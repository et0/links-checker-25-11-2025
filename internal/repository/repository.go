package repository

import "github.com/et0/links-checker-25-11-2025/internal/model"

type LinkRepository interface {
	CreateList(list *model.LinkList) error
	UpdateList(list *model.LinkList) error
	GetList(ids []int) ([]*model.LinkList, error)
	GetAllList() []*model.LinkList
}
