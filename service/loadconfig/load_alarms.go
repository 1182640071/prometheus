package loadconfig

import (
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/prometheus/service/alarm"
	"github.com/prometheus/prometheus/service/common"
	"time"
)

var AlarmTypeConfigs map[string]string

func LoadAlarms(){
	for {
		//60秒加载一次
		time.Sleep(time.Duration(60)*time.Second)
		Load()
	}
}

func Load(){
	var configsAlarm []alarm.ConfigsAlarm
	configsAlarm, err := alarm.SelectAlarmConfigurations()
	if err != nil{
		level.Error(common.Logger).Log("status", "alarm_config表查询信息加载错误", "err", err.Error())
	}

	tmpMap := make(map[string]string, len(configsAlarm))
	for _, alarm := range configsAlarm{
		jobName := alarm.JobName
		tmpMap[jobName] = alarm.Receiver
	}
	AlarmTypeConfigs = tmpMap
}
