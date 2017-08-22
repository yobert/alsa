package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/yobert/alsa"
)

func main() {

	cards, err := alsa.OpenCards()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer alsa.CloseCards(cards)

	for _, card := range cards {
		fmt.Println(card)

		if err := beep_card(card); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func beep_card(card *alsa.Card) error {
	devices, err := card.Devices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if device.Type != alsa.PCM || !device.Play {
			continue
		}
		fmt.Println("───", device)

		if err := beep_device(device); err != nil {
			return err
		}
	}
	return nil
}

func beep_device(device *alsa.Device) error {
	var err error

	if err = device.Open(); err != nil {
		return err
	}
	defer device.Close()

	channels, err := device.NegotiateChannels(1, 2)
	if err != nil {
		return err
	}

	rate, err := device.NegotiateRate(44100)
	if err != nil {
		return err
	}

	format, err := device.NegotiateFormat(alsa.S16_LE, alsa.S32_LE)
	if err != nil {
		return err
	}

	buffer_size, err := device.NegotiateBufferSize(1024, 8192, 16384)
	if err != nil {
		return err
	}

	if err = device.Prepare(); err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	t := 0.0

	fmt.Printf("Negotiated parameters: %d channels, %d hz, %v, %d frame buffer\n",
		channels, rate, format, buffer_size)

	for {
		buf.Reset()

		for i := 0; i < buffer_size; i++ {
			v := math.Sin(t * 2 * math.Pi * 440) // A4
			v *= 0.1                             // make a little quieter

			switch format {
			case alsa.S16_LE:
				sample := int16(v * ((1 << 16) - 1))

				for c := 0; c < channels; c++ {
					binary.Write(buf, binary.LittleEndian, sample)
				}

			case alsa.S32_LE:
				sample := int32(v * ((1 << 32) - 1))

				for c := 0; c < channels; c++ {
					binary.Write(buf, binary.LittleEndian, sample)
				}

			default:
				return fmt.Errorf("Unhandled sample format: %v", format)
			}

			t += 1.0 / float64(rate)
		}

		if err := device.Write(buf.Bytes(), buffer_size); err != nil {
			return err
		}
	}

	return nil
}
