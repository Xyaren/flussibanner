package main

import (
	"flag"
	"github.com/tdewolff/canvas"
	"github.com/xyaren/flussibanner/internal/api"
	"github.com/xyaren/flussibanner/internal/imggen"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
)

const worldId int = 2202

var options = jpeg.Options{Quality: 95}

func main() {
	portPtr := flag.Int("port", 8080, "webserverPort")
	flag.Parse()

	http.HandleFunc("/png", func(w http.ResponseWriter, r *http.Request) {
		var img image.Image = getImage().WriteImage(2.0)
		_ = png.Encode(w, img)
	})
	http.HandleFunc("/jpeg", func(w http.ResponseWriter, r *http.Request) {
		var img image.Image = getImage().WriteImage(2.0)
		_ = jpeg.Encode(w, img, &options)
	})
	http.HandleFunc("/svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml")

		img := getImage()
		svg := canvas.NewSVG(w, img.W, img.H)
		img.Render(svg)
		_ = svg.Close()
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil))
}

func getImage() *canvas.Canvas {
	match, nameMap, stats := api.GetData(worldId)
	return imggen.DrawImage(match, nameMap, stats)
}
