go-mp4
======

Go library for reading and writing MP4 (ISO Base Media File Format)

## Parse MP4 file

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

## Edit MP4 file

Refer to `mp4tool/edit/mp4edit.go`

## Command Line Tool

```
cd mp4tool
make install

mp4tool -help
```
