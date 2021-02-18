package fragment

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
)

const (
	blockSize          = 128 * 1024
	blockHistorySize   = 4
	minSegmentDuration = 15.5
)

func Main(args []string) {
	flagSet := flag.NewFlagSet("fragment", flag.ExitOnError)
	// TODO
	flagSet.Parse(args)

	if len(flagSet.Args()) < 2 {
		fmt.Printf("USAGE: mp4tool alpha fragment [OPTIONS] INPUT.mp4 OUTPUT.mp4\n")
		flagSet.PrintDefaults()
		return
	}

	ipath := flagSet.Args()[0]
	opath := flagSet.Args()[1]

	input, err := os.Open(ipath)
	if err != nil {
		fmt.Println("Failed to open the input file:", err)
		os.Exit(1)
	}
	defer input.Close()

	output, err := os.Create(opath)
	if err != nil {
		fmt.Println("Failed to create the output file:", err)
		os.Exit(1)
	}
	defer output.Close()

	m := &mp4fragment{
		trackMap: make(map[uint32]*track),
		r:        bufseekio.NewReadSeeker(input, blockSize, blockHistorySize),
		w:        mp4.NewWriter(output),
	}
	if err := m.fragment(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type sample struct {
	keyframe   bool
	chunkStart bool
	descIndex  uint32
	dataOffset uint32
	dataSize   uint32
	baseTime   int64
	timeOffset int64
	duration   uint32
}

type chunk struct {
	dataOffset uint32
	samples    []*sample
}

type track struct {
	chunks    []*chunk
	samples   []*sample
	timescale uint32
	startTime int64
	duration  uint64
	avcc      *mp4.AVCDecoderConfiguration
}

type mp4fragment struct {
	trackMap        map[uint32]*track
	videoTrackID    uint32
	videoTrackCount int
	r               io.ReadSeeker
	w               *mp4.Writer
	probe           *mp4.ProbeInfo
}

func (m *mp4fragment) fragment() error {
	probeInfo, err := mp4.Probe(m.r)
	if err != nil {
		return err
	}
	m.probe = probeInfo

	if err := m.buildInit(); err != nil {
		return err
	}

	if m.videoTrackCount == 0 {
		return errors.New("support plain avc video only")
	} else if m.videoTrackCount >= 2 {
		return errors.New("multiple video track found")
	}

	if err := m.markKeyframes(); err != nil {
		return err
	}

	if err := m.buildFragments(); err != nil {
		return err
	}

	return nil
}

func (m *mp4fragment) buildInit() error {
	var track int
	_, err := mp4.ReadBoxStructure(m.r, func(h *mp4.ReadHandle) (interface{}, error) {
		switch h.BoxInfo.Type {
		case mp4.BoxTypeFtyp():
			if _, err := m.w.StartBox(&h.BoxInfo); err != nil {
				return nil, err
			}
			ibox, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			ftyp := ibox.(*mp4.Ftyp)
			ftyp.AddCompatibleBrand(mp4.BrandISO5()) // required for base-data-offset
			if _, err := mp4.Marshal(m.w, ftyp, h.BoxInfo.Context); err != nil {
				return nil, err
			}
			_, err = m.w.EndBox()
			return nil, err

		case mp4.BoxTypeMvhd():
			if _, err := m.w.StartBox(&h.BoxInfo); err != nil {
				return nil, err
			}
			ibox, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}
			mvhd := ibox.(*mp4.Mvhd)
			if mvhd.GetVersion() == 0 {
				mvhd.DurationV0 = 60034 // FIXME
			} else {
				mvhd.DurationV1 = 60034 // FIXME
			}
			if _, err := mp4.Marshal(m.w, mvhd, h.BoxInfo.Context); err != nil {
				return nil, err
			}
			_, err = m.w.EndBox()
			return nil, err

		case mp4.BoxTypeMoov():
			if _, err := m.w.StartBox(&h.BoxInfo); err != nil {
				return nil, err
			}
			if _, err := h.Expand(); err != nil {
				return nil, err
			}
			if err := m.buildMvex(); err != nil {
				return nil, err
			}
			_, err := m.w.EndBox()
			return nil, err

		case mp4.BoxTypeTrak():
			// FIXME
			if track == 1 {
				return nil, nil
			}
			track++
			return nil, m.buildTrak(&h.BoxInfo)

		default:
			return nil, nil
		}
	})
	return err
}

func (m *mp4fragment) buildTrak(trak *mp4.BoxInfo) error {
	iboxes := make(map[mp4.BoxType]mp4.IBox, 8)
	if _, err := mp4.ReadBoxStructureFromInternal(m.r, trak, func(h *mp4.ReadHandle) (interface{}, error) {
		var copyAll bool
		var write bool
		var collect bool
		var editor func(mp4.IBox) mp4.IBox

		switch h.BoxInfo.Type {
		case mp4.BoxTypeTrak(), mp4.BoxTypeEdts(), mp4.BoxTypeMdia(),
			mp4.BoxTypeMinf(), mp4.BoxTypeStbl(), mp4.BoxTypeStsd(),
			mp4.StrToBoxType("avc1"):
			write = true
		case mp4.BoxTypeMdhd():
			write = true
			collect = true
			editor = func(ibox mp4.IBox) mp4.IBox {
				mdhd := ibox.(*mp4.Mdhd)
				mdhd.DurationV0 = 0
				mdhd.DurationV1 = 0
				return mdhd
			}
		case mp4.StrToBoxType("avcC"):
			write = true
			collect = true
		case mp4.BoxTypeTkhd():
			write = true
			collect = true
			editor = func(ibox mp4.IBox) mp4.IBox {
				tkhd := ibox.(*mp4.Tkhd)
				tkhd.AddFlag(0x000007) // track_enabled=0x000001 track_in_movie=0x000002 track_in_preview=0x000004
				return tkhd
			}
		case mp4.BoxTypeElst():
			write = true
			collect = true
			editor = func(ibox mp4.IBox) mp4.IBox {
				elst := ibox.(*mp4.Elst)
				for i := range elst.Entries {
					if elst.GetVersion() == 0 {
						elst.Entries[i].SegmentDurationV0 = 0
					} else {
						elst.Entries[i].SegmentDurationV1 = 0
					}
				}
				return elst
			}
		case mp4.BoxTypeStts():
			write = true
			collect = true
			editor = func(mp4.IBox) mp4.IBox {
				return &mp4.Stts{}
			}
		case mp4.BoxTypeStsz():
			write = true
			collect = true
			editor = func(mp4.IBox) mp4.IBox {
				return &mp4.Stsz{}
			}
		case mp4.BoxTypeStsc():
			write = true
			collect = true
			editor = func(mp4.IBox) mp4.IBox {
				return &mp4.Stsc{}
			}
		case mp4.BoxTypeStco():
			write = true
			collect = true
			editor = func(mp4.IBox) mp4.IBox {
				return &mp4.Stco{}
			}
		case mp4.BoxTypeCtts():
			collect = true
		default:
			// exclude unexpected children of stbl
			if len(h.Path) >= 2 && h.Path[len(h.Path)-2] == mp4.BoxTypeStbl() {
				return nil, nil
			}
			copyAll = true
		}

		if copyAll {
			// copy all
			return nil, m.w.CopyBox(m.r, &h.BoxInfo)
		}
		// read payload
		ibox, _, err := h.ReadPayload()
		if err != nil {
			return nil, err
		}
		if collect {
			// collect box
			if _, exists := iboxes[h.BoxInfo.Type]; exists {
				return nil, fmt.Errorf("multiple %s", h.BoxInfo.Type)
			}
			iboxes[h.BoxInfo.Type] = ibox
		}
		if editor != nil {
			// edit box payload
			ibox = editor(ibox)
		}
		if write {
			// write box header
			if _, err := m.w.StartBox(&h.BoxInfo); err != nil {
				return nil, err
			}
			// write box payload
			if _, err := mp4.Marshal(m.w, ibox, h.BoxInfo.Context); err != nil {
				return nil, err
			}
		}
		// expand children
		if _, err := h.Expand(); err != nil {
			return nil, err
		}
		if write {
			// rewrite box size
			if _, err = m.w.EndBox(); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}); err != nil {
		return err
	}

	track := &track{
		chunks:  make([]*chunk, 0),
		samples: make([]*sample, 0),
	}

	ibox, ok := iboxes[mp4.BoxTypeTkhd()]
	if !ok {
		return errors.New("tkhd not found")
	}
	tkhd := ibox.(*mp4.Tkhd)

	ibox, ok = iboxes[mp4.BoxTypeMdhd()]
	if !ok {
		return errors.New("mdhd not found")
	}
	mdhd := ibox.(*mp4.Mdhd)
	track.timescale = mdhd.Timescale

	ibox, ok = iboxes[mp4.BoxTypeElst()]
	if !ok {
		return errors.New("elst not found")
	}
	elst := ibox.(*mp4.Elst)
	if elst.EntryCount >= 2 {
		return errors.New("not support multiple entries of elst")
	} else if elst.EntryCount == 1 {
		if elst.GetVersion() == 0 {
			track.startTime = int64(elst.Entries[0].MediaTimeV0)
			track.duration = uint64(elst.Entries[0].SegmentDurationV0)
		} else {
			track.startTime = elst.Entries[0].MediaTimeV1
			track.duration = elst.Entries[0].SegmentDurationV1
		}
	}

	ibox, ok = iboxes[mp4.StrToBoxType("avcC")]
	if ok {
		track.avcc = ibox.(*mp4.AVCDecoderConfiguration)
		m.videoTrackID = tkhd.TrackID
		m.videoTrackCount++
	}

	ibox, ok = iboxes[mp4.BoxTypeStco()]
	if !ok {
		return errors.New("stco not found")
	}
	stco := ibox.(*mp4.Stco)
	for _, offset := range stco.ChunkOffset {
		track.chunks = append(track.chunks, &chunk{
			dataOffset: offset,
		})
	}

	ibox, ok = iboxes[mp4.BoxTypeStts()]
	if !ok {
		return errors.New("stts not found")
	}
	stts := ibox.(*mp4.Stts)
	var baseTime int64
	for _, entry := range stts.Entries {
		for i := uint32(0); i < entry.SampleCount; i++ {
			track.samples = append(track.samples, &sample{
				duration: entry.SampleDelta,
				baseTime: baseTime,
			})
			baseTime += int64(entry.SampleDelta)
		}
	}

	ibox, ok = iboxes[mp4.BoxTypeCtts()]
	if ok {
		ctts := ibox.(*mp4.Ctts)
		var si uint32
		for _, entry := range ctts.Entries {
			for i := uint32(0); i < entry.SampleCount; i++ {
				if ctts.GetVersion() == 0 {
					track.samples[si].timeOffset = int64(entry.SampleOffsetV0)
				} else {
					track.samples[si].timeOffset = int64(entry.SampleOffsetV1)
				}
				si++
			}
		}
	}

	ibox, ok = iboxes[mp4.BoxTypeStsz()]
	if !ok {
		return errors.New("stsz not found")
	}
	stsz := ibox.(*mp4.Stsz)
	for i, size := range stsz.EntrySize {
		track.samples[i].dataSize = size
	}

	ibox, ok = iboxes[mp4.BoxTypeStsc()]
	if !ok {
		return errors.New("stsc not found")
	}
	stsc := ibox.(*mp4.Stsc)
	var currSample uint32
	for sci, entry := range stsc.Entries {
		first := entry.FirstChunk - 1
		var end uint32
		if sci != len(stsc.Entries)-1 {
			end = stsc.Entries[sci+1].FirstChunk - 1
		} else {
			end = uint32(len(track.chunks))
		}
		for ci := first; ci < end; ci++ {
			chunk := track.chunks[ci]
			nextSample := currSample + entry.SamplesPerChunk
			chunk.samples = track.samples[currSample:nextSample]
			offset := chunk.dataOffset
			for _, sample := range chunk.samples {
				sample.dataOffset = offset
				sample.descIndex = entry.SampleDescriptionIndex
				offset += sample.dataSize
			}
			chunk.samples[0].chunkStart = true
			currSample = nextSample
		}
	}

	m.trackMap[tkhd.TrackID] = track
	return nil
}

func (m *mp4fragment) buildMvex() error {
	// start mvex
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMvex()}); err != nil {
		return err
	}

	// mehd
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMehd()}); err != nil {
		return err
	}
	if _, err := mp4.Marshal(m.w, &mp4.Mehd{
		FragmentDurationV0: 60034, // FIXME
	}, mp4.Context{}); err != nil {
		return err
	}
	if _, err := m.w.EndBox(); err != nil {
		return err
	}

	// trex
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTrex()}); err != nil {
		return err
	}
	if _, err := mp4.Marshal(m.w, &mp4.Trex{
		TrackID:                       1, // FIXME
		DefaultSampleDescriptionIndex: 1,
		DefaultSampleDuration:         0,
		DefaultSampleSize:             0,
		DefaultSampleFlags:            0x00000000,
	}, mp4.Context{}); err != nil {
		return err
	}
	if _, err := m.w.EndBox(); err != nil {
		return err
	}

	// end mvex
	_, err := m.w.EndBox()
	return err
}

func (m *mp4fragment) markKeyframes() error {
	track := m.trackMap[m.videoTrackID]
	for _, sample := range track.samples {
		if _, err := m.r.Seek(int64(sample.dataOffset), io.SeekStart); err != nil {
			return err
		}
		nal := make([]byte, sample.dataSize)
		if _, err := io.ReadFull(m.r, nal); err != nil {
			return err
		}
		for len(nal) > 0 {
			lengthSize := track.avcc.LengthSizeMinusOne + 1
			var length uint64
			for i := 0; i < int(lengthSize); i++ {
				length = (length << 8) + uint64(nal[i])
			}
			nalHeader := nal[lengthSize]
			nalType := nalHeader & 0x1f
			if nalType == 5 {
				sample.keyframe = true
				break
			}
			nal = nal[uint64(lengthSize)+length:]
		}
	}
	return nil
}

func (m *mp4fragment) buildFragments() error {
	trackID := m.videoTrackID
	track := m.trackMap[trackID]

	type moofInfo struct {
		time   int64
		offset uint64
	}
	moofs := make([]moofInfo, 0)

	sequenceNumber := uint32(1)
	var sampleBegin int
	var sampleEnd int
	for sampleEnd < len(track.samples) {
		sampleBegin = sampleEnd
		var segmentDuration uint32
		for {
			segmentDuration += track.samples[sampleEnd].duration
			sampleEnd++
			if sampleEnd == len(track.samples) {
				break
			}
			if track.samples[sampleEnd].keyframe &&
				float64(segmentDuration)/float64(track.timescale) >= minSegmentDuration {
				break
			}
		}
		samples := track.samples[sampleBegin:sampleEnd]
		constDuration := true
		for i := 1; i < len(samples); i++ {
			if samples[0].duration != samples[i].duration {
				constDuration = false
				break
			}
		}

		// start moof
		bi, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMoof()})
		if err != nil {
			return err
		}
		moofs = append(moofs, moofInfo{
			time:   samples[0].baseTime,
			offset: bi.Offset,
		})
		if _, err := mp4.Marshal(m.w, &mp4.Moof{}, mp4.Context{}); err != nil {
			return err
		}

		// mfhd
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMfhd()}); err != nil {
			return err
		}
		if _, err := mp4.Marshal(m.w, &mp4.Mfhd{
			SequenceNumber: sequenceNumber,
		}, mp4.Context{}); err != nil {
			return err
		}
		if _, err := m.w.EndBox(); err != nil {
			return err
		}
		sequenceNumber++

		// start traf
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTraf()}); err != nil {
			return err
		}

		// tfhd
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTfhd()}); err != nil {
			return err
		}
		tfhd := &mp4.Tfhd{
			FullBox:                mp4.FullBox{Flags: [3]byte{0x02, 0x00, 0x22}},
			TrackID:                trackID,
			SampleDescriptionIndex: samples[0].descIndex,
			DefaultSampleFlags:     0x01010000, // sample_depends_on=1 sample_is_non_sync_sample=1
		}
		if constDuration {
			tfhd.AddFlag(0x000008)
			tfhd.DefaultSampleDuration = samples[0].duration
		}
		if _, err := mp4.Marshal(m.w, tfhd, mp4.Context{}); err != nil {
			return err
		}
		if _, err := m.w.EndBox(); err != nil {
			return err
		}

		// tfdt
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTfdt()}); err != nil {
			return err
		}
		if _, err := mp4.Marshal(m.w, &mp4.Tfdt{
			FullBox:               mp4.FullBox{Version: 1},
			BaseMediaDecodeTimeV1: uint64(samples[0].baseTime),
		}, mp4.Context{}); err != nil {
			return err
		}
		if _, err := m.w.EndBox(); err != nil {
			return err
		}

		// trun
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTrun()}); err != nil {
			return err
		}
		trun := &mp4.Trun{
			FullBox:          mp4.FullBox{Flags: [3]byte{0x00, 0x0a, 0x05}},
			SampleCount:      uint32(len(samples)),
			FirstSampleFlags: 0x02000000, // sample_depends_on=2
			Entries:          make([]mp4.TrunEntry, 0, len(samples)),
		}
		if !constDuration {
			trun.AddFlag(0x000100)
		}
		for _, sample := range samples {
			trunEntry := mp4.TrunEntry{
				SampleDuration:                sample.duration,
				SampleSize:                    sample.dataSize,
				SampleCompositionTimeOffsetV0: uint32(sample.timeOffset),
				SampleCompositionTimeOffsetV1: int32(sample.timeOffset),
			}
			if sample.timeOffset < 0 {
				trun.Version = 1
			}
			trun.Entries = append(trun.Entries, trunEntry)
		}
		if _, err := mp4.Marshal(m.w, trun, mp4.Context{}); err != nil {
			return err
		}
		trunInfo, err := m.w.EndBox()
		if err != nil {
			return err
		}

		// end traf
		if _, err := m.w.EndBox(); err != nil {
			return err
		}

		// end moof
		moofInfo, err := m.w.EndBox()
		if err != nil {
			return err
		}

		// update trun data-offset
		trun.DataOffset = int32(moofInfo.Size + 8)
		if _, err := trunInfo.SeekToPayload(m.w); err != nil {
			return err
		}
		if _, err := mp4.Marshal(m.w, trun, mp4.Context{}); err != nil {
			return err
		}
		if _, err := m.w.Seek(0, io.SeekEnd); err != nil {
			return err
		}

		// mdat
		if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMdat()}); err != nil {
			return err
		}
		for _, sample := range samples {
			if _, err := m.r.Seek(int64(sample.dataOffset), io.SeekStart); err != nil {
				return err
			}
			if _, err := io.CopyN(m.w, m.r, int64(sample.dataSize)); err != nil {
				return err
			}
		}
		if _, err := m.w.EndBox(); err != nil {
			return err
		}
	}

	// start mfra
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMfra()}); err != nil {
		return err
	}

	// tfra
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTfra()}); err != nil {
		return err
	}
	tfra := &mp4.Tfra{
		TrackID:               trackID,
		LengthSizeOfTrafNum:   0x00,
		LengthSizeOfTrunNum:   0x00,
		LengthSizeOfSampleNum: 0x00,
		NumberOfEntry:         4,
		Entries:               make([]mp4.TfraEntry, 0, len(moofs)),
	}
	for _, moofInfo := range moofs {
		tfra.Entries = append(tfra.Entries, mp4.TfraEntry{
			TimeV0:       uint32(moofInfo.time),
			MoofOffsetV0: uint32(moofInfo.offset),
			TrafNumber:   1,
			TrunNumber:   1,
			SampleNumber: 1,
		})
	}
	if _, err := mp4.Marshal(m.w, tfra, mp4.Context{}); err != nil {
		return err
	}
	if _, err := m.w.EndBox(); err != nil {
		return err
	}

	// mfro
	if _, err := m.w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMfro()}); err != nil {
		return err
	}
	mfro := &mp4.Mfro{}
	if _, err := mp4.Marshal(m.w, mfro, mp4.Context{}); err != nil {
		return err
	}
	mfroInfo, err := m.w.EndBox()
	if err != nil {
		return err
	}

	// end mfra
	mfraInfo, err := m.w.EndBox()
	if err != nil {
		return err
	}

	// update mfro
	mfro.Size = uint32(mfraInfo.Size)
	if _, err := mfroInfo.SeekToPayload(m.w); err != nil {
		return err
	}
	_, err = mp4.Marshal(m.w, mfro, mp4.Context{})
	return err
}
