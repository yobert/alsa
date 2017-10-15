package pcm

import (
	"github.com/yobert/alsa/alsatype"
)

type XferI struct {
	Result alsatype.Sframes
	Buf    uintptr
	Frames alsatype.Uframes
}

const XferISize = 24
