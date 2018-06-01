package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sync"
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
			fmt.Printf("error when beeping device: %v\n", err)
		}
	}
	return nil
}

func beepDevice(device *alsa.Device) error {
	var err error

	if err = device.Open(); err != nil {
		return err
	}

	// Cleanup device when done or force cleanup after 3 seconds.
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer wg.Wait()
	childCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()
	go func(ctx context.Context) {
		defer device.Close()
		<-ctx.Done()
		fmt.Println("Closing device.")
		wg.Done()
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

	// A 50ms period is a sensible value to test low-ish latency.
	// We adjust the buffer so it's of minimal size (period * 2) since it appear ALSA won't
	// start playback until the buffer has been filled to a certain degree and the automatic
	// buffer size can be quite large.
	// Some devices only accept even periods while others want powers of 2.
	wantPeriodSize := 2048 // 46ms @ 44100Hz

	periodSize, err := device.NegotiatePeriodSize(wantPeriodSize)
	if err != nil {
		return err
	}

	bufferSize, err := device.NegotiateBufferSize(wantPeriodSize * 2)
	if err != nil {
		return err
	}

	if err = device.Prepare(); err != nil {
		return err
	}

	fmt.Printf("Negotiated parameters: %d channels, %d hz, %v, %d period size, %d buffer size\n",
		channels, rate, format, periodSize, bufferSize)

	// Play 2 seconds of beep.
	duration := 2 * time.Second
	t := time.NewTimer(duration)
	for t := 0.; t < duration.Seconds(); {
		var buf bytes.Buffer

		for i := 0; i < periodSize; i++ {
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

		if err := device.Write(buf.Bytes(), periodSize); err != nil {
			return err
		}
	}
	// Wait for playback to complete.
	<-t.C
	fmt.Printf("Playback should be complete now.\n")
	time.Sleep(1 * time.Second) // To allow a human to compare real playback end with supposed.

	return nil
}
