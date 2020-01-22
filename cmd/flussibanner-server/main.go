package main

import (
	"flag"
	"github.com/kofalt/go-memoize"
	"github.com/tdewolff/canvas"
	"github.com/xyaren/flussibanner/internal/api"
	"github.com/xyaren/flussibanner/internal/imggen"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"
)

const worldId int = 2202

var options = jpeg.Options{Quality: 100}

var apiCache = memoize.NewMemoizer(10*time.Second, 2*time.Minute)

func main() {
	cachePtr := flag.Bool("enableImageCache", false, "enable caching of the generated images (use if not proxy-caches by webserver)")
	portPtr := flag.Int("port", 8080, "webserverPort")

	flag.Parse()

	var cache *memoize.Memoizer
	if *cachePtr {
		cacheDuration := 10 * time.Second
		cache = memoize.NewMemoizer(cacheDuration, cacheDuration*10)
	}

	http.HandleFunc("/png", func(w http.ResponseWriter, r *http.Request) {
		var img image.Image
		if cache != nil {
			cachedImage, _, _ := cache.Memoize("png", func() (interface{}, error) {
				return getImage().WriteImage(2.0), nil
			})
			img = cachedImage.(image.Image)
		} else {
			img = getImage().WriteImage(2.0)
		}
		_ = png.Encode(w, img)
	})
	http.HandleFunc("/jpeg", func(w http.ResponseWriter, r *http.Request) {
		var img image.Image
		if cache != nil {
			cachedImage, _, _ := cache.Memoize("jpeg", func() (interface{}, error) {
				return getImage().WriteImage(2.0), nil
			})
			img = cachedImage.(image.Image)
		} else {
			img = getImage().WriteImage(2.0)
		}
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
	canv, _, _ := apiCache.Memoize("img", func() (interface{}, error) {
		match, nameMap, stats := api.GetData(worldId)

		return imggen.DrawImage(match, nameMap, stats), nil
	})
	return canv.(*canvas.Canvas)
}
