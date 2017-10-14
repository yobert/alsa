package pcm

import (
	"github.com/ironiridis/alsa/misc"
)

type XferI struct {
	Result misc.Sframes
	Buf    uintptr
	Frames misc.Uframes
}

const XferISize = 24
