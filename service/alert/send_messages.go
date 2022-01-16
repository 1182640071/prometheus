package alert

import (
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/service/common"
	"github.com/prometheus/prometheus/service/loadconfig"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)


type DingResult struct {
	ErrCode int `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendMessages 发送告警信息
func SendMessages(w http.ResponseWriter, r *http.Request) {

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	var alertMessage AlertMessage

	defer r.Body.Close()
	con, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(con, &alertMessage)
	if err != nil {
		level.Error(common.Logger).Log("status", "提交信息解析失败", "err", err.Error())
	}

	if len(Messages) < 100 {
		Messages = append(Messages, alertMessage)
	}else{
		level.Warn(common.Logger).Log( "info", "队列长度大于100，告警信息跳过", "msg", alertMessage.CommonAnnotations.Description)
	}
}

func SendAlert(){
	pool := common.NewPool(5)

	for {

		if len(Messages) < 1 {
			time.Sleep(20*time.Second)
			continue
		}

		message := Messages[0]
		if len(Messages) > 1 {
			Messages = Messages[1:]
		}else{
			Messages = make([]AlertMessage, 0)
		}

		pool.Add()

		go func(){
			defer pool.JobDone()
			err := send(message)

			jsonByte, _ := json.Marshal(message)
			jsonStr := string(jsonByte)

			if err != nil{
				level.Error(common.Logger).Log("status", "发送失败", "message", jsonStr)
			}else{
				level.Info(common.Logger).Log("status", "发送成功", "message", jsonStr)
			}
		}()

	}
}

func send(message AlertMessage) error{

	title := ""

	if "firing" == message.Status {
		title = "<font color=#FF0000>告警</font> <br>"
	} else {
		title = "<font color=#00FF00>恢复</font> <br>"
	}

	startTime := message.Alert[0].StartsAt.Format("2006-01-02 15:04:05") //time转string
	sendTime := time.Now().Format("2006-01-02 15:04:05")

	webHook := ""

	jobName := strings.Trim(message.CommonLabels.Job, " ")
	value, ok := loadconfig.AlarmTypeConfigs[jobName]
	if !ok {
		level.Error(common.Logger).Log("status", "不存在" + jobName + "告警方式")
		return errors.New("不存在" + jobName + "告警方式" )
	}else{
		webHook = value
	}

	content := "> 状态：" + title + " \\n\\n " +
				"告警key：" + message.CommonLabels.KeyWord + " \\n\\n " +
				"故障时间：" + startTime + " \\n\\n " +
				"发送时间：" + sendTime + " \\n\\n " +
				"实例：" + message.CommonLabels.Name + " \\n\\n " +
				"节点：" + message.CommonLabels.Instance + " \\n\\n " +
				"平台：" + jobName + " \\n\\n " +
				"告警信息：" + message.CommonAnnotations.Description + " \\n\\n "

	if "firing" == message.Status {
		title = "<font color=#FF0000>告警</font> <br>"
		content += "告警值：" + message.CommonLabels.Value + " \\n\\n "
	}

	textMsg := "{ \"msgtype\": \"markdown\", \"markdown\": {\"title\": \"" + message.CommonAnnotations.Description + "\", \"text\": \"" + content + "\"}}"
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(textMsg))
	if err != nil {
		level.Error(common.Logger).Log("err", err.Error())
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Connection","Close")
	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		level.Error(common.Logger).Log("err", err.Error())
		return err
	}
	//关闭请求
	defer resp.Body.Close()

	body,_ := ioutil.ReadAll(resp.Body)
	//body {"errcode":0,"errmsg":"ok"}
	var dingResult DingResult
	err = json.Unmarshal(body, &dingResult)

	if dingResult.ErrCode != 0 {
		level.Error(common.Logger).Log("status", "发送失败", "msg", dingResult.ErrMsg)
		return err
	}

	if err != nil {
		level.Error(common.Logger).Log("err", err.Error())
		return err
	}else{
		return nil
	}

}