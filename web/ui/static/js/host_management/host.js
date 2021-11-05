jQuery.ajax({
    type: "GET",
    url: "searchHosts",
    dataType: 'json',
    async: false,
    error: function () {
        alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
        return false;
    },
    success: function (result) {
        let code = result.code;
        let html = '';
        $("#hosts").html("");
        if(code == 0) {
            for(var index in result.targets){

                let id = result.targets[index].id;
                let name = result.targets[index].name;
                let ip = result.targets[index].ip;
                let port = result.targets[index].port;
                let groupid = result.targets[index].group_id;
                let groupname = result.targets[index].group_name;
                let label = result.targets[index].label;
                let status = result.targets[index].status;

                let y = index % 2;
                if(y == 1){
                    html = html + "<tr style='background-color: whitesmoke'>";
                }else{
                    html = html + "<tr>";
                }
                html = html + '<td style="display: none">' + id + '</td>';
                html = html + '<td>' + name + '</td>';
                html = html + '<td>' + ip + '</td>';
                html = html + '<td>' + port + '</td>';
                html = html + '<td>' + groupname + '</td>';
                html = html + '<td>' + label + '</td>';
                if(status == "0"){
                    html = html + '<td><a href="#" style="color: green" onClick="changeStatus(\'1\',\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\')">启动</a> | <a href="#" style="color: goldenrod" onClick="deleteTarget(\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\')">删除</a></td>';
                }else{
                    html = html + '<td><a href="#" style="color: red" onClick="changeStatus(\'0\',\'' + id + '\',\'' + groupid + '\', \'' +  groupname + '\')">暂停</a> | <a href="#" style="color: goldenrod" onClick="deleteTarget(\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\')">删除</a></td>';
                }
                html = html + "</tr>";
            }

            $("#hosts").html(html);
        }else {
            alert("监控节点信息查询失败");
        }
    }
})


function changeStatus(status, target_id, group_id, group_name){
    alert(status + "," + target_id);

    jQuery.ajax({
        type: "POST",
        url: "updateTargetStatus",
        dataType: 'json',
        async: false,
        data: JSON.stringify({id: target_id, group_id: group_id, status: status, group_name: group_name}),
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            if(result.code == 0){
                alert("状态修改成功");
            }else{
                alert(result.msg);
            }
        }
    })
}

function deleteTarget(target_id){
    alert(target_id);
}