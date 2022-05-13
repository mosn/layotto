package config

/**
* @Author: azh
* @Date: 2022/5/13 22:12
* @Context:
 */

import (
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
	Postgresql Postgresql `json:"postgresql" yaml:"postgresql"`
}

type Postgresql struct {
	Host     string `json:"host" yaml:"host"`
	Port     int64  `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Db       string `json:"db" yaml:"db"`
}

func (m *Server) Load(filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("read yaml file error %v\n", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, m)
	if err != nil {
		fmt.Printf("unmarshal config error %v\n", err)
		return err
	}
	return nil
}
