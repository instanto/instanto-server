package main

import (
	"database/sql"
	"encoding/json"
	lib "github.com/instanto/instanto-lib"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

type MemberParams struct {
	lib.Member
	Status          int64 `json:"status"`
	FinancedProject int64 `json:"financed_project"`
	Partner         int64 `json:"partner"`
	StudentWork     int64 `json:"student_work"`
	Publication     int64 `json:"publication"`
	ResearchLine    int64 `json:"research_line"`
}

func Members2JSON(members []*lib.Member) (membersJSON []byte, err error) {
	data := make(map[string]interface{})
	data["members"] = members
	membersJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Member2JSON(member *lib.Member) (memberJSON []byte, err error) {
	memberJSON, err = json.Marshal(member)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func MemberGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	members, err := dbProvider.MemberGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	membersJSON, err := Members2JSON(members)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(membersJSON)
}
func MemberCreate(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err := dbProvider.MemberCreate(memberParams.FirstName, memberParams.LastName, memberParams.Degree, memberParams.YearIn, memberParams.YearOut, memberParams.Email, username, memberParams.PrimaryStatus)
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
	member, err := dbProvider.MemberGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	memberJSON, err := Member2JSON(member)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(memberJSON)
}
func MemberGetById(w http.ResponseWriter, r *http.Request) {
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
	member, err := dbProvider.MemberGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	memberJSON, err := Member2JSON(member)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(memberJSON)
	return
}
func MemberUpdate(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err := dbProvider.MemberUpdate(idInt64, memberParams.FirstName, memberParams.LastName, memberParams.Degree, memberParams.YearIn, memberParams.YearOut, memberParams.Email, username, memberParams.PrimaryStatus)
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
	member, err := dbProvider.MemberGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	memberJSON, err := Member2JSON(member)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(memberJSON)
	return
}
func MemberDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err := dbProvider.MemberDelete(idInt64)
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
func MemberGetPrimaryStatus(w http.ResponseWriter, r *http.Request) {
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
	member, err := dbProvider.MemberGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	status, err := dbProvider.StatusGetById(member.PrimaryStatus)
	statusJSON, err := Status2JSON(status)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(statusJSON)
	return
}
func MemberGetSecondaryStatuses(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	statuses, err := dbProvider.MemberGetStatuses(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	statusesJSON, err := Statuses2JSON(statuses)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(statusesJSON)
	return
}
func MemberAddSecondaryStatus(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddStatus(idInt64, memberParams.Status, username)
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
func MemberRemoveSecondaryStatus(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemoveStatus(idInt64, memberParams.Status)
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
func MemberGetFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err := dbProvider.MemberGetFinancedProjects(idInt64)
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
func MemberAddFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddFinancedProject(idInt64, memberParams.FinancedProject, username)
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
func MemberRemoveFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemoveFinancedProject(idInt64, memberParams.FinancedProject)
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
func MemberGetSecondaryFinancedProjectsAsLeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err := dbProvider.MemberGetFinancedProjects(idInt64)
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
func MemberAddSecondaryFinancedProjectAsLeader(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddFinancedProject(idInt64, memberParams.FinancedProject, username)
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
func MemberRemoveSecondaryFinancedProjectAsLeader(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemoveFinancedProject(idInt64, memberParams.FinancedProject)
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
func MemberGetPartners(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	partners, err := dbProvider.MemberGetPartners(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	partnersJSON, err := Partners2JSON(partners)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(partnersJSON)
	return
}
func MemberAddPartner(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddPartner(idInt64, memberParams.Partner, username)
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
func MemberRemovePartner(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemovePartner(idInt64, memberParams.Partner)
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
func MemberGetStudentWorks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	studentWorks, err := dbProvider.MemberGetStudentWorks(idInt64)
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
func MemberGetPrimaryPublications(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	publications, err := dbProvider.PublicationGetByPrimaryAuthor(idInt64)
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
func MemberGetSecondaryPublications(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	publications, err := dbProvider.MemberGetPublications(idInt64)
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
func MemberAddSecondaryPublication(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddPublication(idInt64, memberParams.Publication, username)
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
func MemberRemoveSecondaryPublication(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemovePublication(idInt64, memberParams.Publication)
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
func MemberGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLinea, err := dbProvider.MemberGetResearchLines(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLineaJSON, err := ResearchLines2JSON(researchLinea)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(researchLineaJSON)
	return
}
func MemberAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddResearchLine(idInt64, memberParams.ResearchLine, username)
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
func MemberRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemoveResearchLine(idInt64, memberParams.ResearchLine)
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
func MemberGetPrimaryLeaderedFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err := dbProvider.FinancedProjectGetByPrimaryLeader(idInt64)
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
func MemberGetSecondaryLeaderedFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err := dbProvider.MemberGetFinancedProjectsAsLeader(idInt64)
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
func MemberAddSecondaryLeaderedFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err := dbProvider.MemberAddFinancedProjectAsLeader(idInt64, memberParams.FinancedProject, username)
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
func MemberRemoveSecondaryLeaderedFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	memberParams := &MemberParams{}
	err = json.Unmarshal(body, memberParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err := dbProvider.MemberRemoveFinancedProjectAsLeader(idInt64, memberParams.FinancedProject)
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
func MemberUpdateCv(w http.ResponseWriter, r *http.Request) {
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
	exists, err := dbProvider.MemberExists(idInt64)
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
		verr := &lib.ValidationError{"cv", "filesize must be greater than 0 and less than 2Mb"}
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
		verr := &lib.ValidationError{"cv", "filesize must be greater than 0 and less than 2Mb"}
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
	filename := "member_cv_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err := dbProvider.MemberUpdateCv(idInt64, filename, username)
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

func MemberDeleteCv(w http.ResponseWriter, r *http.Request) {
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
	exists, err := dbProvider.MemberExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err := dbProvider.MemberUpdateCv(idInt64, "", username)
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
func MemberUpdatePhoto(w http.ResponseWriter, r *http.Request) {
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
	exists, err := dbProvider.MemberExists(idInt64)
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
		verr := &lib.ValidationError{"photo", "filesize must be greater than 0 and less than 2Mb"}
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
		verr := &lib.ValidationError{"photo", "filesize must be greater than 0 and less than 2Mb"}
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
	filename := "member_photo_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err := dbProvider.MemberUpdatePhoto(idInt64, filename, username)
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

func MemberDeletePhoto(w http.ResponseWriter, r *http.Request) {
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
	exists, err := dbProvider.MemberExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err := dbProvider.MemberUpdatePhoto(idInt64, "", username)
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
