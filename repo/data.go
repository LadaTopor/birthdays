package repo

import "database/sql"

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetListBirthdays() ([]Birthday, error) {
	rows, err := r.db.Query(`SELECT * FROM birthdays`)
	if err != nil {
		return nil, err
	}

	var result []Birthday
	for rows.Next() {
		var res Birthday
		err = rows.Scan(&res.Id, &res.FullName, &res.Birthdate)
		if err != nil {
			return nil, err
		}
		result = append(result, res)
	}

	return result, nil
}

func (r *Repo) AddNewBirthday(fullName, birthday string) error {
	_, err := r.db.Exec(`INSERT INTO birthdays (full_name, birthdate) VALUES ($1, $2)`, fullName, birthday)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteBirthday(id int) error {
	_, err := r.db.Exec(`DELETE FROM birthdays WHERE id=$1`, id)
	if err != nil {
		return err
	}

	return nil
}
