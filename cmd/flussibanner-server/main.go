package main

import (
	"fmt"
	"github.com/adam-lavrik/go-imath/ix" // int-related functions
	"github.com/kofalt/go-memoize"
	"github.com/tdewolff/canvas"
	"github.com/xyaren/gw2api"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

var robotoFont *canvas.FontFamily

var api *gw2api.GW2Api

var colorBackground = color.RGBA{R: 12, G: 80, B: 127, A: 255}
var colorGreen = color.RGBA{R: 94, G: 185, B: 94, A: 255}
var colorBlue = color.RGBA{R: 14, G: 144, B: 210, A: 255}
var colorRed = color.RGBA{R: 221, G: 81, B: 76, A: 255}

const logoMargin float64 = 10
const worldId int = 2202

const rowSize float64 = 40
const bottomOffset = 10
const layout = "02.01.06 15:04:05 MST"

var cache = memoize.NewMemoizer(10*time.Second, 2*time.Minute)

var mapToHeader = make(map[string]string)

func main() {
	setupMapHeaderMapping()
	loadFonts()

	api = gw2api.NewGW2Api()

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
		matchWorld, _ := api.MatchWorld(worldId)
		worldNameMap := getWorldMap(matchWorld)
		stats, _ := api.MatchWorldStats(worldId)

		return drawImage(matchWorld, worldNameMap, stats), nil
	})
	return canv.(*canvas.Canvas)
}

func setupMapHeaderMapping() {
	mapToHeader["Center"] = "EBG"
	mapToHeader["BlueHome"] = "BBL"
	mapToHeader["GreenHome"] = "GBL"
	mapToHeader["RedHome"] = "RBL"
}

func drawImage(matchWorld gw2api.Match, worldNameMap map[int]string, stats gw2api.MatchStats) *canvas.Canvas {
	result := canvas.New(710, bottomOffset+rowSize*4)
	ctx := canvas.NewContext(result)
	draw(ctx, matchWorld, worldNameMap, stats)
	return result
}

func getWorldMap(matchWorld gw2api.Match) map[int]string {
	var worlds []int
	worlds = append(worlds, matchWorld.AllWorlds.Green...)
	worlds = append(worlds, matchWorld.AllWorlds.Blue...)
	worlds = append(worlds, matchWorld.AllWorlds.Red...)
	ids, _ := api.WorldIds("de", worlds...)

	return toMap(ids)
}

func toMap(ids []gw2api.World) map[int]string {
	var elementMap = make(map[int]string)
	for _, entry := range ids {
		elementMap[entry.ID] = entry.Name
	}
	return elementMap
}

func loadFonts() {
	robotoFont = canvas.NewFontFamily("roboto")
	if err := robotoFont.LoadFontFile("./res/fonts/Roboto-Regular.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}
}

func draw(c *canvas.Context, match gw2api.Match, worldNameMap map[int]string, stats gw2api.MatchStats) {
	fillBackground(c)

	//standardFace := robotoFont.Face(28.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	//
	//
	//
	//textFace := robotoFont.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	//drawText(c, 30.0, canvas.NewTextBox(standardFace, "Document Example", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[0], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))

	imageOffset := drawEmblem(c)
	currentX := imageOffset

	currentX += 10
	currentX += drawServerNames(c, currentX, match, worldNameMap)
	currentX += 10

	scores := match.Skirmishes[len(match.Skirmishes)-1].Scores
	currentX += drawScoreCell(c, currentX, match.VictoryPoints, "Victory Points")
	currentX += 20
	currentX += drawScoreCell(c, currentX, scores, "Current Skirmish Score")
	currentX += 20
	//currentX += drawScoreCell(c, currentX, match, worldNameMap, match.Scores, "Total Score")
	//currentX += 20
	//currentX += drawScoreCell(c, currentX, match.Kills, "Kills")
	//currentX += 20
	//currentX += drawScoreCell(c, currentX, match.Deaths, "Deaths")

	currentX += drawKillDeathRatio(c, currentX, stats)
	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[3], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
	drawTimestamp(c)
}

func drawKillDeathRatio(c *canvas.Context, x float64, stats gw2api.MatchStats) float64 {
	width := float64(35)

	standardFace := robotoFont.Face(32.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	textBox := canvas.NewTextBox(standardFace, "Kill/Death Ratio", width*5, rowSize/2, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(x, rowSize*3+rowSize/2-textBox.Bounds().Y+rowSize/2-textBox.Bounds().H/2, textBox)

	for i, wStats := range stats.Maps {
		drawCell(width, c, x, i, 0, float64(wStats.Kills.Red)/float64(wStats.Deaths.Red))
		drawCell(width, c, x, i, 1, float64(wStats.Kills.Blue)/float64(wStats.Deaths.Blue))
		drawCell(width, c, x, i, 2, float64(wStats.Kills.Green)/float64(wStats.Deaths.Green))
		drawCellHeader(width, c, x, i, 3, mapToHeader[wStats.Type])
	}

	column := 4
	drawCell(width, c, x, column, 0, float64(stats.Kills.Red)/float64(stats.Deaths.Red))
	drawCell(width, c, x, column, 1, float64(stats.Kills.Blue)/float64(stats.Deaths.Blue))
	drawCell(width, c, x, column, 2, float64(stats.Kills.Green)/float64(stats.Deaths.Green))
	drawCellHeader(width, c, x, column, 3, "Ø")

	return 0
}

func drawCell(width float64, c *canvas.Context, x float64, column int, row int, kdRatio float64) {
	c.SetStrokeColor(canvas.Red)
	c.SetFillColor(color.Transparent)
	c.SetStrokeWidth(1)
	cellOffsetY := bottomOffset + rowSize*float64(row)
	cellOffsetX := x + float64(column)*width

	//rectangle := canvas.Rectangle(width, rowSize)
	//c.DrawPath(cellOffsetX, cellOffsetY, rectangle)

	standardFace := robotoFont.Face(30.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	textBox := canvas.NewTextBox(standardFace, fmt.Sprintf("%.2f", kdRatio), width, rowSize, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(cellOffsetX, cellOffsetY-textBox.Bounds().Y+rowSize/2-textBox.Bounds().H/2, textBox)
}

func drawCellHeader(width float64, c *canvas.Context, x float64, column int, row int, text string) {
	thisRowSize := rowSize / 2
	c.SetStrokeColor(canvas.Red)
	c.SetFillColor(color.Transparent)
	c.SetStrokeWidth(1)
	cellOffsetY := bottomOffset + rowSize*float64(row)
	cellOffsetX := x + float64(column)*width

	//rectangle := canvas.Rectangle(width, thisRowSize)
	//c.DrawPath(cellOffsetX, cellOffsetY, rectangle)

	standardFace := robotoFont.Face(32.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	textBox := canvas.NewTextBox(standardFace, text, width, thisRowSize, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(cellOffsetX, cellOffsetY-textBox.Bounds().Y+thisRowSize/2-textBox.Bounds().H/2, textBox)
}

func drawEmblem(c *canvas.Context) float64 {
	emblem, err := os.Open("./res/emblem.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(emblem)
	if err != nil {
		panic(err)
	}
	imageBounds := img.Bounds().Max
	imgDPM := float64(imageBounds.Y) / (c.Height() - logoMargin*2)
	//imgWidth := float64(img.Bounds().Max.X) / imgDPM
	//imgHeight := float64(img.Bounds().Max.Y) / imgDPM
	c.DrawImage(logoMargin, 0+logoMargin, img, imgDPM)
	imageOffset := float64(img.Bounds().Max.X) / imgDPM
	return imageOffset + logoMargin
}

func drawTimestamp(c *canvas.Context) {
	standardFace := robotoFont.Face(20.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	textLine := canvas.NewTextLine(standardFace, "Generated at "+time.Now().Format(layout), canvas.Center)
	c.DrawText(c.Width()-textLine.Bounds().W-textLine.Bounds().X-2, 2, textLine)
}

func drawScoreCell(c *canvas.Context, x float64, scores gw2api.TeamAssoc, title string) float64 {
	const barHeight = rowSize/2 - (2 * 4 /* padding */)
	const barWith = 120
	const radius = 3

	standardFace := robotoFont.Face(45.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	box := canvas.NewTextBox(standardFace, title, barWith, rowSize, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(x, bottomOffset+rowSize*3+rowSize, box)

	maxScore := float64(ix.Maxs(scores.Green, scores.Blue, scores.Red))
	totalScore := float64(scores.Green + scores.Blue + scores.Red)
	width := 0.0
	width = math.Max(width, drawScore(bottomOffset+rowSize*2, float64(scores.Green), maxScore, totalScore, barWith, barHeight, radius, colorGreen, c, x))
	width = math.Max(width, drawScore(bottomOffset+rowSize*1, float64(scores.Blue), maxScore, totalScore, barWith, barHeight, radius, colorBlue, c, x))
	width = math.Max(width, drawScore(bottomOffset+rowSize*0, float64(scores.Red), maxScore, totalScore, barWith, barHeight, radius, colorRed, c, x))
	return width
}

func drawScore(offset float64, score float64, highestScore float64, totalScore float64, barWidth float64, barHeight float64, radius float64, barColor color.Color, c *canvas.Context, x float64) float64 {
	cellHeight := rowSize/2 - (2 * 2)

	//c.SetStrokeColor(barColor)
	//rectangle := canvas.Rectangle(barWidth, cellHeight)
	//c.DrawPath(x, offset+rowSize/2-rectangle.Bounds().H+1, rectangle)
	//c.DrawPath(x, offset+(rowSize/2)*2-rectangle.Bounds().H-1, rectangle)

	backgroundBar := &canvas.Path{}
	backgroundBar = canvas.RoundedRectangle(barWidth, barHeight, radius)

	c.SetStrokeWidth(1)
	c.SetFillColor(color.RGBA{50, 50, 50, 50})
	c.SetStrokeColor(canvas.White)
	c.DrawPath(x, offset+rowSize/2+backgroundBar.Bounds().H/2-1, backgroundBar)

	sizeInset := 0.5

	barGreen := &canvas.Path{}
	barGreen = canvas.RoundedRectangle((score/highestScore)*barWidth-sizeInset*2, barHeight-sizeInset*2, radius)
	barGreen = barGreen.Close()

	c.SetStrokeWidth(0)
	c.SetFillColor(barColor)
	c.SetStrokeColor(canvas.White)
	c.DrawPath(x+sizeInset, offset+rowSize/2+barGreen.Bounds().H/2+sizeInset*2-1, barGreen)

	standardFace := robotoFont.Face(30.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	text := fmt.Sprintf("%.0f (%.1f%%)", score, (score/totalScore)*100)

	box := canvas.NewTextBox(standardFace, text, barWidth, cellHeight, canvas.Center, canvas.Top, 0, 0)

	c.DrawText(x, offset+(rowSize/2)+1, box)

	return barWidth
}

func drawServerNames(c *canvas.Context, currentX float64, match gw2api.Match, worldNameMap map[int]string) float64 {
	maxWidth := float64(120)
	greenText := getName(worldNameMap, match.Worlds.Green, match.AllWorlds.Green).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	blueText := getName(worldNameMap, match.Worlds.Blue, match.AllWorlds.Blue).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	redText := getName(worldNameMap, match.Worlds.Red, match.AllWorlds.Red).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	maxBoxWidth := math.Max(greenText.Bounds().W, math.Max(blueText.Bounds().W, redText.Bounds().W))

	drawServeName(c, currentX, 2, greenText, maxBoxWidth)
	drawServeName(c, currentX, 1, blueText, maxBoxWidth)
	drawServeName(c, currentX, 0, redText, maxBoxWidth)

	return maxBoxWidth
}

func getName(nameMap map[int]string, main int, all []int) *canvas.RichText {
	standardFace := robotoFont.Face(35.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	linkFace := robotoFont.Face(25.0, canvas.White, canvas.FontRegular, canvas.FontNormal)

	text := canvas.NewRichText()
	text.Add(standardFace, nameMap[main])

	for i := range all {
		link := all[i]
		if link != main {
			text.Add(linkFace, "\n + "+nameMap[link])
		}
	}
	return text
}

func drawServeName(c *canvas.Context, currentX float64, row float64, cell *canvas.Text, width float64) {
	offsetY := bottomOffset + (rowSize * row)
	//c.SetFillColor(canvas.Transparent)
	//c.SetStrokeColor(canvas.Red)
	//c.DrawPath(currentX, offsetY, canvas.Rectangle(width, rowSize))

	c.DrawText(currentX-cell.Bounds().X+(width-cell.Bounds().W), offsetY-cell.Bounds().Y+rowSize/2-(cell.Bounds().H/2), cell)
}

func fillBackground(c *canvas.Context) {
	c.SetFillColor(colorBackground)
	p := canvas.Rectangle(c.Width(), c.Height())
	c.DrawPath(0, 0, p)
	c.Fill()
}
