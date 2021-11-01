    function submitGroup(){
        var name = $("#groupName").val();
        var interval = $("#groupInterval").val();
        var uri = $("#groupUri").val();

        var type = $("#requestType").val();
        var check = $("#crtCheck").val();
        var match = $("#groupMatch").val();
        // var federal = $("#gfederal").val();
        var federal = "";
        var honor_labels = $("#honor_labels").val();


        // alert(name + "," + interval + "," + uri + "," + type + "," + check + "," + match + "," + federal + "," + honor_labels )

        if('' == name.trim()){
            alert('名字不能为空!');
            return;
        }

        //检查输入时间是否正确
        var reg = /^[0-9]*[s|m]$/;
        if ('' != interval && !reg.test(interval)){
            alert("请检查监控,匹配,超时时间格式!(时间+s,m,如3s或4m)");
            return;
        }

        reg = /^\/.*/;
        if ('' != uri.trim() && !reg.test(uri)){
            alert("输入的uri格式不对!(/uri)");
            return;
        }

        jQuery.ajax({
            type: "POST",
            url: "addGroupConfiguration",
            dataType: 'json',
            data : JSON.stringify({"name":name, "finterval":interval, "metrics_path":uri, "scheme":type, "insecure_skip_verify": check, "match_regulation": match, "federalid": federal, "honor_labels": honor_labels}),
            // data : {"name":name, "finterval":interval, "metrics_path":uri, "scheme":type, "insecure_skip_verify": check, "match_regulation": match, "federalid": federal, "honor_labels": honor_labels},

            async: false,
            error: function () {
                alert("操作失败,请稍等片刻重新尝试,如仍有问题请联系管理员......");
                return false;
            },
            success: function (result) {
                if (result.code == 0) {
                    alert("添加成功");
                }else{
                    alert(result.msg);
                }
            }
        })
    }