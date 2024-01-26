$('.ui.checkbox')
    .checkbox()
;

$('.http.ui.toggle.checkbox').checkbox({
    onChecked: function () {
        $(".http.two.fields").attr("style", "display: block;")
        config.Http.Enabled = true
        setValue(config)
    },
    onUnchecked: function () {
        $(".http.two.fields").attr("style", "display: none;")
        config.Http.Enabled = false
    }
});

$('.socks.ui.toggle.checkbox').checkbox({
    onChecked: function () {
        $(".socks.two.fields").attr("style", "display: block;")
        config.Socks5.Enabled = true
        setValue(config)
    },
    onUnchecked: function () {
        $(".socks.two.fields").attr("style", "display: none;")
        config.Socks5.Enabled = false
    }
});

sleep = function (fun, time) {
    setTimeout(() => {
        fun();
    }, time);
}

function post() {
    const data = {
        ProxyAddr: $("#proxyAddr").val().split("\n"),
        Control: {
            ConfigAddr: $("#controlAddr").val(),
            LogPath: $("#controlLog").val(),
            TorEnable:  $(".tor.ui.toggle.checkbox").hasClass("checked"),
        },
        Http: {
            Enabled: $(".http.ui.toggle.checkbox").hasClass("checked"),
            ListenAddr: $("#httpAddr").val(),
        },
        Socks5: {
            Enabled: $(".socks.ui.toggle.checkbox").hasClass("checked"),
            ListenAddr: $("#socksAddr").val(),
        }
    };
    $.ajax({
        type: 'POST',
        url: "/",
        data: JSON.stringify(data),
        error: errorMessage,
        success: successMessage,
        dataType: "json",
        contentType: "application/json"
    });

}

function setValue(config) {
    let proxy = "";
    config.ProxyAddr.forEach((e) => {
        proxy += e + "\n"
    })
    proxy = proxy.substring(0, proxy.length - 1)
    $('#proxyAddr').text(proxy)
    console.log(config.Control.ConfigAddr)
    $('#controlAddr').attr("value", config.Control.ConfigAddr)
    $('#controlLog').attr("value", config.Control.LogPath)
    if (config.Control.TorEnable) {
        $('.tor.ui.toggle.checkbox').addClass("checked")
        $("#torEnable").attr("checked", "checked")
    }
    if (config.Http.Enabled) {
        $('.http.ui.toggle.checkbox').addClass("checked")
        $(".http.two.fields").attr("style", "display: block;")
        $("#httpEnable").attr("checked", "checked")
        $("#httpAddr").attr("value", config.Http.ListenAddr)
    }
    if (config.Socks5.Enabled) {
        $('.socks.ui.toggle.checkbox').addClass("checked", "check")
        $(".socks.two.fields").attr("style", "display: block;")
        $("#socksEnable").attr("checked", "checked")
        $("#socksAddr").attr("value", config.Socks5.ListenAddr)
    }
}

function successMessage(e) {
    setValue(e)
    $(".ui.success.icon.message").removeClass("hidden");
    setTimeout(() => {
        $(".ui.success.icon.message").closest(".message").transition("fade");
    }, 2000);
}

function errorMessage(e) {
    $(".error.message.info").text(e.status + ":" + e.responseText)
    $(".ui.error.icon.message").removeClass("hidden");
    setTimeout(() => {
        $(".ui.error.icon.message").closest(".message").transition("fade");
    }, 2000);
}