document.addEventListener("DOMContentLoaded", function () {
    function validateEmail() {
        var email = document.getElementById("email").value;
        var formData = new FormData();

        formData.append("email", email);

        fetch("/validate-email", {
            method: "POST",
            body: formData
        })
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    document.getElementById("error-email").textContent = data.error.email;
                } else {
                    document.getElementById("error-email").textContent = "";
                }
            })
            .catch(function (error) {
                console.error(error);
            });
    }
    document.getElementById("email").addEventListener("input", function () {
        validateEmail();
    });
});

document.getElementById("signup-form").addEventListener("submit", function (event) {
    event.preventDefault();
    var username = document.getElementById("username").value;
    var name = document.getElementById("name").value;
    var email = document.getElementById("email").value;
    var password = document.getElementById("password").value;
    var confirmPassword = document.getElementById("confirm_password").value;

    var formData = new FormData();
    formData.append("username", username);
    formData.append("name", name);
    formData.append("email", email);
    formData.append("password", password);
    formData.append("confirm_password", confirmPassword);


    fetch("/signup", {
        method: "POST",
        body: formData
    })
        // Exibe as mensagens de erro correspondentes nos campos do formulário
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                document.getElementById("error-username").textContent = data.error.username;
                document.getElementById("error-name").textContent = data.error.name;
                document.getElementById("error-email").textContent = data.error.email;
                document.getElementById("error-password").textContent = data.error.password;
                document.getElementById("error-confirm-password").textContent = data.error.confirm_password;

            } else {
                console.log(data.message);
                window.location.replace("/login");

            }
        })
        .catch(error => {
            console.error(error);
        });
});     