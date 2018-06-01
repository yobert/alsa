package alsa

import (
	"fmt"
)

func range_check(params hwParams, p param, v int) error {
	min, max := params.IntervalRange(p)
	if v < 0 || uint32(v) < min || uint32(v) > max {
		if min != max {
			return fmt.Errorf("Requested value %d is out of hardware possible range (min %d max %d)", v, min, max)
		}
		return fmt.Errorf("Requested value %d is not supported by hardware: Must be %d", v, min)
	}
	return nil
}

func format_check(params hwParams, f FormatType) error {
	if params.GetFormatSupport(f) {
		return nil
	}
	list := ""
	count := 0
	for i := FormatTypeFirst; i < FormatTypeLast+1; i++ {
		if !params.GetFormatSupport(i) {
			continue
		}
		if count > 0 {
			list += ", "
		}
		list += i.String()
		count++
	}
	if count == 0 {
		return fmt.Errorf("%s not supported: No possible formats", f)
	}
	if count == 1 {
		return fmt.Errorf("%s not supported: Must be %s", f, list)
	}
	return fmt.Errorf("%s not supported: Must be one of %s", f, list)
}

func (device *Device) NegotiateChannels(channels ...int) (int, error) {
	var err error

	for _, v := range channels {
		err = range_check(device.hwparams, paramChannels, v)
		if err != nil {
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

	return 0, fmt.Errorf("Channel count negotiation failure: %v", err)
}

func (device *Device) NegotiateRate(rates ...int) (int, error) {
	var err error

	for _, v := range rates {
		err = range_check(device.hwparams, paramRate, v)
		if err != nil {
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

	return 0, fmt.Errorf("Rate negotiation failure: %v", err)
}

func (device *Device) NegotiateFormat(formats ...FormatType) (FormatType, error) {
	var err error

	for _, v := range formats {
		err = format_check(device.hwparams, v)
		if err != nil {
			continue
		}

		device.hwparams.Cmask = 0
		device.hwparams.Rmask = 0xffffffff
		device.hwparams.SetFormat(v)

		err = device.refine()
		if err == nil {
			return v, nil
		}
	}

	return 0, fmt.Errorf("Format negotiation failure: %v", err)
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
