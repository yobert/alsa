package alsa

import (
	"fmt"

	"github.com/yobert/alsa/alsatype"
	//	"github.com/yobert/alsa/pcm/state"
)

const (
	cmdWrite = 1
	cmdRead  = 2

	cmdPCMInfo          uintptr = 0x4101
	cmdPCMVersion       uintptr = 0x4100
	cmdPCMTimestamp     uintptr = 0x4102
	cmdPCMTimestampType uintptr = 0x4103
	cmdPCMHwRefine      uintptr = 0x4110
	cmdPCMHwParams      uintptr = 0x4111
	cmdPCMSwParams      uintptr = 0x4113
	cmdPCMStatus        uintptr = 0x4120

	cmdPCMPrepare uintptr = 0x4140
	cmdPCMReset   uintptr = 0x4141
	cmdPCMStart   uintptr = 0x4142
	cmdPCMDrop    uintptr = 0x4143
	cmdPCMDrain   uintptr = 0x4144
	cmdPCMPause   uintptr = 0x4145 // int
	cmdPCMRewind  uintptr = 0x4146 // snd_pcm_uframes_t
	cmdPCMResume  uintptr = 0x4147
	cmdPCMXrun    uintptr = 0x4148
	cmdPCMForward uintptr = 0x4149

	cmdPCMWriteIFrames uintptr = 0x4150 // snd_xferi
	cmdPCMReadIFrames  uintptr = 0x4151 // snd_xferi
	cmdPCMWriteNFrames uintptr = 0x4152 // snd_xfern
	cmdPCMReadNFrames  uintptr = 0x4153 // snd_xfern

	cmdPCMLink   uintptr = 0x4160 // int
	cmdPCMUnlink uintptr = 0x4161

	cmdControlVersion       uintptr = 0x5500
	cmdControlCardInfo      uintptr = 0x5501
	cmdControlPCMNextDevice uintptr = 0x5530
	cmdControlPCMInfo       uintptr = 0x5531
)

const (
	pcmTimestampTypeGettimeofday = iota
	pcmTimestampTypeMonotonic
	pcmTimestampTypeMonotonicRaw
	pcmTimestampTypeLast
)

//const (
//	MapShared     = 0x00000001
//	OffsetData    = 0x00000000
//	OffsetStatus  = 0x80000000
//	OffsetControl = 0x81000000
//)

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
	S8 FormatType = iota
	U8
	S16_LE
	S16_BE
	U16_LE
	U16_BE
	S24_LE
	S24_BE
	U24_LE
	U24_BE
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

//type MmapStatus struct {
//	State          state.State
//	Pad1           int32
//	HWPtr          uint
//	Tstamp         Timespec
//	SuspendedState state.State
//	AudioTstamp    Timespec
//}
//type MmapControl struct {
//	ApplPtr  uint
//	AvailMin uint
//}

type cardInfo struct {
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

func (s cardInfo) String() string {
	return fmt.Sprintf("Card %d %#v", s.Card, gstr(s.Name[:]))
}

type pcmInfo struct {
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

func (s pcmInfo) String() string {
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
	maskMax = 256
)

type mask struct {
	Bits [(maskMax + 31) / 32]uint32
}

type interval struct {
	Min, Max uint32
	Flags    Flags
}

func (i interval) String() string {
	return fmt.Sprintf("Interval(%d/%d 0x%x)", i.Min, i.Max, i.Flags)
}

type hwParams struct {
	Flags     uint32
	Masks     [paramLastMask - paramFirstMask + 1]mask
	_         [5]mask
	Intervals [paramLastInterval - paramFirstInterval + 1]interval
	_         [9]interval
	Rmask     uint32
	Cmask     uint32
	Info      uint32
	Msbits    uint32
	RateNum   uint32
	RateDen   uint32
	FifoSize  alsatype.Uframes
	_         [64]byte
}

func (p *hwParams) SetAccess(a AccessType) {
	p.SetMask(paramAccess, uint32(1<<uint(a)))
}
func (p *hwParams) SetFormat(f FormatType) {
	p.SetMask(paramFormat, uint32(1<<uint(f)))
}
func (p *hwParams) SetMask(param param, v uint32) {
	p.Masks[param-paramFirstMask].Bits[0] = v
}
func (p *hwParams) GetFormatSupport(f FormatType) bool {
	bits := p.Masks[paramFormat-paramFirstMask].Bits[0]
	b := bits & (1 << uint(f))
	if b == 0 {
		return false
	}
	return true
}

func (p *hwParams) SetInterval(param param, min, max uint32, flags Flags) {
	p.Intervals[param-paramFirstInterval].Min = min
	p.Intervals[param-paramFirstInterval].Max = max
	p.Intervals[param-paramFirstInterval].Flags = flags
}
func (p *hwParams) SetIntervalToMin(param param) {
	p.Intervals[param-paramFirstInterval].Max = p.Intervals[param-paramFirstInterval].Min
}
func (p *hwParams) IntervalInRange(param param, v uint32) bool {
	min, max := p.IntervalRange(param)
	if min > v {
		return false
	}
	if max < v {
		return false
	}
	return true
}

func (p *hwParams) IntervalRange(param param) (uint32, uint32) {
	return p.Intervals[param-paramFirstInterval].Min, p.Intervals[param-paramFirstInterval].Max
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

func (s *hwParams) String() string {
	return s.Diff(&hwParams{})
}

func fmt_cmask(v uint32) string {

	s := ""
	o := v
	for p := paramFirstMask; p < paramLastMask; p++ {
		if v&(1<<p) != 0 {
			o ^= (1 << p)
			s += " | " + p.String()
		}
	}
	for p := paramFirstInterval; p < paramLastInterval; p++ {
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

func (s *hwParams) Diff(w *hwParams) string {
	r := ""

	if s.Flags != w.Flags {
		r += fmt.Sprintf("  Flags 0x%x\n", s.Flags)
	}

	for i := range s.Masks {
		for j := range s.Masks[i].Bits {
			if s.Masks[i].Bits[j] != w.Masks[i].Bits[j] {
				v := s.Masks[i].Bits[j]

				sv := ""

				/*				for mv := paramFirstMask; mv < paramLastMask; mv++ {
									if v&(1<<mv) != 0 {
										sv += " " + mv.String()
										//						v ^= (1<<mv)
									}
								}

								for iv := paramFirstInterval; iv < paramLastInterval; iv++ {
									if v&(1<<iv) != 0 {
										sv += " " + iv.String()
										//						v ^= (1 << iv)
									}
								}*/

				if param(i)+paramFirstMask == paramAccess {
					for a := AccessTypeFirst; a <= AccessTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}
				if param(i)+paramFirstMask == paramFormat {
					for a := FormatTypeFirst; a <= FormatTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}
				if param(i)+paramFirstMask == paramSubformat {
					for a := SubformatTypeFirst; a <= SubformatTypeLast; a++ {
						if v&(1<<uint(a)) != 0 {
							sv += " " + a.String()
							if v != 0xffffffff {
								v ^= (1 << uint(a))
							}
						}
					}
				}

				r += fmt.Sprintf("  Mask %d[%d]  %8s  %-12s %s\n", i, j, (param(i) + paramFirstMask).String(), fmt_uint(v), sv)
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

		it := (param(i) + paramFirstInterval).String()
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
