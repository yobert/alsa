package alsa

import (
	"fmt"
	"strings"

	"github.com/yobert/alsa/pcm"
)

type param uint32

const (
	paramAccess    param = 0
	paramFormat    param = 1
	paramSubformat param = 2
	paramFirstMask param = paramAccess
	paramLastMask  param = paramSubformat

	paramSampleBits    param = 8
	paramFrameBits     param = 9
	paramChannels      param = 10
	paramRate          param = 11
	paramPeriodTime    param = 12
	paramPeriodSize    param = 13
	paramPeriodBytes   param = 14
	paramPeriods       param = 15
	paramBufferTime    param = 16
	paramBufferSize    param = 17
	paramBufferBytes   param = 18
	paramTickTime      param = 19
	paramFirstInterval param = paramSampleBits
	paramLastInterval  param = paramTickTime
)

type Flags uint32

const (
	OpenMin Flags = 1 << iota
	OpenMax
	Integer
	Empty
)

func (f Flags) String() string {
	r := ""
	if f&OpenMin != 0 {
		r += "OpenMin "
	}
	if f&OpenMax != 0 {
		r += "OpenMax "
	}
	if f&Integer != 0 {
		r += "Integer "
	}
	if f&Empty != 0 {
		r += "Empty "
	}
	return strings.TrimSpace(r)
}

func (p param) IsMask() bool {
	return p >= paramFirstMask && p <= paramLastMask
}
func (p param) IsInterval() bool {
	return p >= paramFirstInterval && p <= paramLastInterval
}
func (p param) String() string {
	if p.IsMask() {
		return "≡" + p.name()
	}
	if p.IsInterval() {
		return "±" + p.name()
	}
	return "Invalid"
}
func (p param) name() string {
	switch p {
	case paramAccess:
		return "Access"
	case paramFormat:
		return "Format"
	case paramSubformat:
		return "Subfmt"

	case paramSampleBits:
		return "SampleBits"
	case paramFrameBits:
		return "FrameBits"
	case paramChannels:
		return "Channels"
	case paramRate:
		return "Rate"
	case paramPeriodTime:
		return "PeriodTime"
	case paramPeriodSize:
		return "PeriodSize"
	case paramPeriodBytes:
		return "PeriodBytes"
	case paramPeriods:
		return "Periods"
	case paramBufferTime:
		return "BufferTime"
	case paramBufferSize:
		return "BufferSize"
	case paramBufferBytes:
		return "BufferBytes"
	case paramTickTime:
		return "TickTime"
	default:
		return "Invalid"
	}
}

func get_status(fd uintptr) error {
	var status pcm.Status
	err := ioctl(fd, ioctl_encode(cmdRead, pcm.StatusSize, cmdPCMStatus), &status)
	if err != nil {
		return err
	}
	fmt.Println(status)
	return nil
}
