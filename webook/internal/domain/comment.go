package domain

import "time"

type Comment struct {
	Id        int64  `json:"id"`
	Content   string `json:"content"`
	ArticleId int64  `json:"article_id"`
	User      User   `json:"user"`
	Ctime     time.Time `json:"ctime"`
	Utime     time.Time `json:"utime"`
}

type CommentList []Comment

func (c CommentList) Ids() []int64 {
	ids := make([]int64, len(c))
	for i, comment := range c {
		ids[i] = comment.Id
	}
	return ids
}
