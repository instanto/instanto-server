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

type PublicationParams struct {
	lib.Publication
	ResearchLine int64 `json:"research_line"`
	Member       int64 `json:"member"`
}

func Publications2JSON(publications []*lib.Publication) (publicationsJSON []byte, err error) {
	data := make(map[string]interface{})
	data["publications"] = publications
	publicationsJSON, err = json.Marshal(data)
	if err != nil {
		LogError(err)
		return
	}
	return
}
func Publication2JSON(publication *lib.Publication) (publicationJSON []byte, err error) {
	publicationJSON, err = json.Marshal(publication)
	if err != nil {
		LogError(err)
		return
	}
	return
}

func PublicationGetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	publications, err :=dbProvider.PublicationGetAll()
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationsJSON, err := Publications2JSON(publications)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationsJSON)
}
func PublicationCreate(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	id, verr, err :=dbProvider.PublicationCreate(publicationParams.Title, publicationParams.Year, publicationParams.BookTitle, publicationParams.Chapter, publicationParams.City, publicationParams.Country, publicationParams.ConferenceName, publicationParams.Edition, publicationParams.Institution, publicationParams.Isbn, publicationParams.Issn, publicationParams.Journal, publicationParams.Language, publicationParams.Nationality, publicationParams.Number, publicationParams.Organization, publicationParams.Pages, publicationParams.School, publicationParams.Series, publicationParams.Volume, username, publicationParams.PublicationType, publicationParams.Publisher, publicationParams.PrimaryAuthor)
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
	publication, err :=dbProvider.PublicationGetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationJSON, err := Publication2JSON(publication)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(201)
	w.Write(publicationJSON)
}
func PublicationGetById(w http.ResponseWriter, r *http.Request) {
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
	publication, err :=dbProvider.PublicationGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationJSON, err := Publication2JSON(publication)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationJSON)
	return
}
func PublicationUpdate(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	numRows, verr, err :=dbProvider.PublicationUpdate(idInt64, publicationParams.Title, publicationParams.Year, publicationParams.BookTitle, publicationParams.Chapter, publicationParams.City, publicationParams.Country, publicationParams.ConferenceName, publicationParams.Edition, publicationParams.Institution, publicationParams.Isbn, publicationParams.Issn, publicationParams.Journal, publicationParams.Language, publicationParams.Nationality, publicationParams.Number, publicationParams.Organization, publicationParams.Pages, publicationParams.School, publicationParams.Series, publicationParams.Volume, username, publicationParams.PublicationType, publicationParams.Publisher, publicationParams.PrimaryAuthor)
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
	publication, err :=dbProvider.PublicationGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationJSON, err := Publication2JSON(publication)
	if err != nil {
		LogError(err)
		w.WriteHeader(500)
		return
	}
	w.Write(publicationJSON)
	return
}
func PublicationDelete(w http.ResponseWriter, r *http.Request) {
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
	numRows, err :=dbProvider.PublicationDelete(idInt64)
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
func PublicationGetPublisher(w http.ResponseWriter, r *http.Request) {
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
	publication, err :=dbProvider.PublicationGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publisher, err :=dbProvider.PublisherGetById(publication.Publisher)
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
func PublicationGetPublicationType(w http.ResponseWriter, r *http.Request) {
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
	publication, err :=dbProvider.PublicationGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	publicationType, err :=dbProvider.PublicationTypeGetById(publication.PublicationType)
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
func PublicationGetPrimaryAuthor(w http.ResponseWriter, r *http.Request) {
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
	publication, err :=dbProvider.PublicationGetById(idInt64)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
		LogError(err)
		w.WriteHeader(500)
		return
	}
	member, err :=dbProvider.MemberGetById(publication.PrimaryAuthor)
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
func PublicationGetResearchLines(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	researchLines, err :=dbProvider.PublicationGetResearchLines(idInt64)
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
func PublicationAddResearchLine(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.PublicationAddResearchLine(idInt64, publicationParams.ResearchLine, username)
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
func PublicationRemoveResearchLine(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.PublicationRemoveResearchLine(idInt64, publicationParams.ResearchLine)
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
func PublicationGetSecondaryAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	members, err :=dbProvider.PublicationGetMembers(idInt64)
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
func PublicationAddSecondaryAuthor(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	verr, err :=dbProvider.PublicationAddMember(idInt64, publicationParams.Member, username)
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
func PublicationRemoveSecondaryAuthor(w http.ResponseWriter, r *http.Request) {
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
	publicationParams := &PublicationParams{}
	err = json.Unmarshal(body, publicationParams)
	if err != nil {
		w.WriteHeader(415)
		return
	}
	removed, err :=dbProvider.PublicationRemoveMember(idInt64, publicationParams.Member)
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
