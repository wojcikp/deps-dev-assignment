package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
	dependenciesupdater "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_updater"
)

type Api struct {
	db      *database.SQLiteDB
	updater *dependenciesupdater.Updater
}

func NewApi(db *database.SQLiteDB, updater *dependenciesupdater.Updater) *Api {
	return &Api{db, updater}
}

func (a *Api) addDependency(w http.ResponseWriter, r *http.Request) {
	defer log.Print("added new Dependency")
	var dependency dependenciesloader.DependencyDetails
	if err := json.NewDecoder(r.Body).Decode(&dependency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := a.db.AddNewDependencyDetails(dependency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dependency)
}

func (a *Api) updateDependency(w http.ResponseWriter, r *http.Request) {
	defer log.Print("updated Dependency")
	var dependency dependenciesloader.DependencyDetails
	if err := json.NewDecoder(r.Body).Decode(&dependency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := a.db.UpdateDependencyDetails(dependency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dependency)
}

func (a *Api) getDependencyByID(w http.ResponseWriter, r *http.Request) {
	defer log.Print("get Dependency By ID")
	id := r.URL.Query().Get("id")
	log.Print(id)
	dependency, err := a.db.GetDependencyDetailsByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dependency)
}

func (a *Api) deleteDependency(w http.ResponseWriter, r *http.Request) {
	defer log.Print("deleted Dependency")
	id := r.URL.Query().Get("id")
	log.Print(id)
	err := a.db.DeleteDependencyWithDetails(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *Api) getDependencyByScore(w http.ResponseWriter, r *http.Request) {
	defer log.Print("get Dependency By Score")
	scoreParam := mux.Vars(r)["score"]
	score, err := strconv.ParseFloat(scoreParam, 64)
	if err != nil {
		http.Error(w, "Invalid score", http.StatusInternalServerError)
		return
	}
	results, err := a.db.GetDependenciesByOverallScore(score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (a *Api) getAllDependencies(w http.ResponseWriter, r *http.Request) {
	defer log.Print("get all Dependencies")
	results, err := a.db.GetAllDependencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (a *Api) updateAllDependencies(w http.ResponseWriter, r *http.Request) {
	defer log.Print("updating dependencies")
	updatedDependencies, err := a.updater.UpdateDependencies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(updatedDependencies)
}

func (a *Api) Run() {
	r := mux.NewRouter()
	h := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "application/json"}),
	)(r)

	r.HandleFunc("/dependency", a.getDependencyByID).Methods("GET")
	r.HandleFunc("/dependency/score/{score}", a.getDependencyByScore).Methods("GET")
	r.HandleFunc("/dependency/all", a.getAllDependencies).Methods("GET")
	r.HandleFunc("/dependency/update", a.updateAllDependencies).Methods("GET")
	r.HandleFunc("/dependency", a.addDependency).Methods("POST")
	r.HandleFunc("/dependency", a.updateDependency).Methods("PUT")
	r.HandleFunc("/dependency", a.deleteDependency).Methods("DELETE")

	http.ListenAndServe(":3000", h)
}
