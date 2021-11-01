package configuration

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/prometheus/prometheus/service/baseconfig"
	"github.com/prometheus/prometheus/service/db"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

// Configuration proemtheus配置信息
type Configuration  struct{
	Name 		string `json:"name"`
	Interval  	string `json:"finterval"`
	RInterval 	string `json:"rinterval"`
	AUrl 		string `json:"aurl"`
	RPath		string `json:"rpath"`
	JPath		string `json:"jpath"`
	TimetOut	string `json:"timeout"`
}

// GroupConfiguration Group(Job)配置信息
type GroupConfiguration  struct{
	ID                  string `json:"id"`
	Name 				string `json:"name"`
	Interval  			string `json:"finterval"`
	Scheme 				string `json:"scheme"`
	InsecureSkipVerify 	string `json:"insecure_skip_verify"`
	MetricsPath			string `json:"metrics_path"`
	MatchRegulation		string `json:"match_regulation"`
	Federalid			string `json:"federalid"`
	HonorLabels 		string `json:"honor_labels"`
}


// JsonResult 登录认真结果
type JsonResult  struct{
	Code int `json:"code"`
	Msg  string `json:"msg"`
}

//JsonResultConfiguration 获取基础配置信息
type JsonResultConfiguration struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Configuration struct{
		Name 		string `json:"name"`
		Interval  	string `json:"finterval"`
		RInterval 	string `json:"rinterval"`
		AUrl 		string `json:"aurl"`
		RPath		string `json:"rpath"`
		JPath		string `json:"jpath"`
		TimetOut	string `json:"timeout"`
	} `json:"configuration"`
}


// SubmitConfiguration 提交的prometheus.yml配置信息提交入库
func SubmitConfiguration(w http.ResponseWriter, r *http.Request) {

	jsonResult := JsonResult{
		Code: 0,
		Msg: "",
	}

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var configuration Configuration
	err := json.Unmarshal(con, &configuration)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
	}else {
		tx, err := db.DB.Begin()
		if err != nil {
			jsonResult.Code = 1000
			jsonResult.Msg = "db事务开启失败"
			fmt.Println("db事务开启失败")

		}else{
			err := deletePrometheusConfiguration()
			if err != nil{
				jsonResult.Code = 1000
				jsonResult.Msg = "clean history data fail"
				fmt.Println("clean history data fail" + err.Error())
				goto over
			}
			stmt, err := tx.Prepare("insert into configuration (`name`, `finterval`, `rinterval`, `aurl`, `rpath`, `jpath`, `timeout`) values (?, ?, ?, ?, ?, ?, ?)")
			if err != nil{
				fmt.Println("insert prepare fail" + err.Error())
				jsonResult.Code = 1000
				jsonResult.Msg = "insert prepare fail"
				goto over
			}

			_, err = stmt.Exec(
				configuration.Name,
				configuration.Interval,
				configuration.RInterval,
				configuration.AUrl,
				configuration.RPath,
				configuration.JPath,
				configuration.TimetOut,
			)
			if err != nil {
				fmt.Println("insert exec fail")
				jsonResult.Code = 1000
				jsonResult.Msg = "insert exec fail"
			}else{
				tx.Commit()
			}
		}
	}

	over:
	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}


// AddGroupConfiguration 提交的Group组配置入库
func AddGroupConfiguration(w http.ResponseWriter, r *http.Request) {

	jsonResult := JsonResult{
		Code: 0,
		Msg: "",
	}

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var groupConfiguration GroupConfiguration
	err := json.Unmarshal(con, &groupConfiguration)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
	}else {
		tx, err := db.DB.Begin()
		if err != nil {
			jsonResult.Code = 1000
			jsonResult.Msg = "db事务开启失败"
			fmt.Println("db事务开启失败")
		}else{

			if groupConfiguration.MetricsPath == "" {
				groupConfiguration.MetricsPath = "/metrics"
			}

			stmt, err := tx.Prepare("insert into group_config (`name`, `finterval`, `scheme`, `insecure_skip_verify`, `metrics_path`, `match_regulation`, `federalid`, `honor_labels`) values (?, ?, ?, ?, ?, ?, ?, ?)")
			if err != nil{
				fmt.Println("insert prepare fail")
				jsonResult.Code = 1000
				jsonResult.Msg = "insert prepare fail"
				goto over
			}

			_, err = stmt.Exec(
				groupConfiguration.Name,
				groupConfiguration.Interval,
				groupConfiguration.Scheme,
				groupConfiguration.InsecureSkipVerify,
				groupConfiguration.MetricsPath,
				groupConfiguration.MatchRegulation,
				groupConfiguration.Federalid,
				groupConfiguration.HonorLabels,
			)
			if err != nil {
				fmt.Println("insert exec fail")
				jsonResult.Code = 1000
				jsonResult.Msg = "insert exec fail"
				goto over
			}

			err = createJobFile(groupConfiguration.Name)
			if err != nil {
				fmt.Println(err)
				jsonResult.Code = 1000
				jsonResult.Msg = "Job对应的存储节点文件创建失败"
				tx.Rollback()
			}else{
				tx.Commit()
			}
		}
	}

	over:
		msg, _ := json.Marshal(jsonResult)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)

}

// GetConfiguration 获取prometheus.yml基础配置信息
func GetConfiguration(w http.ResponseWriter, r *http.Request) {

	var jsonResult JsonResultConfiguration

	configuration, err := selectPrometheusConfiguration()
	if err != nil{
		jsonResult.Code = 1002
		jsonResult.Msg = "配置信息获取失败"
		goto over
	}

	jsonResult.Configuration.Name = configuration.Name
	jsonResult.Configuration.Interval = configuration.Interval
	jsonResult.Configuration.RInterval = configuration.RInterval
	jsonResult.Configuration.AUrl = configuration.AUrl
	jsonResult.Configuration.RPath = configuration.RPath
	jsonResult.Configuration.JPath = configuration.JPath
	jsonResult.Configuration.TimetOut = configuration.TimetOut

	over:
		msg, _ := json.Marshal(jsonResult)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)
}

// UpdatePrometheusYmlConfig 覆写prometheus.yml文件
func UpdatePrometheusYmlConfig(w http.ResponseWriter, r *http.Request){

	code, desc := RewritePrometheusYmlConfig()

	jsonResult := JsonResult{
		Code: code,
		Msg: desc,
	}

	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}

func RewritePrometheusYmlConfig() (int, string){
	var groupList *list.List
	var prometheuYmlConfig string
	configuration, err := selectPrometheusConfiguration()
	if err != nil{
		fmt.Println(err)
		return 1003, "prometheus.yml基础配置获取失败"
	}

	prometheuYmlConfig = baseconfig.PrometheusYmlConfig
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ name }}", configuration.Name, -1)
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ scrape_interval }}", configuration.Interval, -1)
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ scrape_timeout }}", configuration.TimetOut, -1)
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ evaluation_interval }}", configuration.RInterval, -1)
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ targets }}", configuration.AUrl, -1)
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ rule_path }}", configuration.RPath, -1)

	// 获取所有group信息
	groupList, err = selectGroupConfiguration()
	for p := groupList.Front(); p != nil; p = p.Next() {
		var groupConfiguration = p.Value.(GroupConfiguration)

		if groupConfiguration.Interval == ""{
			groupConfiguration.Interval = configuration.Interval
		}

		groupPath := path.Join(configuration.JPath, groupConfiguration.Name + ".yml")

		groupYmlConfig := baseconfig.JobYmlConfig
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ group_name }}", groupConfiguration.Name, -1)
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ scrape_interval }}", groupConfiguration.Interval, -1)
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ honor }}", groupConfiguration.HonorLabels, -1)
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ metrics }}", groupConfiguration.MetricsPath, -1)
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ scheme }}", groupConfiguration.Scheme, -1)
		groupYmlConfig = strings.Replace(groupYmlConfig, "{{ file }}", groupPath, -1)

		if groupConfiguration.Scheme == "https" && groupConfiguration.InsecureSkipVerify == "true" {
			groupYmlConfig = strings.Replace(groupYmlConfig, "{{ tls }}", baseconfig.Tls, -1)
		}else {
			groupYmlConfig = strings.Replace(groupYmlConfig, "{{ tls }}", "", -1)
		}

		if groupConfiguration.MatchRegulation == "" {
			groupYmlConfig = strings.Replace(groupYmlConfig, "{{ match }}", "", -1)
		}else{
			matchs := strings.Split(groupConfiguration.MatchRegulation, ",")
			matchRegulation := ""
			matchMould := baseconfig.Match
			for _, regulation := range matchs{
				matchRegulation += "        - '" + regulation + "'\n"
			}
			matchMould = strings.Replace(matchMould, "{{ match }}", matchRegulation, -1)


			groupYmlConfig = strings.Replace(groupYmlConfig, "{{ match }}", matchMould, -1)
		}
		prometheuYmlConfig += groupYmlConfig
	}

	//覆写prometheus.yml文件
	err = WriteToFile(baseconfig.BasicConfigs.PrometheusYmlConfigPath, prometheuYmlConfig)
	if err != nil{
		return 1003, "prometheus.yml基础配置文件写入失败"
	}
	return 0, "prometheus.yml基础配置文件写入成功"
}

// WriteToFile /* 文件操作 */
// os.O_TRUNC 覆盖写入，不加则追加写入
func WriteToFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
		return err
	} else {
		// offset
		//os.Truncate(filename, 0) //clear
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

// 创建组(Job)后，需要创建对应的yml文件，用于存储监控target节点信息
func createJobFile(name string) error{
	configuration, err := selectPrometheusConfiguration()
	if err != nil {
		return err
	}
	jobPath := configuration.JPath
	_, err = os.Create(path.Join(jobPath, name + ".yml")) //创建文件
	if err != nil {
		return err
	}
	WriteToFile(path.Join(jobPath, name + ".yml"), "[]")
	return nil
}

// 获取prometheus基础配置
func selectPrometheusConfiguration() (Configuration, error){
	var configuration Configuration
	err := db.DB.QueryRow("select name, finterval, rinterval, aurl, rpath, jpath, timeout from configuration limit 1 ").Scan(
		&configuration.Name, &configuration.Interval, &configuration.RInterval, &configuration.AUrl, &configuration.RPath, &configuration.JPath, &configuration.TimetOut)
	if err != nil{
		fmt.Println(err)
		return configuration, err
	}
	return configuration, nil
}

// 获取group配置
func selectGroupConfiguration() (*list.List, error){

	rows, err := db.DB.Query("select id, name, finterval, scheme, insecure_skip_verify, metrics_path, match_regulation, federalid, honor_labels from group_config ")
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	groupList := list.New()

	for rows.Next() {
		var groupConfiguration GroupConfiguration
		err = rows.Scan(&groupConfiguration.ID, &groupConfiguration.Name, &groupConfiguration.Interval, &groupConfiguration.Scheme , &groupConfiguration.InsecureSkipVerify, &groupConfiguration.MetricsPath, &groupConfiguration.MatchRegulation, &groupConfiguration.Federalid, &groupConfiguration.HonorLabels) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return nil, err
		}
		groupList.PushBack(groupConfiguration)
	}

	return groupList, nil
}

func deletePrometheusConfiguration() error{
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}else {
		stmt, err := tx.Prepare("delete from configuration")
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		tx.Commit()
	}
	return nil
}