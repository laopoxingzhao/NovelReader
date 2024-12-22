package model

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"strings"
)

type NovelBase struct {
	SourceUrl  string
	SourceName string
}

type Chapter struct {
	Title string
	Url   string
}
type Novel struct {
	UrlInfo     string `json:"url_list"`
	UrlImg      string `json:"url_img"`
	Articlename string `json:"articlename"`
	Author      string `json:"author"`
	Intro       string `json:"intro"`
}

type NovelInfo struct {
	BookName       string
	Auther         string
	CoverUrl       string
	LastUpdateTime string
	Status         string
	LastChapter    string
	Introduce      string
}

type Search struct {
	Uri    string
	Novels []Novel
}

func (s *Search) Search() ([]Novel, error) {
	var ee error
	novels := make([]Novel, 0)
	c := colly.NewCollector(
		colly.DetectCharset())
	extensions.Referer(c)
	extensions.RandomUserAgent(c)

	search := RuleNovelSearch{
		Url: "https://bqgbe.cc/user/hm.html?q=,https://bqgbe.cc/user/search.html?q=",
		Name: `@js function(data){
			JSON.parse(data).author
	}`,
		Auther:      "@js JSON.parse(data).author",
		CoverUrl:    "@js JSON.parse(data).url_img}",
		BookInfoUrl: "@js JSON.parse(data).url_list}",
	}
	c.OnResponse(func(r *colly.Response) {
		ss := string(r.Body)
		fmt.Println(ss)
		if len(ss) > 1 {

			vm := goja.New()
			vm.Set("data", ss)
			_, err := vm.RunString(`
				const jsonArray = JSON.parse(data);
bb = []
jsonArray.forEach(a => {
    bb.push({
        UrlInfo: a.url_list,
        UrlImg: a.url_img,
        Articlename: a.articlename,
        Author: a.author,
        Intro: a.intro
    })
});`)
			if err != nil {
				fmt.Println(err)
			}
			novellist := new([]interface{})

			vm.ExportTo(vm.Get("bb"), novellist)

			for _, i := range *novellist {
				novel := Novel{}
				novel.Articlename = i.(map[string]interface{})["Articlename"].(string)
				novel.Author = i.(map[string]interface{})["Author"].(string)
				novel.Intro = i.(map[string]interface{})["Intro"].(string)
				novel.UrlImg = i.(map[string]interface{})["UrlImg"].(string)
				novel.UrlInfo = i.(map[string]interface{})["UrlInfo"].(string)
				novels = append(novels, novel)
			}
			//fmt.Println(novels)
			/*runtime.Set("data", string(r.Body))
			if strings.Contains(search.Name, "@js") {
				replace := strings.Replace(search.Name, "@js", "", -1)
				runString, err := runtime.RunString(replace)
				fmt.Println("name::", runString, err)
			}*/

		}
		s.Novels = novels

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("error:", e)
		ee = e

	})
	split := strings.Split(search.Url, ",")
	for _, i := range split {
		fmt.Println(i + s.Uri)
		c.Visit(i + s.Uri)
	}
	return s.Novels, ee
}
