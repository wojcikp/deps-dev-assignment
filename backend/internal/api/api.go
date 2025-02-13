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
)

type Api struct {
	db *database.SQLiteDB
}

func NewApi(db *database.SQLiteDB) *Api {
	return &Api{db}
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
		http.Error(w, "Invalid score", http.StatusBadRequest)
		return
	}
	results, err := a.db.GetDependenciesByOverallScore(score)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (a *Api) Run() {
	// a.db.CreateTables()

	// if err := a.db.CreateTables(); err != nil {
	// 	log.Fatalf("failed to create db tables due to an error: %v \n exiting...", err)
	// }

	// loader := dependenciesloader.NewDependenciesLoader("https://api.deps.dev/v3/systems/GO/packages/github.com%2Fcli%2Fcli/versions/v1.14.0:dependencies")

	// if err := loader.FetchDepsDevDependencies(); err != nil {
	// 	log.Fatalf("failed to fetch deps.dev dependencies due to an error: %v \n exiting...", err)
	// }

	// detailedDependencies := loader.FetchDetailsForAllDependencies()

	// if err := a.db.LoadDependencies(detailedDependencies); err != nil {
	// 	log.Fatalf("failed to load detailed dependencies into db due to an error: %v \n exiting...", err)
	// }

	r := mux.NewRouter()
	h := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:8080"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "application/json"}),
	)(r)

	r.HandleFunc("/dependency", a.addDependency).Methods("POST")
	r.HandleFunc("/dependency", a.updateDependency).Methods("PUT")
	r.HandleFunc("/dependency", a.getDependencyByID).Methods("GET")
	r.HandleFunc("/dependency", a.deleteDependency).Methods("DELETE")
	r.HandleFunc("/dependency/score/{score}", a.getDependencyByScore).Methods("GET")

	http.ListenAndServe(":3000", h)
}
