package main

import (
	"flag"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/pdf"
	"github.com/tdewolff/canvas/rasterizer"
	"github.com/tdewolff/canvas/svg"
	"github.com/xyaren/flussibanner/internal/api"
	"github.com/xyaren/flussibanner/internal/imggen"
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
		var img = rasterizer.Draw(getImage(), canvas.DPMM(2))
		_ = png.Encode(w, img)
	})
	http.HandleFunc("/jpeg", func(w http.ResponseWriter, r *http.Request) {
		var img = rasterizer.Draw(getImage(), canvas.DPMM(2))
		_ = jpeg.Encode(w, img, &options)
	})
	http.HandleFunc("/svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml")
		_ = svg.Writer(w, getImage())
	})
	http.HandleFunc("/pdf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/pdf")
		_ = pdf.Writer(w, getImage())
	})

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil))
}

func getImage() *canvas.Canvas {
	match, nameMap, stats := api.GetData(worldId)
	return imggen.DrawImage(match, nameMap, stats, worldId)
}
