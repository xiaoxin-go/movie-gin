package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func main() {
	file, err := os.Open("F:\\mv\\ABF-011.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// 解码图片
	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	bounds := img.Bounds()
	fmt.Println("图片格式：", format)
	// 定义要截取的区域 (x, y, width, height)
	rect := image.Rect(441, 0, bounds.Dx(), bounds.Dy()) // 裁剪 50,50 到 200,200 区域
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(rect)
	// 创建输出文件
	outFile, err := os.Create("output.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	// 将裁剪后的图片保存为 JPEG 格式
	err = jpeg.Encode(outFile, croppedImg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
