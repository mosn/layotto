package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Db *sql.DB
)

const (
	defaultTableName    = "tableName"
	tableNameKey        = "tableName"
	connectionStringKey = "connectionString"
	dataBaseName        = "dataBaseName"
)

type MySQLMetadata struct {
	TableName        string
	ConnectionString string
	DataBaseName     string
	Db               *sql.DB
}

func ParseMySQLMetadata(properties map[string]string, db *sql.DB) (MySQLMetadata, error) {
	m := MySQLMetadata{}

	if val, ok := properties[tableNameKey]; ok && val != "" {
		m.TableName = val
	}

	if val, ok := properties[connectionStringKey]; ok && val != "" {
		m.ConnectionString = val
	}

	if val, ok := properties[dataBaseName]; ok && val != "" {
		m.DataBaseName = val
	}

	//m.Db = &sql.DB{}
	m.Db = db

	//Db, err := sql.Open("mysql", m.ConnectionString)
	//
	//m.Db = Db
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("数据库连接成功")

	return m, nil
}

func NewMySQLClient(meta MySQLMetadata) error {

	val := meta

	if val.ConnectionString == "" {
		fmt.Errorf("Missing MySql connection string")
	}

	if val.TableName == "" {
		val.TableName = defaultTableName
	}

	db, err := sql.Open("mysql", val.ConnectionString)

	if err != nil {
		fmt.Errorf("MySql connection fail")
		return err
	}
	meta.Db = db

	return tableExists(meta)
}

func tableExists(meta MySQLMetadata) error {
	exists := ""

	query := `SELECT EXISTS (
		SELECT * FROM ? WHERE TABLE_NAME = ?
		) AS 'exists'`
	err := meta.Db.QueryRow(query, meta.DataBaseName, meta.TableName).Scan(&exists)

	if exists == "1" {
		return nil
	} else {
		return err
	}
}
