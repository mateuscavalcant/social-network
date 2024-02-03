document.addEventListener("DOMContentLoaded", function () {
  var pathParts = window.location.pathname.split("/");
  var username = pathParts[pathParts.length - 1];
  loadProfile(username);

  document.getElementById("home-btn").addEventListener("click", handleHome);
  // Após a conclusão da busca pelos dados do perfil
  fetch("/profile/" + username, {
    method: "POST"
  })
    .then(response => response.json())
    .then(data => {
      // Atribuir o nome e o nome de usuário ao título
      document.title = `${data.profile.name} / (@${data.profile.username})`;
    })
    .catch(error => {
      console.error('Error fetching profile data:', error);
    });
});

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
        </header>
      </div>`;

      var $userHeaderDetails = $(userHeaderDetailsHTML);
      $("#profile-header-container").append($userHeaderDetails);

      var imageData = data.icon;
      var imageUrl = 'data:image/jpeg;base64,' + imageData;

      var userDetailsHTML = `<div class="user">
        <header>
          <img id="profile-icon" src="${imageUrl}" class="user-icon">
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
      if (data.profile.followby) {
        // Se o usuário atual estiver seguindo o usuário-alvo
        $userDetails.find("#follow-btn").hide(); // Esconda o botão Follow
        $userDetails.find("#following-btn").show(); // Mostre o botão Following
      } else {
        // Se o usuário atual não estiver seguindo o usuário-alvo
        $userDetails.find("#follow-btn").show(); // Mostre o botão Follow
        $userDetails.find("#following-btn").hide(); // Esconda o botão Following
      }

      $(document).ready(function () {
        $userDetails.find("#follow-btn").click(function () {
          $(this).prop("disabled", true);
          var pathParts = window.location.pathname.split("/");
          var user_follow_to = pathParts[pathParts.length - 1];

          $.ajax({
            type: "POST",
            url: "/follow",
            data: {
              username: user_follow_to,
            },
            success: function (response) {
              $(this).hide();
              $(this).siblings("#following-btn").show();

              $(this).text("Following").prop("disabled", false);
              console.log("Followed successfully:", response);
              location.reload();
            }.bind(this),
            error: function () {
              $(this).prop("disabled", false);
              console.log("Error following user");
            }.bind(this)
          });
        });

        $userDetails.find("#following-btn").click(function () {
          $(this).prop("disabled", true);
          var pathParts = window.location.pathname.split("/");
          var user_follow_to = pathParts[pathParts.length - 1];
          $.ajax({
            type: "POST",
            url: "/unfollow",
            data: {
              username: user_follow_to,
            },
            success: function (response) {
              $(this).hide();
              $(this).siblings("#follow-btn").show();

              $(this).text("Follow").prop("disabled", false);
              console.log("Unfollowed successfully:", response);
              location.reload();
            }.bind(this),
            error: function () {
              $(this).siblings("#follow-btn").prop("disabled", false);
              console.log("Error following user");
            }.bind(this)
          });
        });
      });
      $("#user-profile-container").append($userDetails);

      $("#posts-container").empty();
      data.posts.forEach(function (post) {
        var postHTML = `<div class="post">
          <header>
            <img id="profile-icon" src="${imageUrl}" class="profile-icon">
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
        $("#posts-container").prepend($post);
      });
    })
    .catch(error => {
      console.error("Error loading profile:", error);
    });
}

function handleHome(event) {
  event.preventDefault();
  window.location.replace("/home");
}

