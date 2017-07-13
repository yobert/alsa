package main

import (
	"strings"
)

type Param uint32

const (
	ParamAccess    Param = 0
	ParamFormat    Param = 1
	ParamSubFormat Param = 2
	ParamFirstMask Param = ParamAccess
	ParamLastMask  Param = ParamSubFormat

	ParamSampleBits    Param = 8
	ParamFrameBits     Param = 9
	ParamChannels      Param = 10
	ParamRate          Param = 11
	ParamPeriodTime    Param = 12
	ParamPeriodSize    Param = 13
	ParamPeriodBytes   Param = 14
	ParamPeriods       Param = 15
	ParamBufferTime    Param = 16
	ParamBufferSize    Param = 17
	ParamBufferBytes   Param = 18
	ParamTickTime      Param = 19
	ParamFirstInterval Param = ParamSampleBits
	ParamLastInterval  Param = ParamTickTime
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
		r += "empty "
	}
	return strings.TrimSpace(r)
}

func (p Param) IsMask() bool {
	return p >= ParamFirstMask && p <= ParamLastMask
}
func (p Param) IsInterval() bool {
	return p >= ParamFirstInterval && p <= ParamLastInterval
}
func (p Param) String() string {
	if p.IsMask() {
		return "≡" + p.name()
	}
	if p.IsInterval() {
		return "±" + p.name()
	}
	return "Invalid"
}
func (p Param) name() string {
	switch p {
	case ParamAccess:
		return "Access"
	case ParamFormat:
		return "Format"
	case ParamSubFormat:
		return "Subfmt"

	case ParamSampleBits:
		return "SampleBits"
	case ParamFrameBits:
		return "FrameBits"
	case ParamChannels:
		return "Channels"
	case ParamRate:
		return "Rate"
	case ParamPeriodTime:
		return "PeriodTime"
	case ParamPeriodSize:
		return "PeriodSize"
	case ParamPeriodBytes:
		return "PeriodBytes"
	case ParamPeriods:
		return "Periods"
	case ParamBufferTime:
		return "BufferTime"
	case ParamBufferSize:
		return "BufferSize"
	case ParamBufferBytes:
		return "BufferBytes"
	case ParamTickTime:
		return "TickTime"
	default:
		return "Invalid"
	}
}
