selectHost();
searchHosts(1);

function searchHosts(pageNow){

    var groupID = $("#groups").val();
    var targetName = $("#target-name").val();
    var targetIp = $("#target-ip").val();
    var targetPort = $("#target-port").val();
    var targetKeyword = $("#target-keyword").val();

    jQuery.ajax({
        type: "Post",
        url: "searchHosts",
        dataType: 'json',
        data: JSON.stringify({ page_now: pageNow, group_id: groupID, target_name: targetName, target_ip: targetIp, target_port: targetPort, target_keyword: targetKeyword} ),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            let html = '';
            $("#hosts").html("");

            // let group_html = '<option selected style="text-align: center" value="">请选择所属组</option>';
            // $("#groups").html("");

            let groups_id = "";

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
                        html = html + '<td><a href="javascript:void(0)" style="color: green" onClick="changeStatus(\'1\',\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\', \'' + pageNow + '\')">启动</a> | <a href="javascript:void(0)" style="color: goldenrod" onClick="deleteTarget(\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\', \'' + pageNow + '\')">删除</a></td>';
                    }else{
                        html = html + '<td><a href="javascript:void(0)" style="color: red" onClick="changeStatus(\'0\',\'' + id + '\',\'' + groupid + '\', \'' +  groupname + '\', \'' + pageNow + '\')">暂停</a> | <a href="javascript:void(0)" style="color: goldenrod" onClick="deleteTarget(\'' + id + '\',\'' + groupid + '\', \'' + groupname + '\', \'' + pageNow + '\')">删除</a></td>';
                    }
                    html = html + "</tr>";

                    // if(groups_id.indexOf(groupid) < 0){
                    //     group_html = group_html + '<option style="text-align: center" value="' + groupid + '">' + groupname + '</option>';
                    //     groups_id = groups_id + "," + groupid;
                    // }
                }

                // let select_groups = $("#groups").html();
                // if(select_groups.trim() == ""){
                //     $("#groups").html(group_html);
                // }
                $("#hosts").html(html);

                let pageCount = result.pages;
                let uppage = 1;
                let downpage = 1;

                html = '<label style="font-size: 14px; font-family: \'Charter\'">共有' + result.count + '条,当前页:' + pageNow + '  ' + '</label>';

                if(pageNow != 1){
                    uppage = pageNow -1;
                    // page = uppage;
                }
                if(pageNow != pageCount){
                    downpage = pageNow +1;
                    // page = downpage;
                }else{
                    downpage = pageNow;
                }

                html = html + '<a href="javascript:searchHosts('+ uppage +')" style="font-size: 14px; margin-left: 3px; padding: 1px; font-size: 12px; padding: 4px; background-color: royalblue; color: white;text-decoration: none ;">上一页</a>';
                if(pageNow - 3 > 0){
                    html = html + '<a href="javascript:searchHosts(1)" style="font-size: 14px; padding: 2px;color: black;">1</a>';
                    html = html + '...';
                    for (let i =pageNow - 3 ; i<= pageNow + 3; i++){
                        if (i == result['pages']){
                            break;
                        }
                        html = html + '<a href="javascript:searchHosts('+ i +')" style="font-size: 14px; padding: 2px;color: black;">'+i+'</a>';
                    }
                    html = html + '...';
                    html = html + '<a href="javascript:searchHosts('+ result['pages'] +')" style="font-size: 14px; padding: 2px;color: black;">'+result['pages']+'</a>';
                }else{
                    let j = false;
                    let e = false;
                    for (var i =1 ; i< pageNow + 3; i++){
                        html = html + '<a href="javascript:searchHosts('+ i +')" style="font-size: 14px; padding: 2px;color: black;">'+i+'</a>';
                        if(i >= result['pages'] - 1){
                            j = true;
                        }
                        if (i == result['pages']){
                            e = true;
                            break;
                        }
                    }
                    if(!j){
                        html = html + '...';
                    }
                    if (!e){
                        html = html + '<a href="javascript:searchHosts('+ result['pages'] +')" style=" padding: 2px;color: black;">'+result['pages']+'</a>';
                    }
                }

                html = html + '<a href="javascript:searchHosts('+ downpage +')" style=" font-size: 12px; padding: 4px; background-color: royalblue; color: white;text-decoration: none ;">下一页</a>';
                // html = html + '<label>当前页:'+pageNow+'</label>';
                $("#showpages").html(html);

            }else {
                alert("监控节点信息查询失败");
            }

        }
    })
}

function searchGroups(){

}


function changeStatus(status, target_id, group_id, group_name, page){
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
                searchHosts(parseInt(page));
                // window.location.reload();
            }else{
                alert(result.msg);
            }
        }
    })
}

function deleteTarget(target_id, group_id, group_name){
    jQuery.ajax({
        type: "POST",
        url: "deleteTarget",
        dataType: 'json',
        async: false,
        data: JSON.stringify({id: target_id, group_id: group_id, group_name: group_name}),
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            if(result.code == 0){
                alert("删除成功");
                window.location.reload();
            }else{
                alert(result.msg);
            }
        }
    })
}


function selectHost(){
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
            $("#groups").html("");
            if(code == 0) {
                for(var index in result.group){
                    let groupID = result.group[index].id;
                    let groupName = result.group[index].name;
                    html = html + '<option style="text-align: center" value="' + groupID + '">' + groupName + '</option>';
                }
                $("#groups").html(html);
            }else {
                alert("组信息加载错误");
            }
        }
    })
}