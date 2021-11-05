package hosts

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/prometheus/service/db"
	"io/ioutil"
	"net/http"
)
import "github.com/prometheus/prometheus/service/configuration"

type JsonResultTargetInfo struct {

	Code int `json:"code"`
	Msg  string `json:"msg"`

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

	sql := "select h.id, h.name, h.ip, h.port, h.group_id, h.label, h.status, g.name from host_config h, group_config g where h.group_id = g.id limit 20 "
	hostConfigurations, err := configuration.SelectHostConfigurations(sql)
	if err != nil{
		jsonResultTargetInfo.Code = 1005
		jsonResultTargetInfo.Msg = "监控节点信息查询错误"
		fmt.Println(err)
	}

	jsonResultTargetInfo.Targets = hostConfigurations

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



