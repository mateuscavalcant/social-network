$(document).ready(function () {
  $("#follow-btn").click(function () {
    $("#follow-btn").prop("disabled", true);
    var pathParts = window.location.pathname.split("/");
    var user_follow_to = pathParts[pathParts.length - 1];
    $.ajax({
      type: "POST",
      url: "/follow",
      data: {
        username: user_follow_to,
      },
      success: function (response) {
        $("#follow-btn").text("Following").prop("disabled", false);
        console.log("Followed successfully:", response);
      },
      error: function () {
        $("#follow-btn").prop("disabled", false);
        console.log("Error following user");
      }
    });
  });
});

