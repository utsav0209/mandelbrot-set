package main

import (
	"flag"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/pkg/profile"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

type Config struct {
	MaxIterations, Width, Height int
	XMin, XMax, YMin, YMax       float64
	Out, Mode                    string
}

var config Config

func main() {
	defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()

	var (
		maxIterations = flag.Int("maxIterations", 100, "Maximum number of number of Mandelbrot iterations")
		width         = flag.Int("width", 1000, "Width of the output image")
		height        = flag.Int("height", 1000, "Height of the output image")
		xMin          = flag.Float64("xMin", -2, "Minimum value of X painted on image")
		xMax          = flag.Float64("xMax", 1, "Maximum value of X painted on image")
		yMin          = flag.Float64("yMin", -2, "Minimum value of Y painted on image")
		yMax          = flag.Float64("yMax", 2, "Maximum value of Y painted on image")
		out           = flag.String("out", "output.png", "Filename of the output image")
		mode          = flag.String("mode", "seq", "Image generation method: `seq|px|row`")
	)
	flag.Parse()

	config = Config{
		MaxIterations: *maxIterations,
		Width:         *width,
		Height:        *height,
		XMin:          *xMin,
		XMax:          *xMax,
		YMin:          *yMin,
		YMax:          *yMax,
		Out:           *out,
		Mode:          *mode,
	}

	img := image.NewRGBA(
		image.Rect(0, 0, config.Width, config.Height),
	)

	switch *mode {
	case "seq":
		drawSequentially(img)
	case "px":
		drawPixelByPixel(img)
	case "row":
		drawRowByRow(img)
	default:
		panic("Unknown mode!")
	}

	err := saveImage(img, *out)
	if err != nil {
		panic(err)
	}
}

func drawSequentially(img *image.RGBA) {
	for row := 0; row < config.Width; row++ {
		for col := 0; col < config.Height; col++ {
			paintPixel(row, col, img)
		}
	}
}

func drawPixelByPixel(img *image.RGBA) {
	var wg sync.WaitGroup
	wg.Add(config.Width * config.Height)

	for row := 0; row < config.Width; row++ {
		for col := 0; col < config.Height; col++ {
			go func(row, col int) {
				paintPixel(row, col, img)
				wg.Done()
			}(row, col)
		}
	}

	wg.Wait()
}

func drawRowByRow(img *image.RGBA) {
	var wg sync.WaitGroup
	wg.Add(config.Width)

	for row := 0; row < config.Width; row++ {
		go func(row int) {
			for col := 0; col < config.Height; col++ {
				paintPixel(row, col, img)
			}
			wg.Done()
		}(row)
	}

	wg.Wait()
}

func paintPixel(row int, col int, img *image.RGBA) {
	c := getComplexNumberFromCoOrdinates(row, col)
	iterations := mandelBrot(c, config.MaxIterations)
	img.Set(row, col, computeColor(iterations, config.MaxIterations))
}

func getComplexNumberFromCoOrdinates(row, col int) complex128 {
	return complex(
		config.XMin+(float64(row)*(config.XMax-config.XMin))/float64(config.Width),
		config.YMin+(float64(col)*(config.YMax-config.YMin))/float64(config.Height),
	)
}

func computeColor(iterations int, maxIterations int) (c color.Color) {
	hue := (float64(iterations) / float64(maxIterations)) * 360
	saturation := 1.0
	value := 0.5
	if iterations == maxIterations {
		value = 0
	}
	c = colorful.Hsl(hue, saturation, value)
	return
}

func saveImage(img *image.RGBA, outputPath string) (err error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	return
}
