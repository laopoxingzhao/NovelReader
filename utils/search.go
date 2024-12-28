package utils

import (
	"NovelReader/model"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func init() {
	GetPublisher().Subscribe("download", func(data any) {
		dir := filepath.Join("E:/", "Novel")

		if err := os.MkdirAll(dir, 0755); err != nil {
			// 处理错误
			log.Println(err)
		}
		/*go func() {

			for _, v := range data.([]model.Chapter) {
				fmt.Println(v.Title)
				content := SearchContent(v.Url)
				filePath := filepath.Join(dir, v.Title+".txt")
				log.Println(filePath)
				file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
				if err != nil {
					log.Println(err)
					return
				}
				defer file.Close()
				fmt.Fprintln(file, content)

			}
		}()*/
		s := SearchContentS(data.([]model.Chapter))
		for ii, chapter := range s {
			filePath := filepath.Join(dir, CleanFileName(strconv.Itoa(ii)+"  "+chapter.Title+".txt"))
			file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				log.Println(err)
				break
			}
			defer file.Close()
			fmt.Fprintln(file, chapter.Url)
		}
	})
}

// CleanFileName 清理文件名中的非法字符
func CleanFileName(fileName string) string {
	// 定义非法字符的正则表达式
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	// 替换非法字符为空字符串
	cleaned := invalidChars.ReplaceAllString(fileName, "")
	return cleaned
}

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

		all = strings.ReplaceAll(e.Text, " ", "\n\t")
		fmt.Println(all)
	})
	c.OnError(func(r *colly.Response, e error) {
		log.Println("loading   ", e)
		//return e
		all = e.Error()
	})

	s := "https://www.bqgbe.cc" + Url
	log.Print(s)
	c.Visit(s)
	return all
}

// SearchContentS 获取章节内容
func SearchContentS(chapters []model.Chapter) []model.Chapter {
	var num int = 10
	j := len(chapters) / num         // 大致每份大小
	remainder := len(chapters) % num // 剩余元素数量

	var chapters1 [10][]model.Chapter
	for i := 0; i < num; i++ {
		start := i * j
		end := start + j
		if i == num-1 {
			end += remainder // 最后一块包含剩余元素
		}
		chapters1[i] = chapters[start:end]
	}
	collector := colly.NewCollector(
		colly.DetectCharset(),
	)
	extensions.Referer(collector)
	extensions.RandomUserAgent(collector)
	var w sync.WaitGroup

	for i, i2 := range chapters1 {
		a := 0
		clone := collector.Clone()
		clone.OnHTML("div#chaptercontent", func(e *colly.HTMLElement) {
			all := strings.ReplaceAll(e.Text, " ", "\n\t")
			i3 := i*j + a
			log.Println(i3, all)
			chapters[i3].Url = all
			a++
		})
		clone.OnError(func(r *colly.Response, e error) {
			log.Println("loading   ", e)
			//return e
			chapters[i*j+a].Url = e.Error()
			a++
		})

		go func() {
			for _, chapter := range i2 {
				w.Add(1)
				clone.Visit("https://www.bqgbe.cc" + chapter.Url)
				w.Done()
			}
		}()
	}
	w.Wait()
	log.Println("success")
	return chapters
}
