package sound

import (
	"bytes"
	"io"
	"log"

	"path/filepath"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 48000
const bytesPerSample = 4

var audioContext = audio.NewContext(sampleRate)

func NewAudioPlayer(path string, loop bool, introLen int64) (*audio.Player, error) {
	data, err := shared.AssetsFS.ReadFile(path)
	if err != nil {
		return &audio.Player{}, err
	}
	if filepath.Ext(path) == ".wav" {
		stream, err := wav.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
		if err != nil {
			return &audio.Player{}, err
		}

		var s io.ReadSeeker
		if loop {
			if introLen > 0 {
				s = audio.NewInfiniteLoopWithIntro(stream, introLen*bytesPerSample*sampleRate, stream.Length()-(introLen*bytesPerSample*sampleRate))
			} else {
				s = audio.NewInfiniteLoop(stream, stream.Length())
			}
		} else {
			s = stream
		}
		audioPlayer, err := audioContext.NewPlayer(s)
		return audioPlayer, nil

	} else if filepath.Ext(path) == ".ogg" {
		stream, err := vorbis.DecodeWithSampleRate(sampleRate, bytes.NewReader(data))
		if err != nil {
			return &audio.Player{}, err
		}
		var s io.ReadSeeker
		if loop {
			if introLen > 0 {
				s = audio.NewInfiniteLoopWithIntro(stream, introLen*bytesPerSample*sampleRate, stream.Length()-(introLen*bytesPerSample*sampleRate))
			} else {
				s = audio.NewInfiniteLoop(stream, stream.Length())
			}
		} else {
			s = stream
		}
		audioPlayer, err := audioContext.NewPlayer(s)
		return audioPlayer, nil
	} else {
		return &audio.Player{}, nil
	}
}

func FadeOut(audioPlayer *audio.Player) bool {
	if audioPlayer.Volume() > 0.001 {
		audioPlayer.SetVolume(audioPlayer.Volume() - 0.004)
		if audioPlayer.Volume() > 0 {
			return false
		}
	}
	err := audioPlayer.Close()
	if err != nil {
		log.Println("audio player already closed: %v", err)
	}
	return true
}
