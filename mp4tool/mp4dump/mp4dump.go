package mp4dump

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/abema/go-mp4"
)

func Main(args []string) {
	flagSet := flag.NewFlagSet("dump", flag.ExitOnError)
	showAll := flagSet.Bool("a", false, "show all contents, excluding mdat")
	mdat := flagSet.Bool("mdat", false, "show content of mdat")
	offset := flagSet.Bool("offset", false, "show offset of box")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 1 {
		fmt.Printf("USAGE: mp4tool dump [OPTIONS] INPUT.mp4\n")
		return
	}

	fpath := flagSet.Args()[0]

	m := &mp4dump{
		showAll: *showAll,
		mdat:    *mdat,
		offset:  *offset,
	}
	err := m.dumpFile(fpath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type mp4dump struct {
	showAll bool
	mdat    bool
	offset  bool
}

func (m *mp4dump) dumpFile(fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	return m.dump(file)
}

func (m *mp4dump) dump(r io.ReadSeeker) error {
	_, err := mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		printIndent(len(h.Path) - 1)

		fmt.Printf("[%s]", h.BoxInfo.Type.String())
		if !h.BoxInfo.Type.IsSupported() {
			fmt.Printf(" (unsupported box type)")
		}
		if m.offset {
			fmt.Printf(" Offset=%d", h.BoxInfo.Offset)
		}
		fmt.Printf(" Size=%d", h.BoxInfo.Size)

		showAll := m.showAll
		switch h.BoxInfo.Type {
		case mp4.BoxTypeMdat():
			if m.mdat {
				showAll = true
			} else {
				fmt.Printf(" Data=[...] (use -mdat option to expand)\n")
				return nil, nil
			}
		}

		// supported box type
		if h.BoxInfo.Type.IsSupported() {
			if !showAll && h.BoxInfo.Size-h.BoxInfo.HeaderSize >= 64 &&
				(h.BoxInfo.Type == mp4.BoxTypeEmsg() ||
					h.BoxInfo.Type == mp4.BoxTypeEsds() ||
					h.BoxInfo.Type == mp4.BoxTypeFtyp() ||
					h.BoxInfo.Type == mp4.BoxTypePssh() ||
					h.BoxInfo.Type == mp4.BoxTypeStco() ||
					h.BoxInfo.Type == mp4.BoxTypeStsc() ||
					h.BoxInfo.Type == mp4.BoxTypeStts() ||
					h.BoxInfo.Type == mp4.BoxTypeStsz() ||
					h.BoxInfo.Type == mp4.BoxTypeTfra() ||
					h.BoxInfo.Type == mp4.BoxTypeTrun()) {
				fmt.Printf(" ... (use -a option to show all)\n")
				return nil, nil
			}

			box, n, err := h.ReadPayload()
			if err != nil {
				panic(err)
			}

			if !showAll && n >= 64 {
				fmt.Printf(" ... (use -a option to show all)\n")
			} else {
				str, err := mp4.Stringify(box)
				if err != nil {
					panic(err)
				}
				fmt.Printf(" %s\n", str)
			}

			_, err = h.Expand()
			return nil, err
		}

		// unsupported box type
		if showAll {
			buf := bytes.NewBuffer(make([]byte, 0, h.BoxInfo.Size-h.BoxInfo.HeaderSize))
			if _, err := h.ReadData(buf); err != nil {
				panic(err)
			}
			fmt.Printf(" Data=[")
			for i, d := range buf.Bytes() {
				if i != 0 {
					fmt.Printf(" ")
				}
				fmt.Printf("0x%02x", d)
			}
			fmt.Println("]")
		} else {
			fmt.Printf(" Data=[...] (use -a option to show all)\n")
		}
		return nil, nil
	})
	return err
}

func printIndent(depth int) {
	for i := 0; i < depth; i++ {
		fmt.Printf("  ")
	}
}
