searchGroups()

function searchGroups(){

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
            let html = '';
            $("#group_list").html("");

            let groups_id = "";

            if(code == 0) {
                for(var index in result.group){
                    let id = result.group[index].id;
                    let name = result.group[index].name;
                    let interval = result.group[index].finterval;
                    let scheme = result.group[index].scheme;
                    let insecureSkipVerify = result.group[index].insecure_skip_verify;
                    let metricsPath = result.group[index].metrics_path;
                    let matchRegulation = result.group[index].match_regulation;
                    let federalid = result.group[index].federalid;
                    let honorLabels = result.group[index].honor_labels;

                    let y = index % 2;
                    if(y == 1){
                        html = html + "<tr style='background-color: whitesmoke'>";
                    }else{
                        html = html + "<tr>";
                    }

                    html = html + '<td style="display: none">' + id + '</td>';
                    html = html + '<td><input disabled="disabled" style="width: 98%; border: 0; background-color: transparent; outline: none;" id="name-' + id + '" value="' + name + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="interval-' + id + '" value="' + interval + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="scheme-' + id + '" value="' + scheme + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="insecureSkipVerify-' + id + '" value="' + insecureSkipVerify + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="metricsPath-' + id + '" value="' + metricsPath + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="matchRegulation-' + id + '" value="' + matchRegulation + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="federalid-' + id + '" value="' + federalid + '"></td>';
                    html = html + '<td><input style="width: 98%; border: 0; background-color: transparent; outline: none;" id="honorLabels-' + id + '" value="' + honorLabels + '"></td>';
                    html = html + '<td><a href="javascript:void(0)" style="color: cadetblue" onClick="updateGroupInfo(\'' + id + '\')">更新</a> | <a href="javascript:void(0)" style="color: darkred" onClick="deleteGroup(\'' + id + '\')">删除</a></td>';

                    html = html + "</tr>";
                }

                $("#group_list").html(html);

            }else {
                alert("JOB信息查询失败");
            }

        }
    })
}


function deleteGroup(group_id){

    jQuery.ajax({
        type: "POST",
        url: "deleteGroups",
        dataType: 'json',
        data: JSON.stringify({id: group_id}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;

            if(code == 0) {
                alert("删除成功，Manage -> General Management -> 更新文件 后修改yaml，reload后生效");
                searchGroups();
            }else {
                alert("删除失败");
            }

        }
    })

}

function updateGroupInfo(group_id){

    let groupName = $("#name-" + group_id).val();
    let groupInterval = $("#interval-" + group_id).val();
    let groupScheme = $("#scheme-" + group_id).val();
    let groupInsecureSkipVerify = $("#insecureSkipVerify-" + group_id).val();
    let groupMetricsPath = $("#metricsPath-" + group_id).val();
    let groupMatchRegulation = $("#matchRegulation-" + group_id).val();
    let groupFederalid = $("#federalid-" + group_id).val();
    let groupHonorLabels = $("#honorLabels-" + group_id).val();


    jQuery.ajax({
        type: "POST",
        url: "updateGroups",
        dataType: 'json',
        data: JSON.stringify({id: group_id, name: groupName, finterval: groupInterval, scheme: groupScheme,
            insecure_skip_verify: groupInsecureSkipVerify, metrics_path: groupMetricsPath,
            match_regulation: groupMatchRegulation, federalid: groupFederalid, honor_labels: groupHonorLabels}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;

            if(code == 0) {
                alert("更新成功，Manage -> General Management -> 更新文件 后修改yaml，reload后生效");
            }else {
                alert("更新失败");
            }

        }
    })

}