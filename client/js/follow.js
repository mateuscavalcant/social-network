document.addEventListener("DOMContentLoaded", function () {
  document.getElementById("follow-btn").addEventListener("click", function () {
    document.getElementById("follow-btn").disabled = true;
    var pathParts = window.location.pathname.split("/");
    var user_follow_to = pathParts[pathParts.length - 1];

    fetch("/follow", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        username: user_follow_to
      })
    })
      .then(response => {
        if (response.ok) {
          document.getElementById("follow-btn").textContent = "Following";
          document.getElementById("follow-btn").disabled = false;
          console.log("Followed successfully");
        } else {
          throw new Error("Error following user");
        }
      })
      .catch(error => {
        document.getElementById("follow-btn").disabled = false;
        console.error(error);
      });
  });
});
