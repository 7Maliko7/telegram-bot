package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/7Maliko7/telegram-bot/pkg/db"
)

type Repository struct {
	db *sqlx.DB
}

func New(connectionString string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln(err)
	}

	return &Repository{
		db: db,
	}, nil
}

func (repo *Repository) Close() error {
	return repo.db.Close()
}

func (repo *Repository) AddUser(ctx context.Context, dialog_id int64) error {
	user := db.User{Dialog_id: dialog_id}
	dbResult, err := repo.db.NamedExec(AddUserQuery, user)
	if err != nil {
		return err
	}
	count, err := dbResult.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("User was not added")
	}
	return nil
}

func (repo *Repository) GetUserByUserId(ctx context.Context, user_id int64) (db.User, error) {
	var user db.User
	err := repo.db.Get(&user, GetUserByUserIdQuery, user_id)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (repo *Repository) GetUserByDialogId(ctx context.Context, dialog_id int64) (*db.User, error) {
	rows, err := repo.db.QueryxContext(ctx, GetUserByDialogIdQuery, dialog_id)
	if err != nil {
		return nil, err
	}

	var user *db.User
	for rows.Next() {
		user = &db.User{}
		err := rows.Scan(&user.User_id, &user.Dialog_id)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (repo *Repository) GetQuestionText(ctx context.Context, slug string) (string, error) {
	var questionText string
	err := repo.db.Get(&questionText, GetQuestionTextQuery, slug)
	if err != nil {
		return questionText, err
	}
	return questionText, nil

}
