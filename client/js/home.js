function loadPosts() {
  fetch("/feed", {
    method: "POST"
  })
    .then(response => response.json())
    .then(data => {
      $("#posts-container").empty();

      data.posts.forEach(post => {
        var imageData = post.iconbase64;
        var imageUrl = 'data:image/jpeg;base64,' + imageData;
        var postHTML = `<div class="post">
              <header>
                  <img id="profile-icon" src="${imageUrl}" class="profile-icon">
                  <div class="post-title">
                      <div class="user-name">
                          <p class="name-user${post.postid}">${post.createdbyname}</p>
                      </div>
                      <div class="user-username">
                          <p class="username-user${post.postid}">@${post.createdby}</p>
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
              <footer></footer>
          </div>`;
        var $post = $(postHTML);
        var username = post.createdby;
        $post.find(".user-name").on("click", function () {
          window.location.href = "/" + username;
        });
        $("#posts-container").prepend($post);
      });
    })
    .catch(error => {
      console.error(error);
    });
}

function handleCreatePostFormSubmit(event) {
  event.preventDefault();

  var formData = $(this).serialize();
  fetch("/create-post", {
    method: "POST",
    body: formData,
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  })
    .then(response => response.json())
    .then(() => {
      $("#post-form")[0].reset();
      $("#create-post-form-container").hide();

      loadPosts();
    })
    .catch(error => {
      console.error(error);
    });
}

function handleProfile(event) {
  event.preventDefault();
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

document.addEventListener("DOMContentLoaded", function () {
  loadPosts();
  $("#post-form").submit(handleCreatePostFormSubmit);
  $("#profile-btn").click(handleProfile);
});