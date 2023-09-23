package util

import "github.com/abema/go-mp4"

func ShouldHasNoChildren(boxType mp4.BoxType) bool {
	return boxType == mp4.BoxTypeEmsg() ||
		boxType == mp4.BoxTypeEsds() ||
		boxType == mp4.BoxTypeFtyp() ||
		boxType == mp4.BoxTypePssh() ||
		boxType == mp4.BoxTypeCtts() ||
		boxType == mp4.BoxTypeCo64() ||
		boxType == mp4.BoxTypeElst() ||
		boxType == mp4.BoxTypeSbgp() ||
		boxType == mp4.BoxTypeSdtp() ||
		boxType == mp4.BoxTypeStco() ||
		boxType == mp4.BoxTypeStsc() ||
		boxType == mp4.BoxTypeStts() ||
		boxType == mp4.BoxTypeStss() ||
		boxType == mp4.BoxTypeStsz() ||
		boxType == mp4.BoxTypeTfra() ||
		boxType == mp4.BoxTypeTrun()
}
