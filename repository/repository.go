package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type IRepository interface {
	GetById(id int64) (interface{}, error)
	GetAll(instanceType interface{}) (interface{}, error)
	Update(query string, args ...interface{}) (result sql.Result, err error)
	Delete(query string, args ...interface{}) (result sql.Result, err error)
	Insert(query string, args ...interface{}) (result sql.Result, err error)
}

type BaseRepository struct {
	dbClient *DBClient
}

type ColorRepository struct {
	BaseRepository
}

type DBClient struct {
	pool *sqlx.DB
}

func NewColorRepository(dbClient *DBClient) ColorRepository {
	return ColorRepository{BaseRepository: BaseRepository{dbClient: dbClient}}
}

type DBClientFactory struct {
	driverName     string
	dataSourceName string
}

func NewDbClientFactory(driverName string, dataSourceName string) DBClientFactory {
	return DBClientFactory{driverName: driverName, dataSourceName: dataSourceName}
}

func (dbCLientFactory DBClientFactory) NewDBClient() *DBClient {
	client := &DBClient{}

	pool, err := Connect(dbCLientFactory.driverName, dbCLientFactory.dataSourceName)

	if err != nil {
		panic(err)
	}

	client.pool = pool

	return client
}
