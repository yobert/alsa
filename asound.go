package main

import (
	"fmt"

	"github.com/yobert/alsa/misc"
	"github.com/yobert/alsa/pcm/state"
)

const (
	CmdWrite = 1
	CmdRead  = 2

	CmdPCMInfo          uintptr = 0x4101
	CmdPCMVersion       uintptr = 0x4100
	CmdPCMTimestamp     uintptr = 0x4102
	CmdPCMTimestampType uintptr = 0x4103
	CmdPCMHwRefine      uintptr = 0x4110
	CmdPCMHwParams      uintptr = 0x4111
	CmdPCMSwParams      uintptr = 0x4113
	CmdPCMStatus        uintptr = 0x4120

	CmdPCMPrepare uintptr = 0x4140
	CmdPCMReset   uintptr = 0x4141
	CmdPCMStart   uintptr = 0x4142
	CmdPCMDrop    uintptr = 0x4143
	CmdPCMDrain   uintptr = 0x4144
	CmdPCMPause   uintptr = 0x4145 // int
	CmdPCMRewind  uintptr = 0x4146 // snd_pcm_uframes_t
	CmdPCMResume  uintptr = 0x4147
	CmdPCMXrun    uintptr = 0x4148
	CmdPCMForward uintptr = 0x4149

	CmdPCMWriteIFrames uintptr = 0x4150 // snd_xferi
	CmdPCMReadIFrames  uintptr = 0x4151 // snd_xferi
	CmdPCMWriteNFrames uintptr = 0x4152 // snd_xfern
	CmdPCMReadNFrames  uintptr = 0x4153 // snd_xfern

	CmdPCMLink   uintptr = 0x4160 // int
	CmdPCMUnlink uintptr = 0x4161

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

const (
	MapShared     = 0x00000001
	OffsetData    = 0x00000000
	OffsetStatus  = 0x80000000
	OffsetControl = 0x81000000
)

type AccessType int

const (
	MmapInterleaved AccessType = iota
	MmapNonInterleaved
	MmapComplex
	RWInterleaved
	RWNonInterleaved
	AccessTypeLast  = RWNonInterleaved
	AccessTypeFirst = MmapInterleaved
)

func (a AccessType) String() string {
	switch a {
	case MmapInterleaved:
		return "MmapInterleaved"
	case MmapNonInterleaved:
		return "MmapNonInterleaved"
	case MmapComplex:
		return "MmapComplex"
	case RWInterleaved:
		return "RWInterleaved"
	case RWNonInterleaved:
		return "RWNonInterleaved"
	default:
		return fmt.Sprintf("Invalid AccessType (%d)", a)
	}
}

type FormatType int

const (
	Unknown FormatType = -1
)
const (
	S8     FormatType = iota // 0
	U8                       // 1
	S16_LE                   // 2
	S16_BE                   // 3
	U16_LE                   // 4
	U16_BE                   // 5
	S24_LE                   // 6
	S24_BE                   // 7
	U24_LE                   // 8
	U24_BE                   // 9
	S32_LE
	S32_BE
	U32_LE
	U32_BE
	FLOAT_LE
	FLOAT_BE
	FLOAT64_LE
	FLOAT64_BE
	// There are so many more...
	FormatTypeLast  = FLOAT64_BE
	FormatTypeFirst = S8
)

func (f FormatType) String() string {
	switch f {
	case S8:
		return "S8"
	case U8:
		return "U8"
	case S16_LE:
		return "S16_LE"
	case S16_BE:
		return "S16_BE"
	case U16_LE:
		return "U16_LE"
	case U16_BE:
		return "U16_BE"
	case S24_LE:
		return "S24_LE"
	case S24_BE:
		return "S24_BE"
	case U24_LE:
		return "U24_LE"
	case U24_BE:
		return "U24_BE"
	case S32_LE:
		return "S32_LE"
	case S32_BE:
		return "S32_BE"
	case U32_LE:
		return "U32_LE"
	case U32_BE:
		return "U32_BE"
	case FLOAT_LE:
		return "FLOAT_LE"
	case FLOAT_BE:
		return "FLOAT_BE"
	case FLOAT64_LE:
		return "FLOAT64_LE"
	case FLOAT64_BE:
		return "FLOAT64_BE"
	default:
		return fmt.Sprintf("Invalid FormatType (%d)", f)
	}
}

type SubformatType int

const (
	StandardSubformat  SubformatType = iota
	SubformatTypeFirst               = StandardSubformat
	SubformatTypeLast                = StandardSubformat
)

func (f SubformatType) String() string {
	switch f {
	case StandardSubformat:
		return "StandardSubformat"
	default:
		return fmt.Sprintf("Invalid SubformatType (%d)", f)
	}
}

type MmapStatus struct {
	State          state.State
	Pad1           int32
	HWPtr          uint
	Tstamp         misc.Timespec
	SuspendedState state.State
	AudioTstamp    misc.Timespec
}
type MmapControl struct {
	ApplPtr  uint
	AvailMin uint
}

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
	return fmt.Sprintf("Protocol %d.%d.%d (%d)", v.Major(), v.Minor(), v.Patch(), uint32(v))
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
	r := fmt.Sprintf("PCM %d/%d/%d ", s.Card, s.Device, s.Subdevice)
	switch s.Stream {
	case 0:
		r += "play"
	case 1:
		r += "capt"
	default:
		r += fmt.Sprintf("unknown stream direction (%d)", s.Stream)
	}
	r += fmt.Sprintf(" %#v", gstr(s.Name[:]))
	if s.SubdevicesCount != 1 || s.SubdevicesAvail != 1 || s.DevClass != 0 || s.DevSubclass != 0 {
		r += fmt.Sprintf(" (%d / %d) cls %d subcls %d", s.SubdevicesCount, s.SubdevicesAvail, s.DevClass, s.DevSubclass)
	}
	return r
}

const (
	MaskMax = 256
)

type Mask struct {
	Bits [(MaskMax + 31) / 32]uint32
}

type Interval struct {
	Min, Max uint32
	Flags    Flags
}

func (i Interval) String() string {
	return fmt.Sprintf("Interval(%d/%d 0x%x)", i.Min, i.Max, i.Flags)
}

type Params struct {
	Flags     uint32
	Masks     [ParamLastMask - ParamFirstMask + 1]Mask
	_         [5]Mask
	Intervals [ParamLastInterval - ParamFirstInterval + 1]Interval
	_         [9]Interval
	Rmask     uint32
	Cmask     uint32
	Info      uint32
	Msbits    uint32
	RateNum   uint32
	RateDen   uint32
	FifoSize  misc.Uframes
	_         [64]byte
}

func (p *Params) SetAccess(a AccessType) {
	p.SetMask(ParamAccess, uint32(1<<uint(a)))
}
func (p *Params) SetFormat(f FormatType) {
	p.SetMask(ParamFormat, uint32(1<<uint(f)))
}
func (p *Params) SetMask(param Param, v uint32) {
	p.Masks[param-ParamFirstMask].Bits[0] = v
}

func (p *Params) SetInterval(param Param, min, max uint32, flags Flags) {
	p.Intervals[param-ParamFirstInterval].Min = min
	p.Intervals[param-ParamFirstInterval].Max = max
	p.Intervals[param-ParamFirstInterval].Flags = flags
}
func (p *Params) SetIntervalToMin(param Param) {
	p.Intervals[param-ParamFirstInterval].Max = p.Intervals[param-ParamFirstInterval].Min
}
func (p *Params) IntervalInRange(param Param, v uint32) bool {
	if p.Intervals[param-ParamFirstInterval].Min > v {
		return false
	}
	if p.Intervals[param-ParamFirstInterval].Max < v {
		return false
	}
	return true
}

func fmt_uint(v uint32) string {
	if v == 0 {
		return "0"
	}
	if v == 0xffffffff {
		//return "λ"
		return "∞"
	}
	return fmt.Sprintf("0x%08x", v)
}

func (s *Params) String() string {
	return s.Diff(&Params{})
}

func fmt_cmask(v uint32) string {

	s := ""
	o := v
	for p := ParamFirstMask; p < ParamLastMask; p++ {
		if v&(1<<p) != 0 {
			o ^= (1 << p)
			s += " | " + p.String()
		}
	}
	for p := ParamFirstInterval; p < ParamLastInterval; p++ {
		if v&(1<<p) != 0 {
			o ^= (1 << p)
			s += " | " + p.String()
		}
	}

	if v == 0 {
		return "0"
	}
	if v == 0xffffffff {
		return "∞"
	}
	return fmt.Sprintf("0x%08x%s (0x%08x left)", v, s, o)
}

func (s *Params) Diff(w *Params) string {
	r := ""

	if s.Flags != w.Flags {
		r += fmt.Sprintf("  Flags 0x%x\n", s.Flags)
	}

	for i := range s.Masks {
		for j := range s.Masks[i].Bits {
			if s.Masks[i].Bits[j] != w.Masks[i].Bits[j] {
				v := s.Masks[i].Bits[j]

				sv := ""

				/*				for mv := ParamFirstMask; mv < ParamLastMask; mv++ {
									if v&(1<<mv) != 0 {
										sv += " " + mv.String()
										//						v ^= (1<<mv)
									}
								}

								for iv := ParamFirstInterval; iv < ParamLastInterval; iv++ {
									if v&(1<<iv) != 0 {
										sv += " " + iv.String()
										//						v ^= (1 << iv)
									}
								}*/

				if Param(i)+ParamFirstMask == ParamAccess {
					for a := AccessTypeFirst; a <= AccessTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}
				if Param(i)+ParamFirstMask == ParamFormat {
					for a := FormatTypeFirst; a <= FormatTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}
				if Param(i)+ParamFirstMask == ParamSubformat {
					for a := SubformatTypeFirst; a <= SubformatTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}

				r += fmt.Sprintf("  Mask %d[%d]  %8s  %-12s %s\n", i, j, (Param(i) + ParamFirstMask).String(), fmt_uint(v), sv)
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

		it := (Param(i) + ParamFirstInterval).String()
		iv := ""

		if s.Intervals[i].Min == 0 && s.Intervals[i].Max == 0xffffffff {
			iv = "0/∞ "
		} else {
			iv = fmt.Sprintf("%d/%d ", s.Intervals[i].Min, s.Intervals[i].Max)
		}

		ix := ""
		if s.Intervals[i].Flags != 0 {
			ix = s.Intervals[i].Flags.String()
		}
		r += fmt.Sprintf("%-20s %20s %-20s\n", it, iv, ix)
	}
	if s.Rmask != w.Rmask {
		r += "  Rmask  " + fmt_cmask(s.Rmask) + "\n"
	}
	if s.Cmask != w.Cmask {
		r += "  Cmask  " + fmt_cmask(s.Cmask) + "\n"
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
