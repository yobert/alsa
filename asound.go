package main

import (
	"fmt"
	"strings"
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

func fmt_mask(i uint32) string {
	switch i {
	case HwParamAccess:
		return "access"
	case HwParamFormat:
		return "format"
	case HwParamSubFormat:
		return "subfmt"
	default:
		return "invalid mask"
	}
}

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

func fmt_interval(i uint32) string {
	switch i {
	case HwParamSampleBits:
		return "SampleBits"
	case HwParamFrameBits:
		return "FrameBits"
	case HwParamChannels:
		return "Channels"
	case HwParamRate:
		return "Rate"
	case HwParamPeriodTime:
		return "PeriodTime"
	case HwParamPeriodSize:
		return "PeriodSize"
	case HwParamPeriodBytes:
		return "PeriodBytes"
	case HwParamPeriods:
		return "Periods"
	case HwParamBufferTime:
		return "BufferTime"
	case HwParamBufferSize:
		return "BufferSize"
	case HwParamBufferBytes:
		return "BufferBytes"
	case HwParamTickTime:
		return "TickTime"
	default:
		return "invalid interval"
	}
}

func fmt_iflags(f uint32) string {
	r := ""
	if f&0x01 != 0 {
		r += "openmin "
	}
	if f&0x02 != 0 {
		r += "openmax "
	}
	if f&0x04 != 0 {
		r += "integer "
	}
	if f&0x08 != 0 {
		r += "empty "
	}
	return strings.TrimSpace(r)
}

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

func fmt_uint(v uint32) string {
	if v == 0 {
		return "0"
	}
	if v == 0xffffffff {
		return "λ"
	}
	return fmt.Sprintf("0x%08x", v)
}

func (s *HwParams) String() string {
	return s.Diff(&HwParams{})
}

func (s *HwParams) Diff(w *HwParams) string {
	r := ""

	if s.Flags != w.Flags {
		r += fmt.Sprintf("  Flags 0x%x\n", s.Flags)
	}

	for i := range s.Masks {
		for j := range s.Masks[i].Bits {
			if s.Masks[i].Bits[j] != w.Masks[i].Bits[j] {
				v := s.Masks[i].Bits[j]

				sv := ""

				for mv := range s.Masks {
					mvv := uint32(mv + HwParamFirstMask)
					if v&(1<<mvv) != 0 {
						sv += " " + fmt_mask(mvv)
						v ^= (1 << mvv)
					}
				}

				for iv := range s.Intervals {
					ivv := uint32(iv + HwParamFirstInterval)
					if v&(1<<ivv) != 0 {
						sv += " " + fmt_interval(ivv)
						v ^= (1 << ivv)
					}
				}

				r += fmt.Sprintf("  Mask %02d  bits %02d  %-12s %8s%s\n", i, j, fmt_uint(v), fmt_mask(uint32(i)), sv)
			}
		}
	}
	for i := range s.Intervals {

		if s.Intervals[i].Min == w.Intervals[i].Min &&
			s.Intervals[i].Max == w.Intervals[i].Max &&
			s.Intervals[i].Flags == w.Intervals[i].Flags {
			continue
		}

		r += fmt.Sprintf("  Interval %d\t", i)

		it := fmt_interval(uint32(i + HwParamFirstInterval))
		iv := ""

		if s.Intervals[i].Min == 0 && s.Intervals[i].Max == 0xffffffff {
			iv = "0/λ "
		} else {
			iv = fmt.Sprintf("%d/%d ", s.Intervals[i].Min, s.Intervals[i].Max)
		}

		ix := ""
		if s.Intervals[i].Flags != 0 {
			ix = fmt_iflags(
				s.Intervals[i].Flags)
		}
		r += fmt.Sprintf("%-20s %20s %-20s\n", it, iv, ix)
	}
	if s.Rmask != w.Rmask {
		r += "  Rmask  " + fmt_uint(s.Rmask) + "\n"
	}
	if s.Cmask != w.Cmask {
		r += "  Cmask  " + fmt_uint(s.Cmask) + "\n"
	}
	if s.Msbits != w.Msbits {
		r += "  Msbits " + fmt_uint(s.Msbits) + "\n"
	}
	if s.RateNum != w.RateNum || s.RateDen != w.RateDen {
		r += fmt.Sprintf("  Rate   %d/%d\n", s.RateNum, s.RateDen)
	}
	if s.FifoSize != w.FifoSize {
		r += fmt.Sprintf("  FifoSz %d\n", s.FifoSize)
	}

	if r == "" {
		r += "  No changes\n"
	}
	return r
}
