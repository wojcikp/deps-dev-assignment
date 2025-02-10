package app

import (
	"log"

	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type App struct {
	dependenciesLoader *dependenciesloader.Loader
}

func NewApp(
	dependenciesLoader *dependenciesloader.Loader,
) *App {
	return &App{dependenciesLoader}
}

func (app App) Run() {
	app.dependenciesLoader.FetchDepsDevDependencies()
	log.Print(app.dependenciesLoader.Dependencies)
}
