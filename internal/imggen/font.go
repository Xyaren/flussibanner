package imggen

import (
	_ "embed"

	"github.com/tdewolff/canvas"
)

//go:embed res/fonts/Roboto-Regular.ttf
var robotoRegular []byte

//go:embed res/fonts/Roboto-Bold.ttf
var robotoBold []byte

var roboto = loadRoboto()

func loadRoboto() *canvas.FontFamily {
	fontFamily := canvas.NewFontFamily("roboto")
	loadFamily(fontFamily, robotoRegular, canvas.FontRegular)
	loadFamily(fontFamily, robotoBold, canvas.FontBold)
	return fontFamily
}

func loadFamily(fontFamily *canvas.FontFamily, data []byte, style canvas.FontStyle) {
	err := fontFamily.LoadFont(data, 0, style)
	if err != nil {
		panic(err)
	}
}
