package model

import (
	"encoding/json"
	"fmt"
	l "log"
	"os"
	"strings"
	"time"

	"github.com/deepaksinghvi/cdule/pkg"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

var cduleConfig *pkg.CduleConfig

var CduleRepos *Repositories

type Repositories struct {
	CduleRepository CduleRepository
	DB              *gorm.DB
}

func ConnectDataBase(param []string) {
	cduleConfig, err := ReadConfig(param)
	if nil != err {
		log.Error(err)
		panic("Failed to read config!")
	}
	printConfig(cduleConfig)
	var db *gorm.DB
	if cduleConfig.Cduletype == string(pkg.DATABASE) {
		if strings.Contains(cduleConfig.Dburl, "postgres") {
			db = postgresConn(cduleConfig.Dburl)
		}
	} else if cduleConfig.Cduletype == string(pkg.MEMORY) {
		db = sqliteConn(cduleConfig.Dburl)
	}
	// Set LogLevel to `logger.Silent` to stop logging sqls
	sqlLogger := logger.New(
		l.New(os.Stdout, "\r\n", l.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	db.Logger = sqlLogger
	Migrate(db)
	DB = db

	// Initialise CduleRepositories
	CduleRepos = &Repositories{
		CduleRepository: NewCduleRepository(db),
		DB:              db,
	}
}

func postgresConn(dbDSN string) (db *gorm.DB) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbDSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Errorf("Error Connecting Postgressql %s, %s", dbDSN, err.Error())
		panic("Failed to connect to database! " + dbDSN)
	}
	return db

}

func sqliteConn(dbDSN string) (db *gorm.DB) {
	//db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	//db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})

	// If you would use file based as mentioned above db
	db, err := gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Error(dbDSN)
		panic("Failed to connect to database! " + dbDSN)
	}
	return db
}
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Job{})
	db.AutoMigrate(&JobHistory{})
	db.AutoMigrate(&Schedule{})
	db.AutoMigrate(&Worker{})
}

func printConfig(config *pkg.CduleConfig) {
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Configuration %s\n", string(configJSON))
}
