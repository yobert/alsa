package alsa

import (
	"fmt"
	"os"
	"unsafe"

	"color"

	"github.com/yobert/alsa/misc"
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

	pversion pVersion

	hwparams      hwParams
	hwparams_prev hwParams

	swparams      swParams
	swparams_prev swParams
}

func (device Device) String() string {
	return device.Title
}

func (card *Card) Devices() ([]*Device, error) {
	var next int32 = -1
	ret := make([]*Device, 0)

	for {
		err := ioctl(card.fh.Fd(), ioctl_encode(cmdRead, 4, cmdControlPCMNextDevice), &next)
		if err != nil {
			return ret, err
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
			err = ioctl(card.fh.Fd(), ioctl_encode(cmdRead|cmdWrite, 288, cmdControlPCMInfo), &pi)
			if err != nil {
				// Probably means that device doesn't match that stream type
			} else {
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
	}

	return ret, nil
}

func (device *Device) Open() error {
	var err error
	device.fh, err = os.OpenFile(device.Path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	err = ioctl(device.fh.Fd(), ioctl_encode(cmdRead, 4, cmdPCMVersion), &device.pversion)
	if err != nil {
		device.fh.Close()
		return err
	}

	ttstamp := uint32(pcmTimestampTypeGettimeofday)
	err = ioctl(device.fh.Fd(), ioctl_encode(cmdWrite, 4, cmdPCMTimestampType), &ttstamp)
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

func (device *Device) Prepare() error {
	if device.debug {
		fmt.Println("Final hardware parameter changes:")
		fmt.Println(color.Text(color.Green))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}
	device.hwparams_prev = device.hwparams

	err := ioctl(device.fh.Fd(), ioctl_encode(cmdRead|cmdWrite, 608, cmdPCMHwParams), &device.hwparams)
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

	device.swparams = swParams{}
	device.swparams_prev = swParams{}

	device.swparams.PeriodStep = 1
	device.swparams.AvailMin = uint(buf_size)
	device.swparams.XferAlign = 1
	device.swparams.StartThreshold = uint(buf_size)
	device.swparams.StopThreshold = uint(buf_size * 2)
	device.swparams.Proto = device.pversion
	device.swparams.TstampType = 1

	if err := device.sw_params(); err != nil {
		return err
	}

	if err := ioctl(device.fh.Fd(), ioctl_encode(0, 0, cmdPCMPrepare), nil); err != nil {
		return err
	}

	return nil
}

func (device *Device) Write(b []byte, frames int) error {
	if len(b) == 0 {
		return nil
	}
	return ioctl(device.fh.Fd(), ioctl_encode(cmdWrite, pcm.XferISize, cmdPCMWriteIFrames), &pcm.XferI{
		Buf:    uintptr(unsafe.Pointer(&b[0])),
		Frames: misc.Uframes(frames),
	})
}

func (device *Device) refine() error {

	if device.debug {
		fmt.Println("Requesting changes:")
		fmt.Println(color.Text(color.Green))
		fmt.Print(device.hwparams.Diff(&device.hwparams_prev))
		fmt.Println(color.Reset())
	}
	device.hwparams_prev = device.hwparams

	err := ioctl(device.fh.Fd(), ioctl_encode(cmdRead|cmdWrite, 608, cmdPCMHwRefine), &device.hwparams)
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

	err := ioctl(device.fh.Fd(), ioctl_encode(cmdRead|cmdWrite, 136, cmdPCMSwParams), &device.swparams)
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
