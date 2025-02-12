package app

import (
	"log"

	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type App struct {
	dependenciesLoader *dependenciesloader.Loader
	db                 *database.SQLiteDB
}

func NewApp(
	dependenciesLoader *dependenciesloader.Loader,
	db *database.SQLiteDB,
) *App {
	return &App{dependenciesLoader, db}
}

func (app App) Run() {
	// if err := app.db.CreateTables(); err != nil {
	// 	log.Fatalf("failed to create db tables due to an error: %v \n exiting...", err)
	// }

	// if err := app.dependenciesLoader.FetchDepsDevDependencies(); err != nil {
	// 	log.Fatalf("failed to fetch deps.dev dependencies due to an error: %v \n exiting...", err)
	// }

	// detailedDependencies := app.dependenciesLoader.FetchDetailsForAllDependencies()

	// if err := app.db.LoadDependencies(detailedDependencies); err != nil {
	// 	log.Fatalf("failed to load detailed dependencies into db due to an error: %v \n exiting...", err)
	// }

	// err := app.db.UpdateDependencyDetails(dependency)
	// if err != nil {
	// 	log.Print("ERROR: ", err)
	// }

	// s, err := app.db.GetDependencyDetailsByID("github.com/briandowns/spinner")
	// if err != nil {
	// 	log.Print("ERROR: ", err)
	// }
	// log.Print(s)

	// app.db.DeleteDependencyWithDetails("github.com/alecaivazis/survey")

	byscore, err := app.db.GetDependenciesByOverallScore(5)
	if err != nil {
		log.Print("ERROR: ", byscore)
	}
	for _, d := range byscore {
		log.Print(d)
		log.Print("---------------")
	}
}
