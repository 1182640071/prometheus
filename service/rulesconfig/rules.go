package rulesconfig

import (
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/service/common"
	"github.com/prometheus/prometheus/service/configuration"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const RULEHEADER = `
groups:
- name: {{ name }}
  rules:
`

const RULECONTENT = `
groups:
- name: {{ name }}
  rules:
  - alert: {{ name }}
    expr: {{ expr }}
    for: {{ for }}
    labels:
      level: "{{ level }}"
      service: "{{ service }}"
      key_word: "{{ keyword }}"
      value: "{{ $value }}"
    annotations:
      summary: "告警主机和端口{{ $labels.instance }} 告警值：{{ $value }}"
      description: "{{ description }}"
`

type JsonResult struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
}

type RuleConfig struct {
	Name string `json:"name"`
	Value string `json:"value"`

	Exper string `json:"exper"`
	RuleFor string `json:"rulefor"`
	KeyWord string `json:"keyword"`
	Level   string `json:"level"`
	Service string `json:"service"`
	Description string `json:"description"`

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


func DeleteRulesConfig(w http.ResponseWriter, r *http.Request){

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

	err = os.Remove(ruleFilePath)
	if err != nil{
		level.Error(common.Logger).Log("status", "rule文件删除失败", "msg", err.Error())
		jsonResult.Code = 1001
		jsonResult.Msg = "rule文件删除失败"
	}else{
		level.Info(common.Logger).Log("status", "rule文件删除成功", "msg", ruleFilePath)
		jsonResult.Code = 0
		jsonResult.Msg = "rule文件删除成功"
	}

over:
	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}


func AddRulesConfig(w http.ResponseWriter, r *http.Request){

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

	value = strings.ReplaceAll(RULECONTENT, "{{ name }}", ruleConfig.Name)
	value = strings.ReplaceAll(value, "{{ expr }}", ruleConfig.Exper)
	value = strings.ReplaceAll(value, "{{ for }}", ruleConfig.RuleFor)
	value = strings.ReplaceAll(value, "{{ level }}", ruleConfig.Level)
	value = strings.ReplaceAll(value, "{{ service }}", ruleConfig.Service)
	value = strings.ReplaceAll(value, "{{ keyword }}", ruleConfig.KeyWord)
	value = strings.ReplaceAll(value, "{{ description }}", ruleConfig.Description)

	err = configuration.WriteToFile(ruleFilePath, value)
	if err != nil{
		level.Error(common.Logger).Log("status", "rule文件创建失败", "msg", err.Error())
		jsonResult.Code = 1001
		jsonResult.Msg = "rule文件创建失败"
	}else{
		level.Info(common.Logger).Log("status", "rule文件创建成功", "msg", ruleFilePath)
		jsonResult.Code = 0
		jsonResult.Msg = "rule文件创建成功"
	}

over:
	msg, _ := json.Marshal(jsonResult)
	w.Header().Set("content-type","text/json")
	w.WriteHeader(200)
	w.Write(msg)

}




