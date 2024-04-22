$(document).ready(function () {
    $("#login-form").on("submit", function (e) {
        e.preventDefault();
        var username = $("#username").val();
        var password = $("#password").val();

        $.ajax({
            url: "/signin",
            method: 'POST',
            data: { username: username, password: password },
            xhrFields: {
                withCredentials: true
            },
            success: function () {
                window.location.href = "/";
            },
            error: function (err) {
                showError(err.responseJSON.error);
            }
        });
    });

    function showError(message) {
        $("#error-text").text(message);
        $("#error-alert").show();
    }
});
