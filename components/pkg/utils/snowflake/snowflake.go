//
// Copyright 2021 Layotto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package snowflake

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DefaultMysqlUrl     = "mysql"
	DefaultDatabaseName = "databaseName"
	DefaultTableName    = "tableName"
	DefaultUserName     = "userName"
	DefaultPassword     = "password"
	Charset             = "utf8"
)

type SnowflakeMetadata struct {
	MysqlMetadata    *MysqlMetadata
	Size             int
	PaddingFactor    int
	ScheduleInterval int
	TimeBits         int64
	WorkerBits       int64
	SeqBits          int64
	StartTime        string
	LogInfo          bool
}

type MysqlMetadata struct {
	UserName     string
	Password     string
	DatabaseName string
	TableName    string
	Db           *sql.DB
	MysqlUrl     string
}

type MysqlConn interface {
	MysqlConnect(meta MysqlMetadata)
}

func ParseSnowflakeMysqlMetadata(properties map[string]string) (MysqlMetadata, error) {
	m := MysqlMetadata{}
	if val, ok := properties[DefaultTableName]; ok && val != "" {
		m.TableName = val
	}
	if val, ok := properties[DefaultDatabaseName]; ok && val != "" {
		m.DatabaseName = val
	}
	if val, ok := properties[DefaultUserName]; ok && val != "" {
		m.UserName = val
	}
	if val, ok := properties[DefaultPassword]; ok && val != "" {
		m.Password = val
	}
	if val, ok := properties[DefaultMysqlUrl]; ok && val != "" {
		m.MysqlUrl = val
	}
	return m, nil
}

func NewMysqlClient(meta MysqlMetadata) (int64, error) {
	if meta.TableName == "" {
		meta.TableName = DefaultTableName
	}

	//for unit test
	if meta.Db == nil {
		mysql := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local", meta.UserName, meta.Password, meta.MysqlUrl, meta.DatabaseName, Charset)
		db, err := sql.Open("mysql", mysql)
		if err != nil {
			return 0, err
		}
		meta.Db = db
	}

	exists, err := tableExists(meta)

	if !exists {
		createTable := fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s(
				ID BIGINT NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
				HOST_NAME VARCHAR(64) NOT NULL COMMENT 'host name',
				PORT VARCHAR(64) NOT NULL COMMENT 'port',
				TYPE INT NOT NULL COMMENT 'node type: ACTUAL or CONTAINER',
				LAUNCH_DATE DATE NOT NULL COMMENT 'launch date',
			    CREATED TIMESTAMP NOT NULL COMMENT 'created time',
				PRIMARY KEY(ID)
			)`, meta.TableName)
		_, err = meta.Db.Exec(createTable)
		if err != nil {
			return 0, err
		}
	}

	workId, err := NewWorkId(meta)
	return workId, err
}

func tableExists(meta MysqlMetadata) (bool, error) {
	var result string
	err := meta.Db.QueryRow(fmt.Sprintf(`SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_NAME = ?`), meta.TableName).Scan(&result)
	if err == sql.ErrNoRows {
		return false, err
	}
	return true, err
}

func NewWorkId(meta MysqlMetadata) (int64, error) {
	defer meta.Db.Close()
	ip, err := externalIP()
	stringIp := ip.String()
	if err != nil {
		return 0, fmt.Errorf("ip get error: %s", err)
	}

	var env int
	if isDocker() {
		env = 0
	} else {
		env = 1
	}

	begin, err := meta.Db.Begin()
	if err != nil {
		return 0, fmt.Errorf("db begin error: %s", err)
	}

	var host_name string
	var port string
	mysqlPort := getMysqlPort()
	tableName := meta.TableName

	err = begin.QueryRow("SELECT HOST_NAME, PORT FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&host_name, &port)

	for err == nil {
		mysqlPort = getMysqlPort()
		err = begin.QueryRow("SELECT HOST_NAME, PORT FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&host_name, &port)
	}

	if err == sql.ErrNoRows {
		_, err = begin.Exec("INSERT INTO "+tableName+"(HOST_NAME, PORT, TYPE, LAUNCH_DATE, CREATED) VALUES(?,?,?,?,?)", stringIp, mysqlPort, env, time.Now(), time.Now())

		if err != nil {
			return 0, fmt.Errorf("mysql insert error: %s", err)
		}
	}

	var workId int64
	err = begin.QueryRow("SELECT ID FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&workId)
	if err != nil {
		return 0, err
	}

	if err = begin.Commit(); err != nil {
		begin.Rollback()
		return 0, fmt.Errorf("mysql commit error: %s", err)
	}
	return workId, nil
}
func getMysqlPort() string {
	currentTimeMills := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	randomData := rand.Int63n(100000)
	mysqlPort := strconv.Itoa(int(currentTimeMills)) + "-" + strconv.Itoa(int(randomData))
	return mysqlPort
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil
	}

	return ip
}
func isDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	} else {
		return false
	}
}
