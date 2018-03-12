package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"time"

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

		if err := beepCard(card); err != nil {
			fmt.Printf("error when beeping card: %v\n", err)
		}
	}
}

func beepCard(card *alsa.Card) error {
	devices, err := card.Devices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		if device.Type != alsa.PCM || !device.Play {
			continue
		}
		fmt.Println("───", device)

		if err := beepDevice(device); err != nil {
			return err
		}
	}
	return nil
}

func beepDevice(device *alsa.Device) error {
	var err error

	if err = device.Open(); err != nil {
		return err
	}

	// Cleanup device when done or force cleanup after 4 seconds.
	childCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(4*time.Second))
	defer cancel()
	go func(ctx context.Context) {
		defer device.Close()
		<-ctx.Done()
		fmt.Println("Closing device.")
	}(childCtx)

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

	// Prefer larger buffer to avoid under-run.
	bufferSize, err := device.NegotiateBufferSize(16384, 8192)
	if err != nil {
		return err
	}

	if err = device.Prepare(); err != nil {
		return err
	}

	fmt.Printf("Negotiated parameters: %d channels, %d hz, %v, %d frame buffer\n",
		channels, rate, format, bufferSize)

	// Play 2 seconds of beep.
	for t := 0.; t < 2; {
		var buf bytes.Buffer

		for i := 0; i < bufferSize; i++ {
			v := math.Sin(t * 2 * math.Pi * 440) // A4
			v *= 0.1                             // make a little quieter

			switch format {
			case alsa.S16_LE:
				sample := int16(v * math.MaxInt16)

				for c := 0; c < channels; c++ {
					binary.Write(&buf, binary.LittleEndian, sample)
				}

			case alsa.S32_LE:
				sample := int32(v * math.MaxInt32)

				for c := 0; c < channels; c++ {
					binary.Write(&buf, binary.LittleEndian, sample)
				}

			default:
				return fmt.Errorf("Unhandled sample format: %v", format)
			}

			t += 1 / float64(rate)
		}

		if err := device.Write(buf.Bytes(), bufferSize); err != nil {
			return err
		}
	}

	return nil
}
