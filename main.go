package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//執行過的
var urlSlice = []string{}

//參數
var masterUrl = os.Args[1]
var projectName = os.Args[2]

var j = 0 //總共跑幾次

func goSpide(url string) bool {
	//實際跑
	res, err := http.Get(masterUrl + url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
		// log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	//紀錄跑的url
	if indexOf(url, urlSlice) != -1 {
		return false
	}
	urlSlice = append(urlSlice, url)
	//讀取HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	check(err)
	//篩選連結出來
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//過濾HTML 抓出連結
		title, ok := s.Attr("href")
		if !ok {
			return
		}
		//連結如果包含主網域，要去除
		title = strings.Replace(title, masterUrl, "", -1)
		//不跑
		if strings.Contains(title, "https://") || strings.Contains(title, "http://") {
			return
		}
		if strings.Contains(title, "javascript:") {
			return
		}
		if strings.Contains(title, "#") {
			return
		}
		if strings.Contains(title, "member") {
			return
		}
		if strings.Contains(title, "shopcart") {
			return
		}
		//重複的網址不跑
		if indexOf(title, urlSlice) != -1 {
			return
		}
		//log
		appendText(title)
		j++
		fmt.Printf("index %d: %s\n", j, title)
		//深度優先
		goSpide(title)
	})
	return true
}

func main() {
	//開始爬
	goSpide("/")
}

//error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//切片是否包含字串
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

//從slice pop一個值
func pop(alist *[]string) string {
	f := len(*alist)
	rv := (*alist)[f-1]
	*alist = (*alist)[:f-1]
	return rv
}

//寫入log
func appendText(logText string) {
	if err := os.Mkdir(projectName, 0755); err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
	var fileName = "./" + projectName + "/sitemap.txt"
	f, err := os.OpenFile(fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.Default()
	logger.SetOutput(f)
	logger.SetFlags(0)
	logger.Println(logText)
	// logger.Println("more text to append")
}
