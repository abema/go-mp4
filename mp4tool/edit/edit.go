package edit

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
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
		flagSet.PrintDefaults()
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

	r := bufseekio.NewReadSeeker(inputFile, 128*1024, 4)
	w := mp4.NewWriter(outputFile)
	_, err = mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		if config.dropBoxes.Exists(h.BoxInfo.Type.String()) {
			// drop
			return uint64(0), nil
		}

		if !h.BoxInfo.IsSupportedType() || h.BoxInfo.Type == mp4.BoxTypeMdat() {
			// copy all data
			return nil, w.CopyBox(r, &h.BoxInfo)
		}

		// write header
		_, err := w.StartBox(&h.BoxInfo)
		if err != nil {
			return nil, err
		}
		// read payload
		box, _, err := h.ReadPayload()
		if err != nil {
			return nil, err
		}
		// edit some fields
		switch h.BoxInfo.Type {
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
		// write payload
		if _, err := mp4.Marshal(w, box, h.BoxInfo.Context); err != nil {
			return nil, err
		}
		// expand all of offsprings
		if _, err := h.Expand(); err != nil {
			return nil, err
		}
		// rewrite box size
		_, err = w.EndBox()
		return nil, err
	})
	return err
}
