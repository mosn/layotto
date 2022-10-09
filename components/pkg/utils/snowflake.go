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
package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

const (
	mysqlHost         = "mysqlHost"
	mysqlDatabaseName = "databaseName"
	mysqlTableName    = "tableName"
	mysqlUserName     = "userName"
	mysqlPassword     = "password"
	mysqlCharset      = "utf8"
	timeBits          = "timeBits"
	workerBits        = "workerBits"
	seqBits           = "seqBits"
	startTime         = "startTime"

	defaultMysqlTableName = "layotto_sequencer_snowflake"
	defalutTimeBits       = 28
	defalutWorkerBits     = 22
	defalutSeqBits        = 13
	defalutStartTime      = "2022-01-01"
)

type SnowflakeMetadata struct {
	// MysqlMetadata
	UserName     string
	Password     string
	DatabaseName string
	TableName    string
	Db           *sql.DB
	//ip:port
	MysqlHost string

	WorkIdBits int64
	TimeBits   int64
	SeqBits    int64
	StartTime  int64
	WorkId     int64
	LogInfo    bool
}

func ParseSnowflakeMetadata(properties map[string]string) (SnowflakeMetadata, error) {
	metadata := SnowflakeMetadata{}
	var err error

	metadata.TableName = defaultTableName
	if val, ok := properties[mysqlTableName]; ok && val != "" {
		metadata.TableName = val
	}

	if val, ok := properties[mysqlDatabaseName]; ok && val != "" {
		metadata.DatabaseName = val
	} else {
		return metadata, errors.New("mysql connect error: missing database name")
	}

	if val, ok := properties[mysqlUserName]; ok && val != "" {
		metadata.UserName = val
	} else {
		return metadata, errors.New("mysql connect error: missing username")
	}

	if val, ok := properties[mysqlPassword]; ok && val != "" {
		metadata.Password = val
	} else {
		return metadata, errors.New("mysql connect error: missing password")
	}

	if val, ok := properties[mysqlHost]; ok && val != "" {
		metadata.MysqlHost = val
	} else {
		return metadata, errors.New("mysql connect error: missing mysqlHost")
	}

	metadata.WorkIdBits = defalutWorkerBits
	if val, ok := properties[workerBits]; ok && val != "" {
		if metadata.WorkIdBits, err = strconv.ParseInt(val, 10, 64); err != nil {
			return metadata, err
		}
	}
	metadata.TimeBits = defalutTimeBits
	if val, ok := properties[timeBits]; ok && val != "" {
		if metadata.TimeBits, err = strconv.ParseInt(val, 10, 64); err != nil {
			return metadata, err
		}
	}
	metadata.SeqBits = defalutSeqBits
	if val, ok := properties[seqBits]; ok && val != "" {
		if metadata.SeqBits, err = strconv.ParseInt(val, 10, 64); err != nil {
			return metadata, err
		}
	}

	if metadata.TimeBits+metadata.WorkIdBits+metadata.SeqBits+1 != 64 {
		return metadata, errors.New("not enough 64bits")
	}

	s := defalutStartTime
	if val, ok := properties[startTime]; ok && val != "" {
		s = val
	}

	var tmp time.Time
	if tmp, err = time.ParseInLocation("2006-01-02", s, time.Local); err != nil {
		return metadata, err
	}
	metadata.StartTime = tmp.Unix()

	return metadata, nil
}

func NewMysqlClient(meta SnowflakeMetadata) (int64, error) {

	var workId int64
	//for unit test
	if meta.Db == nil {
		mysql := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local", meta.UserName, meta.Password, meta.MysqlHost, meta.DatabaseName, mysqlCharset)
		db, err := sql.Open("mysql", mysql)
		if err != nil {
			return workId, err
		}
		meta.Db = db
	}

	var err error
	createTable := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %s
		(
			ID BIGINT NOT NULL AUTO_INCREMENT COMMENT 'auto increment id',
			HOST_NAME VARCHAR(64) NOT NULL COMMENT 'host name',
			PORT VARCHAR(64) NOT NULL COMMENT 'port',
				CREATED TIMESTAMP NOT NULL COMMENT 'created time',
			PRIMARY KEY(ID)
		)`, meta.TableName)
	if _, err = meta.Db.Exec(createTable); err != nil {
		return workId, err
	}

	workId, err = NewWorkId(meta)
	return workId, err
}

//get id from mysql
//host_name = "ip"
//port = "timestamp-random number"
func NewWorkId(meta SnowflakeMetadata) (int64, error) {
	defer meta.Db.Close()

	var workId int64
	ip, err := getIP()
	stringIp := ip.String()
	if err != nil {
		return workId, err
	}

	begin, err := meta.Db.Begin()
	if err != nil {
		return workId, err
	}

	var host_name string
	var port string
	mysqlPort := getMysqlPort()
	tableName := meta.TableName

	err = begin.QueryRow("SELECT HOST_NAME, PORT FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&host_name, &port)

	//insert a new record if the records are duplicated, to avoid clock rollback problems after shutdown
	for err == nil {
		mysqlPort = getMysqlPort()
		err = begin.QueryRow("SELECT HOST_NAME, PORT FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&host_name, &port)
	}
	if err == sql.ErrNoRows {
		_, err = begin.Exec("INSERT INTO "+tableName+"(HOST_NAME, PORT, CREATED) VALUES(?,?,?)", stringIp, mysqlPort, time.Now())

		if err != nil {
			return workId, err
		}
	} else {
		return workId, err
	}

	err = begin.QueryRow("SELECT ID FROM "+tableName+" WHERE HOST_NAME = ? AND PORT = ?", stringIp, mysqlPort).Scan(&workId)
	if err != nil {
		return workId, err
	}

	if err = begin.Commit(); err != nil {
		begin.Rollback()
		return workId, err
	}
	return workId, nil
}

func getMysqlPort() string {
	currentTimeMills := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	randomData := rand.Int63n(100000)
	mysqlPort := strconv.FormatInt(currentTimeMills, 10) + "-" + strconv.FormatInt(randomData, 10)
	return mysqlPort
}

func getIP() (net.IP, error) {
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
	return nil, errors.New("can't find ip")
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
