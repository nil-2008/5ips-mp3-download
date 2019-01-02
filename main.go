package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	iconv "github.com/djimenez/iconv-go"
	"github.com/levigross/grequests"
)

const (
	baseUrl string = "http://www.5ips.net/down_610_221.htm"
)

func main() {
	//url := "http://www.5ips.net/down_610_001.htm"
	//url = getMp3DownloadUrl(url)
	fmt.Print(initDownloadHtmlUrl())

}

func downloadMp3(url string) string {
	ss := strings.Split(url, "/")
	dirname := ss[2] + ss[3] + "/" + ss[4]
	reg := regexp.MustCompile(`\d{3}.mp3`)
	filename := reg.FindString(ss[5])

	fmt.Println("dirname->", dirname)
	fmt.Println("filename->", ss[5])
	fmt.Println("savename->", filename)

	fmt.Print("begin to downloading ...")
	res, _ := grequests.Get(url, &grequests.RequestOptions{
		//结构体可以对指定的类型给值，而不一定都赋值
		Headers: map[string]string{
			"Host":       "www.5ips.net",
			"Referer":    url,
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:52.0) Gecko/20100101 Firefox/52.0"}})

	if res.StatusCode != 200 {
		fmt.Printf("下载失败，:%s\n", url)
		os.Exit(-1)
	}

	//mp3大小比较
	length := res.Header.Get("Content-Length")
	slen, _ := strconv.Atoi(length)
	if slen < 1024*1024*1 {
		fmt.Printf("download failed:%s\n", url)
		os.Exit(-2)
	}

	if _, err := os.Stat(dirname); err != nil {
		fmt.Printf("创建下载文件夹:%s\n", dirname)
		os.MkdirAll(dirname, 0777)
	}
	res.DownloadToFile(dirname + "/" + filename)
	return url
}

func getMp3DownloadUrl(url string) string {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer res.Body.Close()

	utfBody, err := iconv.NewReader(res.Body, "gb2312", "utf-8")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	doc, err := goquery.NewDocumentFromReader(utfBody)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-2)
	}
	dom := doc.Find("script").Eq(7).Text()
	//解析抽取URL
	reg := regexp.MustCompile(`"[http|/pingshu/].*"`)
	url_list := reg.FindAllString(dom, -1)

	return strings.Replace(url_list[0], `"`, "", -1) + strings.Replace(url_list[2], `"`, "", -1)
}

/*
*初始化下载目录
 */
func initDownloadHtmlUrl() []string {
	initDownloadUrl := []string{}

	initalUrls := []string{"http://www.5ips.net/ps/478.htm"}

	for _, url := range initalUrls {
		res, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer res.Body.Close()

		utfBody, err := iconv.NewReader(res.Body, "gb2312", "utf-8")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		doc, err := goquery.NewDocumentFromReader(utfBody)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(-2)
		}

		doc.Find(".displist li").Each(func(index int, sel *goquery.Selection) {
			down_url, exists := sel.Find("a").Attr("href")
			if exists && strings.Contains(down_url, "down") {
				mp3URL := getMp3DownloadUrl(down_url)
				downloadMp3(mp3URL)
			}
		})
	}
	return initDownloadUrl
}
