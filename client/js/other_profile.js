function loadProfile(username) {
  fetch("/profile/" + username, {
    method: "POST"
  })
    .then(response => response.json())
    .then(data => {
      var userHeaderDetailsHTML = `<div class="name">
              <header>
              <div class="user-name">
              <p>${data.profile.name}</p>
              </div>
              <div class="posts-count">
              <p>${data.profile.countposts} posts</p>
              </div>
              <main>
              </main>
              </div>`;
      var $userHeaderDetails = $(userHeaderDetailsHTML);
      $("#profile-header-container").append($userHeaderDetails);

      if (data.profile.followby) {
        var userDetailsHTML = `<div class="user">
                  <header>
                  <img src="client/public/images/golang-icon2.jpeg" class="user-icon">
                  <div class="user-title">
                  <p>@${data.profile.username}</p>
                  </div>
                  <div class="user-bio">
                  <p>${data.profile.bio}</p>
                  </div>
                  </header>
                  <main>
                  <div class="user-followby">
                  <p>${data.profile.followbycount}</p>
                  <p id="followers-name">Followers</p>
                  </div>
                  <div class="user-followto">
                  <p>${data.profile.followtocount}</p>
                  <p id="following-name">Following</p>
                  </div>
                  </main>
                  <footer>
                  <div class="create-btn">
                  <button id="following-btn">Following</button>
                  <button id="follow-btn" style="display: none;">Follow</button>
                  </div>
                  </footer>
                  </div>`;
        var $userDetails = $(userDetailsHTML);

        $("#user-profile-container").append($userDetails);


      } else {
        var userDetailsHTML = `<div class="user">
                  <header>
                  <img src="client/public/images/golang-icon2.jpeg" class="user-icon">
                  <div class="user-title">
                  <p>@${data.profile.username}</p>
                  </div>
                  <div class="user-bio">
                  <p>${data.profile.bio}</p>
                  </div>
                  </header>
                  <main>
                  <div class="user-followby">
                  <p>${data.profile.followbycount}</p>
                  <p id="followers-name">Followers</p>
                  </div>
                  <div class="user-followto">
                  <p>${data.profile.followtocount}</p>
                  <p id="following-name">Following</p>
                  </div>
                  </main>
                  <footer>
                  <div class="create-btn">
                  <button id="follow-btn">Follow</button>
                  <button id="following-btn" style="display: none;">Following</button>
                  </div>
                  </footer>
                  </div>`;
        var $userDetails = $(userDetailsHTML);

        $("#user-profile-container").append($userDetails);

      }

      $("#posts-container").empty();
      data.posts.forEach(function (post) {
        var postHTML = `<div class="post">
                  <header>
                  <img src="client/public/images/golang-icon2.jpeg" class="profile-icon">
                  <div class="post-title">
                  <div class="user-name-post">
                  <p class="name-user${post.postID}">${post.createdbyname}</p>
                  </div>
                  <div class="user-username">
                  <p class="username-user${post.postID}">@${post.createdby}</p>
                  </div>
                  </div>
                  </header>
                  <main>
                  <div class="post-content">
                  <p>${post.content}</p>
                  </div>
                  <div class="post-links">
                  <img src="client/public/images/message.png" alt="Comentário" class="comment-button" data-post-id="${post.postID}">
                  <img src="client/public/images/repost.png" alt="Comment">
                  <img src="client/public/images/like.png" alt="Curtir" class="like-button" data-post-id="${post.postID}">
                  </div>
                  </main>
                  <footer>
                  </footer>
                  </div>`;
        var $post = $(postHTML);
        $post.find(".like-button").on("click", function () {
          var postID = $(this).data("post-id");
          // Alterar a imagem para a nova imagem de curtida
          $(this).attr("src", "");
        });
        $("#posts-container").prepend($post);
      });





      var followBtn = document.getElementById("follow-btn");
      var followingBtn = document.getElementById("following-btn");

      if (data.profile.followby) {
        followBtn.style.display = "none";
        followingBtn.style.display = "inline-block";

        followingBtn.addEventListener("click", function () {
          followingBtn.disabled = true;
          fetch("/unfollow", {
            method: "POST",
            body: JSON.stringify({
              username: username
            }),
            headers: {
              "Content-Type": "application/json"
            }
          })
            .then(response => response.json())
            .then(data => {
              followBtn.style.display = "inline-block";
              followingBtn.style.display = "none";
              followingBtn.textContent = "Follow";
              followingBtn.disabled = false;
              console.log("Unfollowed successfully:", data);
            })
            .catch(error => {
              console.error("Error unfollowing user:", error);
            });
        });
      } else {
        followBtn.style.display = "inline-block";
        followingBtn.style.display = "none";

        followBtn.addEventListener("click", function () {
          followBtn.disabled = true;
          fetch("/follow", {
            method: "POST",
            body: JSON.stringify({
              username: username
            }),
            headers: {
              "Content-Type": "application/json"
            }
          })
            .then(response => response.json())
            .then(data => {
              followBtn.style.display = "none";
              followingBtn.style.display = "inline-block";
              followingBtn.textContent = "Following";
              followBtn.disabled = false;
              console.log("Followed successfully:", data);
            })
            .catch(error => {
              console.error("Error following user:", error);
            });
        });
      }

      // Código para carregar os posts omitido para brevidade...
    })
    .catch(error => {
      console.error("Error loading profile:", error);
    });
}

function handleHome(event) {
  event.preventDefault();
  window.location.replace("/home");
}

document.addEventListener("DOMContentLoaded", function () {
  var pathParts = window.location.pathname.split("/");
  var username = pathParts[pathParts.length - 1];
  loadProfile(username);

  document.getElementById("home-btn").addEventListener("click", handleHome);
});