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
        <main>
        </main>
      </div>`;
      var $userHeaderDetails = $(userHeaderDetailsHTML);
      $("#profile-header-container").append($userHeaderDetails);

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
            <p>${data.profile.followtocount} </p>
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
              <img src="client/public/images/message.png" alt="Comment" class="comment-button" data-post-id="${post.postID}">
              <img src="client/public/images/repost.png" alt="Repost">
              <img src="client/public/images/like.png" alt="Like" class="like-button" data-post-id="${post.postID}">
            </div>
          </main>
          <footer>
          </footer>
        </div>`;
        var $post = $(postHTML);
        $post.find(".like-button").on("click", function () {
          var postID = $(this).data("post-id");
          // Alterar a imagem para a nova imagem de curtida
          $(this).attr("src", "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACQAAAAkCAYAAADhAJiYAAAAAXNSR0IArs4c6QAAAnBJREFUWEft1t+rDGEYB/Dvd3Zn2yNhZ6ToJEWhzp2kKMyk5IILbtwoodw4F4oSLpQokRJulKSUC8UFyc3Z2Vw5f8CpLTdbpyS1s3I4aubsPNpTK2fNj3f2naOjzl7O+7zP83mf931nh1hiPy4xD5ZBWTui1CEB2KnsHRMp7zKAEQqb3bmwaeP9dFwBH3s2GuXqVqFsi4CfZPShFjSmCIgWaAa714Vm9Z4AR0isGEwmIpNGKMdraLR6Y1/hbu6aeE5iZ0zsDME35ZDjqzDRToIldsgvuYdh4AkIO21VIvIdxGlDUI2AhyRXpsfjCyOcsbr113FxsaBO2b0qBq5ntVdrPIouWXONW4M5/gL5JfcQSnirVUxxshHJgTVz3sSf4QtAgh2mX1ndIrBBMadmmLRqwdotxItuP9ECUMd0Tgn5WLNKzunRCStoPIsFtSvOO4IHc2bUChfglR3UjyaBpgmOalXIOVkgTTvwtseDTGeW5EjOnFrhAvlsB976pA59Ivh7UKuS6mSRKSv0xmJBvuk0QO5TzVVEnEBe2oF3LP6W/YsX4sAqROSsHXqP4rcM+0dh8iPJahGrz8wh0pFwdpONyW+xoN5Dv+xchsEbmckKCKBEJ2th42nim7o/0DaduyTPF1AzJYXctALvymBA8r+96d4HcW5RUIIHVlgfj8ud+oHmLwYqBdMDZn4xForKwCiB5g96EZ1SwCiDtFGKmFwgDdQdK6hfVL0cmWdoMJFfcW8DuKBYIBcmd4f6CEVUbszQoPntS+/UUBgtUApqaIw2KAalhSkE1EcJ8MMO6tcUD3tiWO5bplswa/4y6L/r0C89QuAlNQSzNwAAAABJRU5ErkJggg==");
        });
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
  var pathParts = window.location.pathname.split("/");
  var username = pathParts[pathParts.length - 1];
  loadProfile(username);

  $("#home-btn").click(handleHome);
});
