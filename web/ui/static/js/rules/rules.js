function updateRules(name){
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