package db

import (
	"context"
)

type Databaser interface {
	AddUser(ctx context.Context, dialog_id int64) error
	GetUserByUserId(ctx context.Context, user_id int64) (User, error)
	GetUserByDialogId(ctx context.Context, dialog_id int64) (*User, error)
	GetQuestionText(ctx context.Context, slug string) (string, error)
}

type User struct {
	User_id   int64 `db:"user_id"`
	Dialog_id int64 `db:"dialog_id"`
}

type Question struct {
	Id   int    `db:"id"`
	Slug string `db:"slug"`
	Text string `db:"text"`
}
