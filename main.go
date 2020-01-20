package main

import (
	"fmt"
	"github.com/adam-lavrik/go-imath/ix" // int-related functions
	"github.com/tdewolff/canvas"
	"github.com/xyaren/gw2api"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

var robotoFont *canvas.FontFamily
var colorBackground color.RGBA
var api *gw2api.GW2Api

var colorGreen color.RGBA = color.RGBA{R: 94, G: 185, B: 94, A: 255}
var colorBlue color.RGBA = color.RGBA{R: 14, G: 144, B: 210, A: 255}
var colorRed color.RGBA = color.RGBA{R: 221, G: 81, B: 76, A: 255}

const logoMargin int = 10
const worldId int = 2202

const rowSize float64 = 40

func main() {
	api = gw2api.NewGW2Api()
	matchWorld, _ := api.MatchWorld(worldId)
	elementMap := getWorldMap(matchWorld)
	log.Println(matchWorld)
	log.Println(elementMap)

	drawImage(matchWorld, elementMap)
}

func drawImage(matchWorld gw2api.Match, worldNameMap map[int]string) {
	loadFonts()
	colorBackground = ParseHexColor("#0c507f")

	c := canvas.New(1000, rowSize*3+rowSize)
	ctx := canvas.NewContext(c)
	draw(ctx, matchWorld, worldNameMap)
	c.SavePNG("out.png", 2.0)
}

func getWorldMap(matchWorld gw2api.Match) map[int]string {
	ids, _ := api.WorldIds("de", []int{matchWorld.Worlds.Green, matchWorld.Worlds.Blue, matchWorld.Worlds.Red}...)

	elementMap := toMap(ids)
	return elementMap
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
	if err := robotoFont.LoadFontFile("./fonts/Roboto-Regular.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}
}

var lorem = []string{
	`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla malesuada fringilla libero vel ultricies. Phasellus eu lobortis lorem. Phasellus eu cursus mi. Sed enim ex, ornare et velit vitae, sollicitudin volutpat dolor. Sed aliquam sit amet nisi id sodales. Aliquam erat volutpat. In hac habitasse platea dictumst. Pellentesque luctus varius nibh sit amet porta. Vivamus tempus, enim ut sodales aliquet, magna massa viverra eros, nec gravida risus ipsum a erat. Etiam dapibus sem augue, at porta nisi dictum non. Vestibulum quis urna ut ligula dapibus mollis eu vel nisl. Vestibulum lorem dolor, eleifend lacinia fringilla eu, pulvinar vitae metus.`,
	`Morbi dapibus purus vel erat auctor, vehicula tempus leo maximus. Aenean feugiat vel quam sit amet iaculis. Fusce et justo nec arcu maximus porttitor. Cras sed aliquam ipsum. Sed molestie mauris nec dui interdum sollicitudin. Nulla id egestas massa. Fusce congue ante. Interdum et malesuada fames ac ante ipsum primis in faucibus. Praesent faucibus tellus eu viverra blandit. Vivamus mi massa, hendrerit in commodo et, luctus vitae felis.`,
	`Quisque semper aliquet augue, in dignissim eros cursus eu. Pellentesque suscipit consequat nibh, sit amet ultricies risus. Suspendisse blandit interdum tortor, consectetur tristique magna aliquet eu. Aliquam sollicitudin eleifend sapien, in pretium nisi. Sed tempor eleifend velit quis vulputate. Donec condimentum, lectus vel viverra pharetra, ex enim cursus metus, quis luctus est urna ut purus. Donec tempus gravida pharetra. Sed leo nibh, cursus at hendrerit at, ultricies a dui. Maecenas eget elit magna. Quisque sollicitudin odio erat, sed consequat libero tincidunt in. Nullam imperdiet, neque quis consequat pellentesque, metus nisl consectetur eros, ut vehicula dui augue sed tellus.`,
	//` Vivamus varius ex sed nisi vestibulum, sit amet tincidunt ante vestibulum. Nullam et augue blandit dolor accumsan tempus. Quisque at dictum elit, id ullamcorper dolor. Nullam feugiat mauris eu aliquam accumsan.`,
}

var y = 205.0

func drawText(c *canvas.Context, x float64, text *canvas.Text) {
	h := text.Bounds().H
	c.DrawText(x, y, text)
	y -= h + 10.0
}

func draw(c *canvas.Context, match gw2api.Match, worldNameMap map[int]string) {
	fillBackground(c)
	c.SetFillColor(canvas.White)

	//standardFace := robotoFont.Face(28.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	//
	//
	//
	//textFace := robotoFont.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	//drawText(c, 30.0, canvas.NewTextBox(standardFace, "Document Example", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[0], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))

	lenna, err := os.Open("./emblem.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(lenna)
	if err != nil {
		panic(err)
	}

	imageBounds := img.Bounds().Max
	log.Print(imageBounds)
	log.Print(c.Height())

	currentX := 0.0
	currentX += float64(logoMargin)

	imgDPM := float64(imageBounds.Y) / (c.Height() - float64(logoMargin*2))
	//imgWidth := float64(img.Bounds().Max.X) / imgDPM
	//imgHeight := float64(img.Bounds().Max.Y) / imgDPM
	c.DrawImage(currentX, float64(0+logoMargin), img, imgDPM)
	currentX += float64(img.Bounds().Max.X) / imgDPM

	currentX += 10
	currentX += drawServerNames(c, currentX, match, worldNameMap)
	currentX += 10

	scores := match.Skirmishes[len(match.Skirmishes)-1].Scores
	currentX += drawVp(c, currentX, match, worldNameMap, scores, "Skirmish Scores")
	currentX += 20
	currentX += drawVp(c, currentX, match, worldNameMap, match.VictoryPoints, "Victory Points")
	currentX += 20
	//currentX += drawVp(c, currentX, match, worldNameMap, match.Scores, "Total Score")
	//currentX += 20
	currentX += drawVp(c, currentX, match, worldNameMap, match.Kills, "Kills")
	currentX += 20
	currentX += drawVp(c, currentX, match, worldNameMap, match.Deaths, "Deaths")

	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[3], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
}

func drawVp(c *canvas.Context, x float64, match gw2api.Match, nameMap map[int]string, scores gw2api.TeamAssoc, title string) float64 {
	const barHeight = rowSize/2 - (2 * 4 /* padding */)
	const barWith = 120
	const radius = 3

	standardFace := robotoFont.Face(45.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	box := canvas.NewTextBox(standardFace, title, barWith, rowSize, canvas.Center, canvas.Center, 0, 0)
	log.Println(box.Bounds())
	c.DrawText(x, rowSize*3+rowSize, box)

	maxScore := float64(ix.Maxs(scores.Green, scores.Blue, scores.Red))
	totalScore := float64(scores.Green + scores.Blue + scores.Red)
	width := 0.0
	width = math.Max(width, drawScore(rowSize*2, float64(scores.Green), maxScore, totalScore, barWith, barHeight, radius, colorGreen, c, x))
	width = math.Max(width, drawScore(rowSize*1, float64(scores.Blue), maxScore, totalScore, barWith, barHeight, radius, colorBlue, c, x))
	width = math.Max(width, drawScore(rowSize*0, float64(scores.Red), maxScore, totalScore, barWith, barHeight, radius, colorRed, c, x))
	return width
}

func drawScore(offset float64, score float64, highestScore float64,totalScore float64, barWidth float64, barHeight float64, radius float64, barColor color.Color, c *canvas.Context, x float64) float64 {
	cellHeight := rowSize/2 - (2 * 2)
	//
	//rectangle := canvas.Rectangle(barWidth, cellHeight)
	//c.DrawPath(x, offset+(rowSize/2)-rectangle.Bounds().H, rectangle)
	//c.DrawPath(x, offset+(rowSize/2)*2-rectangle.Bounds().H, rectangle)

	backgroundBar := &canvas.Path{}
	backgroundBar = canvas.RoundedRectangle(barWidth, barHeight, radius)

	c.SetStrokeWidth(1)
	c.SetFillColor(color.RGBA{50, 50, 50, 50})
	c.SetStrokeColor(canvas.White)
	c.DrawPath(x, offset+rowSize/2+backgroundBar.Bounds().H/2, backgroundBar)

	sizeOffset := 0.5

	barGreen := &canvas.Path{}
	barGreen = canvas.RoundedRectangle((score/highestScore)*barWidth-sizeOffset*2, barHeight-sizeOffset*2, radius)
	barGreen = barGreen.Close()

	c.SetStrokeWidth(0)
	c.SetFillColor(barColor)
	c.SetStrokeColor(canvas.White)
	c.DrawPath(x+sizeOffset, offset+rowSize/2+sizeOffset*2+barGreen.Bounds().H/2, barGreen)

	standardFace := robotoFont.Face(30.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	text := fmt.Sprintf("%.0f (%.1f%%)", score, (score/totalScore)*100)

	box := canvas.NewTextBox(standardFace, text, barWidth, cellHeight, canvas.Center, canvas.Center, 0, 0)

	c.DrawText(x, offset+(rowSize/2), box)

	return barWidth
}

func drawServerNames(c *canvas.Context, currentX float64, match gw2api.Match, worldNameMap map[int]string) float64 {
	standardFace := robotoFont.Face(45.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	firstRow := canvas.NewTextBox(standardFace, worldNameMap[match.Worlds.Green], 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	c.DrawText(currentX, rowSize*2+(50/2+(firstRow.Bounds().H/2)), firstRow)

	secondRow := canvas.NewTextBox(standardFace, worldNameMap[match.Worlds.Blue], 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	c.DrawText(currentX, rowSize*1+(50/2+(secondRow.Bounds().H/2)), secondRow)

	thirdRow := canvas.NewTextBox(standardFace, worldNameMap[match.Worlds.Red], 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0)
	c.DrawText(currentX, rowSize*0+(50/2+(thirdRow.Bounds().H/2)), thirdRow)

	return math.Max(firstRow.Bounds().W, math.Max(secondRow.Bounds().W, thirdRow.Bounds().W))
}

func fillBackground(c *canvas.Context) {
	c.SetFillColor(colorBackground)
	p := canvas.Rectangle(c.Width(), c.Height())
	c.DrawPath(0, 0, p)
	c.Fill()
}

func ParseHexColor(s string) (c color.RGBA) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, _ = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, _ = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		panic(fmt.Errorf("invalid length, must be 7 or 4"))

	}
	return
}
