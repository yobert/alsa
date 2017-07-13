package main

import (
	"fmt"
	"os"

	"color"
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
		fmt.Println(err)
	}

}

func refine(fd uintptr, params *Params, last *Params) error {

	fmt.Println(color.Text(color.Green))
	fmt.Print(params.Diff(last))
	*last = *params

	err := ioctl(fd, ioctl_encode(CmdRead|CmdWrite, 608, CmdPCMHwRefine), params)
	if err != nil {
		return err
	}

	fmt.Println(color.Text(color.Magenta))
	fmt.Print(params.Diff(last))
	*last = *params

	return nil
}

func boop(path string) error {
	fh, err := os.Open(path)
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

	params.Cmask = 0
	params.Rmask = 0xffffffff
	params.SetInterval(ParamChannels, 2, 2, Integer)
	params.SetInterval(ParamRate, 44100, 44100, Integer)

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
	}

	params.Cmask = 0
	params.Rmask = 0xffffffff
	params.SetInterval(ParamBufferTime, 10, 16000, Integer)
	//params.SetInterval(ParamPeriods, 2, 2, Integer)

	if err := refine(fh.Fd(), params, last); err != nil {
		return err
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

			var pi PCMInfo
			pi.Device = uint32(next)
			pi.Subdevice = 0
			err = ioctl(fh.Fd(), ioctl_encode(CmdRead|CmdWrite, 288, CmdControlPCMInfo), &pi)
			if err != nil {
				//return err
				//fmt.Println(err)
			} else {
				fmt.Println(pi)
			}
		}

		fmt.Println("")
	}
	return nil
}
