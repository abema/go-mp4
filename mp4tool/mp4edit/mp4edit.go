package mp4edit

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/abema/go-mp4"
)

const UNoValue = math.MaxUint64

type Values struct {
	BaseMediaDecodeTime uint64
}

type Boxes []string

func (b Boxes) Exists(boxType string) bool {
	for _, t := range b {
		if t == boxType {
			return true
		}
	}
	return false
}

type Config struct {
	values    Values
	dropBoxes Boxes
}

var config Config

func Main(args []string) {
	flagSet := flag.NewFlagSet("edit", flag.ExitOnError)
	flagSet.Uint64Var(&config.values.BaseMediaDecodeTime, "base_media_decode_time", UNoValue, "set new value to base_media_decode_time")
	dropBoxes := flagSet.String("drop", "", "drop boxes")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 2 {
		fmt.Printf("USAGE: mp4tool edit [OPTIONS] INPUT.mp4 OUTPUT.mp4\n")
		return
	}

	config.dropBoxes = strings.Split(*dropBoxes, ",")

	inputPath := flagSet.Args()[0]
	outputPath := flagSet.Args()[1]

	err := editFile(inputPath, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func editFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = mp4.ReadBoxStructure(inputFile, func(h *mp4.ReadHandle) (interface{}, error) {
		// drop
		if config.dropBoxes.Exists(h.BoxInfo.Type.String()) {
			return uint64(0), nil
		}

		offset, err := outputFile.Seek(0, io.SeekEnd)
		if err != nil {
			return nil, err
		}

		// write header
		bi, err := mp4.WriteBoxInfo(outputFile, &h.BoxInfo)
		if err != nil {
			return nil, err
		}

		bi.Size = bi.HeaderSize

		if bi.Type.IsSupported() && bi.Type != mp4.BoxTypeMdat() {
			box, _, err := h.ReadPayload()
			if err != nil {
				return nil, err
			}

			switch bi.Type {
			case mp4.BoxTypeTfdt():
				tfdt := box.(*mp4.Tfdt)
				if config.values.BaseMediaDecodeTime != UNoValue {
					if tfdt.GetVersion() == 0 {
						tfdt.BaseMediaDecodeTimeV0 = uint32(config.values.BaseMediaDecodeTime)
					} else {
						tfdt.BaseMediaDecodeTimeV1 = config.values.BaseMediaDecodeTime
					}
				}
			}

			n, err := mp4.Marshal(outputFile, box, bi.Context)
			if err != nil {
				return nil, err
			}
			bi.Size += n

			// expand all of offsprings
			vals, err := h.Expand()
			if err != nil {
				return nil, err
			}
			for i := range vals {
				n := vals[i].(uint64)
				bi.Size += n
			}

		} else {
			// write all data
			n, err := h.ReadData(outputFile)
			if err != nil {
				return nil, err
			}
			bi.Size += n
		}

		// rewrite header to fix box size
		_, err = outputFile.Seek(offset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		bi2, err := mp4.WriteBoxInfo(outputFile, bi)
		if err != nil {
			return nil, err
		} else if bi2.HeaderSize != bi.HeaderSize {
			return nil, fmt.Errorf("header size has changed: type=%s, before=%d, after=%d", bi.Type.String(), bi.HeaderSize, bi2.HeaderSize)
		}

		return bi.Size, nil
	})
	return err
}
