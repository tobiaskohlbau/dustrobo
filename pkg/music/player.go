package music

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/pkg/errors"
	"github.com/tobiaskohlbau/dustrobo/pkg/wav"
)

const bufferSize = 4 * 96

type Player struct {
	play   *oto.Player
	input  io.ReadCloser
	buf    []byte
	closer io.Closer // @TODO(koht): close underlying input when closing music source

	mx     sync.RWMutex
	volume float64
	pause  bool
}

func NewMP3Player(r io.ReadCloser) (*Player, error) {
	dec, err := mp3.NewDecoder(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mp3 decoder")
	}
	play, err := oto.NewPlayer(dec.SampleRate(), 2, 2, 8192)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating oto")
	}
	return &Player{
		input:  dec,
		play:   play,
		closer: r,
		buf:    make([]byte, bufferSize),
		volume: 0.5,
	}, nil
}

func NewWAVPlayer(r io.ReadCloser) (*Player, error) {
	dec, err := wav.NewDecoder(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create wav decoder")
	}
	play, err := oto.NewPlayer(44100, 2, 2, 8192)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating oto")
	}
	return &Player{
		input:  dec,
		play:   play,
		closer: r,
		buf:    make([]byte, bufferSize),
		volume: 0.5,
	}, nil
}

func (p *Player) Play() {
	go func() {
		if err := p.loop(); err != nil {
			log.Println(err)
		}
		if err := p.input.Close(); err != nil {
			log.Println("failed to close decoder:", err)
		}
		if err := p.play.Close(); err != nil {
			log.Println("failed to close player:", err)
		}
	}()
}

func (p *Player) Pause() error {
	p.mx.Lock()
	p.pause = !p.pause
	p.mx.Unlock()
	return nil
}

func (p *Player) Stop() error {
	if err := p.play.Close(); err != nil {
		return errors.Wrap(err, "player")
	}
	if err := p.input.Close(); err != nil {
		return errors.Wrap(err, "player")
	}
	if err := p.closer.Close(); err != nil {
		return errors.Wrap(err, "player")
	}
	return nil
}

func (p *Player) Volume(value float64) error {
	p.mx.Lock()
	p.volume = value
	p.mx.Unlock()
	return nil
}

func (p *Player) loop() error {
	sample := make([]int16, 2)
	for {
		p.mx.Lock()
		pause := p.pause
		p.mx.Unlock()
		if pause {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		n, err := p.input.Read(p.buf)
		if err != nil {
			return errors.Wrap(err, "failed to read all samples")
		}
		if n != bufferSize {
			return errors.Wrap(err, "missing bytes while reading samples")
		}

		for i := 0; i < bufferSize; i += 4 {
			sample[0] = int16(p.buf[i]) | (int16(p.buf[i+1]) << 8)
			sample[1] = int16(p.buf[i+2]) | (int16(p.buf[i+3]) << 8)

			p.mx.Lock()
			sample[0] = int16(float64(sample[0]) * p.volume)
			sample[1] = int16(float64(sample[0]) * p.volume)
			p.mx.Unlock()

			p.buf[i] = byte(sample[0])
			p.buf[i+1] = byte(sample[0] >> 8)
			p.buf[i+2] = byte(sample[1])
			p.buf[i+3] = byte(sample[1] >> 8)
		}

		n, err = p.play.Write(p.buf)
		if err != nil {
			return errors.Wrap(err, "failed to write all samples")
		}
		if n != bufferSize {
			return errors.Wrap(err, "missing bytes while writing samples")
		}
	}
	return nil
}
