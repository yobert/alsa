package main

import (
	"fmt"
)

const (
	CmdWrite = 1
	CmdRead  = 2

	CmdPCMInfo              uintptr = 0x4101
	CmdPCMVersion           uintptr = 0x4100
	CmdPCMTimestamp         uintptr = 0x4102
	CmdPCMTimestampType     uintptr = 0x4103
	CmdPCMHwRefine          uintptr = 0x4110
	CmdControlVersion       uintptr = 0x5500
	CmdControlCardInfo      uintptr = 0x5501
	CmdControlPCMNextDevice uintptr = 0x5530
	CmdControlPCMInfo       uintptr = 0x5531
)

const (
	PCMTimestampTypeGettimeofday = iota
	PCMTimestampTypeMonotonic
	PCMTimestampTypeMonotonicRaw
	PCMTimestampTypeLast
)

type CardInfo struct {
	Card       int32
	_          int32
	ID         [16]byte
	Driver     [16]byte
	Name       [32]byte
	LongName   [80]byte
	_          [16]byte
	MixerName  [80]byte
	Components [128]byte
}

func (s CardInfo) String() string {
	return fmt.Sprintf("Card %d %#v", s.Card, gstr(s.Name[:]))
}

type PVersion uint32

func (v PVersion) Major() int {
	return int(v >> 16 & 0xffff)
}
func (v PVersion) Minor() int {
	return int(v >> 8 & 0xff)
}
func (v PVersion) Patch() int {
	return int(v & 0xff)
}
func (v PVersion) String() string {
	return fmt.Sprintf("Protocol %d.%d.%d", v.Major(), v.Minor(), v.Patch())
}

type PCMInfo struct {
	Device          uint32
	Subdevice       uint32
	Stream          int32
	Card            int32
	_               [64]byte
	Name            [80]byte
	Subname         [32]byte
	DevClass        int32
	DevSubclass     int32
	SubdevicesCount uint32
	SubdevicesAvail uint32
	SyncID          [16]byte
	_               [64]byte
}

func (s PCMInfo) String() string {
	return fmt.Sprintf("PCM device %d sub %d stream %d card %d %#v (%d / %d) cls %d subcls %d", s.Device, s.Subdevice, s.Stream, s.Card, gstr(s.Name[:]), s.SubdevicesCount, s.SubdevicesAvail, s.DevClass, s.DevSubclass)
}

const (
	MaskMax = 256
)

type Mask struct {
	Bits [(MaskMax + 31) / 32]uint32
}

const (
	HwParamAccess    = 0
	HwParamFormat    = 1
	HwParamSubFormat = 2
	HwParamFirstMask = HwParamAccess
	HwParamLastMask  = HwParamSubFormat
)

const (
	HwParamSampleBits    = 8
	HwParamFrameBits     = 9
	HwParamChannels      = 10
	HwParamRate          = 11
	HwParamPeriodTime    = 12
	HwParamPeriodSize    = 13
	HwParamPeriodBytes   = 14
	HwParamPeriods       = 15
	HwParamBufferTime    = 16
	HwParamBufferSize    = 17
	HwParamBufferBytes   = 18
	HwParamTickTime      = 19
	HwParamFirstInterval = HwParamSampleBits
	HwParamLastInterval  = HwParamTickTime
)

type Interval struct {
	Min, Max uint32
	Flags    uint32 // bitfield: openmin openmax integer empty
}

func (i Interval) String() string {
	return fmt.Sprintf("Interval(%d/%d 0x%x)", i.Min, i.Max, i.Flags)
}

type UframesType uint64
type SframesType int64

type HwParams struct {
	Flags     uint32
	Masks     [HwParamLastMask - HwParamFirstMask + 1]Mask
	_         [5]Mask
	Intervals [HwParamLastInterval - HwParamFirstInterval + 1]Interval
	_         [9]Interval
	Rmask     uint32
	Cmask     uint32
	Info      uint32
	Msbits    uint32
	RateNum   uint32
	RateDen   uint32
	FifoSize  UframesType
	_         [64]byte
}
