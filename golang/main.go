package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~mendelmaleh/freetype/truetype"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	font_file         = flag.String("font", "Alibaba-PuHuiTi-Regular.ttf", "Enter the Font File's Name which You r Planing Using")
	font_default, err = filepath.Abs(filepath.Join("../Fonts", *font_file))
	letter_size       = flag.Int("letter_size", 10, "Enter Letter Size")
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

// func GenerateText(text string, font []byte, font_size int) {
// 	file := freetype.NewContent()

// 	file.SetFont(font)
// 	file.SetFontSize(font_size)
// }

func main() {

	flag.Parse()
	font, err := filepath.Abs(filepath.Join("../Fonts", *font_file))
	if err == nil {
		font = font_default
	}

	fmt.Printf("\n the using font file is %s\n", font)

	fileByte, err := ioutil.ReadFile(font)
	if err != nil {
		fmt.Println(err)
	}

	f, err := truetype.Parse(fileByte)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Println("\n Successfully Load in font format file\n")
	fmt.Println(f)

}
