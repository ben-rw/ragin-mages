package sound

import (
	"bytes"
	"io"

	"github.com/ben-rw/ragin-mages/cmd/game/internal/shared"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const sampleRate = 48000
const bytesPerSample = 4

var audioContext = audio.NewContext(sampleRate)

func NewAudioPlayer(filepath string, loop bool, introLen int64) (*audio.Player, error) {
	data, err := shared.AssetsFS.ReadFile(filepath)
	if err != nil {
		return &audio.Player{}, err
	}
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
}

func FadeOut(audioPlayer *audio.Player) bool {
	if audioPlayer.Volume() > 0.05 {
		audioPlayer.SetVolume(audioPlayer.Volume() - 0.01)
		return false
	} else {
		audioPlayer.Close()
		return true
	}
}
