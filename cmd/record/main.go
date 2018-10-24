// record some audio and save it as a WAV file
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/yobert/alsa"
)

func main() {
	var (
		channels     int
		rate         int
		duration_str string
		file         string
	)

	flag.IntVar(&channels, "channels", 2, "Channels (1 for mono, 2 for stereo)")
	flag.IntVar(&rate, "rate", 44100, "Frame rate (Hz)")
	flag.StringVar(&duration_str, "duration", "5s", "Recording duration")
	flag.StringVar(&file, "file", "out.wave", "Output file")
	flag.Parse()

	duration, err := time.ParseDuration(duration_str)
	if err != nil {
		fmt.Println("Cannot parse duration:", err)
		os.Exit(1)
	}

	cards, err := alsa.OpenCards()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer alsa.CloseCards(cards)

	// use the first recording device we find
	var recordDevice *alsa.Device

	for _, card := range cards {
		devices, err := card.Devices()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
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
		os.Exit(1)
	}
	fmt.Printf("Recording device: %v\n", recordDevice)

	recording, err := record(recordDevice, duration, channels, rate)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = save(recording, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// success!
	return
}

// record audio for given duration
func record(rec *alsa.Device, duration time.Duration, channels, rate int) (alsa.Buffer, error) {
	var err error

	if err = rec.Open(); err != nil {
		return alsa.Buffer{}, err
	}
	defer rec.Close()

	_, err = rec.NegotiateChannels(channels)
	if err != nil {
		return alsa.Buffer{}, err
	}

	_, err = rec.NegotiateRate(rate)
	if err != nil {
		return alsa.Buffer{}, err
	}

	_, err = rec.NegotiateFormat(alsa.S16_LE, alsa.S32_LE)
	if err != nil {
		return alsa.Buffer{}, err
	}

	bufferSize, err := rec.NegotiateBufferSize(8192, 16384)
	if err != nil {
		return alsa.Buffer{}, err
	}

	if err = rec.Prepare(); err != nil {
		return alsa.Buffer{}, err
	}

	buf := rec.NewBufferDuration(duration)

	fmt.Printf("Negotiated parameters: %v, %d frame buffer, %d bytes/frame\n",
		buf.Format, bufferSize, rec.BytesPerFrame())

	fmt.Printf("Recording for %s (%d frames, %d bytes)...\n", duration, len(buf.Data)/rec.BytesPerFrame(), len(buf.Data))
	err = rec.Read(buf.Data)
	if err != nil {
		return alsa.Buffer{}, err
	}
	fmt.Println("Recording stopped.")
	return buf, nil
}

// save recording to a WAV file
func save(recording alsa.Buffer, file string) error {
	of, err := os.Create(file)
	if err != nil {
		return err
	}
	defer of.Close()

	var sampleBytes int
	switch recording.Format.SampleFormat {
	case alsa.S32_LE:
		sampleBytes = 4
	case alsa.S16_LE:
		sampleBytes = 2
	default:
		return fmt.Errorf("Unhandled ALSA format %v", recording.Format.SampleFormat)
	}

	// normal uncompressed WAV format (I think)
	// https://web.archive.org/web/20080113195252/http://www.borg.com/~jglatt/tech/wave.htm
	wavformat := 1

	enc := wav.NewEncoder(of, recording.Format.Rate, sampleBytes*8, recording.Format.Channels, wavformat)

	sampleCount := len(recording.Data) / sampleBytes
	data := make([]int, sampleCount)

	format := &audio.Format{
		NumChannels: recording.Format.Channels,
		SampleRate:  recording.Format.Rate,
	}

	// Convert into the format go-audio/wav wants
	var off int
	switch recording.Format.SampleFormat {
	case alsa.S32_LE:
		inc := binary.Size(uint32(0))
		for i := 0; i < sampleCount; i++ {
			data[i] = int(binary.LittleEndian.Uint32(recording.Data[off:]))
			off += inc
		}
	case alsa.S16_LE:
		inc := binary.Size(uint16(0))
		for i := 0; i < sampleCount; i++ {
			data[i] = int(binary.LittleEndian.Uint16(recording.Data[off:]))
			off += inc
		}
	default:
		return fmt.Errorf("Unhandled ALSA format %v", recording.Format.SampleFormat)
	}

	intBuf := &audio.IntBuffer{Data: data, Format: format, SourceBitDepth: sampleBytes * 8}

	if err := enc.Write(intBuf); err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return err
	}
	fmt.Printf("Saved recording to %s\n", file)
	return nil
}
