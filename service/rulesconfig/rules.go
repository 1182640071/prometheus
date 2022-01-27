package rulesconfig

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/service/common"
	"github.com/prometheus/prometheus/service/configuration"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const RULEHEADER = `
groups:
- name: {{ name }}
  rules:
`

type JsonResult struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
}

type RuleConfig struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func UpdateRulesConfig(w http.ResponseWriter, r *http.Request){

	rulePath := ""
	ruleFilePath := ""
	value := ""
	var configurations configuration.Configuration

	var jsonResult JsonResult
	var ruleConfig RuleConfig

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(con, &ruleConfig)
	if err != nil {
		fmt.Println(err)
		jsonResult.Code = 1000
		jsonResult.Msg = "提交信息解析失败"
		goto over
	}

	configurations, err = configuration.SelectPrometheusConfiguration()
	if err != nil {
		jsonResult.Code = 1001
		jsonResult.Msg = "prometheus基础配置获取失败"
		level.Error(common.Logger).Log("status", "configuration获取失败", "msg", err.Error())
		goto over
	}

	rulePath = configurations.RPath
	ruleFilePath = filepath.Join(rulePath, ruleConfig.Name + ".yml")

	value = strings.ReplaceAll(RULEHEADER, "{{ name }}", ruleConfig.Name) + "\n"
	value = value + "  - " + strings.ReplaceAll(ruleConfig.Value, "\n", "\n    ")

	err = configuration.WriteToFile(ruleFilePath, value)
	if err != nil{
		level.Error(common.Logger).Log("status", "rule文件修改失败", "msg", err.Error())
		jsonResult.Code = 1001
		jsonResult.Msg = "rule文件修改失败"
	}else{
		level.Info(common.Logger).Log("status", "rule文件修改成功", "msg", ruleFilePath)
		jsonResult.Code = 0
		jsonResult.Msg = "rule文件修改成功"
	}

over:
	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}

