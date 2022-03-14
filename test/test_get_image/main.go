package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
)

func main(){
	uri, err := url2.Parse("http://127.0.0.1:7890")
	url := "https://img2.baidu.com/it/u=1305248331,3698728375&fm=253&fmt=auto&app=138&f=JPEG?w=889&h=500"
	url = "https://www.javbus.com/pics/sample/6erf_2.jpg"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy: http.ProxyURL(uri),
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil{
		log.Fatal(err)
	}

	out, err := os.Create("./static/image/test.jpg")
	if err != nil{
		log.Fatal(err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil{
		log.Println(err)
	}
}
