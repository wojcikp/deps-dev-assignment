package database

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

func GetTestDatabase(t *testing.T) *SQLiteDB {
	dbPath := getDbPath(t)

	db, err := NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatal("failed to create database:", err)
	}

	db.CreateTables()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("database file was not created")
	}

	return db
}

func TestLoadDependencies(t *testing.T) {
	var dependencies dependenciesloader.Dependencies

	cwd, _ := os.Getwd()
	mockFile := path.Join(cwd, "test_data", "dependencies_mock.json")

	data, err := os.ReadFile(mockFile)
	if err != nil {
		t.Fatal("failed to read mock data:", err)
	}

	if err := json.Unmarshal(data, &dependencies); err != nil {
		t.Fatal("failed to parse mock data:", err)
	}

	db := GetTestDatabase(t)

	if err := db.LoadDependencies(dependencies.Nodes); err != nil {
		t.Fatalf("failed to load version keys into test db due to an error: %v", err)
	}
}

func TestLoadDetailedDependencies(t *testing.T) {

	db := GetTestDatabase(t)

	if err := db.LoadDetailedDependencies(getDetailedDependenciesMock(t, "dependencies_details_mock.json")[:5]); err != nil {
		t.Fatalf("failed to load version keys into test db due to an error: %v", err)
	}
}

func TestGetVersionKeys(t *testing.T) {
	db := GetTestDatabase(t)
	keys, err := db.GetVersionKeys()
	if err != nil {
		t.Fatalf("failed to retrieve version keys from test db due to an error: %v", err)
	}
	const want = 5
	got := len(keys)
	if got != want {
		t.Fatalf("version keys not found in test db, want: %d, got: %d", want, got)
	}
}

func TestGetAllDependencies(t *testing.T) {
	db := GetTestDatabase(t)
	checkAllDependencies(t, db, 5)
}

func TestAddNewDependencyDetails(t *testing.T) {
	db := GetTestDatabase(t)

	detailedDependencies := getDetailedDependenciesMock(t, "dependencies_details_mock.json")

	if err := db.AddNewDependencyDetails(detailedDependencies[5]); err != nil {
		t.Fatal("failed to add new dependency details:", err)
	}

	checkAllDependencies(t, db, 6)
}

func TestGetDependencyDetailsByID(t *testing.T) {
	db := GetTestDatabase(t)

	got, err := db.GetDependencyDetailsByID("github.com/cli/cli")
	if err != nil {
		t.Fatal("failed to get dependency details by id:", err)
	}

	want := getDetailedDependenciesMock(t, "dependencies_details_mock.json")[0]

	if !cmp.Equal(*got, want) {
		t.Fatal("dependency details from test db are not equal to mocks: ", cmp.Diff(*got, want))
	}
}

func TestUpdateDependencyDetails(t *testing.T) {
	db := GetTestDatabase(t)

	want := getDetailedDependenciesMock(t, "dependencies_details_update_mock.json")[0]

	if err := db.UpdateDependencyDetails(want); err != nil {
		t.Fatal("failed to update dependency details:", err)
	}

	got, err := db.GetDependencyDetailsByID("github.com/charmbracelet/glamour")
	if err != nil {
		t.Fatal("failed to get dependency details by id:", err)
	}

	if got.Scorecard.Date != want.Scorecard.Date ||
		len(got.Scorecard.Checks) != len(want.Scorecard.Checks) ||
		got.Scorecard.OverallScore != want.Scorecard.OverallScore {
		t.Fatalf("dependencies details are not equal after update:\n dates: %s --- %s\nchecks len: %d --- %d\noverall scores: %f --- %f",
			got.Scorecard.Date, want.Scorecard.Date,
			len(got.Scorecard.Checks), len(want.Scorecard.Checks),
			got.Scorecard.OverallScore, want.Scorecard.OverallScore,
		)
	}

}

func TestGetDependenciesByOverallScore(t *testing.T) {
	db := GetTestDatabase(t)

	got, err := db.GetDependenciesByOverallScore(4)
	if err != nil {
		t.Fatal("failed to get dependency details by overall score:", err)
	}
	const want = 2

	if len(got) != want {
		t.Fatalf("got != want, want: %d, got: %d", want, len(got))
	}
}

func TestDeleteDependencyWithDetails(t *testing.T) {
	db := GetTestDatabase(t)

	if err := db.DeleteDependencyWithDetails("github.com/briandowns/spinner"); err != nil {
		t.Fatal("failed to delete dependency details:", err)
	}

	const want = 5

	got, err := db.GetAllDependencies()
	if err != nil {
		t.Fatalf("failed to retrieve dependencies from test db: %v", err)
	}

	if len(got) != want {
		t.Fatalf("got != want, want: %d, got: %d", want, len(got))
	}
}

func TestCleanupTestDatabase(t *testing.T) {
	p := getDbPath(t)
	if err := os.Remove(p); err != nil {
		t.Fatal("failed to remove test db file after tests")
	}
}

func checkAllDependencies(t *testing.T, db *SQLiteDB, want int) {
	dependencies, err := db.GetAllDependencies()
	if err != nil {
		t.Fatalf("failed to retrieve dependencies from test db: %v", err)
	}

	got := len(dependencies)
	if got != want {
		t.Fatalf("unexpected number of dependencies, want: %d, got: %d", want, got)
	}
}

func getDetailedDependenciesMock(t *testing.T, filename string) []dependenciesloader.DependencyDetails {
	var detailedDependencies struct {
		Dependencies []dependenciesloader.DependencyDetails `json:"dependencies"`
	}

	cwd, _ := os.Getwd()
	mockFile := path.Join(cwd, "test_data", filename)

	data, err := os.ReadFile(mockFile)
	if err != nil {
		t.Fatal("failed to read mock data:", err)
	}

	if err := json.Unmarshal(data, &detailedDependencies); err != nil {
		t.Fatal("failed to parse mock data:", err)
	}

	return detailedDependencies.Dependencies
}

func getDbPath(t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	return path.Join(cwd, "test.db")
}
