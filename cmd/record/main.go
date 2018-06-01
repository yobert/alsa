// record some audio and save it as a WAV file
package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/yobert/alsa"
)

const (
	bitDepth = 32
	intSize  = 4
)

var (
	rate     = 44100
	channels = 2
)

func main() {
	var duration = 5
	var file = "out.wav"

	flag.IntVar(&rate, "rate", 44100, "Frame rate (Hz)")
	flag.IntVar(&duration, "duration", 5, "Recording duration (s)")
	flag.StringVar(&file, "file", "out.wave", "Output file")
	flag.Parse()

	cards, err := alsa.OpenCards()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer alsa.CloseCards(cards)

	// use the first recording device we find
	var recordDevice *alsa.Device

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
			if device.Record && recordDevice == nil {
				recordDevice = device
			}
		}
	}

	if recordDevice == nil {
		fmt.Println("No recording device found")
		return
	}
	fmt.Printf("Recording device: %v\n", recordDevice)

	var recording []byte
	recording, err = record(recordDevice, duration)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = save(recording, file)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// record audio for duration seconds
// Side effects: channels and rate may differ from requested values
func record(rec *alsa.Device, duration int) ([]byte, error) {
	var err error

	if err = rec.Open(); err != nil {
		return nil, err
	}
	defer rec.Close()

	channels, err := rec.NegotiateChannels(channels)
	if err != nil {
		return nil, err
	}

	rate, err := rec.NegotiateRate(rate)
	if err != nil {
		return nil, err
	}

	format, err := rec.NegotiateFormat(alsa.S32_LE)
	if err != nil {
		return nil, err
	}

	bufferSize, err := rec.NegotiateBufferSize(8192, 16384)
	if err != nil {
		return nil, err
	}

	if err = rec.Prepare(); err != nil {
		return nil, err
	}

	fmt.Printf("Negotiated parameters: %d channels, %d hz, %v, %d frame buffer, %d bytes/frame\n",
		channels, rate, format, bufferSize, rec.BytesPerFrame())

	recFrames := rate * duration
	recBytes := recFrames * rec.BytesPerFrame()

	buf := make([]byte, recBytes)
	fmt.Printf("Recording for %d seconds (%d frames, %d bytes)...\n", duration, recFrames, recBytes)
	err = rec.Read(buf, recFrames)
	if err != nil {
		return nil, err
	}
	fmt.Println("Recording stopped.")
	return buf, nil
}

// save recording to a WAV file
func save(recording []byte, file string) error {
	of, err := os.Create(file)
	if err != nil {
		return err
	}
	defer of.Close()

	enc := wav.NewEncoder(of, rate, bitDepth, channels, 1)

	// insert the recording data into an audio.IntBuffer
	// is there a cleaner way?
	var data []int
	dataSize := len(recording) / intSize
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = uintptr(unsafe.Pointer(&recording[0]))
	sh.Len = dataSize
	sh.Cap = dataSize
	format := &audio.Format{
		NumChannels: channels,
		SampleRate:  rate,
	}
	intBuf := &audio.IntBuffer{Data: data, Format: format, SourceBitDepth: bitDepth}

	if err := enc.Write(intBuf); err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return err
	}
	fmt.Printf("Saved recording to %s\n", file)
	return nil
}
