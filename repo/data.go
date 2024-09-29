package repo

import "database/sql"

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetAllOfListBirthdays() ([]Birthday, error) {
	rows, err := r.db.Query(`SELECT * FROM birthdays`)
	if err != nil {
		return nil, err
	}

	var result []Birthday
	for rows.Next() {
		var res Birthday
		err = rows.Scan(&res.Id, &res.FullName, &res.Birthdate, &res.ChatId)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (r *Repo) GetListBirthdays(chatId int64) ([]Birthday, error) {
	rows, err := r.db.Query(`SELECT * FROM birthdays WHERE chat_id = $1`, chatId)
	if err != nil {
		return nil, err
	}

	var result []Birthday
	for rows.Next() {
		var res Birthday
		err = rows.Scan(&res.Id, &res.FullName, &res.Birthdate, &res.ChatId)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (r *Repo) AddNewBirthday(fullName, birthday string, chatId int64) error {
	var maxID int
	err := r.db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM birthdays WHERE chat_id = $1", chatId).Scan(&maxID)
	if err != nil {
		return err
	}
	newID := maxID + 1

	_, err = r.db.Exec(`INSERT INTO birthdays (id, full_name, birthdate, chat_id) VALUES ($1, $2, $3, $4)`, newID, fullName, birthday, chatId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteBirthday(id int, chatId int64) error {
	_, err := r.db.Exec(`DELETE FROM birthdays WHERE id=$1 AND chat_id=$2`, id, chatId)
	if err != nil {
		return err
	}

	return nil
}
