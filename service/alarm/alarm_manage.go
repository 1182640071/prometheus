package alarm

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/service/common"
	"github.com/prometheus/prometheus/service/configuration"
	"github.com/prometheus/prometheus/service/db"
	"io/ioutil"
	"net/http"
	"time"
)

type ConfigsAlarm struct {
	ID                  string `json:"id"`
	JobID 				string `json:"job_id"`
	Receiver  			string `json:"receiver"`
	Describe 			string `json:"describe"`
	JobName				string `json:"job_name"`

	Code int `json:"code"`
	Msg  string `json:"msg"`
}

type AlarmsConfig struct {
	Alarms  []ConfigsAlarm `json:"alarms"`
	Code int `json:"code"`
	Msg  string `json:"msg"`
}

func AddAlarmConfig(w http.ResponseWriter, r *http.Request){
	var jsonResult ConfigsAlarm

	sql := ""
	timeUnix:=time.Now().Unix()   //已知的时间戳
	formatTimeStr:=time.Unix(timeUnix,0).Format("20060102150405")

	alarmID := "A" + formatTimeStr + common.RandChar(6)

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(con, &jsonResult)
	if err != nil {
		level.Error(common.Logger).Log("status", "提交信息解析失败", "err", err.Error())
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}

	sql = "insert into alarm_config (id, job_id, receiver, describe ) values ('" + alarmID + "', '" + jsonResult.JobID +
		"', '" + jsonResult.Receiver + "', '" + jsonResult.Describe + "')"
	err = configuration.DbOper(sql)
	if err != nil{
		level.Error(common.Logger).Log("status", "alarm_config表查询错误", "err", err.Error())
		jsonResult.Code = 1010
		jsonResult.Msg = "告警方式添加失败"
	}else{
		jsonResult.Code = 0
		jsonResult.Msg = "告警方式添加成功"
	}

over:
	msg, _ := json.Marshal(jsonResult)
	jsonStr := string(msg)
	level.Info(common.Logger).Log("status", jsonResult.Msg, "message", jsonStr)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)
}

func GetAlarmConfigs(w http.ResponseWriter, r *http.Request){
	var jsonResult AlarmsConfig

	alarms, err := SelectAlarmConfigurations()
	if err != nil{
		jsonResult.Code = 1010
		jsonResult.Msg = "获取告警方式失败"
	}else{
		jsonResult.Code = 0
		jsonResult.Msg = "获取告警方式成功"
	}
	jsonResult.Alarms = alarms

	msg, _ := json.Marshal(jsonResult)
	jsonStr := string(msg)
	level.Info(common.Logger).Log("status", jsonResult.Msg, "message", jsonStr)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)
}


func DeleteAlarmConfig(w http.ResponseWriter, r *http.Request){
	var jsonResult ConfigsAlarm

	sql := ""

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(con, &jsonResult)
	if err != nil {
		level.Error(common.Logger).Log("status", "提交信息解析失败", "err", err.Error())
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}
	sql = "delete from alarm_config where id='" + jsonResult.ID + "';"

	err = configuration.DbOper(sql)
	if err != nil{
		jsonResult.Code = 1010
		jsonResult.Msg = "告警方式删除失败"
		goto over
	}else{
		jsonResult.Code = 0
		jsonResult.Msg = "告警方式删除成功"
	}

over:
	msg, _ := json.Marshal(jsonResult)
	jsonStr := string(msg)
	level.Info(common.Logger).Log("status", jsonResult.Msg, "message", jsonStr)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)
}


// SelectAlarmConfigurations 获取alarm配置
func SelectAlarmConfigurations() ([]ConfigsAlarm, error){

	var alarms []ConfigsAlarm
	sql := "select g.name as job_name, a.id, a.job_id, a.receiver, a.describe from group_config g, alarm_config a where g.id=a.job_id;"

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
		var configsAlarm ConfigsAlarm
		err = rows.Scan(&configsAlarm.JobName, &configsAlarm.ID, &configsAlarm.JobID, &configsAlarm.Receiver , &configsAlarm.Describe) //不scan会导致连接不释放
		if err != nil {
			level.Error(common.Logger).Log("status", "alarm_config表查询错误", "err", err.Error())
			return nil, err
		}
		alarms = append(alarms, configsAlarm)
	}

	return alarms, nil
}
