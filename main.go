package main

import (
	"awesomeProject/model"
	"awesomeProject/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"image/color"
	"log"
	"net/http"
	"strings"
)

var search = &model.Search{
	Novels: make([]model.Novel, 0),
}

type searchUi struct {
	entry *widget.Entry
	list  *widget.List
}
type NovelInfoUI struct {
	Novelname      *widget.Label
	Auther         *widget.Label
	LastUpdateTime *widget.Label
	Status         *widget.Label
	Introduce      *widget.Label
	List           *widget.List
	DataNovelInfo  model.NovelInfo
	DataChapter    []model.Chapter
	Image          *canvas.Image
}

func makeMainUi(window *fyne.Window, ui *searchUi) {
	entry := ui.entry
	entry = widget.NewEntry()
	entry.SetPlaceHolder("Enter your text here...")
	entry.OnChanged = func(s string) {
		search.Uri = s
	}
	list := widget.NewList(
		func() int {
			return len(search.Novels)
		},
		func() fyne.CanvasObject {
			newText := canvas.NewText("", color.Black)
			newText.TextSize = 20
			label1 := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{
				Bold:      false,
				Italic:    false,
				Monospace: false,
				Symbol:    false,
				TabWidth:  10,
				Underline: true,
			})
			text := canvas.NewText("Text Object", color.Black)
			text.TextStyle = fyne.TextStyle{Italic: true}
			return container.NewHBox(newText, layout.NewSpacer(), label1)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*canvas.Text).Text = search.Novels[i].Articlename
			o.(*fyne.Container).Objects[2].(*widget.Label).SetText(search.Novels[i].Author)
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		publisher.Publish("novelinfo", search.Novels[id].UrlInfo)
	}
	label := widget.NewLabel("书名: ")
	button := widget.NewButton("Submit", func() {
		search.Novels = nil
		if len(search.Uri) > 0 {
			_, e := search.Search()
			if e != nil {
				dialog.NewError(e, *window).Show()
				return
			}
			if len(search.Novels) > 0 {
				list.Refresh()
				list.ScrollToTop()
			}
		}

	})
	box := container.NewBorder(nil, nil, label, button, entry)
	border := container.NewBorder(box, nil, nil, nil, list)
	ui.entry = entry
	ui.list = list
	(*window).SetContent(border)

}
func novelinfo(UrlInfo string) (model.NovelInfo, []model.Chapter) {
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

var publisher = utils.NewPublisher()

func init() {

}
func main() {
	var searchui searchUi

	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")
	novelInfo := myApp.NewWindow("info")
	makeMainUi(&myWindow, &searchui)

	var ui NovelInfoUI
	makeNovelUi(&novelInfo, &ui)
	novelInfo.Resize(fyne.NewSize(400, 750))
	publisher.Subscribe("content", func(data interface{}) {
		chapter := data.(model.Chapter)
		window := myApp.NewWindow(chapter.Title)
		label := widget.NewLabel("loading")
		scroll := container.NewHScroll(label)

		window.SetContent(container.NewVScroll(scroll))
		go func() {
			c := colly.NewCollector(
				colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36"),
				colly.DetectCharset(),
			)
			extensions.Referer(c)
			c.OnHTML("div#chaptercontent", func(e *colly.HTMLElement) {
				split := strings.Split(e.Text, " ")
				var content string
				for _, s := range split {
					var result strings.Builder
					for i, r := range s {
						if i > 0 && i%40 == 0 {
							result.WriteRune('\n')
						}
						result.WriteRune(r)
					}
					content = content + result.String()
				}
				strings.ReplaceAll(content, " ", "\n  ")
				//fmt.Println(all[i])
				label.SetText(content)
			})
			c.Visit("https://www.bqgbe.cc" + chapter.Url)
		}()
		window.Show()
		window.Resize(fyne.NewSize(400, 750))
	})
	myWindow.Show()

	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(400, 750))
	myApp.Run()
}

func (n *NovelInfoUI) Show(novelInfo model.NovelInfo, chapter []model.Chapter) {

	resp, err := http.Get(novelInfo.CoverUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	*(n.Image) = *canvas.NewImageFromReader(resp.Body, "")
	//log.Println("111", n.Image.Size())
	n.Image.FillMode = canvas.ImageFillContain
	n.Image.SetMinSize(fyne.NewSize(150, 200))
	n.DataNovelInfo = novelInfo
	n.Novelname.SetText(novelInfo.BookName)
	n.Introduce.SetText(novelInfo.Introduce)
	n.Auther.SetText(novelInfo.Auther)
	n.LastUpdateTime.SetText(novelInfo.LastUpdateTime)

	n.Status.SetText(novelInfo.Status)
	n.DataChapter = chapter
	n.List.ScrollToTop()
	n.List.UnselectAll()
}

func makeNovelUi(win *fyne.Window, novelinfoui *NovelInfoUI) {

	(*win).Hide()
	(*win).SetCloseIntercept(func() {
		(*win).Hide()
	})

	novelinfoui.Novelname = widget.NewLabel(novelinfoui.DataNovelInfo.BookName)
	novelinfoui.Auther = widget.NewLabel(novelinfoui.DataNovelInfo.Auther)
	novelinfoui.LastUpdateTime = widget.NewLabel(novelinfoui.DataNovelInfo.LastUpdateTime)
	novelinfoui.Status = widget.NewLabel(novelinfoui.DataNovelInfo.Status)

	path, err1 := fyne.LoadResourceFromPath("./1.jpg")
	if err1 != nil {
		//log.Fatal(err)
		fmt.Println(err1)
	}

	novelinfoui.Image = canvas.NewImageFromResource(path)

	novelinfoui.Image.SetMinSize(fyne.NewSize(150, 200))
	//}
	novelinfoui.Introduce = widget.NewLabel(novelinfoui.DataNovelInfo.Introduce)
	novelinfoui.List = widget.NewList(
		func() int {
			chapter := (*novelinfoui).DataChapter
			return len(chapter)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("hello")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(novelinfoui.DataChapter[i].Title)
		},
	)
	novelinfoui.List.OnSelected = func(id widget.ListItemID) {
		publisher.Publish("content", novelinfoui.DataChapter[id])
	}
	box := container.NewVBox(novelinfoui.Novelname, novelinfoui.Auther, novelinfoui.LastUpdateTime, novelinfoui.Status)
	//center := container.NewCenter()
	border := container.NewBorder(layout.NewSpacer(), novelinfoui.Introduce, novelinfoui.Image, nil, box)

	border1 := container.NewBorder(border, nil, nil, nil, novelinfoui.List)
	publisher.Subscribe("novelinfo", func(data any) {
		(*win).Show()
		noinfo, chapter := novelinfo(data.(string))
		if len(noinfo.Introduce) > 20 {
			// 每7个字符插入一个换行符
			var result strings.Builder
			for i, r := range noinfo.Introduce {
				if i > 0 && i%20 == 0 {
					result.WriteRune('\n')
				}
				result.WriteRune(r)
			}
			noinfo.Introduce = result.String()
		}
		novelinfoui.Show(noinfo, chapter)
	})

	(*win).SetContent(border1)
}
