package controllers

import (
	"api/auth"
	"api/models"
	"api/utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var user models.User
	err := json.Unmarshal(body, &user)

	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	userAuthenticate, err := auth.SignIn(user)

	if err != nil {
		utils.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	user.Password = ""
	utils.ToJson(w, userAuthenticate)
}
