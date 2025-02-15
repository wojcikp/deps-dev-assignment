package database

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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

func (s *SQLiteDB) CloseDbConnection() {
	s.db.Close()
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

		`CREATE TABLE IF NOT EXISTS "VersionKeys" (
			name TEXT PRIMARY KEY,
			system TEXT,
			version TEXT
		);`,
	}

	for _, stmt := range tableStatements {
		_, err := s.db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("error executing statement: %s \n error: %w", stmt, err)
		}
	}

	return nil
}

func (s *SQLiteDB) LoadDependencies(nodes []dependenciesloader.Node) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, node := range nodes {
		_, err := tx.Exec(`INSERT INTO "VersionKeys" (name, system, version) VALUES (?, ?, ?) ON CONFLICT(name) DO NOTHING`,
			node.VersionKey.Name,
			node.VersionKey.System,
			node.VersionKey.Version,
		)
		if err != nil {
			return fmt.Errorf("failed to insert into VersionKeys: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLiteDB) GetVersionKeys() ([]dependenciesloader.VersionKey, error) {
	query := `SELECT name, system, version FROM VersionKeys`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}

	var versionKeys []dependenciesloader.VersionKey
	defer rows.Close()

	for rows.Next() {
		var versionKey dependenciesloader.VersionKey
		err := rows.Scan(
			&versionKey.Name,
			&versionKey.System,
			&versionKey.Version,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DependencyDetails: %w", err)
		}
		versionKeys = append(versionKeys, versionKey)
	}

	return versionKeys, nil
}

func (s *SQLiteDB) LoadDetailedDependencies(dependenciesDetails []dependenciesloader.DependencyDetails) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	for _, details := range dependenciesDetails {
		_, err := tx.Exec(`INSERT INTO "ProjectKey" (id) VALUES (?) ON CONFLICT(id) DO NOTHING`, details.ProjectKey.ID)
		if err != nil {
			return fmt.Errorf("failed to insert into ProjectKey: %w", err)
		}

		scorecardResult, err := tx.Exec(`
			INSERT INTO "Scorecard" (date, repositoryName, repositoryCommit, scorecardVersion, scorecardCommit, overallScore, metadata) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
			details.Scorecard.Date,
			details.Scorecard.Repository.Name,
			details.Scorecard.Repository.Commit,
			details.Scorecard.Scorecard.Version,
			details.Scorecard.Scorecard.Commit,
			details.Scorecard.OverallScore,
			fmt.Sprintf("%v", details.Scorecard.Metadata),
		)
		if err != nil {
			return fmt.Errorf("failed to insert into Scorecard: %w", err)
		}

		scorecardId, err := scorecardResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id for Scorecard: %w", err)
		}

		for _, check := range details.Scorecard.Checks {
			docResult, err := tx.Exec(`INSERT INTO "Documentation" (shortDescription, url) VALUES (?, ?)`,
				check.Documentation.ShortDescription, check.Documentation.URL)
			if err != nil {
				return fmt.Errorf("failed to insert into Documentation: %w", err)
			}

			docId, err := docResult.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert id for Documentation: %w", err)
			}

			_, err = tx.Exec(`
                INSERT INTO "Check" (name, documentationId, score, reason, scorecardId) 
                VALUES (?, ?, ?, ?, ?)`,
				check.Name,
				docId,
				check.Score,
				check.Reason,
				scorecardId,
			)
			if err != nil {
				return fmt.Errorf("failed to insert into Check: %w", err)
			}
		}

		_, err = tx.Exec(`
		INSERT INTO "DependencyDetails" (projectKeyId, openIssuesCount, starsCount, forksCount, license, description, homepage, scorecardId) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			details.ProjectKey.ID,
			details.OpenIssuesCount,
			details.StarsCount,
			details.ForksCount,
			details.License,
			details.Description,
			details.Homepage,
			scorecardId)
		if err != nil {
			return fmt.Errorf("failed to insert into DependencyDetails: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLiteDB) AddNewDependencyDetails(details dependenciesloader.DependencyDetails) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO "ProjectKey" (id) VALUES (?) ON CONFLICT(id) DO NOTHING`, details.ProjectKey.ID)
	if err != nil {
		return fmt.Errorf("failed to insert into ProjectKey: %w", err)
	}

	scorecardResult, err := tx.Exec(`
			INSERT INTO "Scorecard" (date, repositoryName, repositoryCommit, scorecardVersion, scorecardCommit, overallScore, metadata) 
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		details.Scorecard.Date,
		details.Scorecard.Repository.Name,
		details.Scorecard.Repository.Commit,
		details.Scorecard.Scorecard.Version,
		details.Scorecard.Scorecard.Commit,
		details.Scorecard.OverallScore,
		fmt.Sprintf("%v", details.Scorecard.Metadata),
	)
	if err != nil {
		return fmt.Errorf("failed to insert into Scorecard: %w", err)
	}

	scorecardId, err := scorecardResult.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id for Scorecard: %w", err)
	}

	for _, check := range details.Scorecard.Checks {
		docResult, err := tx.Exec(`INSERT INTO "Documentation" (shortDescription, url) VALUES (?, ?)`,
			check.Documentation.ShortDescription, check.Documentation.URL)
		if err != nil {
			return fmt.Errorf("failed to insert into Documentation: %w", err)
		}

		docId, err := docResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert id for Documentation: %w", err)
		}

		_, err = tx.Exec(`
                INSERT INTO "Check" (name, documentationId, score, reason, scorecardId) 
                VALUES (?, ?, ?, ?, ?)`,
			check.Name,
			docId,
			check.Score,
			check.Reason,
			scorecardId,
		)
		if err != nil {
			return fmt.Errorf("failed to insert into Check: %w", err)
		}
	}

	_, err = tx.Exec(`
		INSERT INTO "DependencyDetails" (projectKeyId, openIssuesCount, starsCount, forksCount, license, description, homepage, scorecardId) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		details.ProjectKey.ID,
		details.OpenIssuesCount,
		details.StarsCount,
		details.ForksCount,
		details.License,
		details.Description,
		details.Homepage,
		scorecardId)
	if err != nil {
		return fmt.Errorf("failed to insert into DependencyDetails: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *SQLiteDB) UpdateVersionKeys(name, version string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	_, err = tx.Exec(`
		UPDATE "VersionKeys" 
		SET version = ?
		WHERE name = ?
	`, version, name)
	if err != nil {
		return fmt.Errorf("failed to update DependencyDetails: %w", err)
	}

	return nil
}

func (s *SQLiteDB) UpdateDependencyDetails(newDetails dependenciesloader.DependencyDetails) error {
	projectKeyID := newDetails.ProjectKey.ID

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var scorecardID int
	query := `
		SELECT s.id 
		FROM DependencyDetails dd
		JOIN Scorecard s ON dd.scorecardId = s.id
		JOIN ProjectKey pk ON dd.projectKeyId = pk.id
		WHERE pk.id = ?
	`
	err = tx.QueryRow(query, projectKeyID).Scan(&scorecardID)
	if err != nil {
		return fmt.Errorf("failed to get related scorecard ID: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE "DependencyDetails" 
		SET openIssuesCount = ?, starsCount = ?, forksCount = ?, license = ?, description = ?, homepage = ?
		WHERE projectKeyId = ?
	`, newDetails.OpenIssuesCount, newDetails.StarsCount, newDetails.ForksCount, newDetails.License, newDetails.Description, newDetails.Homepage, projectKeyID)
	if err != nil {
		return fmt.Errorf("failed to update DependencyDetails: %w", err)
	}

	metadataJSON, _ := json.Marshal(newDetails.Scorecard.Metadata)
	_, err = tx.Exec(`
		UPDATE "Scorecard" 
		SET date = ?, scorecardVersion = ?, scorecardCommit = ?, overallScore = ?, metadata = ?
		WHERE id = ?
	`, newDetails.Scorecard.Date, newDetails.Scorecard.Scorecard.Version, newDetails.Scorecard.Scorecard.Commit, newDetails.Scorecard.OverallScore, string(metadataJSON), scorecardID)
	if err != nil {
		return fmt.Errorf("failed to update Scorecard: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM "Check"
		WHERE scorecardId = (SELECT scorecardId FROM DependencyDetails WHERE projectKeyId = ?)
	`, projectKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete Check for projectKeyID %s: %w", projectKeyID, err)
	}

	_, err = tx.Exec(`
		DELETE FROM "Documentation"
		WHERE id IN (
			SELECT documentationId 
			FROM "Check"
			WHERE scorecardId = (SELECT scorecardId FROM DependencyDetails WHERE projectKeyId = ?)
		)
	`, projectKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete Documentation for projectKeyID %s: %w", projectKeyID, err)
	}

	for _, check := range newDetails.Scorecard.Checks {
		_, err := tx.Exec(`
			INSERT OR IGNORE INTO "Documentation" (shortDescription, url)
			VALUES (?, ?)
		`, check.Documentation.ShortDescription, check.Documentation.URL)
		if err != nil {
			return fmt.Errorf("failed to insert Documentation for check %s: %w", check.Name, err)
		}

		_, err = tx.Exec(`
			UPDATE "Documentation"
			SET shortDescription = ?, url = ?
			WHERE shortDescription = ? AND url = ?
		`, check.Documentation.ShortDescription, check.Documentation.URL, check.Documentation.ShortDescription, check.Documentation.URL)
		if err != nil {
			return fmt.Errorf("failed to update Documentation for check %s: %w", check.Name, err)
		}

		var documentationID int
		err = tx.QueryRow(`
			SELECT id FROM "Documentation"
			WHERE shortDescription = ? AND url = ?
		`, check.Documentation.ShortDescription, check.Documentation.URL).Scan(&documentationID)
		if err != nil {
			return fmt.Errorf("failed to get Documentation ID for check %s: %w", check.Name, err)
		}

		_, err = tx.Exec(`
			INSERT OR IGNORE INTO "Check" (name, score, reason, documentationId, scorecardId)
			VALUES (?, ?, ?, ?, ?)
		`, check.Name, check.Score, check.Reason, documentationID, scorecardID)
		if err != nil {
			return fmt.Errorf("failed to insert Check for %s: %w", check.Name, err)
		}

		_, err = tx.Exec(`
			UPDATE "Check"
			SET score = ?, reason = ?, documentationId = ?
			WHERE name = ? AND scorecardId = ?
		`, check.Score, check.Reason, documentationID, check.Name, scorecardID)
		if err != nil {
			return fmt.Errorf("failed to update Check for %s: %w", check.Name, err)
		}
	}

	return nil
}

func (s *SQLiteDB) GetDependencyDetailsByID(projectKeyID string) (*dependenciesloader.DependencyDetails, error) {
	var detail dependenciesloader.DependencyDetails

	query := `SELECT pk.id, dd.openIssuesCount, dd.starsCount, dd.forksCount, dd.license,
                     dd.description, dd.homepage, sc.date, sc.repositoryName, sc.repositoryCommit,
                     sc.scorecardVersion, sc.scorecardCommit, sc.overallScore, sc.metadata
              FROM DependencyDetails dd
              JOIN ProjectKey pk ON dd.projectKeyId = pk.id
              JOIN Scorecard sc ON dd.scorecardId = sc.id
              WHERE pk.id = ?`

	var metadataStr string

	err := s.db.QueryRow(query, projectKeyID).Scan(
		&detail.ProjectKey.ID,
		&detail.OpenIssuesCount,
		&detail.StarsCount,
		&detail.ForksCount,
		&detail.License,
		&detail.Description,
		&detail.Homepage,
		&detail.Scorecard.Date,
		&detail.Scorecard.Repository.Name,
		&detail.Scorecard.Repository.Commit,
		&detail.Scorecard.Scorecard.Version,
		&detail.Scorecard.Scorecard.Commit,
		&detail.Scorecard.OverallScore,
		&metadataStr,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get DependencyDetails: %w", err)
	}

	if err := json.Unmarshal([]byte(metadataStr), &detail.Scorecard.Metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	checkQuery := `SELECT c.name, c.score, c.reason, d.shortDescription, d.url
                   FROM "Check" c
                   JOIN Documentation d ON c.documentationId = d.id
                   WHERE c.scorecardId = (SELECT sc.id FROM Scorecard sc
                                          JOIN DependencyDetails dd ON dd.scorecardId = sc.id
                                          JOIN ProjectKey pk ON dd.projectKeyId = pk.id
                                          WHERE pk.id = ?)`

	rows, err := s.db.Query(checkQuery, projectKeyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Checks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var check dependenciesloader.Check
		err := rows.Scan(
			&check.Name,
			&check.Score,
			&check.Reason,
			&check.Documentation.ShortDescription,
			&check.Documentation.URL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan Check: %w", err)
		}
		detail.Scorecard.Checks = append(detail.Scorecard.Checks, check)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over Checks: %w", err)
	}

	return &detail, nil
}

func (s *SQLiteDB) GetAllDependencies() ([]dependenciesloader.DependencyDetails, error) {
	query := `SELECT id FROM ProjectKey`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies: %w", err)
	}

	var dependencies []dependenciesloader.DependencyDetails
	defer rows.Close()

	for rows.Next() {
		var projectKeyId string
		err := rows.Scan(
			&projectKeyId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DependencyDetails: %w", err)
		}
		dependency, err := s.GetDependencyDetailsByID(projectKeyId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan GetDependencyDetailsByID: %w", err)
		}
		dependencies = append(dependencies, *dependency)
	}

	return dependencies, nil
}

func (s *SQLiteDB) GetDependenciesByOverallScore(dependencyScore float64) ([]dependenciesloader.DependencyDetails, error) {
	query := `
        SELECT dd.projectKeyId 
        FROM DependencyDetails dd
        JOIN Scorecard sc ON dd.scorecardId = sc.id
        WHERE sc.overallScore BETWEEN ? AND ?
    `

	rows, err := s.db.Query(query, dependencyScore, dependencyScore+0.99)
	if err != nil {
		return nil, fmt.Errorf("failed to query dependencies by overallScore: %w", err)
	}

	var dependencies []dependenciesloader.DependencyDetails
	defer rows.Close()

	for rows.Next() {
		var projectKeyId string
		err := rows.Scan(
			&projectKeyId,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DependencyDetails: %w", err)
		}
		dependency, err := s.GetDependencyDetailsByID(projectKeyId)
		if err != nil {
			return nil, fmt.Errorf("failed to scan GetDependencyDetailsByID: %w", err)
		}
		dependencies = append(dependencies, *dependency)
	}

	return dependencies, nil
}

func (s *SQLiteDB) DeleteDependencyWithDetails(projectKeyID string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var scorecardID, dependencyDetailsID int
	query := `
        SELECT dd.id, s.id
        FROM DependencyDetails dd
        JOIN Scorecard s ON dd.scorecardId = s.id
        JOIN ProjectKey pk ON dd.projectKeyId = pk.id
        WHERE pk.id = ?
    `
	err = tx.QueryRow(query, projectKeyID).Scan(&dependencyDetailsID, &scorecardID)
	if err != nil {
		return fmt.Errorf("failed to retrieve related IDs: %w", err)
	}

	_, err = tx.Exec(`
        DELETE FROM "Check"
        WHERE scorecardId = ?
    `, scorecardID)
	if err != nil {
		return fmt.Errorf("failed to delete Checks: %w", err)
	}

	_, err = tx.Exec(`
        DELETE FROM "Documentation"
        WHERE id IN (
            SELECT DISTINCT documentationId 
            FROM "Check" 
            WHERE scorecardId = ?
        )
    `, scorecardID)
	if err != nil {
		return fmt.Errorf("failed to delete Documentation: %w", err)
	}

	_, err = tx.Exec(`
        DELETE FROM "Scorecard"
        WHERE id = ?
    `, scorecardID)
	if err != nil {
		return fmt.Errorf("failed to delete Scorecard: %w", err)
	}

	_, err = tx.Exec(`
        DELETE FROM "DependencyDetails"
        WHERE id = ?
    `, dependencyDetailsID)
	if err != nil {
		return fmt.Errorf("failed to delete DependencyDetails: %w", err)
	}

	_, err = tx.Exec(`
        DELETE FROM "ProjectKey"
        WHERE id = ?
    `, projectKeyID)
	if err != nil {
		return fmt.Errorf("failed to delete ProjectKey: %w", err)
	}

	return nil
}
