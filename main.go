package main

import (
	"flag"
	"fmt"
	lib "github.com/instanto/instanto-lib"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Global DbProvider to be able to use it from other files of this main package
var dbProvider *lib.DBProvider
var config *Config

type ServerCORS struct {
	*mux.Router
	config *Config
}

func (s *ServerCORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.config.CORSEnabled {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Add("Access-Control-Allow-Origin", s.config.AccessControlAllowOrigin)
			w.Header().Add("Access-Control-Allow-Methods", s.config.AccessControlAllowMethods)
			w.Header().Add("Access-Control-Allow-Headers", s.config.AccessControlAllowHeaders)
		}
		// Stop here if it is a Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// Lets Gorilla work
	s.Router.ServeHTTP(w, r)
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./config.json", "config=/etc/config.json")
	flag.Parse()

	configProvider, err := NewConfigProvider(configFile)
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
	config, err = configProvider.Parse()
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
	dbProvider, err = lib.NewDBProvider(config.DSN)
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
	router := mux.NewRouter().StrictSlash(true)
	registerCoreAPIs(router, config)
	if config.ServeWebApps == true {
		registerAdminApp(router, config)
		registerPublicApp(router, config)
	}

	LogInfo(fmt.Sprintf("Instanto server listenning on port %d . . .", config.Port))
	http.Handle("/", &ServerCORS{router, config})
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		LogError(err)
		os.Exit(1)
	}
}

func registerPublicApp(router *mux.Router, config *Config) {
	router.PathPrefix(config.WebAppPublicURL).Handler(http.StripPrefix(config.WebAppPublicURL, http.FileServer(http.Dir(config.WebAppPublicDir))))
}

func registerAdminApp(router *mux.Router, config *Config) {
	router.PathPrefix(config.WebAppAdminURL).Handler(http.StripPrefix(config.WebAppAdminURL, http.FileServer(http.Dir(config.WebAppAdminDir))))
}

func registerCoreAPIs(router *mux.Router, config *Config) {

	router.PathPrefix("/api/media/").Handler(http.StripPrefix("/api/media/", http.FileServer(http.Dir(config.MediaDir))))

	// Login
	router.HandleFunc("/api/login", Login).Methods("POST")

	// Status routes OK
	router.HandleFunc("/api/statuses", StatusGetAll).Methods("GET")
	router.HandleFunc("/api/statuses", StatusCreate).Methods("POST")
	router.HandleFunc("/api/statuses/{id}", StatusGetById).Methods("GET")
	router.HandleFunc("/api/statuses/{id}", StatusUpdate).Methods("PUT")
	router.HandleFunc("/api/statuses/{id}", StatusDelete).Methods("DELETE")
	router.HandleFunc("/api/statuses/{id}/primarymembers", StatusGetPrimaryMembers).Methods("GET")
	router.HandleFunc("/api/statuses/{id}/secondarymembers", StatusGetSecondaryMembers).Methods("GET")
	router.HandleFunc("/api/statuses/{id}/secondarymembers", StatusAddSecondaryMember).Methods("POST")
	router.HandleFunc("/api/statuses/{id}/secondarymembers", StatusRemoveSecondaryMember).Methods("DELETE")

	// Newspaper routes OK
	router.HandleFunc("/api/newspapers", NewspaperGetAll).Methods("GET")
	router.HandleFunc("/api/newspapers", NewspaperCreate).Methods("POST")
	router.HandleFunc("/api/newspapers/{id}", NewspaperGetById).Methods("GET")
	router.HandleFunc("/api/newspapers/{id}", NewspaperUpdate).Methods("PUT")
	router.HandleFunc("/api/newspapers/{id}", NewspaperDelete).Methods("DELETE")
	router.HandleFunc("/api/newspapers/{id}/articles", NewspaperGetArticles).Methods("GET")
	router.HandleFunc("/api/newspapers/{id}/logo", NewspaperUpdateLogo).Methods("PUT")
	router.HandleFunc("/api/newspapers/{id}/logo", NewspaperDeleteLogo).Methods("DELETE")

	// Article routes OK
	router.HandleFunc("/api/articles", ArticleGetAll).Methods("GET")
	router.HandleFunc("/api/articles", ArticleCreate).Methods("POST")
	router.HandleFunc("/api/articles/{id}", ArticleGetById).Methods("GET")
	router.HandleFunc("/api/articles/{id}", ArticleUpdate).Methods("PUT")
	router.HandleFunc("/api/articles/{id}", ArticleDelete).Methods("DELETE")
	router.HandleFunc("/api/articles/{id}/researchlines", ArticleGetResearchLines).Methods("GET")
	router.HandleFunc("/api/articles/{id}/researchlines", ArticleAddResearchLine).Methods("POST")
	router.HandleFunc("/api/articles/{id}/researchlines", ArticleRemoveResearchLine).Methods("DELETE")

	// FundingBody routes OK
	router.HandleFunc("/api/fundingbodies", FundingBodyGetAll).Methods("GET")
	router.HandleFunc("/api/fundingbodies", FundingBodyCreate).Methods("POST")
	router.HandleFunc("/api/fundingbodies/{id}", FundingBodyGetById).Methods("GET")
	router.HandleFunc("/api/fundingbodies/{id}", FundingBodyUpdate).Methods("PUT")
	router.HandleFunc("/api/fundingbodies/{id}", FundingBodyDelete).Methods("DELETE")
	router.HandleFunc("/api/fundingbodies/{id}/primaryfinancedprojects", FundingBodyGetPrimaryFinancedProjects).Methods("GET")
	router.HandleFunc("/api/fundingbodies/{id}/secondaryfinancedprojects", FundingBodyGetSecondaryFinancedProjects).Methods("GET")
	router.HandleFunc("/api/fundingbodies/{id}/secondaryfinancedprojects", FundingBodyAddSecondaryFinancedProject).Methods("POST")
	router.HandleFunc("/api/fundingbodies/{id}/secondaryfinancedprojects", FundingBodyRemoveSecondaryFinancedProject).Methods("DELETE")

	// FinancedProject routes OK
	router.HandleFunc("/api/financedprojects", FinancedProjectGetAll).Methods("GET")
	router.HandleFunc("/api/financedprojects", FinancedProjectCreate).Methods("POST")
	router.HandleFunc("/api/financedprojects/{id}", FinancedProjectGetById).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}", FinancedProjectUpdate).Methods("PUT")
	router.HandleFunc("/api/financedprojects/{id}", FinancedProjectDelete).Methods("DELETE")
	router.HandleFunc("/api/financedprojects/{id}/primaryfundingbody", FinancedProjectGetPrimaryFundingBody).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/secondaryfundingbodies", FinancedProjectGetSecondaryFundingBodies).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/secondaryfundingbodies", FinancedProjectAddSecondaryFundingBody).Methods("POST")
	router.HandleFunc("/api/financedprojects/{id}/secondaryfundingbodies", FinancedProjectRemoveSecondaryFundingBody).Methods("DELETE")
	router.HandleFunc("/api/financedprojects/{id}/researchlines", FinancedProjectGetResearchLines).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/researchlines", FinancedProjectAddResearchLine).Methods("POST")
	router.HandleFunc("/api/financedprojects/{id}/researchlines", FinancedProjectRemoveResearchLine).Methods("DELETE")
	router.HandleFunc("/api/financedprojects/{id}/primaryleader", FinancedProjectGetPrimaryLeader).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/secondaryleaders", FinancedProjectGetSecondaryLeaders).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/secondaryleaders", FinancedProjectAddSecondaryLeader).Methods("POST")
	router.HandleFunc("/api/financedprojects/{id}/secondaryleaders", FinancedProjectRemoveSecondaryLeader).Methods("DELETE")
	router.HandleFunc("/api/financedprojects/{id}/members", FinancedProjectGetMembers).Methods("GET")
	router.HandleFunc("/api/financedprojects/{id}/members", FinancedProjectAddMember).Methods("POST")
	router.HandleFunc("/api/financedprojects/{id}/members", FinancedProjectRemoveMember).Methods("DELETE")

	// Partner routes
	router.HandleFunc("/api/partners", PartnerGetAll).Methods("GET")
	router.HandleFunc("/api/partners", PartnerCreate).Methods("POST")
	router.HandleFunc("/api/partners/{id}", PartnerGetById).Methods("GET")
	router.HandleFunc("/api/partners/{id}", PartnerUpdate).Methods("PUT")
	router.HandleFunc("/api/partners/{id}", PartnerDelete).Methods("DELETE")
	router.HandleFunc("/api/partners/{id}/members", PartnerGetMembers).Methods("GET")
	router.HandleFunc("/api/partners/{id}/members", PartnerAddMember).Methods("POST")
	router.HandleFunc("/api/partners/{id}/members", PartnerRemoveMember).Methods("DELETE")
	router.HandleFunc("/api/partners/{id}/researchlines", PartnerGetResearchLines).Methods("GET")
	router.HandleFunc("/api/partners/{id}/researchlines", PartnerAddResearchLine).Methods("POST")
	router.HandleFunc("/api/partners/{id}/researchlines", PartnerRemoveResearchLine).Methods("DELETE")
	router.HandleFunc("/api/partners/{id}/logo", PartnerUpdateLogo).Methods("PUT")
	router.HandleFunc("/api/partners/{id}/logo", PartnerDeleteLogo).Methods("DELETE")

	// ResearchArea routes OK
	router.HandleFunc("/api/researchareas", ResearchAreaGetAll).Methods("GET")
	router.HandleFunc("/api/researchareas", ResearchAreaCreate).Methods("POST")
	router.HandleFunc("/api/researchareas/{id}", ResearchAreaGetById).Methods("GET")
	router.HandleFunc("/api/researchareas/{id}", ResearchAreaUpdate).Methods("PUT")
	router.HandleFunc("/api/researchareas/{id}", ResearchAreaDelete).Methods("DELETE")
	router.HandleFunc("/api/researchareas/{id}/primaryresearchlines", ResearchAreaGetPrimaryResearchLines).Methods("GET")
	router.HandleFunc("/api/researchareas/{id}/secondaryresearchlines", ResearchAreaGetSecondaryResearchLines).Methods("GET")
	router.HandleFunc("/api/researchareas/{id}/secondaryresearchlines", ResearchAreaAddSecondaryResearchLine).Methods("POST")
	router.HandleFunc("/api/researchareas/{id}/secondaryresearchlines", ResearchAreaRemoveSecondaryResearchLine).Methods("DELETE")
	router.HandleFunc("/api/researchareas/{id}/logo", ResearchAreaUpdateLogo).Methods("PUT")
	router.HandleFunc("/api/researchareas/{id}/logo", ResearchAreaDeleteLogo).Methods("DELETE")

	// Member routes
	router.HandleFunc("/api/members", MemberGetAll).Methods("GET")
	router.HandleFunc("/api/members", MemberCreate).Methods("POST")
	router.HandleFunc("/api/members/{id}", MemberGetById).Methods("GET")
	router.HandleFunc("/api/members/{id}", MemberUpdate).Methods("PUT")
	router.HandleFunc("/api/members/{id}", MemberDelete).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/primarystatus", MemberGetPrimaryStatus).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondarystatuses", MemberGetSecondaryStatuses).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondarystatuses", MemberAddSecondaryStatus).Methods("POST")
	router.HandleFunc("/api/members/{id}/secondarystatuses", MemberRemoveSecondaryStatus).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/financedprojects", MemberGetFinancedProjects).Methods("GET")
	router.HandleFunc("/api/members/{id}/financedprojects", MemberAddFinancedProject).Methods("POST")
	router.HandleFunc("/api/members/{id}/financedprojects", MemberRemoveFinancedProject).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/primaryleaderedfinancedprojects", MemberGetPrimaryLeaderedFinancedProjects).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondaryleaderedfinancedprojects", MemberGetSecondaryLeaderedFinancedProjects).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondaryleaderedfinancedprojects", MemberAddSecondaryLeaderedFinancedProject).Methods("POST")
	router.HandleFunc("/api/members/{id}/secondaryleaderedfinancedprojects", MemberRemoveSecondaryLeaderedFinancedProject).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/partners", MemberGetPartners).Methods("GET")
	router.HandleFunc("/api/members/{id}/partners", MemberAddPartner).Methods("POST")
	router.HandleFunc("/api/members/{id}/partners", MemberRemovePartner).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/primarypublications", MemberGetPrimaryPublications).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondarypublications", MemberGetSecondaryPublications).Methods("GET")
	router.HandleFunc("/api/members/{id}/secondarypublications", MemberAddSecondaryPublication).Methods("POST")
	router.HandleFunc("/api/members/{id}/secondarypublications", MemberRemoveSecondaryPublication).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/studentworks", MemberGetStudentWorks).Methods("GET")
	router.HandleFunc("/api/members/{id}/researchlines", MemberGetResearchLines).Methods("GET")
	router.HandleFunc("/api/members/{id}/researchlines", MemberAddResearchLine).Methods("POST")
	router.HandleFunc("/api/members/{id}/researchlines", MemberRemoveResearchLine).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/cv", MemberUpdateCv).Methods("PUT")
	router.HandleFunc("/api/members/{id}/cv", MemberDeleteCv).Methods("DELETE")
	router.HandleFunc("/api/members/{id}/photo", MemberUpdatePhoto).Methods("PUT")
	router.HandleFunc("/api/members/{id}/photo", MemberDeletePhoto).Methods("DELETE")

	// ResearchLine routes
	router.HandleFunc("/api/researchlines", ResearchLineGetAll).Methods("GET")
	router.HandleFunc("/api/researchlines", ResearchLineCreate).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}", ResearchLineGetById).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}", ResearchLineUpdate).Methods("PUT")
	router.HandleFunc("/api/researchlines/{id}", ResearchLineDelete).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/primaryresearcharea", ResearchLineGetPrimaryResearchArea).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/secondaryresearchareas", ResearchLineGetSecondaryResearchAreas).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/secondaryresearchareas", ResearchLineAddSecondaryResearchArea).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/secondaryresearchareas", ResearchLineRemoveSecondaryResearchArea).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/articles", ResearchLineGetArticles).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/articles", ResearchLineAddArticle).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/articles", ResearchLineRemoveArticle).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/financedprojects", ResearchLineGetFinancedProjects).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/financedprojects", ResearchLineAddFinancedProject).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/financedprojects", ResearchLineRemoveFinancedProject).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/partners", ResearchLineGetPartners).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/partners", ResearchLineAddPartner).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/partners", ResearchLineRemovePartner).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/members", ResearchLineGetMembers).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/members", ResearchLineAddMember).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/members", ResearchLineRemoveMember).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/publications", ResearchLineGetPublications).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/publications", ResearchLineAddPublication).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/publications", ResearchLineRemovePublication).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/studentworks", ResearchLineGetStudentWorks).Methods("GET")
	router.HandleFunc("/api/researchlines/{id}/studentworks", ResearchLineAddStudentWork).Methods("POST")
	router.HandleFunc("/api/researchlines/{id}/studentworks", ResearchLineRemoveStudentWork).Methods("DELETE")
	router.HandleFunc("/api/researchlines/{id}/logo", ResearchLineUpdateLogo).Methods("PUT")
	router.HandleFunc("/api/researchlines/{id}/logo", ResearchLineDeleteLogo).Methods("DELETE")

	// Publisher routes
	router.HandleFunc("/api/publishers", PublisherGetAll).Methods("GET")
	router.HandleFunc("/api/publishers", PublisherCreate).Methods("POST")
	router.HandleFunc("/api/publishers/{id}", PublisherGetById).Methods("GET")
	router.HandleFunc("/api/publishers/{id}", PublisherUpdate).Methods("PUT")
	router.HandleFunc("/api/publishers/{id}", PublisherDelete).Methods("DELETE")
	router.HandleFunc("/api/publishers/{id}/publications", PublisherGetPublications).Methods("GET")

	// PublicationType routes
	router.HandleFunc("/api/publicationtypes", PublicationTypeGetAll).Methods("GET")
	router.HandleFunc("/api/publicationtypes", PublicationTypeCreate).Methods("POST")
	router.HandleFunc("/api/publicationtypes/{id}", PublicationTypeGetById).Methods("GET")
	router.HandleFunc("/api/publicationtypes/{id}", PublicationTypeUpdate).Methods("PUT")
	router.HandleFunc("/api/publicationtypes/{id}", PublicationTypeDelete).Methods("DELETE")
	router.HandleFunc("/api/publicationtypes/{id}/publications", PublicationTypeGetPublications).Methods("GET")

	// Publication routes
	router.HandleFunc("/api/publications", PublicationGetAll).Methods("GET")
	router.HandleFunc("/api/publications", PublicationCreate).Methods("POST")
	router.HandleFunc("/api/publications/{id}", PublicationGetById).Methods("GET")
	router.HandleFunc("/api/publications/{id}", PublicationUpdate).Methods("PUT")
	router.HandleFunc("/api/publications/{id}", PublicationDelete).Methods("DELETE")
	router.HandleFunc("/api/publications/{id}/publisher", PublicationGetPublisher).Methods("GET")
	router.HandleFunc("/api/publications/{id}/publicationtype", PublicationGetPublicationType).Methods("GET")
	router.HandleFunc("/api/publications/{id}/primaryauthor", PublicationGetPrimaryAuthor).Methods("GET")
	router.HandleFunc("/api/publications/{id}/secondaryauthors", PublicationGetSecondaryAuthors).Methods("GET")
	router.HandleFunc("/api/publications/{id}/secondaryauthors", PublicationAddSecondaryAuthor).Methods("POST")
	router.HandleFunc("/api/publications/{id}/secondaryauthors", PublicationRemoveSecondaryAuthor).Methods("DELETE")
	router.HandleFunc("/api/publications/{id}/researchlines", PublicationGetResearchLines).Methods("GET")
	router.HandleFunc("/api/publications/{id}/researchlines", PublicationAddResearchLine).Methods("POST")
	router.HandleFunc("/api/publications/{id}/researchlines", PublicationRemoveResearchLine).Methods("DELETE")

	// StudentWorkType routes
	router.HandleFunc("/api/studentworktypes", StudentWorkTypeGetAll).Methods("GET")
	router.HandleFunc("/api/studentworktypes", StudentWorkTypeCreate).Methods("POST")
	router.HandleFunc("/api/studentworktypes/{id}", StudentWorkTypeGetById).Methods("GET")
	router.HandleFunc("/api/studentworktypes/{id}", StudentWorkTypeUpdate).Methods("PUT")
	router.HandleFunc("/api/studentworktypes/{id}", StudentWorkTypeDelete).Methods("DELETE")
	router.HandleFunc("/api/studentworktypes/{id}/studentworks", StudentWorkTypeGetStudentWorks).Methods("GET")

	// StudentWork routes
	router.HandleFunc("/api/studentworks", StudentWorkGetAll).Methods("GET")
	router.HandleFunc("/api/studentworks", StudentWorkCreate).Methods("POST")
	router.HandleFunc("/api/studentworks/{id}", StudentWorkGetById).Methods("GET")
	router.HandleFunc("/api/studentworks/{id}", StudentWorkUpdate).Methods("PUT")
	router.HandleFunc("/api/studentworks/{id}", StudentWorkDelete).Methods("DELETE")
	router.HandleFunc("/api/studentworks/{id}/studentworktype", StudentWorkGetStudentWorkType).Methods("GET")
	router.HandleFunc("/api/studentworks/{id}/publicationtype", PublicationGetPublisher).Methods("GET")
	router.HandleFunc("/api/studentworks/{id}/author", PublicationGetPrimaryAuthor).Methods("GET")
	router.HandleFunc("/api/studentworks/{id}/researchlines", StudentWorkGetResearchLines).Methods("GET")
	router.HandleFunc("/api/studentworks/{id}/researchlines", StudentWorkAddResearchLine).Methods("POST")
	router.HandleFunc("/api/studentworks/{id}/researchlines", StudentWorkRemoveResearchLine).Methods("DELETE")
}
