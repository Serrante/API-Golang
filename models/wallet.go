package models

import (
	"api/utils"
	"errors"
)

var (
	ErrWalletNotFound = errors.New("Carteira não existe")
)

type Wallet struct {
	PublicKey string  `json:"public_key"`
	User      User    `json:"user"`
	Balance   float32 `json:"balance"`
	UpdatedAt string  `json:"updated_at"`
}

func (w *Wallet) GeneratePublicKey() {
	w.PublicKey = utils.Md5(w.User.Nickname + w.User.Password)
}

func GetWallets() ([]Wallet, error) {
	con := Connect()
	defer con.Close()

	sql := `SELECT	u.id,
									u.nickname,
									u.email,
									u.status,
									u.created_at,
									u.updated_at,
									w.public_key,
									w.balance,
									w.updated_at
						FROM	WALLETS AS w
						JOIN  USERS AS u ON u.id = w.user
						ORDER BY w.user asc`
	rs, err := con.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rs.Close()

	var wallets []Wallet

	for rs.Next() {
		var wallet Wallet

		err := rs.Scan(
			&wallet.User.UID,
			&wallet.User.Nickname,
			&wallet.User.Email,
			&wallet.User.Status,
			&wallet.User.CreatedAt,
			&wallet.User.UpdatedAt,
			&wallet.PublicKey,
			&wallet.Balance,
			&wallet.UpdatedAt)

		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func GetWalletByPublicKey(publicKey string) (Wallet, error) {
	con := Connect()
	defer con.Close()

	sql := `SELECT	u.id,
									u.nickname,
									u.email,
									u.status,
									u.created_at,
									u.updated_at,
									w.public_key,
									w.balance,
									w.updated_at
						FROM	WALLETS AS w
						JOIN  USERS AS u ON u.id = w.user
						WHERE w.public_key::text = $1`
	rs, err := con.Query(sql, publicKey)

	if err != nil {
		return Wallet{}, err
	}

	defer rs.Close()

	var wallet Wallet

	for rs.Next() {
		err := rs.Scan(
			&wallet.User.UID,
			&wallet.User.Nickname,
			&wallet.User.Email,
			&wallet.User.Status,
			&wallet.User.CreatedAt,
			&wallet.User.UpdatedAt,
			&wallet.PublicKey,
			&wallet.Balance,
			&wallet.UpdatedAt)

		if err != nil {
			return Wallet{}, err
		}
	}

	if wallet.PublicKey == "" {
		return Wallet{}, ErrWalletNotFound
	}

	return wallet, nil
}

func UpdateWallet(w Wallet) (int64, error) {
	con := Connect()
	defer con.Close()

	sql := `UPDATE WALLETS SET
							balance = $1
					WHERE public_key = $2`
	stmt, err := con.Prepare(sql)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	rs, err := stmt.Exec(w.Balance, w.PublicKey)

	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}

func AddBalance(w Wallet) (int64, error) {
	con := Connect()
	defer con.Close()

	sql := `UPDATE WALLETS SET
							balance = (balance + $1)
					WHERE public_key = $2`
	stmt, err := con.Prepare(sql)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	rs, err := stmt.Exec(w.Balance, w.PublicKey)

	if err != nil {
		return 0, err
	}

	return rs.RowsAffected()
}
