package models

import (
	"api/utils"
	"errors"
)

var (
	ErrUserNotFound = errors.New("Usuário não encontrado")
)

type User struct {
	UID       string `json:"id"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
	Status    int8   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewUser(user User) (bool, error) {
	con := Connect()
	defer con.Close()

	tx, err := con.Begin()

	if err != nil {
		tx.Rollback()
		return false, err
	}

	sql := `INSERT INTO USERS (
						nickname,
						email,
						password
					) VALUES (
						$1,
						$2,
						$3
					) returning id`
	{
		stmt, err := tx.Prepare(sql)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		defer stmt.Close()

		hashedPassword, err := utils.Bcrypt(user.Password)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		err = stmt.QueryRow(user.Nickname, user.Email, hashedPassword).Scan(&user.UID)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	sql = "INSERT INTO WALLETS (public_key, \"user\") VALUES ($1, $2)"
	wallet := Wallet{User: user}
	wallet.GeneratePublicKey()
	{
		stmt, err := tx.Prepare(sql)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = stmt.Exec(wallet.PublicKey, wallet.User.UID)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}
	return true, tx.Commit()
}

func GetUsers() ([]User, error) {
	con := Connect()
	defer con.Close()

	sql := `SELECT	id,
									nickname,
									email,
									status,
									created_at,
									updated_at
				 	FROM 		USERS`
	rs, err := con.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rs.Close()

	var users []User

	for rs.Next() {
		var user User
		err := rs.Scan(
			&user.UID,
			&user.Nickname,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUser(id string) (User, error) {
	con := Connect()
	defer con.Close()

	sql := `SELECT	id,
									nickname,
									email,
									status,
									created_at,
									updated_at
					FROM 		USERS
					WHERE   id::text = $1`
	rs, err := con.Query(sql, id)

	if err != nil {
		return User{}, err
	}

	defer rs.Close()

	var user User

	for rs.Next() {
		err := rs.Scan(
			&user.UID,
			&user.Nickname,
			&user.Email,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt)

		if err != nil {
			return User{}, err
		}
	}

	if user.UID == "" || len(user.UID) == 0 {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

func UpdateUser(user User) (int64, error) {
	con := Connect()
	defer con.Close()

	sql := `UPDATE USERS SET
							nickname = $1,
							email = $2,
							status = $3
					WHERE id::text = $4`
	stmt, err := con.Prepare(sql)

	if err != nil {
		return 0, nil
	}

	defer stmt.Close()
	rs, err := stmt.Exec(user.Nickname, user.Email, user.Status, user.UID)

	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}

func DeleteUser(id string) (int64, error) {
	con := Connect()
	defer con.Close()

	sql := `DELETE FROM USERS
					WHERE id::text = $1`
	stmt, err := con.Prepare(sql)

	if err != nil {
		return 0, nil
	}

	defer stmt.Close()
	rs, err := stmt.Exec(id)

	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}
