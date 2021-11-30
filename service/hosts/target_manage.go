package hosts

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/prometheus/service/configuration"
	"github.com/prometheus/prometheus/service/db"
	"io/ioutil"
	"net/http"
)

type JsonResultTargetInfo struct {

	Code int `json:"code"`
	Msg  string `json:"msg"`
	PageNow int `json:"page_now"`
	Count int `json:"count"`
	Pages int `json:"pages"`

	GroupName string `json:"group_name"`
	GroupID string `json:"group_id"`
	TargetName string `json:"target_name"`
	TargetIp string `json:"target_ip"`
	TargetPort string `json:"target_port"`
	TargetKeyword string `json:"target_keyword"`

	Targets []configuration.HostConfiguration `json:"targets"`
}

type TargetStatus struct {
	ID 			string `json:"id"`
	GroupID 	string `json:"group_id"`
	GroupName   string `json:"group_name"`
	Status 		string `json:"status"`
}


func SearchTargets(w http.ResponseWriter, r *http.Request){

	jsonResultTargetInfo := JsonResultTargetInfo{
		Code: 0,
		Msg: "OK",
	}

	var sql string
	var hostConfigurations []configuration.HostConfiguration
	var count int

	var groupID string
	var targetName string
	var targetIp string
	var targetPort string
	var targetKeyword string

	pageNow := 0
	start := 0
	end := 0
	size := 20
	pages := 0

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(con, &jsonResultTargetInfo)
	if err != nil {
		fmt.Println(err)
		jsonResultTargetInfo.Code = 1000
		jsonResultTargetInfo.Msg = "提交信息解析失败"
		goto over
	}

	sql = "select h.id, h.name, h.ip, h.port, h.group_id, h.label, h.status, g.name from host_config h, group_config g where h.group_id = g.id "

	groupID = jsonResultTargetInfo.GroupID
	targetName = jsonResultTargetInfo.TargetName
	targetIp = jsonResultTargetInfo.TargetIp
	targetPort = jsonResultTargetInfo.TargetPort
	targetKeyword = jsonResultTargetInfo.TargetKeyword

	if groupID != "" {
		sql += " and h.group_id = " + groupID
	}

	if targetName != "" {
		sql += " and h.name like '%" + targetName + "%'"
	}

	if targetIp != "" {
		sql += " and h.ip = '" + targetIp + "'"
	}

	if targetPort != "" {
		sql += " and h.port = " + targetPort
	}

	if targetKeyword != "" {
		sql += " and h.label like '%" + targetKeyword + "%'"
	}

	sql += " order by h.id desc"

	hostConfigurations, err = configuration.SelectHostConfigurations(sql)
	if err != nil{
		jsonResultTargetInfo.Code = 1005
		jsonResultTargetInfo.Msg = "监控节点信息查询错误"
		fmt.Println(err)
	}

	count = len(hostConfigurations)

	pages = count / 20
	if count % 20 != 0{
		pages += 1
	}

	jsonResultTargetInfo.Pages = pages
	pageNow = jsonResultTargetInfo.PageNow
	start = size * (pageNow - 1)
	end = size * pageNow

	if end > count {
		end = count
	}

	jsonResultTargetInfo.Targets = hostConfigurations[start:end]
	jsonResultTargetInfo.Count = count

	over:
		msg, _ := json.Marshal(jsonResultTargetInfo)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)
}


func UpdateTargetsStatus(w http.ResponseWriter, r *http.Request){

	jsonResult := JsonResultTargetInfo{
		Code: 0,
		Msg: "OK",
	}

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var targetStatus TargetStatus
	err := json.Unmarshal(con, &targetStatus)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
	}else{
		sql := "update host_config set status = '" + targetStatus.Status + "' where id='" + targetStatus.ID + "'"
		tx, err := db.DB.Begin()
		if err != nil {
			tx.Rollback()
			jsonResult.Code = 1002
			jsonResult.Msg = "状态更新失败，事务开启失败"
			goto over
		}else {
			stmt, err := tx.Prepare(sql)
			if err != nil {
				tx.Rollback()
				jsonResult.Code = 1002
				jsonResult.Msg = "状态更新失败，sql预备失败"
				goto over
			}
			_, err = stmt.Exec()
			if err != nil {
				tx.Rollback()
				jsonResult.Code = 1002
				jsonResult.Msg = "状态更新失败，sql执行失败"
				goto over
			}

			if targetStatus.GroupName == ""{
				fmt.Println("组名为空")
				jsonResult.Code = 1002
				jsonResult.Msg = "组名为空"
				tx.Rollback()
				goto over
			}

			tx.Commit()
			err = configuration.RewriteJobFile(targetStatus.GroupID, targetStatus.GroupName)
			if err != nil{
				fmt.Println(err)
				jsonResult.Code = 1002
				jsonResult.Msg = "Job文件更新失败"
			}
		}
	}

	over:
		msg, _ := json.Marshal(jsonResult)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)
}


func DeleteTargetsStatus(w http.ResponseWriter, r *http.Request){
	jsonResult := JsonResultTargetInfo{
		Code: 0,
		Msg: "OK",
	}

	var sql string

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var targetStatus TargetStatus
	err := json.Unmarshal(con, &targetStatus)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}

	sql = "delete from host_config where id = '" + targetStatus.ID + "'"

	err = configuration.DbOper(sql)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "监控节点删除失败"
		goto over
	}

	err = configuration.RewriteJobFile(targetStatus.GroupID, targetStatus.GroupName)
	if err != nil{
		fmt.Println(err)
		jsonResult.Code = 1002
		jsonResult.Msg = "Job文件更新失败"
	}

	over:
		msg, _ := json.Marshal(jsonResult)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)

}


