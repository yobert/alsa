package pcm

import (
	"fmt"

	"github.com/yobert/alsa/misc"
	"github.com/yobert/alsa/pcm/state"
)

type Status struct {
	State               state.State
	TriggerTstamp       misc.Timespec
	Tstamp              misc.Timespec
	ApplPtr             misc.Uframes
	HwPtr               misc.Uframes
	Delay               misc.Sframes
	Avail               misc.Uframes
	AvailMax            misc.Uframes
	Overrange           misc.Uframes
	SuspendedState      state.State
	AudioTstampData     uint32
	AudioTstamp         misc.Timespec
	DriverTstamp        misc.Timespec
	AudioTstampAccuracy uint32 // nanoseconds
	Reserved            [20]byte
}

const StatusSize = 152

func (status Status) String() string {
	return fmt.Sprintf("Status{%s, delay %d avail %d max %d over %d suspended %s}",
		status.State, status.Delay, status.Avail, status.AvailMax, status.Overrange, status.SuspendedState)
}
