# deps dev app

#### Running the app:
1. Clone this repository to your local machine
2. Make sure to have docker compose installed on your computer
3. Make sure to have **:3000** and **:8080** ports on your machine **available**
4. In the main app directory run command "docker-compose up --build"

When the docker build process is ready backend of the app will be available at **localhost:3000**, frontend will be available at localhost:8080. 
To see the application interface go to **localhost:8080** address in your web browser.

#### Available endpoints:
1. "/dependency", Methods("GET"), example: `curl -X GET "http://localhost:3000/dependency?id=github.com/briandowns/spinner"`
2. "/dependency/score/{score}", Methods("GET"), example: `curl -X GET "http://localhost:3000/dependency/score/4"`
3. "/dependency/all", Methods("GET"), example: `curl -X GET "http://localhost:3000/dependency/all"`
4. "/dependency/update", Methods("GET"), example: `curl -X GET "http://localhost:3000/dependency/update"`
5. "/dependency", Methods("DELETE"), example: `curl -X DELETE "http://localhost:3000/dependency?id=github.com/briandowns/spinner"`
6. "/dependency", Methods("POST"), example: 
```
curl --location 'http://localhost:3000/dependency' \
--header 'Content-Type: application/json' \
--data ''
```
7. "/dependency", Methods("PUT"), example:
```
curl --location --request PUT 'http://localhost:3000/dependency' \
--header 'Content-Type: application/json' \
--data ''
```
**NOTE**: In data field provide a valid json structured like response from deps.dev api, for example result of: `curl -s 'https://api.deps.dev/v3/projects/github.com%2Fcharmbracelet%2Fglamour'`

#### SQLite database schema:
```
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
```
