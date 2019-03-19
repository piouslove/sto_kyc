package models

import (
	"database/sql"
	"errors"
	"log"

	"Sto_kyc/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB_mysql *sql.DB

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Selector string `json:"selector"`
	Passport string `json:"passport"`
	Status   int    `json:"status"`
}

func CreateUser(user *User) (int, error) {
	result, err := DB_mysql.Exec(`insert into users(name, email, address, 
	selector, passport, status) values (?, ?, ?, ?, ?, ?)`,
		user.Name,
		user.Email,
		user.Address,
		user.Selector,
		user.Passport,
		user.Status)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func UpdateUser(status, id int) error {
	result, err := DB_mysql.Exec("update users set status = ? where id = ?",
		status,
		id)

	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func ReadUser(id int, selector string) (int, *User, error) {
	user := new(User)
	row := DB_mysql.QueryRow(`select id, name, address, email, passport 
	from users where id > ? and status = 2 and selector = ?`,
		id,
		selector)

	err := row.Scan(&id, &user.Name, &user.Address, &user.Email, &user.Passport)

	if err != nil {
		return 0, nil, err
	}

	user.Status = 2
	user.Selector = selector

	return id, user, nil
}

func CheckUserCertified(address, selector string) (bool, error) {
	var tmp int
	row := DB_mysql.QueryRow(`select status from users where address = ?
	 and selector = ?`,
		address,
		selector)
	err := row.Scan(&tmp)

	if err != nil {
		return false, err
	}

	if tmp != 1 {
		return false, errors.New("此用户没有通过KYC认证！")
	}
	return true, nil
}

func CheckUserExists(user *User) (bool, error) {
	var tmp string
	row := DB_mysql.QueryRow(`select selector from users where address = ?`,
		user.Address)
	err := row.Scan(&tmp)
	if err != nil {
		return false, err
	}
	if tmp == user.Selector {
		return true, errors.New("此地址已被注册！")
	}

	row = DB_mysql.QueryRow(`select selector from users where name = ?`,
		user.Name)
	err = row.Scan(&tmp)
	if err != nil {
		return false, err
	}
	if tmp == user.Selector {
		return true, errors.New("此用户已被注册！")
	}

	row = DB_mysql.QueryRow(`select selector from users where email = ?`,
		user.Email)
	err = row.Scan(&tmp)
	if err != nil {
		return false, err
	}
	if tmp == user.Selector {
		return true, errors.New("此邮箱已被注册！")
	}

	return false, nil
}

func init() {
	var err error
	username := config.V.Mysql.Username
	password := config.V.Mysql.Password
	host := config.V.Mysql.Host
	port := config.V.Mysql.Port
	dbname := config.V.Mysql.Dbname
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	DB_mysql, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
