package domain

import "time"

type Comment struct {
	Id        int64
	Content   string
	ArticleId int64
	User      User
	Ctime     time.Time
	Utime     time.Time
}

type CommentList []Comment

func (c CommentList) Ids() []int64 {
	ids := make([]int64, len(c))
	for i, comment := range c {
		ids[i] = comment.Id
	}
	return ids
}
