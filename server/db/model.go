package db

import (
	"log"
	"time"
)

type User struct {
	ID         int
	Username   string
	Password   string
	CreateDate time.Time
}

func Query() ([]User, error) {
	var result []User
	sql := "select * from t_user"
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Fatalln("ERROR:", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatalln("ERROR:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreateDate)
		result = append(result, u)
		if err != nil {
			log.Fatalln("ERROR:", err)
			continue
		}
	}

	return result, nil
}
