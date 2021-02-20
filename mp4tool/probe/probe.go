package probe

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
	"gopkg.in/yaml.v2"
)

func Main(args []string) {
	flagSet := flag.NewFlagSet("fragment", flag.ExitOnError)
	format := flagSet.String("format", "json", "output format (yaml|json)")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 1 {
		fmt.Printf("USAGE: mp4tool beta probe [OPTIONS] INPUT.mp4\n")
		flagSet.PrintDefaults()
		return
	}

	ipath := flagSet.Args()[0]
	input, err := os.Open(ipath)
	if err != nil {
		fmt.Println("Failed to open the input file:", err)
		os.Exit(1)
	}
	defer input.Close()

	r := bufseekio.NewReadSeeker(input, 1024, 4)
	rep, err := buildReport(r)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	switch *format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err = enc.Encode(rep)
	default:
		err = yaml.NewEncoder(os.Stdout).Encode(rep)
	}
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type report struct {
	MajorBrand       string   `yaml:"major_brand"`
	MinorVersion     uint32   `yaml:"minor_version"`
	CompatibleBrands []string `yaml:"compatible_brands"`
	FastStart        bool     `yaml:"fast_start"`
	Timescale        uint32   `yaml:"timescale"`
	Duration         uint64   `yaml:"duration"`
	DurationSeconds  float32  `yaml:"duration_seconds"`
	Tracks           []*track `yaml:"tracks"`
}

type track struct {
	TrackID         uint32  `yaml:"track_id"`
	Timescale       uint32  `yaml:"timescale"`
	Duration        uint64  `yaml:"duration"`
	DurationSeconds float32 `yaml:"duration_seconds"`
	Codec           string  `yaml:"codec"`
	Encrypted       bool    `yaml:"encrypted"`
	Width           uint16  `json:",omitempty" yaml:"width,omitempty"`
	Height          uint16  `json:",omitempty" yaml:"height,omitempty"`
	SampleNum       int     `json:",omitempty" yaml:"sample_num,omitempty"`
	ChunkNum        int     `json:",omitempty" yaml:"chunk_num,omitempty"`
	IDRFrameNum     int     `json:",omitempty" yaml:"idr_frame_num,omitempty"`
	Bitrate         uint64  `json:",omitempty" yaml:"bitrate,omitempty"`
	MaxBitrate      uint64  `json:",omitempty" yaml:"max_bitrate,omitempty"`
}

func buildReport(r io.ReadSeeker) (*report, error) {
	info, err := mp4.Probe(r)
	if err != nil {
		return nil, err
	}

	rep := &report{
		MajorBrand:       string(info.MajorBrand[:]),
		MinorVersion:     info.MinorVersion,
		CompatibleBrands: make([]string, 0, len(info.CompatibleBrands)),
		FastStart:        info.FastStart,
		Timescale:        info.Timescale,
		Duration:         info.Duration,
		DurationSeconds:  float32(info.Duration) / float32(info.Timescale),
		Tracks:           make([]*track, 0, len(info.Tracks)),
	}
	for _, brand := range info.CompatibleBrands {
		rep.CompatibleBrands = append(rep.CompatibleBrands, string(brand[:]))
	}
	for _, tr := range info.Tracks {
		bitrate := tr.Samples.GetBitrate(tr.Timescale)
		maxBitrate := tr.Samples.GetMaxBitrate(tr.Timescale, uint64(tr.Timescale))
		if bitrate == 0 || maxBitrate == 0 {
			bitrate = info.Segments.GetBitrate(tr.TrackID, tr.Timescale)
			maxBitrate = info.Segments.GetMaxBitrate(tr.TrackID, tr.Timescale)
		}
		t := &track{
			TrackID:         tr.TrackID,
			Timescale:       tr.Timescale,
			Duration:        tr.Duration,
			DurationSeconds: float32(tr.Duration) / float32(tr.Timescale),
			Encrypted:       tr.Encrypted,
			Bitrate:         bitrate,
			MaxBitrate:      maxBitrate,
			SampleNum:       len(tr.Samples),
			ChunkNum:        len(tr.Chunks),
		}
		switch tr.Codec {
		case mp4.CodecAVC1:
			if tr.AVC != nil {
				t.Codec = fmt.Sprintf("avc1.%02X%02X%02X",
					tr.AVC.Profile,
					tr.AVC.ProfileCompatibility,
					tr.AVC.Level,
				)
				t.Width = tr.AVC.Width
				t.Height = tr.AVC.Height
			} else {
				t.Codec = "avc1"
			}
			idxs, err := mp4.FindIDRFrames(r, tr)
			if err != nil {
				return nil, err
			}
			t.IDRFrameNum = len(idxs)
		case mp4.CodecMP4A:
			if tr.MP4A == nil || tr.MP4A.OTI == 0 {
				t.Codec = "mp4a"
			} else if tr.MP4A.AudOTI == 0 {
				t.Codec = fmt.Sprintf("mp4a.%X", tr.MP4A.OTI)
			} else {
				t.Codec = fmt.Sprintf("mp4a.%X.%d", tr.MP4A.OTI, tr.MP4A.AudOTI)
			}
		default:
			t.Codec = "unknown"
		}
		rep.Tracks = append(rep.Tracks, t)
	}
	return rep, nil
}
