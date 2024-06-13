package model

import (
	"GOEND/common"
	"errors"
	"strings"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title can not be empty")
)

type TodoItem struct {
	common.SQLModel
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
	Status      string `json:"status" gorm:"column:status;"`
}

func (TodoItem) TableName() string { return "todo_items" }

type TodoItemCreation struct {
	Id          int    `json:"id" gorm:"column:id;"`
	Title       string `json:"title" gorm:"column:title;"`
	Description string `json:"description" gorm:"column:description;"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	if len(i.Title) == 0 {
		return ErrTitleCannotBeEmpty
	}
	return nil
}

// Su dung con tro: khi update chuoi rong thi db moi update theo, neu khong de thi GORM hieu la ko thay doi gi
// Giai thich: neu co tro tro toi gia tri nil, fale, "", thi nghia la no van co gia tri, nen GORM tien hanh update
type TodoItemUpdate struct {
	Title       *string `json:"title" gorm:"column:title;"`
	Description *string `json:"description" gorm:"column:description;"`
	Status      *string `json:"status" gorm:"column:status;"`
}

func (TodoItemUpdate) TableName() string { return TodoItem{}.TableName() }

func (TodoItemCreation) TableName() string { return TodoItem{}.TableName() }
