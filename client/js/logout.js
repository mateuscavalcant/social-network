function handleLogout(event) {
  event.preventDefault();
  fetch("/logout", {
    method: "POST"
  })
    .then(response => {
      if (response.ok) {
        window.location.replace("/login");
      } else {
        throw new Error('Logout request failed');
      }
    })
    .catch(error => {
      console.error('Error:', error);
    });
}

document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("logout-btn").addEventListener("click", handleLogout);
});
