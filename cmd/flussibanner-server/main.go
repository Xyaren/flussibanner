package main

import (
	"github.com/kofalt/go-memoize"
	"github.com/tdewolff/canvas"
	"github.com/xyaren/flussibanner/internal/api"
	"github.com/xyaren/flussibanner/internal/imggen"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"time"
)

const worldId int = 2202

var cache = memoize.NewMemoizer(10*time.Second, 2*time.Minute)

func main() {

	http.HandleFunc("/png", func(w http.ResponseWriter, r *http.Request) {
		img, _, _ := cache.Memoize("png", func() (interface{}, error) {
			return getImage().WriteImage(2.0), nil
		})
		img2 := img.(image.Image)
		_ = png.Encode(w, img2)
	})
	http.HandleFunc("/jpeg", func(w http.ResponseWriter, r *http.Request) {
		img, _, _ := cache.Memoize("jpeg", func() (interface{}, error) {
			return getImage().WriteImage(3.0), nil
		})
		img2 := img.(image.Image)
		_ = jpeg.Encode(w, img2, nil)
	})
	http.HandleFunc("/svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml")

		i := getImage()
		svg := canvas.NewSVG(w, i.W, i.H)
		i.Render(svg)
		_ = svg.Close()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getImage() *canvas.Canvas {
	canv, _, _ := cache.Memoize("img", func() (interface{}, error) {
		match, nameMap, stats := api.GetData(worldId)

		return imggen.DrawImage(match, nameMap, stats), nil
	})
	return canv.(*canvas.Canvas)
}
