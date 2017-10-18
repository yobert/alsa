package pcm

import (
	"fmt"

	"github.com/yobert/alsa/alsatype"
	"github.com/yobert/alsa/pcm/state"
)

type Status struct {
	State               state.State
	TriggerTstamp       alsatype.Timespec
	Tstamp              alsatype.Timespec
	ApplPtr             alsatype.Uframes
	HwPtr               alsatype.Uframes
	Delay               alsatype.Sframes
	Avail               alsatype.Uframes
	AvailMax            alsatype.Uframes
	Overrange           alsatype.Uframes
	SuspendedState      state.State
	AudioTstampData     uint32
	AudioTstamp         alsatype.Timespec
	DriverTstamp        alsatype.Timespec
	AudioTstampAccuracy uint32 // nanoseconds
	Reserved            [20]byte
}

const StatusSize = 152

func (status Status) String() string {
	return fmt.Sprintf("Status{%s, delay %d avail %d max %d over %d suspended %s}",
		status.State, status.Delay, status.Avail, status.AvailMax, status.Overrange, status.SuspendedState)
}
