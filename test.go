package main

import (
	"fmt"
	"os"
)

// _, _, errnop := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
//errno := int(errno)

//func ioctl(fd int, request, argp uintptr) error {
//	_, _, errorp := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), request, argp)
//	return os.NewSyscallError("ioctl", int(errorp))
//}

func main() {
	//	if err := list_the_things(); err != nil {
	//		fmt.Println(err)
	//	}

	if err := boop("/dev/snd/pcmC1D0p"); err != nil {
		fmt.Println(err)
	}

}

func boop(path string) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	var pv PVersion
	err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead, 4, CmdPCMVersion), &pv)
	if err != nil {
		return err
	}

	ttstamp := uint32(PCMTimestampTypeGettimeofday)
	err = ioctl(path, fh.Fd(), ioctl_encode(CmdWrite, 4, CmdPCMTimestampType), &ttstamp)
	if err != nil {
		return err
	}

	var params HwParams
	err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead|CmdWrite, 608, CmdPCMHwRefine), &params)
	if err != nil {
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
		err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead, 4, CmdControlVersion), &pv)
		if err != nil {
			return err
		}

		var ci CardInfo
		err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead, 376, CmdControlCardInfo), &ci)
		if err != nil {
			return err
		}

		fmt.Println(ci, pv)

		var next int32 = -1

		for {
			err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead, 4, CmdControlPCMNextDevice), &next)
			if err != nil {
				return err
			}

			if next == -1 {
				break
			}

			var pi PCMInfo
			pi.Device = uint32(next)
			pi.Subdevice = 0
			err = ioctl(path, fh.Fd(), ioctl_encode(CmdRead|CmdWrite, 288, CmdControlPCMInfo), &pi)
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
