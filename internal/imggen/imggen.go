package imggen

import (
	"fmt"
	"github.com/adam-lavrik/go-imath/ix"
	"github.com/tdewolff/canvas"
	"github.com/xyaren/gw2api"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
	"time"
)

var robotoFont = loadRoboto()

func loadRoboto() *canvas.FontFamily {
	fontFamily := canvas.NewFontFamily("roboto")
	loadFamily(fontFamily, "./res/fonts/Roboto-Regular.ttf", canvas.FontRegular)
	loadFamily(fontFamily, "./res/fonts/Roboto-Bold.ttf", canvas.FontBold)
	return fontFamily
}

func loadFamily(fontFamily *canvas.FontFamily, path string, style canvas.FontStyle) {
	err := fontFamily.LoadFontFile(path, style)
	if err != nil {
		panic(err)
	}
}

var colorBackground = color.RGBA{R: 12, G: 80, B: 127, A: 255}

var colorGreen = color.RGBA{R: 94, G: 185, B: 94, A: 255}
var colorBlue = color.RGBA{R: 14, G: 144, B: 210, A: 255}
var colorRed = color.RGBA{R: 221, G: 81, B: 76, A: 255}

const logoMargin float64 = 10

const rowSize float64 = 40
const bottomOffset = 10
const layout = "02.01.2006 15:04:05 MST"

var mapToHeader = setupMapHeaderMapping()

func DrawImage(matchWorld gw2api.Match, worldNameMap map[int]string, stats gw2api.MatchStats, worldId int) *canvas.Canvas {
	result := canvas.New(710, bottomOffset+rowSize*4)
	ctx := canvas.NewContext(result)
	draw(ctx, matchWorld, worldNameMap, stats, worldId)
	return result
}

func draw(c *canvas.Context, match gw2api.Match, worldNameMap map[int]string, stats gw2api.MatchStats, worldId int) {
	fillCanvas(c, colorBackground)

	//standardFace := robotoFont.Face(28.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	//
	//
	//
	//textFace := robotoFont.Face(12.0, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	//drawText(c, 30.0, canvas.NewTextBox(standardFace, "Document Example", 0.0, 0.0, canvas.Left, canvas.Top, 0.0, 0.0))
	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[0], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
	split := strings.Split(match.ID, "-")
	//region := split[0];
	tier := split[1]

	imageOffset := drawEmblem(c)
	currentX := imageOffset

	currentX += 5

	currentX += drawServerNames(c, currentX, match, worldNameMap, tier, worldId)
	currentX += 10
	currentX += drawBars(c, currentX, match.VictoryPoints, "Victory Points")
	currentX += 20
	currentX += drawBars(c, currentX, match.Skirmishes[len(match.Skirmishes)-1].Scores, "Current Skirmish Score")
	currentX += 20
	currentX += drawKillDeathRatio(c, currentX, stats)

	//drawText(c, 30.0, canvas.NewTextBox(textFace, lorem[3], 140.0, 0.0, canvas.Justify, canvas.Top, 5.0, 0.0))
	drawTimestamp(c)
}

func drawKillDeathRatio(c *canvas.Context, x float64, stats gw2api.MatchStats) float64 {
	cellWidth := float64(35)

	headerFace := robotoFont.Face(40.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	headerBox := canvas.NewTextBox(headerFace, "Kill/Death Ratio", cellWidth*5, rowSize/2, canvas.Center, canvas.Center, 0, 0)
	c.DrawText(x, rowSize*3+rowSize/2-headerBox.Bounds().Y+rowSize/2-headerBox.Bounds().H/2, headerBox)

	for i, wStats := range stats.Maps {
		kdGreen := float64(wStats.Kills.Green) / float64(wStats.Deaths.Green)
		kdBlue := float64(wStats.Kills.Blue) / float64(wStats.Deaths.Blue)
		kdRed := float64(wStats.Kills.Red) / float64(wStats.Deaths.Red)

		bestKdMap := math.Max(kdGreen, math.Max(kdBlue, kdRed))

		drawCellHeader(cellWidth, c, x, i, 3, mapToHeader[wStats.Type])
		drawCell(cellWidth, c, x, i, 2, kdGreen, color.White, kdGreen == bestKdMap)
		drawCell(cellWidth, c, x, i, 1, kdBlue, color.White, kdBlue == bestKdMap)
		drawCell(cellWidth, c, x, i, 0, kdRed, color.White, kdRed == bestKdMap)
	}
	{
		kdGreen := float64(stats.Kills.Green) / float64(stats.Deaths.Green)
		kdBlue := float64(stats.Kills.Blue) / float64(stats.Deaths.Blue)
		kdRed := float64(stats.Kills.Red) / float64(stats.Deaths.Red)

		bestKdMap := math.Max(kdGreen, math.Max(kdBlue, kdRed))

		column := 4

		drawCellHeader(cellWidth, c, x, column, 3, "Ø")
		drawCell(cellWidth, c, x, column, 2, kdGreen, color.White, kdGreen == bestKdMap)
		drawCell(cellWidth, c, x, column, 1, kdBlue, color.White, kdBlue == bestKdMap)
		drawCell(cellWidth, c, x, column, 0, kdRed, color.White, kdRed == bestKdMap)
	}
	return 0
}

func drawCell(width float64, c *canvas.Context, x float64, column int, row int, kdRatio float64, textColor color.Color, isBestOnMap bool) {
	c.SetStrokeColor(canvas.Red)
	c.SetFillColor(color.Transparent)
	c.SetStrokeWidth(1)
	cellOffsetY := bottomOffset + rowSize*float64(row)
	cellOffsetX := x + float64(column)*width

	//rectangle := canvas.Rectangle(width, rowSize)
	//c.DrawPath(cellOffsetX, cellOffsetY, rectangle)

	var fontStyle = canvas.FontRegular
	decorators := make([]canvas.FontDecorator, 0)
	if !math.IsNaN(kdRatio) && isBestOnMap {
		//fontStyle = canvas.FontBold
		decorators = append(decorators, canvas.FontUnderline)
	}
	standardFace := robotoFont.Face(35.0, textColor, fontStyle, canvas.FontNormal, decorators...)

	var text string
	if math.IsNaN(kdRatio) {
		text = "-"
	} else {
		text = fmt.Sprintf("%.2f", kdRatio)
	}

	textBox := canvas.NewTextBox(standardFace, text, width, rowSize, canvas.Center, canvas.Center, 0, 0)
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

	standardFace := robotoFont.Face(40.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
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
	loc, _ := time.LoadLocation("Europe/Berlin")
	textLine := canvas.NewTextLine(standardFace, "Generated at "+time.Now().In(loc).Format(layout)+"    © Tobi", canvas.Center)
	c.DrawText(c.Width()-textLine.Bounds().W-textLine.Bounds().X-2, 2, textLine)
}

func drawBars(c *canvas.Context, x float64, scores gw2api.TeamAssoc, title string) float64 {
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

	var text string
	if score == 0 || totalScore == 0 {
		text = fmt.Sprintf("%.0f", score)
	} else {
		text = fmt.Sprintf("%.0f (%.1f%%)", score, (score/totalScore)*100)
	}
	box := canvas.NewTextBox(standardFace, text, barWidth, cellHeight, canvas.Center, canvas.Top, 0, 0)

	c.DrawText(x, offset+(rowSize/2)+1, box)

	return barWidth
}

func drawServerNames(c *canvas.Context, currentX float64, match gw2api.Match, worldNameMap map[int]string, tier string, worldId int) float64 {
	maxWidth := float64(105)
	greenText := getName(worldNameMap, match.Worlds.Green, match.AllWorlds.Green, worldId).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	blueText := getName(worldNameMap, match.Worlds.Blue, match.AllWorlds.Blue, worldId).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	redText := getName(worldNameMap, match.Worlds.Red, match.AllWorlds.Red, worldId).ToText(maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	//maxBoxWidth := math.Max(greenText.Bounds().W, math.Max(blueText.Bounds().W, redText.Bounds().W))

	drawTier(c, currentX, tier, maxWidth)

	drawServeName(c, currentX, 2, greenText, maxWidth)
	drawServeName(c, currentX, 1, blueText, maxWidth)
	drawServeName(c, currentX, 0, redText, maxWidth)
	return maxWidth
}

func drawTier(c *canvas.Context, currentX float64, tier string, maxWidth float64) {
	standardFace := robotoFont.Face(45.0, canvas.White, canvas.FontBold, canvas.FontNormal)
	text := "Tier " + tier
	box := canvas.NewTextBox(standardFace, text, maxWidth, rowSize, canvas.Right, canvas.Center, 0, 0)
	c.DrawText(currentX, bottomOffset+rowSize*3+rowSize, box)
}

func getName(nameMap map[int]string, main int, all []int, worldId int) *canvas.RichText {
	standardFace := robotoFont.Face(35.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	standardFaceTargetWorld := robotoFont.Face(35.0, canvas.White, canvas.FontBold, canvas.FontNormal)
	linkFace := robotoFont.Face(25.0, canvas.White, canvas.FontRegular, canvas.FontNormal)
	linkFaceTargetWorld := robotoFont.Face(25.0, canvas.White, canvas.FontBold, canvas.FontNormal)

	text := canvas.NewRichText()

	if main == worldId {
		text.Add(standardFaceTargetWorld, nameMap[main])
	} else {
		text.Add(standardFace, nameMap[main])
	}

	for i := range all {
		link := all[i]
		if link != main {
			if link == worldId {
				text.Add(linkFaceTargetWorld, "\n + "+nameMap[link])
			} else {
				text.Add(linkFace, "\n + "+nameMap[link])
			}
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

func fillCanvas(c *canvas.Context, color color.Color) {
	c.Push()
	c.SetFillColor(color)
	p := canvas.Rectangle(c.Width(), c.Height())
	c.DrawPath(0, 0, p)
	c.Fill()
	c.Pop()
}

func setupMapHeaderMapping() map[string]string {
	mapToHeader := make(map[string]string)
	mapToHeader["Center"] = "EBG"
	mapToHeader["BlueHome"] = "BBL"
	mapToHeader["GreenHome"] = "GBL"
	mapToHeader["RedHome"] = "RBL"
	return mapToHeader
}
