package main

import (
	"database/sql"
	"encoding/json"
	lib "github.com/instanto/instanto-lib"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ArticleParams struct {
	lib.Article
	ResearchLine int64 `json:"research_line"`
}

func Articles2JSON(articles []*lib.Article) (articlesJSON []byte, err error) {
	data := make(map[string]interface{})
	data["articles"] = articles
	articlesJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Article2JSON(article *lib.Article) (articleJSON []byte, err error) {
	articleJSON, err = json.Marshal(article)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func ArticleGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	articles, err := dbProvider.ArticleGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articlesJSON, err := Articles2JSON(articles)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(articlesJSON)
}
func ArticleCreate(w http.ResponseWriter, r *http.Request) {
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
	articleParams := &ArticleParams{}
	err = json.Unmarshal(body, articleParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err := dbProvider.ArticleCreate(articleParams.Title, articleParams.Web, articleParams.Date, username, articleParams.Newspaper)
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
	article, err := dbProvider.ArticleGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articleJSON, err := Article2JSON(article)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(articleJSON)
}
func ArticleGetById(w http.ResponseWriter, r *http.Request) {
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
	article, err := dbProvider.ArticleGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articleJSON, err := Article2JSON(article)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(articleJSON)
	return
}
func ArticleUpdate(w http.ResponseWriter, r *http.Request) {
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
	articleParams := &ArticleParams{}
	err = json.Unmarshal(body, articleParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err := dbProvider.ArticleUpdate(idInt64, articleParams.Title, articleParams.Web, articleParams.Date, username, articleParams.Newspaper)
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
	article, err := dbProvider.ArticleGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articleJSON, err := Article2JSON(article)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(articleJSON)
	return
}
func ArticleDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err := dbProvider.ArticleDelete(idInt64)
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
func ArticleGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err := dbProvider.ArticleGetResearchLines(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLinesJSON, err := ResearchLines2JSON(researchLines)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(researchLinesJSON)
	return
}
func ArticleAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	articleParams := &ArticleParams{}
	err = json.Unmarshal(body, articleParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.ArticleAddResearchLine(idInt64, articleParams.ResearchLine, username)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
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
	return
}
func ArticleRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articleParams := &ArticleParams{}
	err = json.Unmarshal(body, articleParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.ArticleRemoveResearchLine(idInt64, articleParams.ResearchLine)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !removed {
		w.WriteHeader(404)
		return
	}
	return
}
