package config

import (
	"io/ioutil"

	"github.com/olebedev/config"
)

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	Dbname   string
}

type Config struct {
	Mysql
	Selectors []string
	Managers  []string
	Port      string
	ImagesDir string
}

var V Config

func init() {
	file, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	env, err := cfg.String("environment")
	if err != nil {
		panic(err)
	}

	cfg, err = cfg.Get(env)
	if err != nil {
		panic(err)
	}

	V.Mysql.Username, err = cfg.String("mysql.username")
	if err != nil {
		panic(err)
	}
	V.Mysql.Password, err = cfg.String("mysql.password")
	if err != nil {
		panic(err)
	}
	V.Mysql.Host, err = cfg.String("mysql.host")
	if err != nil {
		panic(err)
	}
	V.Mysql.Port, err = cfg.String("mysql.port")
	if err != nil {
		panic(err)
	}
	V.Mysql.Dbname, err = cfg.String("mysql.dbname")
	if err != nil {
		panic(err)
	}

	V.Port, err = cfg.String("port")
	if err != nil {
		panic(err)
	}

	V.ImagesDir, err = cfg.String("imagesDir")
	if err != nil {
		panic(err)
	}

	interfaceList, err := cfg.List("selectors")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(interfaceList); i++ {
		if v, ok := interfaceList[i].(string); ok {
			V.Selectors = append(V.Selectors, v)
		}
	}

	interfaceList, err = cfg.List("managers")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(interfaceList); i++ {
		if v, ok := interfaceList[i].(string); ok {
			V.Managers = append(V.Managers, v)
		}
	}
}
