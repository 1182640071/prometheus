getAlarms();
getJobs();

function getJobs() {
    jQuery.ajax({
        type: "GET",
        url: "getGroups",
        dataType: 'json',
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            let html = '<option selected style="text-align: center" value="">请选择所属组</option>';
            $("#jobs").html("");
            if(code == 0) {
                for(var index in result.group){
                    let groupID = result.group[index].id;
                    let groupName = result.group[index].name;
                    html = html + '<option style="text-align: center" value="' + groupID + '">' + groupName + '</option>';
                }
                $("#jobs").html(html);
            }else {
                alert("job信息加载错误");
            }
        }
    })
}

function addAlarm(){
    let job_id = $("#jobs").val().trim();
    let receiver = $("#receiver").val().trim();
    let describe = $("#describe").val().trim();
    jQuery.ajax({
        type: "POST",
        url: "addAlarm",
        dataType: 'json',
        contentType: 'json/application',
        data : JSON.stringify({job_id: job_id, receiver: receiver, describe: describe}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {
                alert("添加成功");
                getAlarms();
            }else{
                alert("添加失败");
            }
        }
    })
}

function getAlarms(){
    jQuery.ajax({
        type: "GET",
        url: "getAlarm",
        dataType: 'json',
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {

                let alarms = result.alarms;

                let html = "";
                for(let index in alarms){
                    let alarm = alarms[index];

                    let job_name = alarm.job_name;
                    let job_id = alarm.job_id;
                    let id = alarm.id;
                    let receiver = alarm.receiver
                    let describe = alarm.describe

                    html = html + '<tr>';
                    html = html + '<td className = "alarm-value" style = "display: none" ></td>';
                    html = html + '<td className="alarm-value">' + job_name + '</td>';
                    html = html + '<td className="alarm-value">' + receiver + '</td>';
                    html = html + '<td className="alarm-value">' + describe + '</td>';
                    html = html + '<td className="alarm-value"><button onClick="deleteAlarm(\'' + id + '\')">删除</button></td>';
                }
                $("#alarm_list").html(html);
            }else{
                alert("添加失败");
            }
        }
    })
}

function deleteAlarm(alarm_id){
    let rs = warning();
    if (!rs){
        return;
    }

    jQuery.ajax({
        type: "POST",
        url: "deleteAlarm",
        dataType: 'json',
        contentType: 'json/application',
        data : JSON.stringify({id: alarm_id}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {
                alert("删除成功");
                getAlarms();
            }else{
                alert("删除失败");
            }
        }
    })
}