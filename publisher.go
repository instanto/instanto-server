package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	lib "github.com/instanto/instanto-lib"

	"github.com/gorilla/mux"
)

type PublisherParams struct {
	lib.Publisher
	Publication int64 `json:"publication"`
}

func Publishers2JSON(publishers []*lib.Publisher) (publishersJSON []byte, err error) {
	data := make(map[string]interface{})
	data["publisher_collection"] = publishers
	publishersJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Publisher2JSON(publisher *lib.Publisher) (publisherJSON []byte, err error) {
	publisherJSON, err = json.Marshal(publisher)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func PublisherGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	publishers, err :=dbProvider.PublisherGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publishersJSON, err := Publishers2JSON(publishers)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publishersJSON)
}
func PublisherCreate(w http.ResponseWriter, r *http.Request) {
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
	publisherParams := &PublisherParams{}
	err = json.Unmarshal(body, publisherParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.PublisherCreate(publisherParams.Name, username)
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
	publisher, err :=dbProvider.PublisherGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publisherJSON, err := Publisher2JSON(publisher)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(publisherJSON)
}
func PublisherGetById(w http.ResponseWriter, r *http.Request) {
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
	publisher, err :=dbProvider.PublisherGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publisherJSON, err := Publisher2JSON(publisher)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publisherJSON)
	return
}
func PublisherUpdate(w http.ResponseWriter, r *http.Request) {
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
	publisherParams := &PublisherParams{}
	err = json.Unmarshal(body, publisherParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.PublisherUpdate(idInt64, publisherParams.Name, username)
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
	publisher, err :=dbProvider.PublisherGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publisherJSON, err := Publisher2JSON(publisher)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publisherJSON)
	return
}
func PublisherDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.PublisherDelete(idInt64)
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
func PublisherGetPublications(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	publications, err :=dbProvider.PublicationGetByPublisher(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationsJSON, err := Publications2JSON(publications)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(publicationsJSON)
	return
}
