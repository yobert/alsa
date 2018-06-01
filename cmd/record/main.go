// record some audio and save it as a WAV file
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/yobert/alsa"
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

	recording, err := record(recordDevice, duration)
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
func record(rec *alsa.Device, duration int) (alsa.Buffer, error) {
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

	buf := rec.NewBufferSeconds(duration)

	fmt.Printf("Negotiated parameters: %v, %d frame buffer, %d bytes/frame\n",
		buf.Format, bufferSize, rec.BytesPerFrame())

	fmt.Printf("Recording for %d seconds (%d frames, %d bytes)...\n", duration, len(buf.Data)/rec.BytesPerFrame(), len(buf.Data))
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
