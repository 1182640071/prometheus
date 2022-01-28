function updateRules(name){

    let warn = warning()
    if (!warn) {
        return
    }

    let value = $("#" + name + "-rules").html().trim();
    if (value == ""){
        alert("告警规则不能为空！");
        return
    }

    jQuery.ajax({
        type: "POST",
        url: "updateRules",
        dataType: 'json',
        data: JSON.stringify({name: name, value: value}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {
                alert("修改成功，点击reload重载后生效");
            }else{
                alert("修改失败");
            }
        }
    })

}


function deleteRules(name){

    let warn = warning()
    if (!warn) {
        return
    }

    jQuery.ajax({
        type: "POST",
        url: "deleteRules",
        dataType: 'json',
        data: JSON.stringify({name: name}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {
                alert("删除成功，点击reload重载后生效");
            }else{
                alert("删除失败");
            }
        }
    })

}


function openAddRuleWin(){
    let features = "height=500, width=800, top=100, left=100, toolbar=no, menubar=no, " +
        "scrollbars=no,resizable=no, location=no, status=no";  //设置新窗口的特性
    window.open("/openAddRuleWin", "newW", features, false);  //打开新窗口
}

function addRule(){
    alert(1)
    let name = $("#rule_name").html().trim();
    let exper = $("#rule_exper").html().trim();
    let rulefor = $("#rule_for").html().trim();
    let keyword = $("#rule_keyword").html().trim();
    let level = $("#rule_level").html().trim();
    let service = $("#rule_service").html().trim();
    let description = $("#rule_description").html().trim();

    if (name == "" || exper == "" || rulefor == "" || keyword == "" || level == "" || service == "" || description == ""){
        alert("参数不能有空值!");
        return;
    }

    let warn = warning();
    if (!warn) {
        return;
    }

    jQuery.ajax({
        type: "POST",
        url: "addRules",
        dataType: 'json',
        data: JSON.stringify({name: name, exper: exper, rulefor: rulefor, keyword: keyword, level: level, service: service, description: description}),
        async: false,
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            let code = result.code;
            if(code == 0) {
                alert("添加成功，点击reload重载后生效");
            }else{
                alert("添加失败");
            }
        }
    })


}
