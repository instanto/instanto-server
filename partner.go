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

type PartnerParams struct {
	lib.Partner
	Member       int64 `json:"member"`
	ResearchLine int64 `json:"research_line"`
}

func Partners2JSON(partners []*lib.Partner) (partnersJSON []byte, err error) {
	data := make(map[string]interface{})
	data["partners"] = partners
	partnersJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Partner2JSON(partner *lib.Partner) (partnerJSON []byte, err error) {
	partnerJSON, err = json.Marshal(partner)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func PartnerGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	partners, err :=dbProvider.PartnerGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	partnersJSON, err := Partners2JSON(partners)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(partnersJSON)
}
func PartnerCreate(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.PartnerCreate(partnerParams.Name, partnerParams.Web, partnerParams.SameDepartment, partnerParams.Scope, username)
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
	partner, err :=dbProvider.PartnerGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	partnerJSON, err := Partner2JSON(partner)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(partnerJSON)
}
func PartnerGetById(w http.ResponseWriter, r *http.Request) {
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
	partner, err :=dbProvider.PartnerGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	partnerJSON, err := Partner2JSON(partner)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(partnerJSON)
	return
}
func PartnerUpdate(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.PartnerUpdate(idInt64, partnerParams.Name, partnerParams.Web, partnerParams.SameDepartment, partnerParams.Scope, username)
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
	partner, err :=dbProvider.PartnerGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	partnerJSON, err := Partner2JSON(partner)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(partnerJSON)
	return
}
func PartnerDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.PartnerDelete(idInt64)
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
func PartnerGetMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	members, err :=dbProvider.PartnerGetMembers(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	membersJSON, err := Members2JSON(members)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(membersJSON)
	return
}
func PartnerAddMember(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.PartnerAddMember(idInt64, partnerParams.Member, username)
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
func PartnerRemoveMember(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.PartnerRemoveMember(idInt64, partnerParams.Member)
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
func PartnerGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err :=dbProvider.PartnerGetResearchLines(idInt64)
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
func PartnerAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.PartnerAddResearchLine(idInt64, partnerParams.ResearchLine, username)
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
func PartnerRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
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
	partnerParams := &PartnerParams{}
	err = json.Unmarshal(body, partnerParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.PartnerRemoveResearchLine(idInt64, partnerParams.ResearchLine)
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

func PartnerUpdateLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.PartnerExists(idInt64)
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
	filename := "partner_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err :=dbProvider.PartnerUpdateLogo(idInt64, filename, username)
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

func PartnerDeleteLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.PartnerExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err :=dbProvider.PartnerUpdateLogo(idInt64, "", username)
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
