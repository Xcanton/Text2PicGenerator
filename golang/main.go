package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~mendelmaleh/freetype"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var (
	font_file         = flag.String("font", "Alibaba-PuHuiTi-Regular.otf", "Enter the Font File's Name which You r Planing Using")
	font_default, err = filepath.Abs(filepath.Join("../Fonts", *font_file))
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

func main() {

	flag.Parse()
	font, err := filepath.Abs(filepath.Join("../Fonts", *font_file))
	if err == nil {
		font = font_default
	}

	fmt.Printf("\n the using font file is %s\n", font)

	uft16Bytes, _ := ReadFileUTF16(font)
	fmt.Println("读取font字体文件完成")

	fmt.Println("转换font字体文件为utf-8完成")

	f, err := freetype.ParseFont(uft16Bytes)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Print(f)
}
