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

type FundingBodyParams struct {
	lib.FundingBody
	FinancedProject int64  `json:"financed_project"`
	Record          string `json:"record"`
}

func FundingBodies2JSON(fundingBodies []*lib.FundingBody) (fundingBodiesJSON []byte, err error) {
	data := make(map[string]interface{})
	data["fundingbody_collection"] = fundingBodies
	fundingBodiesJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func FundingBody2JSON(fundingBody *lib.FundingBody) (fundingBodyJSON []byte, err error) {
	fundingBodyJSON, err = json.Marshal(fundingBody)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func FundingBodyGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fundingBodies, err :=dbProvider.FundingBodyGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBodiesJSON, err := FundingBodies2JSON(fundingBodies)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(fundingBodiesJSON)
}
func FundingBodyCreate(w http.ResponseWriter, r *http.Request) {
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
	fundingBodyParams := &FundingBodyParams{}
	err = json.Unmarshal(body, fundingBodyParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.FundingBodyCreate(fundingBodyParams.Name, fundingBodyParams.Web, fundingBodyParams.Scope, username)
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
	fundingBody, err :=dbProvider.FundingBodyGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBodyJSON, err := FundingBody2JSON(fundingBody)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(fundingBodyJSON)
}
func FundingBodyGetById(w http.ResponseWriter, r *http.Request) {
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
	fundingBody, err :=dbProvider.FundingBodyGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBodyJSON, err := FundingBody2JSON(fundingBody)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(fundingBodyJSON)
	return
}
func FundingBodyUpdate(w http.ResponseWriter, r *http.Request) {
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
	fundingBodyParams := &FundingBodyParams{}
	err = json.Unmarshal(body, fundingBodyParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.FundingBodyUpdate(idInt64, fundingBodyParams.Name, fundingBodyParams.Web, fundingBodyParams.Scope, username)
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
	fundingBody, err :=dbProvider.FundingBodyGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBodyJSON, err := FundingBody2JSON(fundingBody)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(fundingBodyJSON)
	return
}
func FundingBodyDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.FundingBodyDelete(idInt64)
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
func FundingBodyGetPrimaryFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err :=dbProvider.FinancedProjectGetByPrimaryFundingBody(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectsJSON, err := FinancedProjects2JSON(financedProjects)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(financedProjectsJSON)
	return
}
func FundingBodyGetSecondaryFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err :=dbProvider.FundingBodyGetFinancedProjects(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectsJSON, err := FinancedProjects2JSON(financedProjects)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(financedProjectsJSON)
	return
}
func FundingBodyAddSecondaryFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	fundingBodyParams := &FundingBodyParams{}
	err = json.Unmarshal(body, fundingBodyParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.FundingBodyAddFinancedProject(idInt64, fundingBodyParams.FinancedProject, fundingBodyParams.Record, username)
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
	w.WriteHeader(201)
	return
}
func FundingBodyRemoveSecondaryFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	fundingBodyParams := &FundingBodyParams{}
	err = json.Unmarshal(body, fundingBodyParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.FundingBodyRemoveFinancedProject(idInt64, fundingBodyParams.FinancedProject)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !removed {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)
	return
}
