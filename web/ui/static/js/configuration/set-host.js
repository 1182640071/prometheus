function addLabel() {

    let key = $("#label-key").val().trim();
    let value = $("#label-value").val().trim();

    if(key == "" || value == ""){
        alert("label的key和value不能为空");
        return;
    }

    let label = "<div id='" + key + "'>";
    label = label + '<input id="' + key + '-key" type="button" style="float: left; margin-left: 10px; height: 25px; margin-top: 10px; border: 1px solid lightgrey; border-radius: 0" disabled="disabled" value="' + key + '">'
    label = label + '<input id="' + key + '-value" type="button" style="float: left; margin-left: 10px; height: 25px; margin-top: 10px; border: 1px solid lightgrey; display: none; border-radius: 0" disabled="disabled" value="' + value + '">'
    label = label + '<input type="button" onClick="closeLabel(\'' + key + '\')" style="float: left; height: 25px; width: 25px; margin-top: 10px; border: 1px solid lightgrey; border-left: 0; border-radius: 0" value="×">';
    label = label + "</div>";

    $("#extra_label").append(label);

    $("#label-key").val("");
    $("#label-value").val("");

    var labels = $("#labels").html();
    $("#labels").html(labels + "|!" + key);
}

function closeLabel(node){
    $("#" + node).html("");
    var labels = $("#labels").html();
    $("#labels").html(labels.replace("|!" + node, ""));
}


function submitHosts(){

    let group = $("#groups").val();
    let name = $("#host-name").val();
    let ip = $("#host-ip").val();
    let port = $("#host-port").val();

    if("" == group){
        alert("请选择属组!");
        return;
    }

    if("" == name){
        alert("请填写主机名称!");
        return;
    }


    if("" == ip){
        alert("请填写主机IP!");
        return;
    }

    if("" == port){
        alert("请填写主机端口!");
        return;
    }

    let extraLabels = {}

    let labelKeys = $("#labels").html();
    let keys = labelKeys.split("|!");
    for(let index in keys){
        let key = keys[index].trim();
        if(key == ""){
            continue;
        }
        let value = $("#" + key + "-value").val();
        extraLabels[key] = value;
    }
    extraLabels["name"] = name + "_" + ip;

    jQuery.ajax({
        type: "POST",
        url: "addHostConfig",
        dataType: 'json',
        contentType: 'json/application',
        data : JSON.stringify({name: name, group_id: group, ip: ip, port: port, label: JSON.stringify(extraLabels).toString()}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            if(result.code == 0){
                alert("添加成功");
            }else{
                alert(result.msg);
            }
        }
    })

}

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