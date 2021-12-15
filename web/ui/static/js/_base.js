function reload(){
    jQuery.ajax({
        type: "POST",
        url: "-/reload",
        dataType: 'json',
        async: false,
        // data: JSON.stringify({id: target_id, group_id: group_id, status: status, group_name: group_name}),
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {
            alert(result.result);
        }
    })

}