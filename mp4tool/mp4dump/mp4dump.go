package mp4dump

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/abema/go-mp4"
)

func Main(args []string) {
	flagSet := flag.NewFlagSet("dump", flag.ExitOnError)
	full := flagSet.String("full", "", "Show full content of specified box types\nFor example: -full free,ctts,stts")
	showAll := flagSet.Bool("a", false, "Deprecated: see -full")
	mdat := flagSet.Bool("mdat", false, "Deprecated: see -full")
	free := flagSet.Bool("free", false, "Deprecated: see -full")
	offset := flagSet.Bool("offset", false, "Show offset of box")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 1 {
		fmt.Printf("USAGE: mp4tool dump [OPTIONS] INPUT.mp4\n")
		return
	}

	fpath := flagSet.Args()[0]

	fmap := make(map[string]struct{})
	for _, tname := range strings.Split(*full, ",") {
		fmap[tname] = struct{}{}
	}
	if *mdat {
		fmap["mdat"] = struct{}{}
	}
	if *free {
		fmap["free"] = struct{}{}
		fmap["skip"] = struct{}{}
	}

	m := &mp4dump{
		full:    fmap,
		showAll: *showAll,
		offset:  *offset,
	}
	err := m.dumpFile(fpath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

type mp4dump struct {
	full    map[string]struct{}
	showAll bool
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

		_, full := m.full[h.BoxInfo.Type.String()]
		if !full &&
			(h.BoxInfo.Type == mp4.BoxTypeMdat() ||
				h.BoxInfo.Type == mp4.BoxTypeFree() ||
				h.BoxInfo.Type == mp4.BoxTypeSkip()) {
			fmt.Printf(" Data=[...] (use \"-full %s\" to show all)\n", h.BoxInfo.Type)
			return nil, nil
		}
		full = full || m.showAll

		// supported box type
		if h.BoxInfo.Type.IsSupported() {
			if !full && h.BoxInfo.Size-h.BoxInfo.HeaderSize >= 64 &&
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
				fmt.Printf(" ... (use \"-full %s\" to show all)\n", h.BoxInfo.Type)
				return nil, nil
			}

			box, n, err := h.ReadPayload()
			if err != nil {
				panic(err)
			}

			if !full && n >= 64 {
				fmt.Printf(" ... (use \"-full %s\" to show all)\n", h.BoxInfo.Type)
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
		if full {
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
			fmt.Printf(" Data=[...] (use \"-full %s\" to show all)\n", h.BoxInfo.Type)
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
