package pcm

import (
	"github.com/yobert/alsa/misc"
)

type XferI struct {
	Result misc.Sframes
	Buf    uintptr
	Frames misc.Uframes
}

const XferISize = 24
