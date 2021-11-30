package job

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/prometheus/service/configuration"
	"io/ioutil"
	"net/http"
)

func UpdateJobInfo(w http.ResponseWriter, r *http.Request){
	var jsonResult configuration.JsonResult

	var sql string

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var groupConfiguration configuration.GroupConfiguration
	err := json.Unmarshal(con, &groupConfiguration)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}

	sql = "update group_config set name='" + groupConfiguration.Name + "', finterval='" + groupConfiguration.Interval +
		"', scheme='" + groupConfiguration.Scheme + "', insecure_skip_verify='" + groupConfiguration.InsecureSkipVerify +
		"', metrics_path='" + groupConfiguration.MetricsPath + "', match_regulation='" + groupConfiguration.MatchRegulation +
		"', federalid='" + groupConfiguration.Federalid + "', honor_labels='" + groupConfiguration.HonorLabels +
		"' where id=" + groupConfiguration.ID

	err = configuration.DbOper(sql)
	if err != nil{
		fmt.Println(err.Error())
		jsonResult.Code = 1010
		jsonResult.Msg = "更新失败"
	}else{
		jsonResult.Code = 0
		jsonResult.Msg = "更新成功"
	}

	over:
		msg, _ := json.Marshal(jsonResult)
		w.Header().Set("content-type","text/json")
		w.WriteHeader(200)
		w.Write(msg)
}


func DeleteGroups(w http.ResponseWriter, r *http.Request){
	var jsonResult configuration.JsonResult

	var sql string

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	var groupConfiguration configuration.GroupConfiguration
	err := json.Unmarshal(con, &groupConfiguration)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}

	sql = "delete from group_config where id=" + groupConfiguration.ID

	err = configuration.DbOper(sql)
	if err != nil{
		fmt.Println(err.Error())
		jsonResult.Code = 1010
		jsonResult.Msg = "删除失败"
		goto over
	}else{
		jsonResult.Code = 0
		jsonResult.Msg = "删除成功"
	}

	sql = "delete from host_config where group_id=" + groupConfiguration.ID
	err = configuration.DbOper(sql)
	if err != nil{
		fmt.Println(err.Error())
	}

over:
	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}