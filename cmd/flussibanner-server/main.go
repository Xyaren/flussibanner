package main

import (
	"flag"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
	"github.com/tdewolff/canvas/renderers/pdf"
	"github.com/tdewolff/canvas/renderers/svg"
	"github.com/xyaren/flussibanner/internal/api"
	"github.com/xyaren/flussibanner/internal/imggen"
)

func main() {
	portPtr := flag.Int("port", 8080, "webserverPort")
	worldId := flag.Int("worldId", 2202, "world Id")
	flag.Parse()

	colorSpace := canvas.SRGBColorSpace{}

	resolution := canvas.DPMM(2)

	imager := imggen.NewImager()
	http.HandleFunc("/png", handler(*worldId, imager, renderers.PNG(resolution, colorSpace), "image/png", false))
	http.HandleFunc("/jpeg", handler(*worldId, imager, renderers.JPEG(resolution, colorSpace, &jpeg.Options{Quality: 95}), "image/jpeg", false))
	http.HandleFunc("/svg", handler(*worldId, imager, renderers.SVG(svgOptions()), "image/svg+xml", true))
	http.HandleFunc("/pdf", handler(*worldId, imager, renderers.PDF(pdfOptions()), "application/pdf", false))

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil))
}

func pdfOptions() *pdf.Options {
	opt := pdf.DefaultOptions
	opt.Compress = true
	opt.SubsetFonts = true
	return &opt
}

func svgOptions() *svg.Options {
	opt := svg.DefaultOptions
	opt.Compression = 2
	opt.EmbedFonts = true
	opt.SubsetFonts = true
	return &opt
}

func handler(worldId int, imager *imggen.Imager, writer canvas.Writer, contentType string, gzipEncoded bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		if gzipEncoded {
			w.Header().Add("Content-Encoding", "gzip")
		}
		image := getImage(worldId, imager)
		err := writer(w, image)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getImage(id int, imager *imggen.Imager) *canvas.Canvas {
	match, nameMap, stats := api.GetData(id)
	return imager.DrawImage(match, nameMap, stats, id)
}
