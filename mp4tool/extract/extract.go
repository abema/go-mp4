package extract

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/abema/go-mp4"
	"github.com/abema/go-mp4/mp4tool/util"
	"github.com/sunfish-shogi/bufseekio"
)

const (
	blockSize        = 128 * 1024
	blockHistorySize = 4
)

func Main(args []string) {
	flagSet := flag.NewFlagSet("extract", flag.ExitOnError)
	flagSet.Parse(args)

	if len(flagSet.Args()) < 2 {
		println("USAGE: mp4tool extract [OPTIONS] BOX_TYPE INPUT.mp4")
		flagSet.PrintDefaults()
		os.Exit(1)
	}

	boxType := flagSet.Args()[0]
	inputPath := flagSet.Args()[1]

	if len(boxType) != 4 {
		println("Error:", "invalid argument:", boxType)
		println("BOX_TYPE must be 4 characters.")
		os.Exit(1)
	}

	input, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer input.Close()

	r := bufseekio.NewReadSeeker(input, blockSize, blockHistorySize)
	if err := extract(r, mp4.StrToBoxType(boxType)); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func extract(r io.ReadSeeker, boxType mp4.BoxType) error {
	_, err := mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		if h.BoxInfo.Type == boxType {
			h.BoxInfo.SeekToStart(r)
			if _, err := io.CopyN(os.Stdout, r, int64(h.BoxInfo.Size)); err != nil {
				return nil, err
			}
		}
		if !h.BoxInfo.IsSupportedType() {
			return nil, nil
		}
		if h.BoxInfo.Size >= 256 && util.ShouldHasNoChildren(h.BoxInfo.Type) {
			return nil, nil
		}
		_, err := h.Expand()
		if err == mp4.ErrUnsupportedBoxVersion {
			return nil, nil
		}
		return nil, err
	})
	return err
}
