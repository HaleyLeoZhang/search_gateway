package conf

import (
	"github.com/HaleyLeoZhang/go-component/driver/db"
	"github.com/HaleyLeoZhang/go-component/driver/xelastic"
	"github.com/HaleyLeoZhang/go-component/driver/xkafka"
	"github.com/HaleyLeoZhang/go-component/driver/xredis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var (
	Conf = &Config{}
	//confPath string
)

// Config struct
type Config struct {
	DB    *db.Config       `yaml:"db" json:"db"`
	Redis *xredis.Config   `yaml:"redis" json:"redis"`
	Es    *xelastic.Config `yaml:"elastic" json:"elastic"`
	Kafka *xkafka.Config   `yaml:"kafka" json:"kafka"`
}

func init() {
	//flag.StringVar(&confPath, "conf", "", "conf values")
}

func Init() (err error) {
	var yamlFile string
	//if confPath != "" {
	//	yamlFile, err = filepath.Abs(confPath)
	//} else {
	yamlFile, err = filepath.Abs("../conf/debug.yaml") // 可以参照 debug.example.yaml
	//}
	if err != nil {
		return
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlRead, Conf)
	if err != nil {
		return
	}
	return
}
