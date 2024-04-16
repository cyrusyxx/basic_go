package domain

type Article struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Status  ArticleStatus
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
