//NOTE: Cannot get working on my machine, please fix! Not kept up-to-date!

package impl

import (
	"math"
	"math/rand"

	"github.com/averseabfun/gochip8/engine/interfaces"
	"github.com/gordonklaus/portaudio"
)

type PortAudio struct {
	tones   map[interfaces.AudioID]*portaudio.Stream
	playing map[interfaces.AudioID]bool
}

func (pa *PortAudio) InitAudio() error {
	return portaudio.Initialize()
}

func (pa *PortAudio) DeinitAudio() error {
	return portaudio.Terminate()
}

func (pa *PortAudio) PlayTone(freqHz float64) (interfaces.AudioID, error) {
	var toneFunc = func(out [][]float32, timeInfo portaudio.StreamCallbackTimeInfo, flags portaudio.StreamCallbackFlags) {
		for i := range out[0] {
			out[0][i] = float32(math.Sin(2 * math.Pi * freqHz))
			out[1][i] = float32(math.Sin(2 * math.Pi * freqHz))
		}
	}
	var stream *portaudio.Stream
	var err error
	if stream, err = portaudio.OpenDefaultStream(0, 2, 44100, 20, toneFunc); err != nil {
		return 0, err
	}
	if err := stream.Start(); err != nil {
		return 0, err
	}
	var id = interfaces.AudioID(rand.Uint64())
	pa.tones[id] = stream
	pa.playing[id] = true
	return id, nil
}

func (pa *PortAudio) Playing(tone interfaces.AudioID) bool {
	var playing, ok = pa.playing[tone]
	if !ok {
		return false
	}
	return playing
}

func (pa *PortAudio) PauseTone(tone interfaces.AudioID) error {
	pa.playing[tone] = false
	return pa.tones[tone].Stop()
}

func (pa *PortAudio) ResumeTone(tone interfaces.AudioID) error {
	pa.playing[tone] = true
	return pa.tones[tone].Start()
}

func (pa *PortAudio) StopTone(tone interfaces.AudioID) error {
	if err := pa.tones[tone].Stop(); err != nil {
		return err
	}
	if err := pa.tones[tone].Close(); err != nil {
		return err
	}
	delete(pa.tones, tone)
	delete(pa.playing, tone)
	return nil
}
