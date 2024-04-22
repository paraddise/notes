var logintab = document.getElementById("loginTab");
if (logintab != null) {
    logintab.addEventListener("click", function() {
        document.getElementById("loginForm").style.display = "block";
        document.getElementById("registerForm").style.display = "none";
        document.getElementById("loginTab").classList.add("is-active");
        document.getElementById("registerTab").classList.remove("is-active");
    });
}

var regtab = document.getElementById("registerTab");
if (regtab != null) {
    regtab.addEventListener("click", function() {
        document.getElementById("registerForm").style.display = "block";
        document.getElementById("loginForm").style.display = "none";
        document.getElementById("registerTab").classList.add("is-active");
        document.getElementById("loginTab").classList.remove("is-active");
    });
}

function loginUser() {
    var username = document.getElementById("loginUsername").value;
    var password = document.getElementById("loginPassword").value;
    var url = `/api/login?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`;
    fetch(url, {method: "POST"})
        .then(response => response.json())
        .then(data => {
            if (data.Success) {
                window.location.href = "clicker.html";
            } else {
                alert("Ошибка входа: Неправильный логин или пароль");
            }
        });
}

function registerUser() {
    var username = document.getElementById("registerUsername").value;
    var password = document.getElementById("registerPassword").value;
    if (username.length < 10 || password.length < 10) {
        alert("Логин и пароль не могут быть меньше 10 символов");
        return;
    }
    var url = `/api/register?username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`;
    fetch(url, {method: "POST"})
        .then(response => response.json())
        .then(data => {
            if (data.Success) {
                window.location.reload();
            } else {
                alert("Ошибка регистрации: " + data.Error);
            }
        });
}

function makeClick() {
    fetch('/api/step', {method: "POST"})
        .then(response => response.json())
        .then(data => {
            updateCounter()
        });
}

function buyColor(colorName) {
    fetch(`/api/buyColor?color=${encodeURIComponent(colorName)}`, {method: "POST"})
        .then(response => response.json())
        .then(data => {
            if (!data.Success) {
                alert("Ошибка при смене цвета: " + data.Error);
            }
            updateCounter()
        });
}

function buyPremium() {
    fetch('/api/buyPremium', {method: "POST"})
        .then(response => response.json())
        .then(data => {
            if (!data.Success) {
                alert("Ошибка при покупке премиума: " + data.Error);
            }
            updateCounter()
        });
}

function getNoun(number, one, two, five) {
    let n = Math.abs(number);
    n %= 100;
    if (n >= 5 && n <= 20) {
        return "кликов";
    }
    n %= 10;
    if (n === 1) {
        return "клик";
    }
    if (n >= 2 && n <= 4) {
        return "клика";
    }
    return "кликов";
}

function updateCounter() {
    if (document.location.href.indexOf("clicker") > 0 && document.cookie.indexOf("session") < 0) {
        window.location.href = "/";
    }

    fetch('/api/getInfo', {method: "POST"}).then(response => response.json()).then(data => {
        var text = "";

        if (!data.Success) {
            text = "Ошибка в получении данных, попробуйте обновить страницу или перезайти в аккаунт"
        } else {
            if (data.Prize != "") {
                alert("Поздравляем, вы получаете приз: " + data.Prize);
                text = "Кажется, вы закатили камень."
            } else {
                var left = 413370 - data.Offset;
                text = "Вы сделали " + data.Offset + " " + getNoun(data.Offset) + ", до получения приза осталось еще " + left + " " + getNoun(data.Multiplier);
            }

            if (data.Color == 0) {
                document.getElementsByTagName("html")[0].style = "background: white"; 
                document.getElementsByTagName("body")[0].style = ""; 
            } else if (data.Color == 1) {
                document.getElementsByTagName("html")[0].style = "background: black";
                document.getElementsByTagName("body")[0].style = "color: white"; 
            } else if (data.Color == 2) {
                document.getElementsByTagName("html")[0].style = "background: aqua";
                document.getElementsByTagName("body")[0].style = ""; 
            }

            document.getElementById("clickMultiplicator").innerText = "Одно нажатие на кнопку дает " + data.Multiplier + " " + getNoun(data.Multiplier) + ".";
        }
        
        document.getElementById("clickCount").innerText = text;
    });
}

function logout() {
    document.cookie = 'session=; Max-Age=-99999999;'; 
    window.location.href = "/";
}

updateCounter();