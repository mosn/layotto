package utils

/**
* @Author: azh
* @Date: 2022/5/13 22:20
* @Context:
 */

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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
	postgresqlHost    = "host"
	postgresqlPort      = "port"
	postgresqlUsername      = "username"
	postgresqlPassword = "password"
	postgresqlDBname    = "db"
)

var postDB *sql.DB
var err error

type product struct {
	ProductNo string
	Name      string
	Price     float64
}

//func DBopen() error {
//	info := fmt.Sprintf("host=%s port=%d user=%s "+
//		"password=%s dbname=%s sslmode=disable",
//		postgresqlHost, postgresqlPort, postgresqlUsername, postgresqlPassword, postgresqlDBname)
//	postDB, err = sql.Open("postgres", info)
//	if err != nil {
//		return err
//	}
//	err = postDB.Ping()
//	if err != nil {
//		fmt.Println("conected faild")
//		return err
//	}
//	fmt.Println("success connected")
//	return nil
//}

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

func InitPostgresql(proterties map[string]string) (*service.PostgresqlService, error) {
	s := &PostgresqlServer{}
	err := s.init(proterties)
	return s.NewPostgresqlServer(), err
}

func NewPostgresqlClient(conf *config.Server) *sql.DB {
	info := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Postgresql.Host, conf.Postgresql.Port, conf.Postgresql.Username, conf.Postgresql.Password, conf.Postgresql.Db)
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

func (s *PostgresqlServer) init(proterties map[string]string)  error {
		conf := &config.Server{}

		if val, ok := proterties[postgresqlHost]; ok && val != "" {
			conf.Postgresql.Host = val
		}else {
			return  errors.New("posgresql store error: missing host address")
		}

		if val, ok := proterties[postgresqlPort]; ok && val != "0" {
			conf.Postgresql.Port, err = strconv.ParseInt(val, 10, 64)
		}else {
			return err
		}

		if val, ok := proterties[postgresqlUsername]; ok {
			conf.Postgresql.Username = val
		}else {
			return errors.New("posgresql store error: username error")
		}

		if val, ok := proterties[postgresqlPassword]; ok && val != "" {
			conf.Postgresql.Password = val
		}else {
			return errors.New("posgresql store error: password error")
		}

		if val, ok := proterties[postgresqlDBname]; ok && val != "" {
			conf.Postgresql.Db = val
		}else {
			return errors.New("posgresql store error: db not found error")
		}
		s.conf = conf
		return nil
}
