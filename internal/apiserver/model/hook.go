package model

import (
	"github.com/xiahuaxiahua0616/miniblog/internal/pkg/rid"
	"github.com/xiahuaxiahua0616/miniblog/pkg/auth"
	"gorm.io/gorm"
)

// AfterCreate 在创建数据库记录之后生成 postID.
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.New(uint64(m.ID))

	return tx.Save(m).Error
}

// AfterCreate 在创建数据库记录之后生成 userID.
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))

	return tx.Save(m).Error
}

func (m *UserM) BeforeCreate(tx *gorm.DB) error {
	var err error
	m.Password, err = auth.Encrypt(m.Password)
	if err != nil {
		return err
	}

	return nil
}
