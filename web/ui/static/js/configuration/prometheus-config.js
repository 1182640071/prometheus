
// 获取当前prometheus的配置信息
$(document).ready(function () {
    jQuery.ajax({
        type: "GET",
        url: "getConfiguration",
        dataType: 'json',
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {

            var status = result.code;
            if(status == 0){
                $("#prometheusName").val(result.configuration.name);
                $("#prometheusInterval").val(result.configuration.finterval);
                $("#prometheusRInterval").val(result.configuration.rinterval);
                $("#prometheusAUrl").val(result.configuration.aurl);
                $("#prometheusRPath").val(result.configuration.rpath);
                $("#prometheusJPath").val(result.configuration.jpath);
                $("#prometheusTimeout").val(result.configuration.timeout);
            }else {
                alert('数据加载异常,请稍后尝试或联系管理员');
            }

        }
    })
});

function submitConfiguration() {
    var name = $("#prometheusName").val();
    var interval = $("#prometheusInterval").val();
    var RInterval = $("#prometheusRInterval").val();
    var AUrl = $("#prometheusAUrl").val();
    var RPath = $("#prometheusRPath").val();
    var JPath = $("#prometheusJPath").val();
    var timetOut = $("#prometheusTimeout").val();

    //检查名字中是否含有特殊字符
    if (checkQuote(name) || "" == name.trim()) {
        alert("输入名为空或存在特殊字符!");
        return;
    }

    //检查输入时间是否正确
    var reg = /^[0-9]*[s|m]$/;
    if (!reg.test(interval) || !reg.test(RInterval) || !reg.test(timetOut)) {
        alert("请检查监控,匹配,超时时间格式!(时间+s,m,如3s或4m)");
        return;
    }

    //检查url
    if (checkURL(AUrl)) {
        alert("url为空或参数格式存在问题!");
        return;
    }

    //检查路径是否为空
    if ("" == RPath || "" == JPath) {
        alert("路径参数不能为空");
        return;
    }

    var x = confirm("是否确认提交");
    if (x == false) {
        return;
    }

    jQuery.ajax({
        type: "POST",
        url: "submitConfiguration",
        contentType: 'application/json;charset=UTF-8',
        dataType: 'json',
        data : JSON.stringify({"name":name, "finterval":interval, "rinterval":RInterval, "aurl":AUrl, "rpath": RPath, "jpath": JPath, "timeout": timetOut}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            if(result.code == 0){
                alert("提交成功！");
            }else{
                alert(result.msg);
            }
        }
    })
}

//判断是否含有特殊字符
function checkQuote(str) {
    var items = new Array("~", "`", "!", "@", "#", "$", "%", "^", "&", "{", "}", "[", "]", "(", ")");
    items.push(":", ";", "'", "|", "\\", "<", ">", "?", "/", "<<", ">>", "||", "/", "*", ",", ".");
    //str = str.toLowerCase();
    for (var i = 0; i < items.length; i++) {
        if (str.indexOf(items[i]) >= 0) {
            return true;
        }
    }
    return false;
}

//检测url地址是否合法
function checkURL(str) {
    str = 'http://' + str;
    if (str.match(/(http[s]?|ftp):\/\/[^\/\.]+?..+\w$/i) == null) {
        return true
    }else {
        return false;
    }
}

function changeNode(node){
    if (node == "host"){
        document.getElementById("host").style.borderBottom="0px solid white";

        document.getElementById("group").style.borderBottom="1px solid lightgrey";
        document.getElementById("consul").style.borderBottom="1px solid lightgrey";

        $("#change_iframe").attr("src","/toHost");

    }else if(node == "group"){
        document.getElementById("group").style.borderBottom="0px solid white";

        document.getElementById("host").style.borderBottom="1px solid lightgrey";
        document.getElementById("consul").style.borderBottom="1px solid lightgrey";

        $("#change_iframe").attr("src","/toGroup");

    }else if(node == "consul"){
        document.getElementById("consul").style.borderBottom="0px solid white";

        document.getElementById("host").style.borderBottom="1px solid lightgrey";
        document.getElementById("group").style.borderBottom="1px solid lightgrey";

        $("#change_iframe").attr("src","toConsul");
    }
}

function showPConfig(){
    jQuery.ajax({
        type: "POST",
        url: "showPConfig",
        dataType: 'json',
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            if(result['status'] == 0){
                alert('数据错误,请检查配置信息!');
            }else{
                $("#content").html(result['baseconfig']);
                document.getElementById("showconfig").style.display="block";
                document.getElementById("backshow").style.display="block";
            }
        }
    })
}

function updatePConfig(){
    var x = confirm("是否确认刷新prometheus配置文件");
    if (x == false) {
        return;
    }

    jQuery.ajax({
        type: "GET",
        url: "updatePrometheusYmlConfig",
        dataType: 'json',
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            alert(result);
        }
    })

}


function closeBack(){
    document.getElementById("showconfig").style.display="none";
    document.getElementById("backshow").style.display="none";
}
