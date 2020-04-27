go-mp4
------

go-mp4 is Go library for reading and writing MP4 (ISO Base Media File Format).

## Integration with your Go application

You can parse MP4 file as follows:

```go
_, err := mp4.ReadBoxStructure(file, func(h *mp4.ReadHandle) (interface{}, error) {
	fmt.Println("depth", len(h.Path))

	// Box Type (e.g. "mdhd", "tfdt", "mdat")
	fmt.Println(h.BoxInfo.Type.String())

	// Box Size
	fmt.Println(h.BoxInfo.Size)

	// Payload
	if h.BoxInfo.Type.IsSupported() {
		box, _, _ := h.ReadPayload()
		fmt.Println(mp4.Stringify(box))
	}

	// Expands sibling boxes
	return h.Expand()
})
```

You can create additional box definition as follows:

```go
func BoxTypeXxxx() BoxType { return StrToBoxType("xxxx") }

func init() {
	AddBoxDef(&Xxxx{}, 0)
}

type Xxxx struct {
	FullBox  `mp4:"extend"`
	UI32      uint32 `mp4:"size=32"`
	ByteArray []byte `mp4:"size=8,len=dynamic"`
}

func (*Xxxx) GetType() BoxType {
	return BoxTypeXxxx()
}
```

## Command Line Tool

Install mp4tool as follows:

```sh
go get github.com/abema/go-mp4/mp4tool

mp4tool -help
```

For example, `mp4tool dump MP4_FILE_NAME` command prints MP4 box tree as follows:

```
[moof] Size=504
  [mfhd] Size=16 Version=0 Flags=0x000000 SequenceNumber=1
  [traf] Size=480
    [tfhd] Size=28 Version=0 Flags=0x020038 TrackID=1 DefaultSampleDuration=9000 DefaultSampleSize=33550 DefaultSampleFlags=0x1010000
    [tfdt] Size=20 Version=1 Flags=0x000000 BaseMediaDecodeTimeV1=0
    [trun] Size=424 ... (use -a option to show all)
[mdat] Size=44569 Data=[...] (use -mdat option to expand)
```
