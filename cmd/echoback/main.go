package main

import (
	"fmt"
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

	var (
		record_device   *alsa.Device
		playback_device *alsa.Device
	)

	// just use the "first" playback and record devices

	for _, card := range cards {
		devices, err := card.Devices()
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, device := range devices {
			if device.Type != alsa.PCM {
				continue
			}
			if device.Play && playback_device == nil {
				playback_device = device
			}
			if device.Record && record_device == nil {
				record_device = device
			}
		}
	}

	err = echoback(record_device, playback_device, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func echoback(rec, play *alsa.Device, delay time.Duration) error {
	var err error

	if err = rec.Open(); err != nil {
		return err
	}
	defer rec.Close()

	if err = play.Open(); err != nil {
		return err
	}
	defer play.Close()

	play.Debug(true)

	channels, err := rec.NegotiateChannels(2)
	if err != nil {
		return err
	}
	_, err = play.NegotiateChannels(channels)
	if err != nil {
		return err
	}

	rate, err := rec.NegotiateRate(44100)
	if err != nil {
		return err
	}
	_, err = play.NegotiateRate(rate)
	if err != nil {
		return err
	}

	format, err := rec.NegotiateFormat(alsa.S32_LE)
	if err != nil {
		return err
	}
	_, err = play.NegotiateFormat(format)
	if err != nil {
		return err
	}

	buffer_size, err := rec.NegotiateBufferSize(8192, 16384)
	if err != nil {
		return err
	}
	_, err = play.NegotiateBufferSize(buffer_size)
	if err != nil {
		return err
	}

	if err = rec.Prepare(); err != nil {
		return err
	}
	if err = play.Prepare(); err != nil {
		return err
	}

	bytes_per_frame := rec.BytesPerFrame()

	fmt.Printf("Negotiated parameters: %d channels, %d hz, %v, %d frame buffer, %d bytes per frame\n",
		channels, rate, format, buffer_size, bytes_per_frame)

	delay_bytes := int(float64(rate) * float64(bytes_per_frame) * delay.Seconds())
	buffer_bytes := bytes_per_frame * buffer_size

	// a sloppy circular buffer shared between reader and
	// writer goroutines
	buf := make([]byte,
		delay_bytes+
			// enough additional space for a read and a write
			(buffer_size*2*bytes_per_frame))

	fmt.Println("buffer_bytes\t", buffer_bytes)
	fmt.Println("delay_bytes\t", delay_bytes)
	fmt.Println("len(buf)\t", len(buf))

	var cursor_mu sync.Mutex
	rec_cursor := 0
	play_cursor := 0

	go func() {
		for {
			cursor_mu.Lock()
			avail := play_cursor - rec_cursor
			if avail <= 0 {
				avail = len(buf) - rec_cursor
				if avail == 0 && play_cursor > 0 {
					// loop around
					rec_cursor = 0
					avail = play_cursor
				}
			}
			if avail > buffer_bytes {
				// cap read size to alsa buffer size
				avail = buffer_bytes
			}
			rc := rec_cursor
			cursor_mu.Unlock()

			if avail == 0 {
				fmt.Println(" rec: waiting")
				time.Sleep(time.Millisecond * 100)
				continue
			}

			fmt.Println("reading...")
			err := rec.Read(buf[rc : rc+avail])
			if err != nil {
				fmt.Println(err)
				return
			}

			cursor_mu.Lock()
			rec_cursor += avail
			if rec_cursor == len(buf) {
				rec_cursor = 0
			} else if rec_cursor > len(buf) {
				panic("wtf, rec cursor overrun")
			}
			fmt.Printf("%d\t%d\n", rec_cursor, play_cursor)
			cursor_mu.Unlock()
		}
	}()

	go func() {
		for {
			cursor_mu.Lock()
			avail := rec_cursor - play_cursor
			actual := avail
			if avail < 0 {
				avail = len(buf) - play_cursor
				actual = avail
				avail += rec_cursor
			}
			pc := play_cursor
			cursor_mu.Unlock()

			if avail <= delay_bytes {
				fmt.Println("play: waiting")
				time.Sleep(time.Millisecond * 100)
				continue
			}

			frames_actual := actual / bytes_per_frame

			fmt.Println("writing...")
			err = play.Write(buf[pc:pc+actual], frames_actual)
			if err != nil {
				fmt.Println(err)
				return
			}

			cursor_mu.Lock()
			play_cursor += actual
			if play_cursor == len(buf) {
				play_cursor = 0
			} else if play_cursor > len(buf) {
				panic("wtf, play cursor overrun")
			}
			fmt.Printf("%d\t%d\n", rec_cursor, play_cursor)
			cursor_mu.Unlock()
		}
	}()

	select {}
}
