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

type PublicationTypeParams struct {
	lib.PublicationType
}

func PublicationTypes2JSON(publicationTypes []*lib.PublicationType) (publicationTypesJSON []byte, err error) {
	data := make(map[string]interface{})
	data["publication_types"] = publicationTypes
	publicationTypesJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func PublicationType2JSON(publicationType *lib.PublicationType) (publicationTypeJSON []byte, err error) {
	publicationTypeJSON, err = json.Marshal(publicationType)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func PublicationTypeGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	publicationTypes, err :=dbProvider.PublicationTypeGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationTypesJSON, err := PublicationTypes2JSON(publicationTypes)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationTypesJSON)
}
func PublicationTypeCreate(w http.ResponseWriter, r *http.Request) {
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
	publicationTypeParams := &PublicationTypeParams{}
	err = json.Unmarshal(body, publicationTypeParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.PublicationTypeCreate(publicationTypeParams.Name, username)
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
	publicationType, err :=dbProvider.PublicationTypeGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationTypeJSON, err := PublicationType2JSON(publicationType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(publicationTypeJSON)
}
func PublicationTypeGetById(w http.ResponseWriter, r *http.Request) {
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
	publicationType, err :=dbProvider.PublicationTypeGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationTypeJSON, err := PublicationType2JSON(publicationType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationTypeJSON)
	return
}
func PublicationTypeUpdate(w http.ResponseWriter, r *http.Request) {
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
	publicationTypeParams := &PublicationTypeParams{}
	err = json.Unmarshal(body, publicationTypeParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.PublicationTypeUpdate(idInt64, publicationTypeParams.Name, username)
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
	publicationType, err :=dbProvider.PublicationTypeGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationTypeJSON, err := PublicationType2JSON(publicationType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationTypeJSON)
	return
}
func PublicationTypeDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.PublicationTypeDelete(idInt64)
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

func PublicationTypeGetPublications(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	publications, err :=dbProvider.PublicationGetByPublicationType(idInt64)
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
