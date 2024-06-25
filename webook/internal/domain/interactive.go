package domain

type InteractiveCount struct {
	BizId      int64
	ViewCnt    int64
	LikeCnt    int64
	CollectCnt int64

	Liked     bool
	Collected bool
}
