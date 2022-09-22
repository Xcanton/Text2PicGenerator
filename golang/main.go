package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	sunicode "unicode"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	font_file         = flag.String("font", "Alibaba-PuHuiTi-Regular.ttf", "Enter the Font File's Name which You r Planing Using")
	font_default, err = filepath.Abs(filepath.Join("../Fonts", *font_file))
	letter_size       = flag.Int("letter_size", 10, "Enter Letter Size")

	// 通过 goroutine 将后续参数加载并生成图片
	// text2generate     = flag.String("text", "", "Enter the text that you need to generate pic")
)

func ReadFileUTF16(filename string) ([]byte, error) {

	// Read the file into a []byte:
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, err := ioutil.ReadAll(unicodeReader)
	return decoded, err
}

func ReadFileUTF8(filename string) ([]byte, error) {
	// Read the file into a []byte:
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	gbkReader := bytes.NewReader(raw)
	utf8Reader := transform.NewReader(gbkReader, simplifiedchinese.GBK.NewDecoder())
	utf8Bytes, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		log.Println(err)
	}
	return utf8Bytes, err
}

func getActuralLen(text string) float64 {
	var length float64
	for _, char := range text {
		if sunicode.Is(sunicode.Han, char) {
			length++
		} else {
			length = length + 0.5
		}
	}
	return length
}

func GenerateText(text string, font_file *truetype.Font, font_size int) {

	file := freetype.NewContext()

	file.SetDPI(72)
	file.SetFont(font_file)
	file.SetFontSize(float64(font_size))
	file.SetHinting(font.HintingFull)

	size_y := int(float64(file.PointToFixed(float64(font_size))) / 48)
	size_x := int(float64(getActuralLen(text) * float64(size_y)))

	img := image.NewNRGBA(image.Rect(0, 0, size_x, size_y))

	file.SetSrc(image.White)
	file.SetClip(img.Bounds())
	file.SetDst(img)

	_, err = file.DrawString(text, freetype.Pt(0, int(file.PointToFixed(float64(font_size))>>6)))
	if err != nil {
		log.Println(err)
		return
	}

	imgfile, _ := os.Create(fmt.Sprintf("./%s.png", text))
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	flag.Parse()
	font_file, err := filepath.Abs(filepath.Join("../Fonts", *font_file))
	if err == nil {
		font_file = font_default
	}

	fmt.Printf("\n the using font file is %s\n", font_file)

	fileByte, err := ioutil.ReadFile(font_file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := freetype.ParseFont(fileByte)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("\n Successfully Load in font format file\n")

	for _, text := range os.Args[1:] {
		if !strings.Contains(text, "--") {
			GenerateText(text, f, *letter_size)
			fmt.Println(getActuralLen(text))
		}
	}
}
