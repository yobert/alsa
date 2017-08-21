package alsa

import (
	"fmt"
)

func (device *Device) NegotiateChannels(channels ...int) (int, error) {
	var err error

	for _, v := range channels {

		if !device.hwparams.IntervalInRange(ParamChannels, uint32(v)) {
			err = fmt.Errorf("Channels %d out of range")
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(ParamChannels, uint32(v), uint32(v), Integer)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func (device *Device) NegotiateRate(rates ...int) (int, error) {
	var err error

	for _, v := range rates {
		if !device.hwparams.IntervalInRange(ParamRate, uint32(v)) {
			err = fmt.Errorf("Rate %d out of range")
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(ParamRate, uint32(v), uint32(v), Integer)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func (device *Device) NegotiateFormat(formats ...FormatType) (FormatType, error) {
	var err error

	for _, v := range formats {
		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetFormat(v)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func (device *Device) NegotiateBufferSize(buffer_sizes ...int) (int, error) {
	var err error

	for _, v := range buffer_sizes {
		if !device.hwparams.IntervalInRange(ParamBufferSize, uint32(v)) {
			err = fmt.Errorf("Buffer size %d out of range")
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(ParamBufferSize, uint32(v), uint32(v), Integer)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}
