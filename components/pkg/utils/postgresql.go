package utils

/**
* @Author: azh
* @Date: 2022/5/13 22:20
* @Context:
 */

import (
	"database/sql"
	"fmt"
	//_ "github.com/bmizerany/pq"
	"log"
	"mosn.io/layotto/components/sequencer/postgresql/config"
	dao2 "mosn.io/layotto/components/sequencer/postgresql/dao"
	"mosn.io/layotto/components/sequencer/postgresql/model"
	"mosn.io/layotto/components/sequencer/postgresql/service"
	"sync"
)

type PostgrepsqlMetaData struct {
	Host     string
	Username string
	Password string
	Port     int
}

const (
	ipHost    = "localhost"
	port      = 5432
	user      = "postgres"
	passwords = "213213"
	dbname    = "test_db"
)

var postDB *sql.DB
var err error

type product struct {
	ProductNo string
	Name      string
	Price     float64
}

func DBopen() error {
	info := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		ipHost, port, user, passwords, dbname)
	fmt.Println(info)
	postDB, err = sql.Open("postgres", info)
	if err != nil {
		return err
	}
	err = postDB.Ping()
	if err != nil {
		fmt.Println("conected faild")
		return err
	}
	fmt.Println("success connected")
	return nil
}

type PostgresqlServer struct {
	conf *config.Server
	once sync.Once
}

func (s *PostgresqlServer) NewPostgresqlServer() *service.PostgresqlService {
	db := NewPostgresqlClient(s.conf)
	postgresqlDB := dao2.NewPostgresqlDB(db)
	postgresqlDAO := dao2.NewPostgresqlDAO(postgresqlDB, db)
	seq := model.NewPostgresqlSeq()
	postgresqlService := service.NewPostgresqlService(postgresqlDAO, seq)
	return postgresqlService
}

func InitPostgresql() *service.PostgresqlService {
	s := &PostgresqlServer{}
	s.init()
	return s.NewPostgresqlServer()
}

func NewPostgresqlClient(conf *config.Server) *sql.DB {
	info := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgresql.Host, conf.Postgresql.Port, conf.Postgresql.Username, conf.Postgresql.Password, conf.Postgresql.Db)
	fmt.Println(info)
	postDB, err = sql.Open("postgres", info)
	if err != nil {
		return nil
	}
	err = postDB.Ping()
	if err != nil {
		fmt.Println("conected faild")
		return nil
	}
	fmt.Println("success connected")
	return postDB
}

func (s *PostgresqlServer) init() {
	s.once.Do(func() {
		conf := &config.Server{}
		err := conf.Load("D:\\goTest\\test\\layotto\\components\\sequencer\\postgresql\\conf\\postgresql.yaml")
		if err != nil {
			log.Panic("load conf file failed", err)
		}
		s.conf = conf
	})
}
