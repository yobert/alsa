package alsatype

import (
	"fmt"
	"reflect"
)

func (s *SwParams) String() string {
	return s.Diff(&SwParams{})
}

func (s *SwParams) Diff(w *SwParams) string {
	r := ""

	v1 := reflect.ValueOf(*s)
	v2 := reflect.ValueOf(*w)

	typ := reflect.TypeOf(*s)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Name == "Reserved" || field.Name == "padding_for_c" {
			continue
		}

		v1v := v1.Field(i)
		v2v := v2.Field(i)
		d := ""

		switch v1v.Type().Kind() {
		case reflect.Uint32:
			if v1v.Uint() != v2v.Uint() {
				d = fmt_uint(uint32(v1v.Uint())) + fmt.Sprintf(" (%d)", v1v.Uint())
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64:
			if v1v.Uint() != v2v.Uint() {
				d = fmt.Sprintf("%d", v1v.Uint())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v1v.Int() != v2v.Int() {
				d = fmt.Sprintf("%d", v1v.Int())
			}
		case reflect.String:
			if v1v.String() != v2v.String() {
				d = v1v.String()
			}
		default:
			d = v1v.Type().Kind().String()
		}
		if d != "" {
			r += fmt.Sprintf("%20s %s\n", field.Name, d)
		}
	}

	if r == "" {
		r += "  No changes\n"
	}
	return r
}
