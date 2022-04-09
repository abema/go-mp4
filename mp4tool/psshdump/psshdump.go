package psshdump

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	mp4 "github.com/abema/go-mp4"
	"github.com/sunfish-shogi/bufseekio"
)

func Main(args []string) {
	if len(args) < 1 {
		println("USAGE: mp4tool psshdump INPUT.mp4")
		os.Exit(1)
	}

	if err := dump(args[0]); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func dump(inputFilePath string) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	r := bufseekio.NewReadSeeker(inputFile, 1024, 4)

	bs, err := mp4.ExtractBoxWithPayload(r, nil, mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypePssh()})
	if err != nil {
		return err
	}

	for i := range bs {
		pssh := bs[i].Payload.(*mp4.Pssh)

		sysid, _ := pssh.StringifyField("SystemID", "", 0, bs[i].Info.Context)

		if _, err := bs[i].Info.SeekToStart(r); err != nil {
			return err
		}
		rawData := make([]byte, bs[i].Info.Size)
		if _, err := io.ReadFull(r, rawData); err != nil {
			return err
		}

		fmt.Printf("%d:\n", i)
		fmt.Printf("  offset: %d\n", bs[i].Info.Offset)
		fmt.Printf("  size: %d\n", bs[i].Info.Size)
		fmt.Printf("  version: %d\n", pssh.Version)
		fmt.Printf("  flags: 0x%x\n", pssh.Flags)
		fmt.Printf("  systemId: %s\n", sysid)
		fmt.Printf("  dataSize: %d\n", pssh.DataSize)
		fmt.Printf("  base64: \"%s\"\n", base64.StdEncoding.EncodeToString(rawData))
		fmt.Println()
	}

	return nil
}
