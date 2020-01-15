go-mp4
======

Go library for reading and writing MP4 (ISO Base Media File Format)

# Usage

## Integrate to your Go application

### Parse MP4 file

```
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

### Edit MP4 file

Refer to `mp4tool/edit/mp4edit.go`

### Add box definition

```
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

### Install

```
cd mp4tool
make install

mp4tool -help
```

### Dump

```
mp4tool dump sample.mp4
```

Output example:

```
[moof] Size=504
  [mfhd] Size=16 Version=0 Flags=0x000000 SequenceNumber=1
  [traf] Size=480
    [tfhd] Size=28 Version=0 Flags=0x020038 TrackID=1 DefaultSampleDuration=9000 DefaultSampleSize=33550 DefaultSampleFlags=0x1010000
    [tfdt] Size=20 Version=1 Flags=0x000000 BaseMediaDecodeTimeV1=0
    [trun] Size=424 ... (use -a option to show all)
[mdat] Size=44569 Data=[...] (use -mdat option to expand)
```
