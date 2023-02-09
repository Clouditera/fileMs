package picture

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
)

func ReadPic(imgFile string) image.Image {
	// 打开图片文件
	f, err := os.Open(imgFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	// 解析图片
	img, fmtName, err := image.Decode(f)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Printf("Name: %v, Bounds: %+v, Color: %+v", fmtName, img.Bounds(), img.ColorModel())
	x := img.Bounds().Size().X
	y := img.Bounds().Size().Y

	fmt.Println(fmtName, x, y)

	out := image.NewRGBA(img.Bounds())

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			//fmt.Println(uint8(r), uint8(g), uint8(b))
			gray := Gray(uint8(r), uint8(g), uint8(b))
			printPic(gray)
			out.Set(i, j, color.RGBA{R: gray, G: gray, B: gray, A: 255})
		}
		fmt.Println("")
	}

	SavePic(out)

	return img
}

/*
反转，把每个像素点的每个rgb值都与255作差（alpha的值不改变）
r, g, b = 255-r , 255-g , 255-b
*/

func Invert(r, g, b uint8) (uint8, uint8, uint8) {
	return 255 - r, 255 - g, 255 - b
}

func Gray(r, g, b uint8) uint8 {
	//fmt.Println("\ngray", (r+g+b)/3)
	//return (r + g + b) / 3
	return g
}

func SavePic(img *image.RGBA) {
	f, err := os.Create("/static/out/uims1.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, img)
	if err != nil {
		log.Fatalln(err)
	}
	b.Flush()
}

// 输出 V2
func printPic(gray uint8) {
	charArray := []string{"K", "S", "P", "k", "s", "p", ";", "."}
	index := math.Round(float64(int(gray) * (len(charArray) + 1) / 255))
	if int(index) >= len(charArray) {
		fmt.Print(" ")
	} else {
		fmt.Printf("%s", charArray[int(index)])
	}
}
