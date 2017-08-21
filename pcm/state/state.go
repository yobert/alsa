package state

import (
	"fmt"
)

type State int32

const (
	Open State = iota
	Setup
	Prepared
	Running
	Xrun
	Draining
	Paused
	Suspended
	Disconnected
	Last  = Disconnected
	First = Open
)

func (s State) String() string {
	switch s {
	case Open:
		return "Open"
	case Setup:
		return "Setup"
	case Prepared:
		return "Prepared"
	case Running:
		return "Running"
	case Xrun:
		return "Xrun"
	case Draining:
		return "Draining"
	case Paused:
		return "Paused"
	case Suspended:
		return "Suspended"
	case Disconnected:
		return "Disconnected"
	default:
		return fmt.Sprintf("Invalid PCM State (%d)", s)
	}
}
