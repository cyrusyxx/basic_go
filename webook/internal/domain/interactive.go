package domain

type InteractiveCount struct {
	ViewCnt    int64
	LikeCnt    int64
	CollectCnt int64

	Liked     bool
	Collected bool
}
