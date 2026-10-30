package shared

import (
	"bytes"
	"log"

	"github.com/ben-rw/ragin-mages/internal/gameAssets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var FontSrc = getFontSrc()

var GameFontPath = "assets/fonts/pico-8-mono.otf"

func getFontSrc() *text.GoTextFaceSource {
	fontBytes, err := gameAssets.AssetsFS.ReadFile(GameFontPath)
	if err != nil {
		log.Printf("error reading font file: %v\n", err)
	}

	fontSrc, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Printf("error created font face src: %v\n", err)
	}

	return fontSrc
}
