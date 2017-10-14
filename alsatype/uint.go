package alsatype

import "fmt"

func fmt_uint(v uint32) string {
	if v == 0 {
		return "0"
	}
	if v == 0xffffffff {
		//return "λ"
		return "∞"
	}
	return fmt.Sprintf("0x%08x", v)
}
