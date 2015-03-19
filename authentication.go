package main

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	LoginParams := &LoginParams{}
	err = json.Unmarshal(body, LoginParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.UserCheckLogin(LoginParams.Username, LoginParams.Password)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if verr != nil {
		verrJSON, err := HelperValidationError2JSON(verr)
		if err != nil {
			LogError(err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(verrJSON)
		return
	}
	// Create token
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["iss"] = LoginParams.Username
	token.Claims["exp"] = time.Now().Add(time.Minute * 480).Unix()
	token.Claims["permissions"] = []string{"status_list", "status_add", "status_update", "status_delete"}
	tokenString, err := token.SignedString([]byte(config.Secret))
	data := make(map[string]string)
	data["token"] = tokenString
	dataJSON, err := json.Marshal(data)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(dataJSON)
}

func AuthenticateApp(r *http.Request) {
	// TODO check API Key and API secret
	// and apply bussiness validations like
	// rate limiting and time trials
	// Maybe create our custom proxy like other product
}

func AuthenticateUser(r *http.Request) (username string, ok bool) {
	tokenHeader := strings.Split(r.Header.Get("Authorization"), " ")
	if len(tokenHeader) < 2 {
		return
	}
	token, err := jwt.Parse(string(tokenHeader[1]), func(token *jwt.Token) (key interface{}, err error) {
		return []byte(config.Secret), nil
	})
	if err != nil {
		return
	}
	username = token.Claims["iss"].(string)
	ok = true
	return
}
func AuthorizeUser(username string, permission string) (err error) {
	// check permissions of the user
	return
}
