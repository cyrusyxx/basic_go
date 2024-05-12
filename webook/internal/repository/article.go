package repository

import (
	"context"
	"time"
	"webook/webook/internal/domain"
	"webook/webook/internal/repository/cache"
	"webook/webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, arti domain.Article) (int64, error)
	Update(ctx context.Context, arti domain.Article) error
	Sync(ctx context.Context, arti domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	GetByAuthor(ctx context.Context, uid int64, offset int64, limit int64) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache

	userRepo UserRepository

	readerDAO dao.ArticleReaderDAO
	authorDAO dao.ArticleAuthorDAO
}

func NewCachedArticleRepository(dao dao.ArticleDAO,
	cache cache.ArticleCache,
	userRepo UserRepository) ArticleRepository {
	return &CachedArticleRepository{
		dao:      dao,
		cache:    cache,
		userRepo: userRepo,
	}
}

func (r *CachedArticleRepository) Create(ctx context.Context,
	arti domain.Article) (int64, error) {
	id, err := r.dao.Insert(ctx, r.toEntity(arti))
	if err != nil {
		return 0, err
	}

	// Delete cache
	err = r.cache.DelFirstPage(ctx, arti.Author.Id)
	if err != nil {
		// log
	}
	return id, nil
}

func (r *CachedArticleRepository) Update(ctx context.Context,
	arti domain.Article) error {
	err := r.dao.UpdateById(ctx, r.toEntity(arti))
	if err != nil {
		return err
	}

	// Delete cache
	err = r.cache.DelFirstPage(ctx, arti.Author.Id)
	if err != nil {
		// log
	}
	return nil
}

func (r *CachedArticleRepository) Sync(ctx context.Context,
	arti domain.Article) (int64, error) {
	// Sync article
	id, err := r.dao.Sync(ctx, r.toEntity(arti))
	if err != nil {
		return 0, err
	}

	// Delete cache
	err = r.cache.DelFirstPage(ctx, arti.Author.Id)
	if err != nil {
		// log
	}

	// Set cache
	user, err := r.userRepo.FindByID(ctx, arti.Author.Id)
	if err != nil {
		// log
		return id, nil
	}
	arti.Author.Name = user.NickName
	err = r.cache.SetPub(ctx, arti)
	if err != nil {
		// log
	}

	return id, nil
}

func (r *CachedArticleRepository) SyncStatus(ctx context.Context,
	uid int64, id int64, status domain.ArticleStatus) error {
	err := r.dao.SyncStatus(ctx, uid, id, uint8(status))
	if err != nil {
		return err
	}

	// Delete cache
	err = r.cache.DelFirstPage(ctx, uid)
	if err != nil {
		// log
	}
	return nil
}

func (r *CachedArticleRepository) GetByAuthor(ctx context.Context,
	uid int64, offset int64, limit int64) ([]domain.Article, error) {
	// Get page from cache
	if offset == 0 && limit == 100 {
		res, err := r.cache.GetFirstPage(ctx, uid)
		if err == nil {
			return res, nil
		} else {
			// log
		}
	}

	// Get page from database
	artis, err := r.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	res := r.toDomains(artis)

	// Set cache
	err = r.cache.SetFirstPage(ctx, uid, res)
	if err != nil {
		// log
	}

	r.preCache(ctx, res)

	return res, nil
}

func (r *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	// Get from cache
	res, err := r.cache.GetById(ctx, id)
	if err == nil {
		return res, nil
	}

	// Get from database
	arti, err := r.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}

	// Set cache
	err = r.cache.Set(ctx, r.toDomain(arti))
	if err != nil {
		// log
	}

	return r.toDomain(arti), nil
}

func (r *CachedArticleRepository) GetPubById(ctx context.Context,
	id int64) (domain.Article, error) {
	// Get from cache
	res, err := r.cache.GetPubById(ctx, id)
	if err == nil {
		return res, nil
	}

	// Get article from database
	arti, err := r.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	res = r.toDomain(dao.Article(arti))

	// Get author's name from user repository
	author, err := r.userRepo.FindByID(ctx, arti.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res.Author.Name = author.NickName

	// Set cache
	err = r.cache.SetPub(ctx, res)
	if err != nil {
		// log
	}

	return res, nil
}

func (r *CachedArticleRepository) preCache(ctx context.Context,
	arti []domain.Article) {
	const maxlen = 1024 * 1024
	if len(arti) > 0 && len(arti[0].Content) < maxlen {
		err := r.cache.Set(ctx, arti[0])
		if err != nil {
			// log
		}
	}
}

// toEntity convert domain.Article to dao.Article
func (r *CachedArticleRepository) toEntity(arti domain.Article) dao.Article {
	return dao.Article{
		Id:       arti.Id,
		AuthorId: arti.Author.Id,
		Title:    arti.Title,
		Content:  arti.Content,
		Status:   uint8(arti.Status),
	}
}

// toDomain convert dao.Article to domain.Article
func (r *CachedArticleRepository) toDomain(arti dao.Article) domain.Article {
	return domain.Article{
		Id:      arti.Id,
		Title:   arti.Title,
		Content: arti.Content,
		Author:  domain.Author{Id: arti.AuthorId},
		Status:  domain.ArticleStatus(arti.Status),

		Ctime: time.UnixMilli(arti.Ctime),
		Utime: time.UnixMilli(arti.Utime),
	}
}

// toDomains convert []dao.Article to []domain.Article
func (r *CachedArticleRepository) toDomains(artis []dao.Article) []domain.Article {
	var res []domain.Article
	for _, arti := range artis {
		res = append(res, r.toDomain(arti))
	}
	return res
}
