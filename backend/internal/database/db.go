package database

import (
	"database/sql"
	"fmt"
	"log"

	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &SQLiteDB{db}, nil
}

func (s *SQLiteDB) CreateTables() error {
	tableStatements := []string{
		`CREATE TABLE IF NOT EXISTS "ProjectKey" (
			id TEXT PRIMARY KEY
		);`,

		`CREATE TABLE IF NOT EXISTS "Documentation" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			shortDescription TEXT,
			url TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS "Check" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			documentationId INTEGER,
			score INTEGER,
			reason TEXT,
			FOREIGN KEY (documentationId) REFERENCES "Documentation"(id)
		);`,

		`CREATE TABLE IF NOT EXISTS "Scorecard" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT,
			repositoryName TEXT,
			repositoryCommit TEXT,
			scorecardVersion TEXT,
			scorecardCommit TEXT,
			overallScore REAL,
			metadata TEXT
		);`,

		`CREATE TABLE IF NOT EXISTS "DependencyDetails" (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			projectKeyId TEXT,
			openIssuesCount INTEGER,
			starsCount INTEGER,
			forksCount INTEGER,
			license TEXT,
			description TEXT,
			homepage TEXT,
			scorecardId INTEGER,
			FOREIGN KEY (projectKeyId) REFERENCES "ProjectKey"(id),
			FOREIGN KEY (scorecardId) REFERENCES "Scorecard"(id)
		);`,
	}

	for _, stmt := range tableStatements {
		_, err := s.db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement: %s \n error: %w", stmt, err)
		}
		log.Printf("executed db statement: %s", stmt)
	}

	return nil
}

func (s *SQLiteDB) LoadDependencies(detailsSlice []dependenciesloader.DependencyDetails) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, details := range detailsSlice {
		_, err := tx.Exec("INSERT OR REPLACE INTO \"ProjectKey\" (id) VALUES (?)", details.ProjectKey.ID)
		if err != nil {
			return err
		}

		result, err := tx.Exec(`INSERT INTO "Scorecard" (date, repositoryName, repositoryCommit, scorecardVersion, scorecardCommit, overallScore, metadata)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			details.Scorecard.Date,
			details.Scorecard.Repository.Name,
			details.Scorecard.Repository.Commit,
			details.Scorecard.Scorecard.Version,
			details.Scorecard.Scorecard.Commit,
			details.Scorecard.OverallScore,
			fmt.Sprintf("%v", details.Scorecard.Metadata))
		if err != nil {
			return err
		}

		scorecardID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, check := range details.Scorecard.Checks {
			documentationResult, err := tx.Exec(`INSERT INTO "Documentation" (shortDescription, url) VALUES (?, ?)`,
				check.Documentation.ShortDescription, check.Documentation.URL)
			if err != nil {
				return err
			}
			documentationId, err := documentationResult.LastInsertId()
			if err != nil {
				return err
			}

			_, err = tx.Exec(`INSERT INTO "Check" (name, documentationId, score, reason) 
				VALUES (?, ?, ?, ?)`, check.Name, documentationId, check.Score, check.Reason)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`INSERT INTO "DependencyDetails" (projectKeyId, openIssuesCount, starsCount, forksCount, license, description, homepage, scorecardId) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			details.ProjectKey.ID,
			details.OpenIssuesCount,
			details.StarsCount,
			details.ForksCount,
			details.License,
			details.Description,
			details.Homepage,
			scorecardID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
