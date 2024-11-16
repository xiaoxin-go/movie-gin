package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func main() {
	url1 := "https://www.javbus.com/pics/cover/ausm_b.jpg"
	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		fmt.Println("request error:", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://www.javbus.com")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	//req.Header.Set("Host", "https://www.javbus.com")
	//uri, err := url.Parse("127.0.0.1:7897")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		//Proxy:           http.ProxyURL(uri),
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 10}
	res, e := client.Do(req)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(res.StatusCode)
}
