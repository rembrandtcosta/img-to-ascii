package cmd

import (
	"fmt"
	"image"
	"log"
	"os"

	"image/color"
	_ "image/png"
	_ "image/jpeg"

	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
)

var (
	asciiGreyScale = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	width = 0
	height = 0
)

var rootCmd = &cobra.Command{
	Use: "img-to-ascii",
	Short: "generate ascii art from image files",
	Long: "generate ascii art from image files",
	Run: func(cmd *cobra.Command, args []string) {
		printAscii(args[0])	
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&width, "width", "W", 0, "Set width in characters")
        rootCmd.PersistentFlags().IntVarP(&height, "height", "H", 0, "Set height in characters")
}


func printAscii(imagePath string) {
	reader, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err);
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	if width == 0 && height == 0 {
		width = m.Bounds().Dx();
		height = m.Bounds().Dy();
	} else {
		imgWidth := float64(m.Bounds().Dx())
		imgHeight := float64(m.Bounds().Dy())
		aspectRatio := imgWidth / imgHeight

		if width == 0 {
			width = int(float64(height) / aspectRatio)
			width = int(2 * float64(height))
		} else if height == 0 {
			height = int(float64(width) / aspectRatio)
			height = int(0.5 * float64(height))
		}
	}

	m = imaging.Resize(m, width, height, imaging.Lanczos)
	bounds := m.Bounds()

	var imgSet [][]uint32

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		var tmp []uint32
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.GrayModel.Convert(m.At(x, y))
			charDepth, _, _, _ := pixel.RGBA()
			tmp = append(tmp, charDepth)
		}
		imgSet = append(imgSet, tmp)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			charDepth := imgSet[y][x]/257
			idxF := float64(charDepth) / 255.0 * float64(len(asciiGreyScale))
			if charDepth == 255 {
				idxF = float64(len(asciiGreyScale) - 1)
			}
			idx := int(idxF)
			asciiChar := asciiGreyScale[idx]
			fmt.Printf("%c", asciiChar)
		}
		fmt.Println()
	}
}
