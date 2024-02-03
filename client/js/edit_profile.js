function handleProfile() {
    fetch("/profile", {
        method: "POST"
    })
        .then(response => response.json())
        .then(data => {
            var username = data.profile.username;
            window.location.replace("/" + username);
        })
        .catch(error => {
            console.error(error);
        });
}

document.getElementById("edit-profile-form").addEventListener("submit", function (event) {
    event.preventDefault();
    const fileInput = document.getElementById('icon');
    const file = fileInput.files[0];
    const name = document.getElementById("name").value;
    const bio = document.getElementById("bio").value;

    const formData = new FormData();
    formData.append("icon", file);
    formData.append("name", name);
    formData.append("bio", bio);


    fetch("/edit-profile", {
        method: "POST",
        body: formData
    })
        // Exibe as mensagens de erro correspondentes nos campos do formulário
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                document.getElementById("error-icon").textContent = data.error.icon;
                document.getElementById("error-name").textContent = data.error.name;
                document.getElementById("error-bio").textContent = data.error.bio;
            } else {
                console.log(data.message);
                handleProfile();
            }
        })
        .catch(error => {
            console.error(error);
        });
});

