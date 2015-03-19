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

type StudentWorkTypeParams struct {
	lib.StudentWorkType
}

func StudentWorkTypes2JSON(studentWorkTypes []*lib.StudentWorkType) (studentWorkTypesJSON []byte, err error) {
	data := make(map[string]interface{})
	data["student_work_types"] = studentWorkTypes
	studentWorkTypesJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func StudentWorkType2JSON(studentWorkType *lib.StudentWorkType) (studentWorkTypeJSON []byte, err error) {
	studentWorkTypeJSON, err = json.Marshal(studentWorkType)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func StudentWorkTypeGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	studentWorkTypes, err :=dbProvider.StudentWorkTypeGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkTypesJSON, err := StudentWorkTypes2JSON(studentWorkTypes)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorkTypesJSON)
}
func StudentWorkTypeCreate(w http.ResponseWriter, r *http.Request) {
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
	studentWorkTypeParams := &StudentWorkTypeParams{}
	err = json.Unmarshal(body, studentWorkTypeParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.StudentWorkTypeCreate(studentWorkTypeParams.Name, username)
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
	studentWorkType, err :=dbProvider.StudentWorkTypeGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkTypeJSON, err := StudentWorkType2JSON(studentWorkType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(studentWorkTypeJSON)
}
func StudentWorkTypeGetById(w http.ResponseWriter, r *http.Request) {
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
	studentWorkType, err :=dbProvider.StudentWorkTypeGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkTypeJSON, err := StudentWorkType2JSON(studentWorkType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorkTypeJSON)
	return
}
func StudentWorkTypeUpdate(w http.ResponseWriter, r *http.Request) {
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
	studentWorkTypeParams := &StudentWorkTypeParams{}
	err = json.Unmarshal(body, studentWorkTypeParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.StudentWorkTypeUpdate(idInt64, studentWorkTypeParams.Name, username)
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
	studentWorkType, err :=dbProvider.StudentWorkTypeGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkTypeJSON, err := StudentWorkType2JSON(studentWorkType)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorkTypeJSON)
	return
}
func StudentWorkTypeDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.StudentWorkTypeDelete(idInt64)
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

func StudentWorkTypeGetStudentWorks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	studentWorks, err :=dbProvider.StudentWorkGetByStudentWorkType(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorksJSON, err := StudentWorks2JSON(studentWorks)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(studentWorksJSON)
	return
}
