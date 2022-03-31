package dump

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	indentSize       = 2
	blockSize        = 128 * 1024
	blockHistorySize = 4
)

var terminalWidth = 180

func init() {
	if width, _, err := terminal.GetSize(syscall.Stdin); err == nil {
		terminalWidth = width
	}
}

func Main(args []string) {
	flagSet := flag.NewFlagSet("dump", flag.ExitOnError)
	full := flagSet.String("full", "", "Show full content of specified box types\nFor example: -full free,ctts,stts")
	showAll := flagSet.Bool("a", false, "Show full content of boxes excepting mdat, free and styp")
	mdat := flagSet.Bool("mdat", false, "Deprecated: use \"-full mdat\"")
	free := flagSet.Bool("free", false, "Deprecated: use \"-full free,styp\"")
	offset := flagSet.Bool("offset", false, "Show offset of box")
	hex := flagSet.Bool("hex", false, "Use hex for size and offset")
	flagSet.Parse(args)

	if len(flagSet.Args()) < 1 {
		fmt.Fprintf(os.Stderr, "USAGE: mp4tool dump [OPTIONS] INPUT.mp4\n")
		flagSet.PrintDefaults()
		os.Exit(1)
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
		hex:     *hex,
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
	hex     bool
}

func (m *mp4dump) dumpFile(fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer file.Close()

	return m.dump(bufseekio.NewReadSeeker(file, blockSize, blockHistorySize))
}

func (m *mp4dump) dump(r io.ReadSeeker) error {
	_, err := mp4.ReadBoxStructure(r, func(h *mp4.ReadHandle) (interface{}, error) {
		line := bytes.NewBuffer(make([]byte, 0, terminalWidth))

		printIndent(line, len(h.Path)-1)

		fmt.Fprintf(line, "[%s]", h.BoxInfo.Type.String())
		if !h.BoxInfo.IsSupportedType() {
			fmt.Fprintf(line, " (unsupported box type)")
		}
		sizeFormat := "%d"
		if m.hex {
			sizeFormat = "0x%x"
		}
		if m.offset {
			fmt.Fprintf(line, " Offset="+sizeFormat, h.BoxInfo.Offset)
		}
		fmt.Fprintf(line, " Size="+sizeFormat, h.BoxInfo.Size)

		_, full := m.full[h.BoxInfo.Type.String()]
		if !full &&
			(h.BoxInfo.Type == mp4.BoxTypeMdat() ||
				h.BoxInfo.Type == mp4.BoxTypeFree() ||
				h.BoxInfo.Type == mp4.BoxTypeSkip()) {
			fmt.Fprintf(line, " Data=[...] (use \"-full %s\" to show all)", h.BoxInfo.Type)
			fmt.Println(line.String())
			return nil, nil
		}
		full = full || m.showAll

		// supported box type
		if h.BoxInfo.IsSupportedType() {
			if !full && h.BoxInfo.Size-h.BoxInfo.HeaderSize >= 64 &&
				(h.BoxInfo.Type == mp4.BoxTypeEmsg() ||
					h.BoxInfo.Type == mp4.BoxTypeEsds() ||
					h.BoxInfo.Type == mp4.BoxTypeFtyp() ||
					h.BoxInfo.Type == mp4.BoxTypePssh() ||
					h.BoxInfo.Type == mp4.BoxTypeCtts() ||
					h.BoxInfo.Type == mp4.BoxTypeCo64() ||
					h.BoxInfo.Type == mp4.BoxTypeElst() ||
					h.BoxInfo.Type == mp4.BoxTypeSbgp() ||
					h.BoxInfo.Type == mp4.BoxTypeSdtp() ||
					h.BoxInfo.Type == mp4.BoxTypeStco() ||
					h.BoxInfo.Type == mp4.BoxTypeStsc() ||
					h.BoxInfo.Type == mp4.BoxTypeStts() ||
					h.BoxInfo.Type == mp4.BoxTypeStss() ||
					h.BoxInfo.Type == mp4.BoxTypeStsz() ||
					h.BoxInfo.Type == mp4.BoxTypeTfra() ||
					h.BoxInfo.Type == mp4.BoxTypeTrun()) {
				fmt.Fprintf(line, " ... (use \"-full %s\" to show all)", h.BoxInfo.Type)
				fmt.Println(line.String())
				return nil, nil
			}

			box, _, err := h.ReadPayload()
			if err != mp4.ErrUnsupportedBoxVersion {
				if err != nil {
					return nil, err
				}
				str, err := mp4.Stringify(box, h.BoxInfo.Context)
				if err != nil {
					return nil, err
				}
				if !full && line.Len()+len(str)+2 > terminalWidth {
					fmt.Fprintf(line, " ... (use \"-full %s\" to show all)", h.BoxInfo.Type)
				} else if str != "" {
					if h.BoxInfo.Type == mp4.BoxTypeMdat() {
						p := strings.Index(str, "[")
						q := strings.Index(str, "]")
						hexset := strings.Split(str[p + 1:q], ", ")
						hexseq := ""
						hexln  := ""
						for cnt, val := range hexset {
							ret, _ := strconv.ParseUint(val, 0, 16)
							if cnt % 16 == 0 {
								hexseq = fmt.Sprintf("    %02x", ret)
							} else if cnt % 16 == 15 {
								hexseq = fmt.Sprintf("%s %02x\n", hexseq, ret)
								hexln += hexseq
							} else {
								hexseq = fmt.Sprintf("%s %02x", hexseq, ret)
							}
						}
						hexln += hexseq
						fmt.Fprintf(line, " %s\n%s]", str[:p + 1], hexln)
					} else {
						fmt.Fprintf(line, " %s", str)
					}

				}
				fmt.Println(line.String())
				_, err = h.Expand()
				return nil, err
			}
			fmt.Fprintf(line, " (unsupported box version)")
		}

		// unsupported box type
		if full {
			buf := bytes.NewBuffer(make([]byte, 0, h.BoxInfo.Size-h.BoxInfo.HeaderSize))
			if _, err := h.ReadData(buf); err != nil {
				return nil, err
			}
			fmt.Fprintf(line, " Data=[")
			for i, d := range buf.Bytes() {
				if i != 0 {
					fmt.Fprintf(line, " ")
				}
				fmt.Fprintf(line, "0x%02x", d)
			}
			fmt.Fprintf(line, "]")
		} else {
			fmt.Fprintf(line, " Data=[...] (use \"-full %s\" to show all)", h.BoxInfo.Type)
		}
		fmt.Println(line.String())
		return nil, nil
	})
	return err
}

func printIndent(w io.Writer, depth int) {
	for i := 0; i < depth*indentSize; i++ {
		fmt.Fprintf(w, " ")
	}
}
