package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/deepaksinghvi/cdule/pkg"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_ConnectDatabase(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping PostgreSQL test in -short mode")
	}
	param := []string{"../../resources", "config_in_memory", "Info"} // default path for resource
	cduleConfig, err := ConnectDataBase(param)
	require.NoError(t, err)
	require.NotEqual(t, pkg.EMPTYSTRING, cduleConfig.Dburl)
	_ = os.Remove("./sqlite.db")
}

func Test_ConnectDatabaseFailedToReadConfig(t *testing.T) {
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			log.Warning("Recovered in Test_ConnectPostgresDBPanic ", r)
			recovered = true
		}
	}()
	param := []string{"./resources", "config_in_memory", "Info"} // default path for resource
	_, err := ConnectDataBase(param)
	require.Error(t, err)
	require.EqualValues(t, true, recovered)
}

func Test_ConnectPostgresDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping PostgreSQL test in -short mode")
	}
	db := postgresConn("postgres://cduleuser:cdulepassword@localhost:5432/cdule?sslmode=disable")
	require.NotNil(t, db)
}

func Test_ConnectPostgresDBPanic(t *testing.T) {
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			log.Warning("Recovered in Test_ConnectPostgresDBPanic ", r)
			recovered = true
		}
	}()
	db := postgresConn("postgres://abc:abc@localhost:5432/cdule?sslmode=disable")
	require.Nil(t, db)
	require.EqualValues(t, true, recovered)
}

func Test_ConnectSqlite(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Remove(dirname + "/sqlite.db")

	db := sqliteConn(dirname + "/sqlite.db")
	require.NotNil(t, db)
}

// Test_ConnectSqliteDBPanic tests the panic when the database file is not found
func Test_ConnectSqliteDBPanic(t *testing.T) {
	recovered := false
	defer func() {
		if r := recover(); r != nil {
			log.Warning("Recovered in Test_ConnectSqliteDBPanic ", r)
			recovered = true
		}
	}()
	// Use a guaranteed invalid path (non-existent nested directory)
	invalidPath := filepath.Join(os.TempDir(), "cdule-test-nonexistent", "nested", "db.sqlite")
	db := sqliteConn(invalidPath)
	require.Nil(t, db)
	require.EqualValues(t, true, recovered)
}
