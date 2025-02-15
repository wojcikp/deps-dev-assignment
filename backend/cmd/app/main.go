package main

import (
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wojcikp/deps-dev-assignment/backend/internal/api"
	"github.com/wojcikp/deps-dev-assignment/backend/internal/app"
	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
	dependenciesupdater "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_updater"
)

const repositoryApiUrl = "https://api.deps.dev/v3/systems/GO/packages/github.com%2Fcli%2Fcli/versions/v1.14.0:dependencies"

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	p := path.Join(cwd, "data", "production.db")
	db, err := database.NewSQLiteDB(p)
	if err != nil {
		log.Fatal("failed to establish database connection, exiting...")
	}

	dependenciesLoader := dependenciesloader.NewDependenciesLoader(repositoryApiUrl)
	dependenciesUpdater := dependenciesupdater.NewDependenciesUpdater(dependenciesLoader, db)
	api := api.NewApi(db, dependenciesUpdater)
	app := app.NewApp(dependenciesLoader, db, api)

	app.Run()
}
