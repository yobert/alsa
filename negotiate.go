package alsa

import (
	"fmt"
)

func (device *Device) NegotiateChannels(channels ...int) (int, error) {
	var err error

	for _, v := range channels {

		if !device.hwparams.IntervalInRange(paramChannels, uint32(v)) {
			err = fmt.Errorf("Channels %d out of range", v)
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(paramChannels, uint32(v), uint32(v), Integer)

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
		if !device.hwparams.IntervalInRange(paramRate, uint32(v)) {
			err = fmt.Errorf("Rate %d out of range", v)
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(paramRate, uint32(v), uint32(v), Integer)

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
		if !device.hwparams.IntervalInRange(paramBufferSize, uint32(v)) {
			err = fmt.Errorf("Buffer size %d out of range", v)
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(paramBufferSize, uint32(v), uint32(v), Integer)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func (device *Device) NegotiatePeriodSize(period_sizes ...int) (int, error) {
	var err error

	for _, v := range period_sizes {
		if !device.hwparams.IntervalInRange(paramPeriodSize, uint32(v)) {
			err = fmt.Errorf("Period size %d out of range", v)
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetInterval(paramPeriodSize, uint32(v), uint32(v), Integer)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, err
}

func (device *Device) BytesPerFrame() int {
	sample_size := int(device.hwparams.Intervals[paramSampleBits-paramFirstInterval].Max) / 8
	channels := int(device.hwparams.Intervals[paramChannels-paramFirstInterval].Max)
	return sample_size * channels
}
