package resource

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/joaoteixeira/go-auth-service-v2/core"
	"github.com/joaoteixeira/go-auth-service-v2/helper"
	"gopkg.in/mgo.v2/bson"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user User

	json.NewDecoder(r.Body).Decode(&user)

	if user.Password == "" || user.Email == "" {
		helper.ResponseWithError(w, 401, "user or password is empty")
		return
	}

	password := user.Password

	filter := bson.M{"email": user.Email}

	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		helper.ResponseWithError(w, 404, "user or password incorrect")
		return
	}

	if errV := core.VerifyPassword(password, user.Password, user.Salt); errV != nil {
		helper.ResponseWithError(w, 401, "user or password incorrect")
		return
	}

	tokenDetail, errToken := core.CreateToken(user.ID.Hex())

	if errToken != nil {
		helper.ResponseWithError(w, 401, errToken.Error())

	}

	userJson, _ := json.Marshal(user)

	errAuth := core.CreateAuth(userJson, tokenDetail)

	if errAuth != nil {
		helper.ResponseWithError(w, 400, errToken.Error())

	}

	helper.Response(w, 200, tokenDetail)

}
