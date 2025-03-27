package domain

import "time"

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus

	Ctime time.Time
	Utime time.Time
}

type ArticleList []Article

func (a ArticleList) Ids() []int64 {
	ids := make([]int64, len(a))
	for i, article := range a {
		ids[i] = article.Id
	}
	return ids
}

func (a Article) Abstract() string {
	str := []rune(a.Content)
	if len(str) > 100 {
		return string(str[:100])
	}
	return string(str)
}

type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown     = iota
	ArticleStatusUnpublished = iota
	ArticleStatusPublished   = iota
	ArticleStatusPrivate     = iota
)
