package alsa

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/yobert/alsa/alsatype"
	"github.com/yobert/alsa/color"
	"github.com/yobert/alsa/pcm"
)

type DeviceType int

const (
	UnknownDeviceType DeviceType = iota
	PCM
)

func (t DeviceType) String() string {
	switch t {
	case PCM:
		return "PCM"
	default:
		return fmt.Sprintf("UnknownDeviceType(%d)", t)
	}
}

type Device struct {
	Type         DeviceType
	Number       int
	Play, Record bool

	Path  string
	Title string

	debug bool

	fh      *os.File
	pcminfo pcmInfo

	pversion alsatype.PVersion

	hwparams      hwParams
	hwparams_prev hwParams

	swparams      alsatype.SwParams
	swparams_prev alsatype.SwParams
}

func (device Device) String() string {
	return device.Title
}

func (card *Card) Devices() ([]*Device, error) {
	var ret []*Device
	next := int32(-1)

	for {
		err := ioctl(card.fh.Fd(), ioctl_encode_ptr(cmdRead, &next, cmdControlPCMNextDevice), &next)
		if err != nil {
			return nil, err
		}
		if next == -1 {
			// No more devices
			break
		}

		for stream := int32(0); stream < 2; stream++ {
			var pi pcmInfo
			pi.Device = uint32(next)
			pi.Subdevice = 0
			pi.Stream = stream
			err = ioctl(card.fh.Fd(), ioctl_encode_ptr(cmdRead|cmdWrite, &pi, cmdControlPCMInfo), &pi)
			if err != nil {
				// Probably means that device doesn't match that stream type.
				continue
			}

			play := true
			record := false
			sstr := "p"
			if stream == 1 {
				play = false
				record = true
				sstr = "c"
			}

			ret = append(ret, &Device{
				Type:    PCM,
				Path:    fmt.Sprintf("/dev/snd/pcmC%dD%d%s", card.Number, next, sstr),
				Play:    play,
				Record:  record,
				Number:  int(next),
				Title:   gstr(pi.Name[:]),
				pcminfo: pi,
			})
		}
	}

	return ret, nil
}

func (device *Device) Open() error {
	var err error
	device.fh, err = os.OpenFile(device.Path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	err = ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdRead, &device.pversion, cmdPCMVersion), &device.pversion)
	if err != nil {
		device.fh.Close()
		return err
	}

	ttstamp := uint32(pcmTimestampTypeGettimeofday)
	err = ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdWrite, &ttstamp, cmdPCMTimestampType), &ttstamp)
	if err != nil {
		device.fh.Close()
		return err
	}

	device.hwparams = hwParams{}
	device.hwparams_prev = hwParams{}

	for i := range device.hwparams.Masks {
		for ii := 0; ii < 2; ii++ {
			device.hwparams.Masks[i].Bits[ii] = 0xffffffff
		}
	}
	for i := range device.hwparams.Intervals {
		device.hwparams.Intervals[i].Max = 0xffffffff
	}
	device.hwparams.Rmask = 0xffffffff

	if err := device.refine(); err != nil {
		return err
	}

	device.hwparams.Cmask = 0
	device.hwparams.Rmask = 0xffffffff
	device.hwparams.SetAccess(RWInterleaved)

	if err := device.refine(); err != nil {
		return err
	}

	return nil
}

func (device *Device) Close() {
	if device.fh != nil {
		device.fh.Close()
	}
}

func (device *Device) Debug(v bool) {
	device.debug = v
}

func (device *Device) Prepare() error {
	if device.debug {
		fmt.Println("Final hardware parameter changes:")
		fmt.Println(color.Text(color.Green))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}
	device.hwparams_prev = device.hwparams

	err := ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdRead|cmdWrite, &device.hwparams, cmdPCMHwParams), &device.hwparams)
	if err != nil {
		return err
	}

	if device.debug {
		fmt.Println("Final hardware parameter results:")
		fmt.Println(color.Text(color.Magenta))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}

	device.hwparams_prev = device.hwparams

	// final buf size
	buf_size := int(device.hwparams.Intervals[paramBufferSize-paramFirstInterval].Max)

	device.swparams = alsatype.SwParams{}
	device.swparams_prev = alsatype.SwParams{}

	device.swparams.PeriodStep = 1
	device.swparams.AvailMin = alsatype.Uframes(buf_size)
	device.swparams.XferAlign = 1
	device.swparams.StartThreshold = alsatype.Uframes(buf_size)
	device.swparams.StopThreshold = alsatype.Uframes(buf_size * 2)
	device.swparams.Proto = device.pversion
	device.swparams.TstampType = 1

	if err := device.sw_params(); err != nil {
		return err
	}

	if err := ioctl(device.fh.Fd(), ioctl_encode(0, 0, cmdPCMPrepare), nil); err != nil {
		return fmt.Errorf("Device prepare failure: %v", err)
	}

	return nil
}

// BufferFormat() is not valid until after Prepare() is called
func (device *Device) BufferFormat() BufferFormat {
	bf := BufferFormat{}

	for i := FormatTypeFirst; i < FormatTypeLast+1; i++ {
		// there should only be 1 format bit set at this point
		if device.hwparams.GetFormatSupport(i) {
			bf.SampleFormat = i
			break
		}
	}

	v, _ := device.hwparams.IntervalRange(paramRate)
	bf.Rate = int(v)
	v, _ = device.hwparams.IntervalRange(paramChannels)
	bf.Channels = int(v)

	return bf
}

// This function is deprecated and will be removed at some point.
// Please use NewBufferDuration
func (device *Device) NewBufferSeconds(seconds int) Buffer {
	return device.NewBufferDuration(time.Second * time.Duration(seconds))
}

func (device *Device) NewBufferDuration(d time.Duration) Buffer {
	bf := device.BufferFormat()

	frames := int(float64(bf.Rate)*d.Seconds() + 0.5)
	bytecount := frames * device.BytesPerFrame()
	data := make([]byte, bytecount)

	return Buffer{Format: bf, Data: data}
}

func (device *Device) Read(buf []byte) error {
	frames := len(buf) / device.BytesPerFrame()
	x := pcm.XferI{
		Buf:    uintptr(unsafe.Pointer(&buf[0])),
		Frames: alsatype.Uframes(frames),
	}
	return ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdRead, &x, cmdPCMReadIFrames), &x)
}

func (device *Device) Write(buf []byte, frames int) error {
	x := pcm.XferI{
		Buf:    uintptr(unsafe.Pointer(&buf[0])),
		Frames: alsatype.Uframes(frames),
	}
	return ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdWrite, &x, cmdPCMWriteIFrames), &x)
}

func (device *Device) refine() error {

	if device.debug {
		fmt.Println("Requesting changes:")
		fmt.Println(color.Text(color.Green))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}
	device.hwparams_prev = device.hwparams

	err := ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdRead|cmdWrite, &device.hwparams, cmdPCMHwRefine), &device.hwparams)
	if err != nil {
		return err
	}

	if device.debug {
		fmt.Println("Results:")
		fmt.Println(color.Text(color.Magenta))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}
	device.hwparams_prev = device.hwparams

	return nil
}

func (device *Device) sw_params() error {
	if device.debug {
		fmt.Println("Requesting soft parameters:")
		fmt.Println(color.Text(color.Green))
		fmt.Print(device.swparams.Diff(&device.swparams_prev))
		fmt.Println(color.Reset())
	}
	device.swparams_prev = device.swparams

	err := ioctl(device.fh.Fd(), ioctl_encode_ptr(cmdRead|cmdWrite, &device.swparams, cmdPCMSwParams), &device.swparams)
	if err != nil {
		return err
	}

	if device.debug {
		fmt.Println("Results:")
		fmt.Println(color.Text(color.Magenta))
		fmt.Print(device.swparams.Diff(&device.swparams_prev))
		fmt.Println(color.Reset())
	}
	device.swparams_prev = device.swparams

	return nil
}
