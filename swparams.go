package main

import (
	"fmt"
	"reflect"

	"color"
)

type SwParams struct {
	TstampMode int32
	PeriodStep uint32
	SleepMin   uint32

	AvailMin         uint
	XferAlign        uint
	StartThreshold   uint
	StopThreshold    uint
	SilenceThreshold uint
	SilenceSize      uint
	Boundary         uint

	Proto         PVersion
	TstampType    uint32
	Reserved      [56]byte
	padding_for_c [4]byte
}

func sw_params(fd uintptr, params *SwParams, last *SwParams) error {

	fmt.Println(color.Text(color.Green))
	fmt.Print(params.Diff(last))
	*last = *params

	err := ioctl(fd, ioctl_encode(CmdRead|CmdWrite, 136, CmdPCMSwParams), params)
	if err != nil {
		return err
	}

	fmt.Println(color.Text(color.Magenta))
	fmt.Print(params.Diff(last))
	*last = *params

	return nil
}

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
		case reflect.Uint:
			fallthrough
		case reflect.Uint8:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint64:
			if v1v.Uint() != v2v.Uint() {
				d = fmt.Sprintf("%d", v1v.Uint())
			}
		case reflect.Int:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
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
			r += field.Name + "\t" + d + "\n"
		}
	}

	if r == "" {
		r += "  No changes\n"
	}
	return r
}
