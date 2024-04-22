$(document).ready(function() {
    $("#register-form").on("submit", function(e) {
        e.preventDefault();
        var username = $("#username").val();
        var password = $("#password").val();

        $.post("/signup", { username: username, password: password })
            .done(function() {
                window.location.href = "/signin";
            })
            .fail(function(err) {
                showError(err.responseJSON.error);
            });
    });

    function showError(message) {
        $("#error-text").text(message);
        $("#error-alert").show();
    }
});
