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

type StudentWorkParams struct {
	lib.StudentWork
	ResearchLine    int64 `json:"research_line"`
	StudentWorkType int64 `json:"student_work_type"`
}

func StudentWorks2JSON(studentWorks []*lib.StudentWork) (studentWorksJSON []byte, err error) {
	data := make(map[string]interface{})
	data["studentwork_collection"] = studentWorks
	studentWorksJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func StudentWork2JSON(studentWork *lib.StudentWork) (studentWorkJSON []byte, err error) {
	studentWorkJSON, err = json.Marshal(studentWork)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func StudentWorkGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	studentWorks, err :=dbProvider.StudentWorkGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorksJSON, err := StudentWorks2JSON(studentWorks)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorksJSON)
}
func StudentWorkCreate(w http.ResponseWriter, r *http.Request) {
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
	studentWorkParams := &StudentWorkParams{}
	err = json.Unmarshal(body, studentWorkParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.StudentWorkCreate(studentWorkParams.Title, studentWorkParams.Year, studentWorkParams.School, studentWorkParams.Volume, username, studentWorkParams.StudentWorkType, studentWorkParams.Author)
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
	studentWork, err :=dbProvider.StudentWorkGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkJSON, err := StudentWork2JSON(studentWork)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(studentWorkJSON)
}
func StudentWorkGetById(w http.ResponseWriter, r *http.Request) {
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
	studentWork, err :=dbProvider.StudentWorkGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkJSON, err := StudentWork2JSON(studentWork)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorkJSON)
	return
}
func StudentWorkUpdate(w http.ResponseWriter, r *http.Request) {
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
	studentWorkParams := &StudentWorkParams{}
	err = json.Unmarshal(body, studentWorkParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.StudentWorkUpdate(idInt64, studentWorkParams.Title, studentWorkParams.Year, studentWorkParams.School, studentWorkParams.Volume, username, studentWorkParams.StudentWorkType, studentWorkParams.Author)
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
	studentWork, err :=dbProvider.StudentWorkGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkJSON, err := StudentWork2JSON(studentWork)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(studentWorkJSON)
	return
}
func StudentWorkDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.StudentWorkDelete(idInt64)
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
func StudentWorkGetStudentWorkType(w http.ResponseWriter, r *http.Request) {
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
	studentWork, err :=dbProvider.StudentWorkGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	studentWorkType, err :=dbProvider.StudentWorkTypeGetById(studentWork.StudentWorkType)
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
func StudentWorkGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err :=dbProvider.StudentWorkGetResearchLines(idInt64)
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
func StudentWorkAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	studentWorkParams := &StudentWorkParams{}
	err = json.Unmarshal(body, studentWorkParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.StudentWorkAddResearchLine(idInt64, studentWorkParams.ResearchLine, username)
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
func StudentWorkRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
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
	studentWorkParams := &StudentWorkParams{}
	err = json.Unmarshal(body, studentWorkParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.StudentWorkRemoveResearchLine(idInt64, studentWorkParams.ResearchLine)
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
