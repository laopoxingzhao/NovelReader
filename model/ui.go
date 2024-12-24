package model

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"log"
	"net/http"
)

type SearchUi struct {
	Entry *widget.Entry
	List  *widget.List
}
type NovelInfoUI struct {
	Novelname      *widget.Label
	Auther         *widget.Label
	LastUpdateTime *widget.Label
	Status         *widget.Label
	Introduce      *widget.Label
	List           *widget.List
	DataNovelInfo  NovelInfo
	DataChapter    []Chapter
	Image          *canvas.Image
}

func (n *NovelInfoUI) Show(novelInfo NovelInfo, chapter []Chapter) {

	resp, err := http.Get(novelInfo.CoverUrl)
	if err != nil {
		log.Fatal(err)
		*(n.Image) = *canvas.NewImageFromFile("./1.jpg")

	} else {
		*(n.Image) = *canvas.NewImageFromReader(resp.Body, "")

	}
	defer resp.Body.Close()
	n.Image.FillMode = canvas.ImageFillContain
	n.Image.SetMinSize(fyne.NewSize(150, 200))
	//log.Println("111", n.Image.Size())

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
