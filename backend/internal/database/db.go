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
			scorecardId INTEGER,
			FOREIGN KEY (documentationId) REFERENCES "Documentation"(id),
			FOREIGN KEY (scorecardId) REFERENCES "Scorecard"(id)
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

func (s *SQLiteDB) LoadDependencies(dependencyDetails []dependenciesloader.DependencyDetails) error {

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, detail := range dependencyDetails {
		_, err := tx.Exec(`INSERT INTO "ProjectKey" (id) VALUES (?) ON CONFLICT(id) DO NOTHING`, detail.ProjectKey.ID)
		if err != nil {
			return fmt.Errorf("failed to insert into ProjectKey: %v", err)
		}

		scorecardResult, err := tx.Exec(`
			INSERT INTO "Scorecard" (date, repositoryName, repositoryCommit, scorecardVersion, scorecardCommit, overallScore, metadata) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			detail.Scorecard.Date,
			detail.Scorecard.Repository.Name,
			detail.Scorecard.Repository.Commit,
			detail.Scorecard.Scorecard.Version, detail.Scorecard.Scorecard.Commit,
			detail.Scorecard.OverallScore,
			fmt.Sprintf("%v", detail.Scorecard.Metadata),
		)
		if err != nil {
			return fmt.Errorf("failed to insert into Scorecard: %v", err)
		}

		scorecardId, err := scorecardResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id for Scorecard: %v", err)
		}

		for _, check := range detail.Scorecard.Checks {
			docResult, err := tx.Exec(`INSERT INTO "Documentation" (shortDescription, url) VALUES (?, ?)`,
				check.Documentation.ShortDescription, check.Documentation.URL)
			if err != nil {
				return fmt.Errorf("failed to insert into Documentation: %v", err)
			}

			docId, err := docResult.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert id for Documentation: %v", err)
			}

			_, err = tx.Exec(`
                INSERT INTO "Check" (name, documentationId, score, reason, scorecardId) 
                VALUES (?, ?, ?, ?, ?)
            `, check.Name, docId, check.Score, check.Reason, scorecardId)
			if err != nil {
				return fmt.Errorf("failed to insert into Check: %v", err)
			}
		}

		_, err = tx.Exec(`
		INSERT INTO "DependencyDetails" (projectKeyId, openIssuesCount, starsCount, forksCount, license, description, homepage, scorecardId) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			detail.ProjectKey.ID,
			detail.OpenIssuesCount,
			detail.StarsCount,
			detail.ForksCount,
			detail.License,
			detail.Description,
			detail.Homepage,
			scorecardId)
		if err != nil {
			return fmt.Errorf("failed to insert into DependencyDetails: %v", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
