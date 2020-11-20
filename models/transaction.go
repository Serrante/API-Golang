package models

type Transaction struct {
	UID       string  `json:"id"`
	Origin    Wallet  `json:"origin"`
	Target    Wallet  `json:"target"`
	Cash      float32 `json:"cash"`
	Message   string  `json:"message"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func NewTransaction(transaction Transaction) (bool, error) {
	con := Connect()
	defer con.Close()

	tx, err := con.Begin()

	if err != nil {
		return false, err
	}

	sql := `UPDATE WALLETS SET
							balance = (balance - $1)
					WHERE public_key = $2`
	{
		stmt, err := tx.Prepare(sql)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = stmt.Exec(transaction.Origin.Balance, transaction.Origin.PublicKey)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	sql = `UPDATE WALLETS SET
							balance = (balance + $1)
				 WHERE public_key = $2`
	{
		stmt, err := tx.Prepare(sql)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = stmt.Exec(transaction.Origin.Balance, transaction.Target.PublicKey)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	sql = `INSERT INTO TRANSACTIONS (
							origin,
							target,
							cash,
							message
				 ) VALUES (
							$1,
							$2,
							$3,
							$4
				 )`
	{
		stmt, err := tx.Prepare(sql)

		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = stmt.Exec(
			transaction.Origin.PublicKey,
			transaction.Target.PublicKey,
			transaction.Cash,
			transaction.Message)

		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	return true, tx.Commit()
}

func GetTransactions() ([]Transaction, error) {
	con := Connect()
	defer con.Close()

	sql := `SELECT  *
					FROM    TRANSACTIONS`
	rs, err := con.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rs.Close()
	var transactions []Transaction

	for rs.Next() {
		var transaction Transaction
		err := rs.Scan(
			&transaction.UID,
			&transaction.Origin.PublicKey,
			&transaction.Target.PublicKey,
			&transaction.Cash,
			&transaction.Message,
			&transaction.CreatedAt,
			&transaction.UpdatedAt)

		if err != nil {
			return nil, err
		}

		origin, err := GetWalletByPublicKey(transaction.Origin.PublicKey)

		if err != nil {
			return nil, err
		}

		target, err := GetWalletByPublicKey(transaction.Target.PublicKey)

		if err != nil {
			return nil, err
		}

		transaction.Origin = origin
		transaction.Target = target
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
