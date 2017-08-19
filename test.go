package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"time"
	"unsafe"
	//"github.com/edsrzf/mmap-go"
	"github.com/yobert/alsa/pcm"
)

// _, _, errnop := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
//errno := int(errno)

//func ioctl(fd int, request, argp uintptr) error {
//	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
//	return os.NewSyscallError("ioctl", int(errorp))
//}

func main() {
	if err := list_the_things(); err != nil {
		fmt.Println(err)
	}

	if err := boop("/dev/snd/pcmC0D0p"); err != nil {
		//if err := boop("/dev/snd/pcmC2D0p"); err != nil {
		fmt.Println(err)
	}

}

func boop(path string) error {
	fh, err := os.OpenFile(path, os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer fh.Close()

	var pv PVersion
	err = ioctl(fh.Fd(), ioctl_encode(CmdRead, 4, CmdPCMVersion), &pv)
	if err != nil {
		return err
	}

	ttstamp := uint32(PCMTimestampTypeGettimeofday)
	err = ioctl(fh.Fd(), ioctl_encode(CmdWrite, 4, CmdPCMTimestampType), &ttstamp)
	if err != nil {
		return err
	}

	params := &Params{}
	last := &Params{}

	for i := range params.Masks {
		for ii := 0; ii < 2; ii++ {
			params.Masks[i].Bits[ii] = 0xffffffff
		}
	}
	for i := range params.Intervals {
		params.Intervals[i].Max = 0xffffffff
	}
	params.Rmask = 0xffffffff

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
	}

	if !params.IntervalInRange(ParamChannels, 2) {
		return fmt.Errorf("Stereo not available")
	}

	params.Cmask = 0
	params.Rmask = 0xffffffff
	params.SetInterval(ParamChannels, 2, 2, Integer)

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
	}

	rate := uint32(44100)

	if !params.IntervalInRange(ParamRate, rate) {
		return fmt.Errorf("%d Khz not available", rate)
	}

	params.Cmask = 0
	params.Rmask = 0xffffffff
	params.SetInterval(ParamRate, rate, rate, Integer)

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
	}

	params.Cmask = 0
	params.Rmask = 0xffffffff
	params.SetAccess(RWInterleaved)
	params.SetFormat(U24_BE)
	//params.SetFormat(U8)

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
	}

	//	params.Cmask = 0
	//	params.Rmask = 0xffffffff
	//	params.SetIntervalToMin(ParamBufferTime)
	//	if err := refine(fh.Fd(), params, last); err != nil {
	//		return err
	//	}

	if err := hw_params(fh.Fd(), params, last); err != nil {
		return err
	}

	swparams := &SwParams{}
	swlast := &SwParams{}

	swparams.PeriodStep = 1
	swparams.AvailMin = 1024
	swparams.XferAlign = 1
	swparams.StartThreshold = 1
	swparams.StopThreshold = 16384
	swparams.Proto = pv
	swparams.TstampType = 1

	if err := sw_params(fh.Fd(), swparams, swlast); err != nil {
		return err
	}

	if err := get_status(fh.Fd()); err != nil {
		return err
	}
	if err := ioctl(fh.Fd(), ioctl_encode(0, 0, CmdPCMPrepare), nil); err != nil {
		return err
	}
	if err := get_status(fh.Fd()); err != nil {
		return err
	}

	/*	if err := ioctl(fh.Fd(), ioctl_encode(0, 0, CmdPCMStart), nil); err != nil {
			return err
		}
		if err := get_status(fh.Fd()); err != nil {
			return err
		}*/

	buf_size := int(params.Intervals[ParamBufferSize-ParamFirstInterval].Max)

	buf_bytes := int(params.Intervals[ParamBufferBytes-ParamFirstInterval].Max)
	buf := bytes.NewBuffer(make([]byte, 0, buf_bytes))

	fmt.Println("buf", buf_bytes, "/", buf_size, "frames")
	fmt.Println("rate", rate)

	t := 0.0

	xfer := pcm.XferI{}

	for i := 0; i < buf_size; i++ {

		v := math.Sin(t * 2 * math.Pi * 440)
		//v += math.Sin(t * 2 * math.Pi * 261.63)
		//v += math.Sin(t * 2 * math.Pi * 349.23)
		v *= 0.1

		//buf[i] = uint8((v*0.5+0.5)*255)
		sample := uint32((v*0.5 + 0.5) * 16777215)

		binary.Write(buf, binary.BigEndian, sample)
		binary.Write(buf, binary.BigEndian, sample)

		t += 1.0 / float64(rate)

		xfer.Frames++
	}

	xfer.Buf = uintptr(unsafe.Pointer(&buf.Bytes()[0]))

	if err := get_status(fh.Fd()); err != nil {
		return err
	}
	err = ioctl(fh.Fd(), ioctl_encode(CmdWrite, pcm.XferISize, CmdPCMWriteIFrames), &xfer)
	if err != nil {
		return err
	}
	fmt.Println("xfer complete")
	for i := 0; i < 10; i++ {
		if err := get_status(fh.Fd()); err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

func list_the_things() error {
	for i := 0; i < 10; i++ {
		path := fmt.Sprintf("/dev/snd/controlC%d", i)
		_, err := os.Stat(path)
		if err != nil {
			continue
		}
		fh, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fh.Close()

		var pv PVersion
		err = ioctl(fh.Fd(), ioctl_encode(CmdRead, 4, CmdControlVersion), &pv)
		if err != nil {
			return err
		}

		var ci CardInfo
		err = ioctl(fh.Fd(), ioctl_encode(CmdRead, 376, CmdControlCardInfo), &ci)
		if err != nil {
			return err
		}

		fmt.Println(ci, pv)

		var next int32 = -1

		for {
			err = ioctl(fh.Fd(), ioctl_encode(CmdRead, 4, CmdControlPCMNextDevice), &next)
			if err != nil {
				return err
			}

			if next == -1 {
				break
			}

			for stream := int32(0); stream < 2; stream++ {
				var pi PCMInfo
				pi.Device = uint32(next)
				pi.Subdevice = 0
				pi.Stream = stream
				err = ioctl(fh.Fd(), ioctl_encode(CmdRead|CmdWrite, 288, CmdControlPCMInfo), &pi)
				if err != nil {
					//return err
					//fmt.Println(err)
				} else {
					fmt.Println(pi)
				}
			}
		}

		fmt.Println("")
	}
	return nil
}
