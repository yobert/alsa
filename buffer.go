package alsa

import (
	"fmt"
)

type BufferFormat struct {
	SampleFormat FormatType
	Rate         int
	Channels     int
}

type Buffer struct {
	Format BufferFormat
	Data   []byte
}

func (bp BufferFormat) String() string {
	return fmt.Sprintf("%d channels, %d hz, %v", bp.Channels, bp.Rate, bp.SampleFormat)
}
