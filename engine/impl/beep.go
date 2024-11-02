package impl

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/averseabfun/gochip8/engine/interfaces"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
)

type Beep struct {
	tones      map[interfaces.AudioID]beep.Streamer
	playing    map[interfaces.AudioID]bool
	initalized bool
}

func (b *Beep) InitAudio() error {
	if b.initalized {
		return nil
	}
	sr := beep.SampleRate(44100)

	err := speaker.Init(sr, sr.N(time.Second/10))
	if err != nil {
		return err
	}
	b.tones = make(map[interfaces.AudioID]beep.Streamer)
	b.playing = make(map[interfaces.AudioID]bool)
	b.initalized = true
	return nil
}

func (b *Beep) DeinitAudio() error {
	speaker.Clear()
	speaker.Close()
	return nil
}

type sineGenerator struct {
	dt float64
	t  float64
}

// Creates a streamer which will procude an infinite sine wave with the given frequency.
// use other wrappers of this package to change amplitude or add time limit.
// sampleRate must be at least two times grater then frequency, otherwise this function will return an error.
func sineTone(sr beep.SampleRate, freq float64) (beep.Streamer, error) {
	dt := freq / float64(sr)

	if dt >= 1.0/2.0 {
		return nil, errors.New("gopxl sine tone generator: samplerate must be at least 2 times grater then frequency")
	}

	return &sineGenerator{dt, 0}, nil
}

func (g *sineGenerator) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		v := math.Sin(g.t * 2.0 * math.Pi)
		samples[i][0] = v
		samples[i][1] = v
		_, g.t = math.Modf(g.t + g.dt)
	}

	return len(samples), true
}

func (*sineGenerator) Err() error {
	return nil
}

func (b *Beep) PlayTone(freqHz float64) (interfaces.AudioID, error) {
	b.InitAudio()
	var ts, err = sineTone(beep.SampleRate(44100), freqHz)
	if err != nil {
		return 0, err
	}
	speaker.Play(ts)

	var id = interfaces.AudioID(rand.Uint64())
	b.tones[id] = ts
	b.playing[id] = true

	return id, nil
}

func (b *Beep) Playing(tone interfaces.AudioID) bool {
	b.InitAudio()
	var playing, ok = b.playing[tone]
	if !ok {
		return false
	}
	return playing
}

func (b *Beep) PauseTone(tone interfaces.AudioID) error {
	b.InitAudio()
	b.playing[tone] = false
	speaker.Clear()

	for id, t := range b.tones {
		if id == tone {
			continue
		}
		speaker.Play(t)
	}
	return nil
}

func (b *Beep) ResumeTone(tone interfaces.AudioID) error {
	b.InitAudio()
	b.playing[tone] = true
	speaker.Play(b.tones[tone])
	return nil
}

func (b *Beep) StopTone(tone interfaces.AudioID) error {
	b.InitAudio()
	b.PauseTone(tone)
	delete(b.tones, tone)
	delete(b.playing, tone)
	return nil
}

func (b *Beep) StopAll() error {
	speaker.Clear()
	return nil
}
