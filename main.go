package main

import (
	"NovelReader/model"
	"NovelReader/utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"image/color"
	"log"
	"strings"
)

var search = &model.Search{
	Novels: make([]model.Novel, 0),
}

var publisher = utils.NewPublisher()

func init() {

}
func main() {
	var searchui model.SearchUi

	myApp := app.New()
	myWindow := myApp.NewWindow("Entry Widget")
	novelInfo := myApp.NewWindow("info")
	MakeMainUi(&myWindow, &searchui)

	var ui model.NovelInfoUI
	MakeNovelUi(&novelInfo, &ui)
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
			log.Print("loading")
			/*	func() {

				}()*/
			c.OnHTML("div#chaptercontent", func(e *colly.HTMLElement) {
				log.Print("loading_____")

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
				label.SetText(content)
			})

			c.OnError(func(r *colly.Response, e error) {
				if window != nil {
					dialog.NewError(e, window).Show()
				}
				log.Println("loading   ", e)
			})
			s := "https://www.bqgbe.cc" + chapter.Url
			log.Print(s)
			c.Visit(s)
		}()
		window.Show()
		window.Resize(fyne.NewSize(400, 750))
	})
	myWindow.Show()

	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(400, 750))
	myApp.Run()
}

func MakeMainUi(window *fyne.Window, ui *model.SearchUi) {
	entry := ui.Entry
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
	ui.Entry = entry
	ui.List = list
	(*window).SetContent(border)

}

func MakeNovelUi(win *fyne.Window, novelinfoui *model.NovelInfoUI) {

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
		noinfo, chapter := utils.Novelinfo(data.(string))
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

func SearchAllContent() {

}
