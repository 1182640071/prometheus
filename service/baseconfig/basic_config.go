package baseconfig

import (
	"github.com/Unknwon/goconfig"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"strconv"
	"strings"
)

var BasicConfigs BasicConfig

// PrometheusYmlConfig prometheus基础配置
const PrometheusYmlConfig = "# name: {{ name }}\nglobal:\n  scrape_interval:       {{ scrape_interval }} # By default, scrape targets every 15 seconds.\n  scrape_timeout:       {{ scrape_timeout }}\n  evaluation_interval: {{ evaluation_interval }} # Evaluate rules every 15 seconds.\n\nalerting:\n  alertmanagers:\n  - static_configs:\n    - targets:\n      - {{ targets }}\n\nrule_files:\n  - {{ rule_path }}\n\nscrape_configs:\n"
// JobYmlConfig Job配置模板
const JobYmlConfig = "\n\n  - job_name: '{{ group_name }}'\n    scrape_interval: {{ scrape_interval }}\n    honor_labels: {{ honor }}\n    metrics_path: '{{ metrics }}'\n    {{ match }}\n    scheme: {{ scheme }}\n    {{ tls }}\n    file_sd_configs:\n    - files:\n      - {{ file }}"
// Match 匹配规则模板
const Match = "params:\n      match[]:\n{{ match }}"
// Tls 跳过证书模板
const Tls = "tls_config:\n      insecure_skip_verify: true"


type BasicConfig struct {
	Username 		string
	Password 		string
	Ip       		string
	Port     		string
	Dbname   		string
	DbDriverName 	string

	MaxIdleConns 	int
	MaxOpenConns 	int

	DingUrl 		string

	PrometheusYmlConfigPath string

	NoNeedLoginMetrics []string
}

func InitBasicConfig(configPath string, logger log.Logger) {

	level.Info(logger).Log("msg", "Loading Basic configuration file", "filename", configPath)


	cfg, err := goconfig.LoadConfigFile(configPath)
	if err != nil{
		panic(err)
	}

	BasicConfigs.DbDriverName, err = cfg.GetValue("db", "drivername")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Username, err = cfg.GetValue("db", "username")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Username, err = cfg.GetValue("db", "username")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Password, err = cfg.GetValue("db", "password")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Ip, err = cfg.GetValue("db", "ip")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Port, err = cfg.GetValue("db", "port")
	if err != nil{
		panic(err)
	}

	BasicConfigs.Dbname, err = cfg.GetValue("db", "dbname")
	if err != nil{
		panic(err)
	}

	maxIdleConns, err := cfg.GetValue("db", "maxIdleConns")
	if err != nil{
		panic(err)
	}
	BasicConfigs.MaxIdleConns, err = strconv.Atoi(maxIdleConns)
	if err != nil{
		panic(err)
	}

	maxOpenConns, err := cfg.GetValue("db", "maxOpenConns")
	if err != nil{
		panic(err)
	}
	BasicConfigs.MaxOpenConns, err = strconv.Atoi(maxOpenConns)
	if err != nil{
		panic(err)
	}

	metrics, err := cfg.GetValue("login", "metrics")
	if err != nil{
		panic(err)
	}
	BasicConfigs.NoNeedLoginMetrics = strings.Split(metrics, ",")
}