package configuration

import (
	"container/list"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/prometheus/prometheus/service/baseconfig"
	"github.com/prometheus/prometheus/service/common"
	"github.com/prometheus/prometheus/service/db"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
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

type HostConfiguration struct {
	ID                  string `json:"id"`
	Name 				string `json:"name"`
	Ip  				string `json:"ip"`
	Port 				string `json:"port"`
	GroupID 			string `json:"group_id"`
	Label				string `json:"label"`
	Status				string `json:"status"`
	GroupName			string `json:"group_name"`
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

//JsonResultGroup 获取组配置信息
type JsonResultGroup struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Group []GroupConfiguration `json:"group"`
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
			//stmt, err := tx.Prepare("insert into configuration (\"name\", finterval, rinterval, aurl, rpath, jpath, \"timeout\") values ('11341', '50s', '30s', '/tat/', '/tmp', '/tmp', '50s');")
			stmt, err := tx.Prepare(`insert into "configuration" ("name", "finterval", "rinterval", "aurl", "rpath", "jpath", "timeout") values ($1, $2, $3, $4, $5, $6, $7)`)
			if err != nil{
				fmt.Println("insert prepare fail" + err.Error())
				jsonResult.Code = 1000
				jsonResult.Msg = "insert prepare fail" + err.Error()
				tx.Rollback()
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

			stmt, err := tx.Prepare(`insert into group_config ("name", "finterval", "scheme", "insecure_skip_verify", "metrics_path", "match_regulation", "federalid", "honor_labels") values ($1, $2, $3, $4, $5, $6, $7, $8)`)
			if err != nil{
				fmt.Println("insert prepare fail" + err.Error())
				jsonResult.Code = 1000
				jsonResult.Msg = "insert prepare fail" + err.Error()
				tx.Rollback()
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
				tx.Rollback()
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


// GetGroups 获取所有group信息
func GetGroups(w http.ResponseWriter, r *http.Request){
	code := 0
	desc := ""

	groupList, err := SelectGroupConfigurations()
	if err != nil {
		code = 1004
		desc = "组信息查询失败"
		fmt.Println(err.Error())
	}

	jsonResult := JsonResultGroup{
		Code: code,
		Msg: desc,
		Group: groupList,
	}

	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}


// ToHostManagement 获取所有group信息
func ToHostManagement(w http.ResponseWriter, r *http.Request){

}


// AddHostConfig 提交的Host配置入库
func AddHostConfig(w http.ResponseWriter, r *http.Request) {

	jsonResult := JsonResult{
		Code: 0,
		Msg: "",
	}

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var hostConfiguration HostConfiguration
	err := json.Unmarshal(con, &hostConfiguration)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
	}else {
		hostConfiguration.Status = "0"
		tx, err := db.DB.Begin()
		if err != nil {
			jsonResult.Code = 1000
			jsonResult.Msg = "db事务开启失败"
			fmt.Println("db事务开启失败")
		}else{
			timeUnix:=time.Now().Unix()   //已知的时间戳
			formatTimeStr:=time.Unix(timeUnix,0).Format("20060102150405")

			hostID := "H" + formatTimeStr + common.RandChar(6)

			stmt, err := tx.Prepare(`insert into host_config ("id", "name", "ip", "port", "group_id", "label", "status") values ($1, $2, $3, $4, $5, $6, $7)`)
			if err != nil{
				fmt.Println(err.Error())
				jsonResult.Code = 1000
				jsonResult.Msg = "insert prepare fail" + err.Error()
				tx.Rollback()
				goto over
			}

			_, err = stmt.Exec(
				hostID,
				hostConfiguration.Name,
				hostConfiguration.Ip,
				hostConfiguration.Port,
				hostConfiguration.GroupID,
				hostConfiguration.Label,
				hostConfiguration.Status,
			)
			if err != nil {
				fmt.Println(err.Error())
				jsonResult.Code = 1000
				jsonResult.Msg = "insert exec fail" + err.Error()
				tx.Rollback()
				goto over
			}
			groupName, err := SelectGroupName(hostConfiguration.GroupID)
			if err != nil {
				tx.Rollback()
				fmt.Println(err)
				jsonResult.Code = 1000
				jsonResult.Msg = "group name select error"
				goto over
			}
			tx.Commit()

			err = RewriteJobFile(hostConfiguration.GroupID, groupName)
			if err != nil{
				fmt.Println(err)
				jsonResult.Code = 1000
				jsonResult.Msg = "Job file rewrite fail"
			}
		}
	}

over:
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
	prometheuYmlConfig = strings.Replace(prometheuYmlConfig, "{{ rule_path }}", path.Join(configuration.RPath, "*.yml"), -1)

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

	// TODO 添加文件存在判断，不存在时创建文件

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("file create failed. err: " + err.Error())
		return err
	} else {
		n, _ := f.Seek(0, os.SEEK_END)
		_, err = f.WriteAt([]byte(content), n)
		fmt.Println("write succeed!")
		defer f.Close()
	}
	return err
}

// RewriteJobFile 重写Job文件，更新监控target信息
func RewriteJobFile(groupID string, groupName string) error{
	sql := "select h.id, h.name, h.ip, h.port, h.group_id, h.label, h.status, g.name from host_config h, group_config g where h.group_id = g.id and g.id = " + groupID
	hostConfigurations, err := SelectHostConfigurations(sql)
	if err != nil{
		fmt.Println("SelectHostConfigurations 查询失败")
		return err
	}

	var contents []string
	for _, hostConfiguration := range hostConfigurations{
		if hostConfiguration.Status != "0"{
			continue
		}
		content := "{\"labels\":" + hostConfiguration.Label + "," + "\"targets\":[\"" + hostConfiguration.Ip + ":" + hostConfiguration.Port + "\"]}"
		contents = append(contents, content)
	}

	contentWriteToFile := "[" + strings.Join(contents, ",") + "]"
	//groupName := hostConfigurations[0].GroupName

	configuration, err := selectPrometheusConfiguration()
	if err != nil {
		return err
	}
	jobPath := configuration.JPath
	filename := path.Join(jobPath, groupName + ".yml")

	err = WriteToFile(filename, contentWriteToFile)
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
		if err == sql.ErrNoRows {
			return configuration, nil
		}
		fmt.Println(err.Error())
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

// SelectGroupConfigurations 获取group配置
func SelectGroupConfigurations() ([]GroupConfiguration, error){

	var groups []GroupConfiguration

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

	for rows.Next() {
		var groupConfiguration GroupConfiguration
		err = rows.Scan(&groupConfiguration.ID, &groupConfiguration.Name, &groupConfiguration.Interval, &groupConfiguration.Scheme , &groupConfiguration.InsecureSkipVerify, &groupConfiguration.MetricsPath, &groupConfiguration.MatchRegulation, &groupConfiguration.Federalid, &groupConfiguration.HonorLabels) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return nil, err
		}
		groups = append(groups, groupConfiguration)
	}

	return groups, nil
}


// SelectHostConfigurations 获取host配置
func SelectHostConfigurations(sql string) ([]HostConfiguration, error){

	var hosts []HostConfiguration

	rows, err := db.DB.Query(sql)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil{
		fmt.Println(err)
		return nil, err
	}

	for rows.Next() {
		var hostConfiguration HostConfiguration
		err = rows.Scan(&hostConfiguration.ID, &hostConfiguration.Name, &hostConfiguration.Ip, &hostConfiguration.Port , &hostConfiguration.GroupID, &hostConfiguration.Label, &hostConfiguration.Status, &hostConfiguration.GroupName) //不scan会导致连接不释放
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return nil, err
		}
		hosts = append(hosts, hostConfiguration)
	}

	return hosts, nil
}


func deletePrometheusConfiguration() error{
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}else {
		stmt, err := tx.Prepare("delete from configuration")
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return err
}

func DbOper(sql string) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}else {
		stmt, err := tx.Prepare(sql)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}
	return err
}


// SelectGroupName 获取group配置
func SelectGroupName(groupID string) (string, error){

	var groupName string

	err := db.DB.QueryRow("select name from group_config where id=" + groupID + " limit 1").Scan(&groupName)
	return groupName, err
}
