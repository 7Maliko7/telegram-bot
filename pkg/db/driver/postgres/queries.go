package postgres

const (
	AddUserQuery           = `/*NO LOAD BALANCE*/ insert into users.list (dialog_id) values (:dialog_id);`
	GetUserByUserIdQuery   = `select * from users.list where user_id=$1;`
	GetUserByDialogIdQuery = `select * from users.list where dialog_id=$1;`
	GetQuestionTextQuery = `select question_text from questions.list where question_slug=$1;`
)
