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

type ResearchAreaParams struct {
	lib.ResearchArea
	ResearchLine int64 `json:"research_line"`
}

func ResearchAreas2JSON(researchAreas []*lib.ResearchArea) (researchAreasJSON []byte, err error) {
	data := make(map[string]interface{})
	data["researcharea_collection"] = researchAreas
	researchAreasJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func ResearchArea2JSON(researchArea *lib.ResearchArea) (researchAreaJSON []byte, err error) {
	researchAreaJSON, err = json.Marshal(researchArea)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func ResearchAreaGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	researchAreas, err :=dbProvider.ResearchAreaGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchAreasJSON, err := ResearchAreas2JSON(researchAreas)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchAreasJSON)
}
func ResearchAreaCreate(w http.ResponseWriter, r *http.Request) {
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
	researchAreaParams := &ResearchAreaParams{}
	err = json.Unmarshal(body, researchAreaParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.ResearchAreaCreate(researchAreaParams.Name, username)
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
	researchArea, err :=dbProvider.ResearchAreaGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchAreaJSON, err := ResearchArea2JSON(researchArea)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(researchAreaJSON)
}
func ResearchAreaGetById(w http.ResponseWriter, r *http.Request) {
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
	researchArea, err :=dbProvider.ResearchAreaGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchAreaJSON, err := ResearchArea2JSON(researchArea)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchAreaJSON)
	return
}
func ResearchAreaUpdate(w http.ResponseWriter, r *http.Request) {
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
	researchAreaParams := &ResearchAreaParams{}
	err = json.Unmarshal(body, researchAreaParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.ResearchAreaUpdate(idInt64, researchAreaParams.Name, username)
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
	researchArea, err :=dbProvider.ResearchAreaGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchAreaJSON, err := ResearchArea2JSON(researchArea)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchAreaJSON)
	return
}
func ResearchAreaDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.ResearchAreaDelete(idInt64)
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
func ResearchAreaGetPrimaryResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err :=dbProvider.ResearchLineGetByPrimaryResearchArea(idInt64)
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
func ResearchAreaGetSecondaryResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err :=dbProvider.ResearchAreaGetResearchLines(idInt64)
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
func ResearchAreaAddSecondaryResearchLine(w http.ResponseWriter, r *http.Request) {
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
	researchAreaParams := &ResearchAreaParams{}
	err = json.Unmarshal(body, researchAreaParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchAreaAddResearchLine(idInt64, researchAreaParams.ResearchLine, username)
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
func ResearchAreaRemoveSecondaryResearchLine(w http.ResponseWriter, r *http.Request) {
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
	researchAreaParams := &ResearchAreaParams{}
	err = json.Unmarshal(body, researchAreaParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchAreaRemoveResearchLine(idInt64, researchAreaParams.ResearchLine)
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
func ResearchAreaUpdateLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.ResearchAreaExists(idInt64)
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
	filename := "researcharea_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err :=dbProvider.ResearchAreaUpdateLogo(idInt64, filename, username)
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

func ResearchAreaDeleteLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.ResearchAreaExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err :=dbProvider.ResearchAreaUpdateLogo(idInt64, "", username)
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
