package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

func main(){
	//clip("3DSVR-0832.jpg")
	//return
	files, err := os.ReadDir("E:\\FFOutput\\static\\images")
	if err != nil{
		log.Fatalln(err)
	}
	for _, file := range files{
		if len(strings.Split(file.Name(), "-")) > 2{
			continue
		}
		fmt.Println(file.Name())
		clip(file.Name())
	}
}

func clip(name string){
	src := "E:\\FFOutput\\static\\images\\" + name
	dst := "E:\\FFOutput\\static\\logo\\" + name
	_, err := os.Stat(dst)
	if err == nil{
		fmt.Printf("%s is already exists", name)
		return
	}
	fIn, err := os.Open(src)
	if err != nil{
		log.Fatalln(err)
	}

	defer fIn.Close()
	fOut, _ := os.Create(dst)
	defer fOut.Close()
	origin, fm, err := image.Decode(fIn)
	fmt.Println(fm, err)
	img := origin.(*image.YCbCr)
	subImg := img.SubImage(image.Rect(422, 0, 800, 537)).(*image.YCbCr)
	jpeg.Encode(fOut, subImg, &jpeg.Options{Quality: 100})
}
