package model

type RuleNovelSearch struct {
	Url         string
	Name        string
	Auther      string
	CoverUrl    string
	BookInfoUrl string
	Introduce   string
}

type RuleNovelInfo struct {
	BookName       string
	Auther         string
	CoverUrl       string
	LastUpdateTime string
	Status         string
	LastChapter    string
	ChapterUrl     string
	Introduce      string
}

type RuleChapter struct {
	Title string
	Url   string
}

type RuleContent struct {
	NextPageUrl string
	Content     string
	Title       string
}
