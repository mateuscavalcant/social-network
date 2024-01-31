document.getElementById("login-form").addEventListener("submit", function (event) {
    event.preventDefault();
    var identifier = document.getElementById("identifier").value;
    var password = document.getElementById("password").value;
    var confirmPassword = document.getElementById("confirm_password").value;

    var formData = new FormData();
    formData.append("identifier", identifier);
    formData.append("password", password);
    formData.append("confirm_password", confirmPassword);

    fetch("/delete_account", {
        method: "POST",
        body: formData
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                // Exibe as mensagens de erro correspondentes nos campos do formulário
                document.getElementById("error-identifier").textContent = data.error.identifier;
                document.getElementById("error-password").textContent = data.error.password;
                document.getElementById("error-confirm_password").textContent = data.error.confirm_password;
            } else {
                console.log(data.message);
                window.location.replace("/login");
            }
        })
        .catch(error => {
            console.error(error);
        });
});