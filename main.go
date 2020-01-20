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
	"time"
)

var robotoFont *canvas.FontFamily
var colorBackground color.RGBA
var api *gw2api.GW2Api

var colorGreen = color.RGBA{R: 94, G: 185, B: 94, A: 255}
var colorBlue = color.RGBA{R: 14, G: 144, B: 210, A: 255}
var colorRed = color.RGBA{R: 221, G: 81, B: 76, A: 255}

const logoMargin float64 = 10
const worldId int = 2202

const rowSize float64 = 40
const layout = "02.01.06 15:04:05 MST"

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

	c := canvas.New(810, rowSize*3+rowSize)
	ctx := canvas.NewContext(c)
	draw(ctx, matchWorld, worldNameMap)
	c.SavePNG("out.png", 2.0)
}

func getWorldMap(matchWorld gw2api.Match) map[int]string {
	var worlds []int
	worlds = append(worlds, matchWorld.AllWorlds.Green...)
	worlds = append(worlds, matchWorld.AllWorlds.Blue...)
	worlds = append(worlds, matchWorld.AllWorlds.Red...)
	ids, _ := api.WorldIds("de", worlds...)

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

	drawTimestamp(c)
}

func drawEmblem(c *canvas.Context) float64 {
	emblem, err := os.Open("./emblem.png")
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
	log.Println(textLine.Bounds())
	c.DrawText(c.Width()-textLine.Bounds().W-textLine.Bounds().X-2, 2, textLine)
}

func drawVp(c *canvas.Context, x float64, match gw2api.Match, nameMap map[int]string, scores gw2api.TeamAssoc, title string) float64 {
	const barHeight = rowSize/2 - (2 * 4 /* padding */)
	const barWith = 120
	const radius = 3

	standardFace := robotoFont.Face(45.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	box := canvas.NewTextBox(standardFace, title, barWith, rowSize, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(x, rowSize*3+rowSize, box)

	maxScore := float64(ix.Maxs(scores.Green, scores.Blue, scores.Red))
	totalScore := float64(scores.Green + scores.Blue + scores.Red)
	width := 0.0
	width = math.Max(width, drawScore(rowSize*2, float64(scores.Green), maxScore, totalScore, barWith, barHeight, radius, colorGreen, c, x))
	width = math.Max(width, drawScore(rowSize*1, float64(scores.Blue), maxScore, totalScore, barWith, barHeight, radius, colorBlue, c, x))
	width = math.Max(width, drawScore(rowSize*0, float64(scores.Red), maxScore, totalScore, barWith, barHeight, radius, colorRed, c, x))
	return width
}

func drawScore(offset float64, score float64, highestScore float64, totalScore float64, barWidth float64, barHeight float64, radius float64, barColor color.Color, c *canvas.Context, x float64) float64 {
	cellHeight := rowSize/2 - (2 * 2)

	//c.SetFillColor(barColor)
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

	box := canvas.NewTextBox(standardFace, text, barWidth, cellHeight, canvas.Center, canvas.Top, 0, 0)

	c.DrawText(x, offset+(rowSize/2), box)

	return barWidth
}

func drawServerNames(c *canvas.Context, currentX float64, match gw2api.Match, worldNameMap map[int]string) float64 {
	standardFace := robotoFont.Face(40.0, canvas.White, canvas.FontRegular, canvas.FontNormal)

	firstRow := drawServerName(standardFace, getName(worldNameMap, match.Worlds.Green, match.AllWorlds.Green), c, currentX, 2)
	secondRow := drawServerName(standardFace, getName(worldNameMap, match.Worlds.Blue, match.AllWorlds.Blue), c, currentX, 1)
	thirdRow := drawServerName(standardFace, getName(worldNameMap, match.Worlds.Red, match.AllWorlds.Red), c, currentX, 0)

	log.Print(math.Max(firstRow.Bounds().W, math.Max(secondRow.Bounds().W, thirdRow.Bounds().W)))
	return 120
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

func drawServerName(standardFace canvas.FontFace, text *canvas.RichText, c *canvas.Context, currentX float64, row float64) *canvas.Text {
	cell := text.ToText(120, rowSize, canvas.Right, canvas.Center, 0, 0)
	log.Println(cell.Bounds())

	c.SetFillColor(canvas.Transparent)
	c.SetStrokeColor(canvas.Red)
	c.DrawPath(currentX, rowSize*row, canvas.Rectangle(120, rowSize))

	c.DrawText(0-cell.Bounds().X, 0-cell.Bounds().Y, cell)
	c.DrawText(currentX-cell.Bounds().X+(120-cell.Bounds().W), rowSize*row-cell.Bounds().Y+rowSize/2-(cell.Bounds().H/2), cell)
	return cell
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
