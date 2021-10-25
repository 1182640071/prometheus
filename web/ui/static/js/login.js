function login(){
    var username = $("#username").val().trim();
    var password = $("#password").val().trim();

    if(username == "" || password == ""){
        $(".status-login").html("用户名、密码不能为空");
        return;
    }

    jQuery.ajax({
        type: "POST",
        url: "userAuthentication",
        contentType: 'application/json;charset=UTF-8',
        dataType: 'json',
        async: false,
        data : JSON.stringify({username: username, password: password}),
        error: function () {
            alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
            return false;
        },
        success: function (result) {

            if(result.code == 0){
                window.location.href = "/welcome";
            }else{
                alert("信息密码错误！");
            }
        }
    })
}