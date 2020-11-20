package controllers

import (
	"api/models"
	"api/utils"
	"api/validations"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrInvalidCash = errors.New("Valor transferido é inválido")
)

func GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := models.GetTransactions()

	if err != nil {
		utils.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	utils.ToJson(w, transactions)
}

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	transaction, err := verifyTransaction(r)

	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}

	_, err = models.NewTransaction(transaction)

	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnprocessableEntity)
		return
	}

	utils.ToJson(w, utils.DefaultResponse{Data: "Transação concluída com sucesso!", Status: http.StatusCreated})
}

func verifyTransaction(r *http.Request) (models.Transaction, error) {
	params := mux.Vars(r)
	targetKey := params["public_key"]

	target, err := models.GetWalletByPublicKey(targetKey)

	if err != nil {
		return models.Transaction{}, err
	}

	body, _ := ioutil.ReadAll(r.Body)
	var origin models.Wallet
	err = json.Unmarshal(body, &origin)

	if err != nil {
		return models.Transaction{}, err
	}

	originVerify, err := models.GetWalletByPublicKey(origin.PublicKey)

	if err != nil {
		return models.Transaction{}, err
	}

	if validations.IsEmpty(target.PublicKey) || validations.IsEmpty(originVerify.PublicKey) {
		return models.Transaction{}, models.ErrWalletNotFound
	}

	if origin.Balance > originVerify.Balance || origin.Balance < 0 {
		return models.Transaction{}, ErrInvalidCash
	}

	var transaction models.Transaction
	transaction.Origin = origin
	transaction.Target = target
	transaction.Cash = origin.Balance
	transaction.Message = fmt.Sprintf("%s transferiu $%.2f, para %s", originVerify.User.Nickname, origin.Balance, target.User.Nickname)
	return transaction, nil
}
