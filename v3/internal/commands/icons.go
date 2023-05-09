package commands

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/jackmordaunt/icns/v2"
	"github.com/leaanthony/winicon"
)

type IconsOptions struct {
	Example         bool   `description:"Generate example icon file (appicon.png) in the current directory"`
	Input           string `description:"The input image file"`
	Sizes           string `description:"The sizes to generate in .ico file (comma separated)" default:"256,128,64,48,32,16"`
	WindowsFilename string `description:"The output filename for the Windows icon" default:"icon.ico"`
	MacFilename     string `description:"The output filename for the Mac icon bundle" default:"icons.icns"`
}

func GenerateIcons(options *IconsOptions) error {

	if options.Example {
		return generateExampleIcon()
	}

	if options.Input == "" {
		return fmt.Errorf("input is required")
	}

	if options.WindowsFilename == "" && options.MacFilename == "" {
		return fmt.Errorf("at least one output filename is required")
	}

	// Parse sizes
	var sizes = []int{256, 128, 64, 48, 32, 16}
	var err error
	if options.Sizes != "" {
		sizes, err = parseSizes(options.Sizes)
		if err != nil {
			return err
		}
	}
	iconData, err := os.ReadFile(options.Input)
	if err != nil {
		return err
	}

	if options.WindowsFilename != "" {
		err := generateWindowsIcon(iconData, sizes, options)
		if err != nil {
			return err
		}
	}

	if options.MacFilename != "" {
		err := generateMacIcon(iconData, options)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateExampleIcon() error {
	return os.WriteFile("appicon.png", []byte(AppIcon), 0644)
}

func parseSizes(sizes string) ([]int, error) {
	// split the input string by comma and confirm that each one is an integer
	parsedSizes := strings.Split(sizes, ",")
	var result []int
	for _, size := range parsedSizes {
		s, err := strconv.Atoi(size)
		if err != nil {
			return nil, err
		}
		if s == 0 {
			continue
		}
		result = append(result, s)
	}

	// put all integers in a slice and return
	return result, nil
}

func generateMacIcon(iconData []byte, options *IconsOptions) error {

	srcImg, _, err := image.Decode(bytes.NewBuffer(iconData))
	if err != nil {
		return err
	}

	dest, err := os.Create(options.MacFilename)
	if err != nil {
		return err

	}
	defer func() {
		err = dest.Close()
		if err == nil {
			return
		}
	}()
	return icns.Encode(dest, srcImg)
}

func generateWindowsIcon(iconData []byte, sizes []int, options *IconsOptions) error {

	var output bytes.Buffer

	err := winicon.GenerateIcon(bytes.NewBuffer(iconData), &output, sizes)
	if err != nil {
		return err
	}

	err = os.WriteFile(options.WindowsFilename, output.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenerateTemplateIcon(data []byte, outputFilename string) error {
	// Decode the input file as a PNG
	buffer := bytes.NewBuffer(data)
	img, err := png.Decode(buffer)
	if err != nil {
		return fmt.Errorf("failed to decode input file as PNG: %w", err)
	}

	// Create a new image with the same dimensions and RGBA color model
	bounds := img.Bounds()
	iconImg := image.NewRGBA(bounds)

	// Iterate over each pixel of the input image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get the alpha of the pixel
			_, _, _, a := img.At(x, y).RGBA()
			iconImg.SetRGBA(x, y, color.RGBA{R: 0, G: 0, B: 0, A: uint8(a)})
		}
	}

	// Create the output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Encode the template icon image as a PNG and write it to the output file
	if err := png.Encode(outFile, iconImg); err != nil {
		return fmt.Errorf("failed to encode output image as PNG: %w", err)
	}

	return nil
}
