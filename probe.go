package mp4

import "io"

type TrackInfo struct {
	TrackID   uint32
	Timescale uint32
}

type SegmentInfo struct {
	TrackID               uint32
	MoofOffset            uint64
	BaseMediaDecodeTime   uint64
	DefaultSampleDuration uint32
	SampleCount           uint32
	Duration              uint32
}

type FraProbeInfo struct {
	Tracks   []TrackInfo
	Segments []SegmentInfo
}

// ProbeFra probes fragmented MP4 file
func ProbeFra(r io.ReadSeeker) (FraProbeInfo, error) {
	info := FraProbeInfo{
		Tracks:   make([]TrackInfo, 0, 8),
		Segments: make([]SegmentInfo, 0, 8),
	}

	boxes, err := ExtractBoxes(r, nil, []BoxPath{
		{BoxTypeMoov(), BoxTypeTrak()},
		{BoxTypeMoof()},
		{BoxTypeMfra(), BoxTypeTfra()},
	})
	if err != nil {
		return info, err
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeTrak() {
			track, err := probeTrak(r, boxes[bi])
			if err != nil {
				return info, err
			}
			info.Tracks = append(info.Tracks, track)
		}
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeMoof() {
			segment, err := probeMoof(r, boxes[bi])
			if err != nil {
				return info, err
			}
			info.Segments = append(info.Segments, segment)
		}
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeTfra() {
			err := probeTfra(r, boxes[bi], &info)
			if err != nil {
				return info, err
			}
		}
	}

	return info, nil
}

func probeTrak(r io.ReadSeeker, bi *BoxInfo) (TrackInfo, error) {
	track := TrackInfo{}

	boxes, err := ExtractBoxes(r, bi, []BoxPath{
		{BoxTypeTkhd()},
		{BoxTypeMdia(), BoxTypeMdhd()},
	})
	if err != nil {
		return track, err
	}

	for bi := range boxes {
		switch boxes[bi].Type {
		case BoxTypeTkhd():
			probeTkhd(r, boxes[bi], &track)
			if err != nil {
				return track, err
			}

		case BoxTypeMdhd():
			probeMdhd(r, boxes[bi], &track)
			if err != nil {
				return track, err
			}
		}
	}

	return track, nil
}

func probeTkhd(r io.ReadSeeker, bi *BoxInfo, info *TrackInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	tkhd := Tkhd{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &tkhd)
	if err != nil {
		return err
	}

	if tkhd.Version == 0 {
		info.TrackID = tkhd.TrackIDV0
	} else {
		info.TrackID = tkhd.TrackIDV1
	}

	return nil
}

func probeMdhd(r io.ReadSeeker, bi *BoxInfo, info *TrackInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	mdhd := Mdhd{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &mdhd)
	if err != nil {
		return err
	}

	info.Timescale = mdhd.Timescale
	return nil
}

func probeMoof(r io.ReadSeeker, bi *BoxInfo) (SegmentInfo, error) {
	segment := SegmentInfo{}

	boxes, err := ExtractBoxes(r, bi, []BoxPath{
		{BoxTypeTraf(), BoxTypeTfhd()},
		{BoxTypeTraf(), BoxTypeTfdt()},
		{BoxTypeTraf(), BoxTypeTrun()},
	})
	if err != nil {
		return segment, err
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeTfhd() {
			probeTfhd(r, boxes[bi], &segment)
			if err != nil {
				return segment, err
			}
		}
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeTfdt() {
			probeTfdt(r, boxes[bi], &segment)
			if err != nil {
				return segment, err
			}
		}
	}

	for bi := range boxes {
		if boxes[bi].Type == BoxTypeTrun() {
			probeTrun(r, boxes[bi], &segment)
			if err != nil {
				return segment, err
			}
		}
	}

	return segment, nil
}

func probeTfhd(r io.ReadSeeker, bi *BoxInfo, segment *SegmentInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	tfhd := Tfhd{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &tfhd)
	if err != nil {
		return err
	}

	segment.TrackID = tfhd.TrackID
	segment.DefaultSampleDuration = tfhd.DefaultSampleDuration

	return nil
}

func probeTfdt(r io.ReadSeeker, bi *BoxInfo, segment *SegmentInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	tfdt := Tfdt{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &tfdt)
	if err != nil {
		return err
	}

	if tfdt.Version == 0 {
		segment.BaseMediaDecodeTime = uint64(tfdt.BaseMediaDecodeTimeV0)
	} else {
		segment.BaseMediaDecodeTime = tfdt.BaseMediaDecodeTimeV1
	}

	return nil
}

func probeTrun(r io.ReadSeeker, bi *BoxInfo, segment *SegmentInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	trun := Trun{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &trun)
	if err != nil {
		return err
	}

	segment.SampleCount = trun.SampleCount

	if trun.CheckFlag(0x000100) {
		segment.Duration = 0
		for ei := range trun.Entries {
			segment.Duration += trun.Entries[ei].SampleDuration
		}
	} else {
		segment.Duration = segment.DefaultSampleDuration * segment.SampleCount
	}

	return nil
}

func probeTfra(r io.ReadSeeker, bi *BoxInfo, info *FraProbeInfo) error {
	if _, err := bi.SeekToPayload(r); err != nil {
		return err
	}

	tfra := Tfra{}
	_, err := Unmarshal(r, bi.Size-bi.HeaderSize, &tfra)
	if err != nil {
		return err
	}

	si := 0
	ei := 0
	for si < len(info.Segments) && ei < len(tfra.Entries) {
		if info.Segments[si].TrackID != tfra.TrackID {
			si++
			continue
		}

		if tfra.Version == 0 {
			info.Segments[si].MoofOffset = uint64(tfra.Entries[ei].MoofOffsetV0)
		} else {
			info.Segments[si].MoofOffset = tfra.Entries[ei].MoofOffsetV1
		}

		si++
		ei++
	}

	return nil
}
