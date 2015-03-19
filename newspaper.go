package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	lib "github.com/instanto/instanto-lib"

	"github.com/gorilla/mux"
)

type NewspaperParams struct {
	lib.Newspaper
}

func Newspapers2JSON(newspapers []*lib.Newspaper) (newspapersJSON []byte, err error) {
	data := make(map[string]interface{})
	data["newspapers"] = newspapers
	newspapersJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Newspaper2JSON(newspaper *lib.Newspaper) (newspaperJSON []byte, err error) {
	newspaperJSON, err = json.Marshal(newspaper)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func NewspaperGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	newspapers, err :=dbProvider.NewspaperGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspapersJSON, err := Newspapers2JSON(newspapers)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(newspapersJSON)
}
func NewspaperCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	username, ok := AuthenticateUser(r)
	if !ok {
		w.WriteHeader(401)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspaperParams := &NewspaperParams{}
	err = json.Unmarshal(body, newspaperParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.NewspaperCreate(newspaperParams.Name, newspaperParams.Web, username)
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
	newspaper, err :=dbProvider.NewspaperGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspaperJSON, err := Newspaper2JSON(newspaper)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(newspaperJSON)
}
func NewspaperGetById(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		w.WriteHeader(404)
		return
	}
	newspaper, err :=dbProvider.NewspaperGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspaperJSON, err := Newspaper2JSON(newspaper)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(newspaperJSON)
	return
}
func NewspaperUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	username, ok := AuthenticateUser(r)
	if !ok {
		w.WriteHeader(401)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspaperParams := &NewspaperParams{}
	err = json.Unmarshal(body, newspaperParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.NewspaperUpdate(idInt64, newspaperParams.Name, newspaperParams.Web, username)
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
	if numRows == 0 {
		w.WriteHeader(404)
		return
	}
	newspaper, err :=dbProvider.NewspaperGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	newspaperJSON, err := Newspaper2JSON(newspaper)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(newspaperJSON)
	return
}
func NewspaperDelete(w http.ResponseWriter, r *http.Request) {
	_, ok := AuthenticateUser(r)
	if !ok {
		w.WriteHeader(401)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	numRows, err :=dbProvider.NewspaperDelete(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if numRows == 0 {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
	return
}
func NewspaperGetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	articles, err :=dbProvider.ArticleGetByNewspaper(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articlesJSON, err := Articles2JSON(articles)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(articlesJSON)
	return
}
func NewspaperUpdateLogo(w http.ResponseWriter, r *http.Request) {
	username, ok := AuthenticateUser(r)
	if !ok {
		w.WriteHeader(401)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	exists, err :=dbProvider.NewspaperExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	fileSize := r.ContentLength
	if fileSize == 0 || fileSize > 2*1024*1024 {
		verr := &lib.ValidationError{"logo", "filesize must be greater than 0 and less than 2Mb"}
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
	var buffer = make([]byte, fileSize)
	_, err = io.ReadFull(r.Body, buffer)
	if err != nil {
		LogError(err)
		verr := &lib.ValidationError{"logo", "filesize must be greater than 0 and less than 2Mb"}
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
	filename := "newspaper_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err :=dbProvider.NewspaperUpdateLogo(idInt64, filename, username)
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
	fd, err := os.Create(config.MediaDir + filename)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	_, err = fd.Write(buffer)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	return
}

func NewspaperDeleteLogo(w http.ResponseWriter, r *http.Request) {
	username, ok := AuthenticateUser(r)
	if !ok {
		w.WriteHeader(401)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	exists, err :=dbProvider.NewspaperExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err :=dbProvider.NewspaperUpdateLogo(idInt64, "", username)
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
}
