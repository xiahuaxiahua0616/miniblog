package store

import (
	"context"
	"errors"

	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/xiahuaxiahua0616/miniblog/internal/apiserver/model"
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/errno"
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/log"
	"gorm.io/gorm"
)

// PostStore 定义了 post 模块在 store 层所实现的方法.
type PostStore interface {
	Create(ctx context.Context, obj *model.PostM) error
	Update(ctx context.Context, obj *model.PostM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.PostM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.PostM, error)

	PostExpansion
}

// PostExpansion 定义了帖子操作的附加方法.
type PostExpansion any

// postStore 是 PostStore 接口的实现.
type postStore struct {
	store *datastore
}

// 确保 postStore 实现了 PostStore 接口.
var _ PostStore = (*postStore)(nil)

// newPostStore 创建 postStore 的实例.
func newPostStore(store *datastore) *postStore {
	return &postStore{store}
}

// Create 插入一条帖子记录.
func (s *postStore) Create(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		log.Errorw("Failed to insert post into database", "err", err, "post", obj)
		return errno.ErrDBWrite.WithMessage("%s", err.Error())
	}

	return nil
}

// Update 更新帖子数据库记录.
func (s *postStore) Update(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update post in database", "err", err, "post", obj)
		return errno.ErrDBWrite.WithMessage("%s", err.Error())
	}

	return nil
}

// Delete 根据条件删除帖子记录.
func (s *postStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.PostM)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete post from database", "err", err, "conditions", opts)
		return errno.ErrDBWrite.WithMessage("%s", err.Error())
	}

	return nil
}

// Get 根据条件查询帖子记录.
func (s *postStore) Get(ctx context.Context, opts *where.Options) (*model.PostM, error) {
	var obj model.PostM
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve post from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrPostNotFound
		}
		return nil, errno.ErrDBRead.WithMessage("%s", err.Error())
	}

	return &obj, nil
}

// List 返回帖子列表和总数.
// nolint: nonamedreturns
func (s *postStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.PostM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list posts from database", "err", err, "conditions", opts)
		err = errno.ErrDBRead.WithMessage("%s", err.Error())
	}
	return
}
