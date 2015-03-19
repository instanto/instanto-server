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

type FinancedProjectParams struct {
	lib.FinancedProject
	FundingBody  int64  `json:"funding_body"`
	Record       string `json:"record"`
	Member       int64  `json:"member"`
	Leader       int64  `json:"leader"`
	ResearchLine int64  `json:"research_line"`
}

func FinancedProjects2JSON(financedProjects []*lib.FinancedProject) (financedProjectsJSON []byte, err error) {
	data := make(map[string]interface{})
	data["financed_projects"] = financedProjects
	financedProjectsJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func FinancedProject2JSON(financedProject *lib.FinancedProject) (financedProjectJSON []byte, err error) {
	financedProjectJSON, err = json.Marshal(financedProject)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func FinancedProjectGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	financedProjects, err := dbProvider.FinancedProjectGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectsJSON, err := FinancedProjects2JSON(financedProjects)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(financedProjectsJSON)
}
func FinancedProjectCreate(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err := dbProvider.FinancedProjectCreate(financedProjectParams.Title, financedProjectParams.Started, financedProjectParams.Ended, financedProjectParams.Budget, financedProjectParams.Scope, username, financedProjectParams.PrimaryFundingBody, financedProjectParams.PrimaryRecord, financedProjectParams.PrimaryLeader)
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
	financedProject, err := dbProvider.FinancedProjectGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectJSON, err := FinancedProject2JSON(financedProject)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(financedProjectJSON)
}
func FinancedProjectGetById(w http.ResponseWriter, r *http.Request) {
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
	financedProject, err := dbProvider.FinancedProjectGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectJSON, err := FinancedProject2JSON(financedProject)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(financedProjectJSON)
	return
}
func FinancedProjectUpdate(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err := dbProvider.FinancedProjectUpdate(idInt64, financedProjectParams.Title, financedProjectParams.Started, financedProjectParams.Ended, financedProjectParams.Budget, financedProjectParams.Scope, username, financedProjectParams.PrimaryFundingBody, financedProjectParams.PrimaryRecord, financedProjectParams.PrimaryLeader)
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
	financedProject, err := dbProvider.FinancedProjectGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	financedProjectJSON, err := FinancedProject2JSON(financedProject)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(financedProjectJSON)
	return
}
func FinancedProjectDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err := dbProvider.FinancedProjectDelete(idInt64)
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
func FinancedProjectGetPrimaryFundingBody(w http.ResponseWriter, r *http.Request) {
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
	financedProject, err := dbProvider.FinancedProjectGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBody, err := dbProvider.FundingBodyGetById(financedProject.PrimaryFundingBody)
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
func FinancedProjectGetSecondaryFundingBodies(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	fundingBodies, err := dbProvider.FinancedProjectGetFundingBodies(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	fundingBodiesJSON, err := FundingBodies2JSON(fundingBodies)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(fundingBodiesJSON)
	return
}
func FinancedProjectAddSecondaryFundingBody(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.FinancedProjectAddFundingBody(idInt64, financedProjectParams.FundingBody, financedProjectParams.Record, username)
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
func FinancedProjectRemoveSecondaryFundingBody(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.FinancedProjectRemoveFundingBody(idInt64, financedProjectParams.FundingBody)
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
func FinancedProjectGetMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	members, err := dbProvider.FinancedProjectGetMembers(idInt64)
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
func FinancedProjectAddMember(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.FinancedProjectAddMember(idInt64, financedProjectParams.Member, username)
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
func FinancedProjectRemoveMember(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.FinancedProjectRemoveMember(idInt64, financedProjectParams.Member)
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
func FinancedProjectGetPrimaryLeader(w http.ResponseWriter, r *http.Request) {
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
	financedProject, err := dbProvider.FinancedProjectGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	leader, err := dbProvider.MemberGetById(financedProject.PrimaryLeader)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	leaderJSON, err := Member2JSON(leader)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(leaderJSON)
	return
}
func FinancedProjectGetSecondaryLeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	leaders, err := dbProvider.FinancedProjectGetLeaders(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	leadersJSON, err := Members2JSON(leaders)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(leadersJSON)
	return
}
func FinancedProjectAddSecondaryLeader(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.FinancedProjectAddLeader(idInt64, financedProjectParams.Leader, username)
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
func FinancedProjectRemoveSecondaryLeader(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.FinancedProjectRemoveLeader(idInt64, financedProjectParams.Leader)
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
func FinancedProjectGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err := dbProvider.FinancedProjectGetResearchLines(idInt64)
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
func FinancedProjectAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.FinancedProjectAddResearchLine(idInt64, financedProjectParams.ResearchLine, username)
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
func FinancedProjectRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
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
	financedProjectParams := &FinancedProjectParams{}
	err = json.Unmarshal(body, financedProjectParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.FinancedProjectRemoveResearchLine(idInt64, financedProjectParams.ResearchLine)
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
