package alsatype

import "fmt"

func (v PVersion) Major() int {
	return int(v >> 16 & 0xffff)
}
func (v PVersion) Minor() int {
	return int(v >> 8 & 0xff)
}
func (v PVersion) Patch() int {
	return int(v & 0xff)
}
func (v PVersion) String() string {
	return fmt.Sprintf("Protocol %d.%d.%d (%d)", v.Major(), v.Minor(), v.Patch(), uint32(v))
}
