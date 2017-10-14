package alsa

import (
	"fmt"
	"reflect"
	"syscall"
)

type ioctl_e uintptr

func (c ioctl_e) String() string {
	mode := c >> 30 & 0x03
	size := c >> 16 & 0x3fff
	cmd := c & 0xffff
	mode_str := ""
	if mode&cmdWrite > 0 {
		mode_str += " write"
	}
	if mode&cmdRead > 0 {
		mode_str += " read "
	}
	return fmt.Sprintf("ioctl%s (%d bytes) 0x%04x", mode_str, size, uintptr(cmd))
}

func ioctl(fd uintptr, c ioctl_e, ptr interface{}) error {
	var p uintptr

	if ptr != nil {
		v := reflect.ValueOf(ptr)
		p = v.Pointer()
	}

	//fmt.Printf("%s :: %d bytes\n", c, reflect.TypeOf(ptr).Elem().Size())
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(c), p)
	if e != 0 {
		return fmt.Errorf("%s failed: %v", c, e)
	}
	return nil
}

func gstr(c []byte) string {
	for i, v := range c {
		if v == 0 {
			return string(c[:i])
		}
	}
	return string(c)
}

func ioctl_encode(mode byte, size uint16, cmd uintptr) ioctl_e {
	return ioctl_e(mode)<<30 | ioctl_e(size)<<16 | ioctl_e(cmd)
}

func ioctl_encode_ptr(mode byte, ref interface{}, cmd uintptr) ioctl_e {
	return ioctl_encode(mode, uint16(reflect.TypeOf(ref).Elem().Size()), cmd)
}
