// Package color is a very basic ANSI escape sequence colorer.
// Use it like this:
//
//     fmt.Println(color.Text(color.Red) + "I AM SO RED" + color.Reset())
//     fmt.Println(color.All(color.Red, false, color.White) + "Red on white." + color.Reset())
//
package color

import (
	"strconv"
)

type Color int

const (
	Black Color = iota + 1
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func Text(c Color) string {
	return "\x1b[" + strconv.Itoa(int(c)+29) + "m"
}
func All(foreground Color, bright bool, background Color) string {
	if bright {
		return "\x1b[1;" + strconv.Itoa(int(foreground)+29) + ";" + strconv.Itoa(int(background)+39) + "m"
	} else {
		return "\x1b[" + strconv.Itoa(int(foreground)+29) + ";" + strconv.Itoa(int(background)+39) + "m"
	}
}
func Reset() string {
	return "\x1b[m"
}

// Error() returns an error string colored red
func Error(e error) string {
	return Text(Red) + e.Error() + Reset()
}

// Pass() returns a cute little green check mark ✓
func Pass() string {
	return Text(Green) + "✓" + Reset()
}

// Fail() returns a cute little red cross mark ✗
func Fail() string {
	return Text(Red) + "✗" + Reset()
}
