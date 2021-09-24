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

const worldId int = 2202

func main() {
	portPtr := flag.Int("port", 8080, "webserverPort")
	//worldId := flag.Int("worldId", 2202, "Gw2 World id (e.g. 2202)")
	flag.Parse()

	colorSpace := canvas.SRGBColorSpace{}

	resolution := canvas.DPMM(2)

	imager := imggen.NewImager()
	http.HandleFunc("/png", handler(imager, renderers.PNG(resolution, colorSpace), "image/png", false))
	http.HandleFunc("/jpeg", handler(imager, renderers.JPEG(resolution, colorSpace, &jpeg.Options{Quality: 95}), "image/jpeg", false))
	http.HandleFunc("/svg", handler(imager, renderers.SVG(svgOptions()), "image/svg+xml", true))
	http.HandleFunc("/pdf", handler(imager, renderers.PDF(pdfOptions()), "application/pdf", false))

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

func handler(imager *imggen.Imager, writer canvas.Writer, contentType string, gzipEncoded bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		if gzipEncoded {
			w.Header().Add("Content-Encoding", "gzip")
		}
		image := getImage(imager)
		err := writer(w, image)
		handleError(err)
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func getImage(imager *imggen.Imager) *canvas.Canvas {
	match, nameMap, stats := api.GetData(worldId)
	return imager.DrawImage(match, nameMap, stats, worldId)
}
