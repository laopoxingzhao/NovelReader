package utils

import (
	"NovelReader/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
	"strings"
)

// novelinfo 获取小说信息 和章节
func Novelinfo(UrlInfo string) (model.NovelInfo, []model.Chapter) {
	info := model.RuleNovelInfo{
		BookName:       "div.info h1",
		Auther:         "div.info div.small span:first-child",
		CoverUrl:       "div.info div.cover img",
		LastUpdateTime: "div.info div.small>span:nth-of-type(3)",
		Status:         "div.info div.small>span:nth-of-type(2)",
		LastChapter:    "div.info div.small>span[class=last] a",
		ChapterUrl:     "",
		Introduce:      "div.info div.intro dd",
	}
	chapterr := model.RuleChapter{
		Title: "div.listmain dd a",
		Url:   "div.listmain dd a",
	}
	chapter := []model.Chapter{}

	s := "https://www.bqgbe.cc" + UrlInfo
	resp, err := http.Get(s)
	if err != nil {
		log.Println(err)
		return model.NovelInfo{}, chapter
	}
	defer resp.Body.Close()

	reader, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	n := model.NovelInfo{}

	reader.Find(info.BookName).Each(func(i int, selection *goquery.Selection) {
		n.BookName = selection.Text()
	})

	reader.Find(info.Auther).Each(func(i int, selection *goquery.Selection) {
		n.Auther = selection.Text()
	})
	reader.Find(info.CoverUrl).Each(func(i int, selection *goquery.Selection) {
		val, _ := selection.Attr("src")
		n.CoverUrl = val
	})
	reader.Find(info.LastUpdateTime).Each(func(i int, selection *goquery.Selection) {
		n.LastUpdateTime = selection.Text()
	})
	reader.Find(info.Status).Each(func(i int, selection *goquery.Selection) {
		n.Status = selection.Text()
	})
	reader.Find(info.Introduce).Each(func(i int, selection *goquery.Selection) {
		n.Introduce = selection.Text()
	})
	chapter1 := model.Chapter{}
	reader.Find(chapterr.Title).Each(func(i int, selection *goquery.Selection) {
		val, e := selection.Attr("rel")
		if e == true && val == "nofollow" {
			return
		}
		chapter1.Title = selection.Text()
		chapter1.Url = selection.AttrOr("href", "")

		chapter = append(chapter, chapter1)
	})

	return n, chapter
}

// SearchContent 获取小说内容
func SearchContent(Url string) string {
	var all string
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"),
		colly.DetectCharset(),
	)
	extensions.Referer(c)
	log.Print("loading")
	c.OnHTML("div#chaptercontent", func(e *colly.HTMLElement) {
		log.Print("loading_____")

		all = strings.ReplaceAll(e.Text, "  ", "\n  ")
		fmt.Println(all)
	})
	c.OnError(func(r *colly.Response, e error) {
		log.Println("loading   ", e)
	})

	s := "https://www.bqgbe.cc" + Url
	log.Print(s)
	c.Visit(s)
	return all
}
