package store

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tehcyx/simple-db-api/pkg/logging"
)

// Postgres simply holds data for the runtime of the application
type Postgres struct {
	user string
	pass string
	host string
	port string
	db   string

	Handle *gorm.DB
}

// NewPostgresStore returns an instance of an postgres store
func NewPostgresStore(log *logrus.Logger, user, pass, host, port, db string) *Postgres {
	info := &Postgres{user: user, pass: pass, host: host, port: port, db: db}

	maxReconnects := 3
	var err error
	var dbHandle *gorm.DB
	for try := 0; try < maxReconnects; try++ {
		dbHandle, err = gorm.Open("postgres", info.psqlInfo())
		if err != nil && try <= maxReconnects {
			backoff := rand.Intn(30) + 1
			log.Error(fmt.Errorf("Could not connect to database: %w, retrying in %d seconds", err, backoff))
			time.Sleep(time.Second * time.Duration(backoff))
		}
	}
	if err != nil {
		log.Fatal(fmt.Errorf("Tried %d times but couldn't connect to database. Reason: %w", maxReconnects, err))
	}
	info.Handle = dbHandle

	dbHandle = dbHandle.Debug()
	// dbHandle.AutoMigrate(&StorageData{}) // Create a table for commerce orders

	dbHandle.DropTableIfExists(&StorageData{}) // drops table and old data
	dbHandle.CreateTable(&StorageData{})       // recreates tables

	return info
}

func (p *Postgres) psqlInfo() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", p.host, p.port, p.user, p.pass, p.db)
}

// Write writes the storage object to the postgres store
func (p *Postgres) Write(ctx context.Context, data StorageData) error {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	log.Debugf("writing: %+v", data)

	createErrs := p.Handle.Create(&data).GetErrors()
	for _, e := range createErrs {
		err := fmt.Errorf("Couldn't create record: %w", e)
		log.Error(err)
		return err
	}
	log.Info("Succesfully persisted in DB")
	log.Debugf("wrote: %+v", data)

	return nil
}

// ReadAll returns all data stored postgres
func (p *Postgres) ReadAll(ctx context.Context) ([]StorageData, error) {
	log := ctx.Value(logging.CtxKeyLog{}).(logrus.FieldLogger)
	var res []StorageData
	err := p.Handle.Find(&res).Error
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve data from database: %w", err)
	}
	log.Debugf("reading: %+v", res)
	return res, nil
}
