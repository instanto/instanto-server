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

type ResearchLineParams struct {
	lib.ResearchLine
	ResearchArea    int64 `json:"research_area"`
	FinancedProject int64 `json:"financed_project"`
	Article         int64 `json:"article"`
	Member          int64 `json:"member"`
	Partner         int64 `json:"partner"`
	Publication     int64 `json:"publication"`
	StudentWork     int64 `json:"student_work"`
}

func ResearchLines2JSON(researchLines []*lib.ResearchLine) (researchLinesJSON []byte, err error) {
	data := make(map[string]interface{})
	data["research_lines"] = researchLines
	researchLinesJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func ResearchLine2JSON(researchLine *lib.ResearchLine) (researchLineJSON []byte, err error) {
	researchLineJSON, err = json.Marshal(researchLine)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func ResearchLineGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	researchLines, err :=dbProvider.ResearchLineGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLinesJSON, err := ResearchLines2JSON(researchLines)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchLinesJSON)
}
func ResearchLineCreate(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.ResearchLineCreate(researchLineParams.Title, researchLineParams.Finished, researchLineParams.Description, username, researchLineParams.PrimaryResearchArea)
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
	researchLine, err :=dbProvider.ResearchLineGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLineJSON, err := ResearchLine2JSON(researchLine)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(researchLineJSON)
}
func ResearchLineGetById(w http.ResponseWriter, r *http.Request) {
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
	researchLine, err :=dbProvider.ResearchLineGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLineJSON, err := ResearchLine2JSON(researchLine)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchLineJSON)
	return
}
func ResearchLineUpdate(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.ResearchLineUpdate(idInt64, researchLineParams.Title, researchLineParams.Finished, researchLineParams.Description, username, researchLineParams.PrimaryResearchArea)
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
	researchLine, err :=dbProvider.ResearchLineGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchLineJSON, err := ResearchLine2JSON(researchLine)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(researchLineJSON)
	return
}
func ResearchLineDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.ResearchLineDelete(idInt64)
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
func ResearchLineGetPrimaryResearchArea(w http.ResponseWriter, r *http.Request) {
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
	researchLine, err :=dbProvider.ResearchLineGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchArea, err :=dbProvider.ResearchAreaGetById(researchLine.PrimaryResearchArea)
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
func ResearchLineGetSecondaryResearchAreas(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchAreas, err :=dbProvider.ResearchLineGetResearchAreas(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	researchAreasJSON, err := ResearchAreas2JSON(researchAreas)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(researchAreasJSON)
	return
}
func ResearchLineAddSecondaryResearchArea(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddResearchArea(idInt64, researchLineParams.ResearchArea, username)
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
func ResearchLineRemoveSecondaryResearchArea(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemoveResearchArea(idInt64, researchLineParams.ResearchArea)
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

func ResearchLineGetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	articles, err :=dbProvider.ResearchLineGetArticles(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	articlesJSON, err := Articles2JSON(articles)
	if err != nil {
		LogError(err)
		return
	}
	w.Write(articlesJSON)
	return
}
func ResearchLineAddArticle(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddArticle(idInt64, researchLineParams.Article, username)
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
func ResearchLineRemoveArticle(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemoveArticle(idInt64, researchLineParams.Article)
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
func ResearchLineGetFinancedProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	financedProjects, err :=dbProvider.ResearchLineGetFinancedProjects(idInt64)
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
func ResearchLineAddFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddFinancedProject(idInt64, researchLineParams.FinancedProject, username)
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
func ResearchLineRemoveFinancedProject(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemoveFinancedProject(idInt64, researchLineParams.FinancedProject)
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
func ResearchLineGetPartners(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	partners, err :=dbProvider.ResearchLineGetPartners(idInt64)
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
func ResearchLineAddPartner(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddPartner(idInt64, researchLineParams.Partner, username)
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
func ResearchLineRemovePartner(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemovePartner(idInt64, researchLineParams.Partner)
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
func ResearchLineGetMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	members, err :=dbProvider.ResearchLineGetMembers(idInt64)
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
func ResearchLineAddMember(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddMember(idInt64, researchLineParams.Member, username)
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
func ResearchLineRemoveMember(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemoveMember(idInt64, researchLineParams.Member)
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

func ResearchLineGetPublications(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	publications, err :=dbProvider.ResearchLineGetPublications(idInt64)
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
func ResearchLineAddPublication(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddPublication(idInt64, researchLineParams.Publication, username)
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
func ResearchLineRemovePublication(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemovePublication(idInt64, researchLineParams.Publication)
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

func ResearchLineGetStudentWorks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	studentWorks, err :=dbProvider.ResearchLineGetStudentWorks(idInt64)
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
func ResearchLineAddStudentWork(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.ResearchLineAddStudentWork(idInt64, researchLineParams.StudentWork, username)
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
func ResearchLineRemoveStudentWork(w http.ResponseWriter, r *http.Request) {
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
	researchLineParams := &ResearchLineParams{}
	err = json.Unmarshal(body, researchLineParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.ResearchLineRemoveStudentWork(idInt64, researchLineParams.StudentWork)
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
func ResearchLineUpdateLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.ResearchLineExists(idInt64)
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
	filename := "researchline_" + id + path.Ext(r.URL.Query().Get("filename"))
	_, verr, err :=dbProvider.ResearchLineUpdateLogo(idInt64, filename, username)
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

func ResearchLineDeleteLogo(w http.ResponseWriter, r *http.Request) {
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
	exists, err :=dbProvider.ResearchLineExists(idInt64)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	if !exists {
		w.WriteHeader(404)
		return
	}
	_, verr, err :=dbProvider.ResearchLineUpdateLogo(idInt64, "", username)
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
