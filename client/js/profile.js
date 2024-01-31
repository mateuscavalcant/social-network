
function loadPostsProfile() {
  fetch("/profile", {
    method: "POST"
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      var userHeaderDetailsHTML = `<div class="name">
        <header>
          <div class="user-name">
            <p>${data.profile.name}</p>
          </div>
          <div class="posts-count">
            <p>${data.profile.countposts} posts</p>
          </div>
        </header>
        <main>
        </main>
      </div>`;

      var $userHeaderDetails = $(userHeaderDetailsHTML);
      $("#profile-header-container").append($userHeaderDetails);

      var imageData = data.icon;
      var imageUrl = 'data:image/jpeg;base64,' + imageData;
      $("#profile-icon").attr("src", imageUrl);
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
        <div class="create-btn">
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
            <button id="edit-profile-btn">Edit Profile</button>
          </footer>
        </div>
      </div>`;
      var $userDetails = $(userDetailsHTML);
      $userDetails.find("#edit-profile-btn").click(function () {
        window.location.href = "/edit-profile";
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
              <img src="client/public/images/message.png" alt="Comment" class="comment-button" data-post-id="${post.postID}">
              <img src="client/public/images/repost.png" alt="Repost">
              <img src="client/public/images/like.png" alt="Like" class="like-button" data-post-id="${post.postID}">
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
      console.error('Error:', error);
    });
}

function handleHome(event) {
  event.preventDefault();
  window.location.replace("/home");
}

$(document).ready(function () {
  loadPostsProfile();
  $("#home-btn").click(handleHome);

});