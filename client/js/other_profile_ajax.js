function loadProfile(username) {
    $.ajax({
        url: "/profile/" + username,
        method: "POST",
        success: function (response) {
            var userHeaderDetailsHTML = '<div class="name">' +
                '<header>' +
                '<div class="user-name">' +
                '<p>' + response.profile.name + '</p>' +
                '</div>' +
                '<div class="posts-count">' +
                '<p>' + response.profile.countposts + ' posts</p>' +
                '</div>' +
                '<main>' +
                '</main>' +
                '</div>';
            var $userHeaderDetails = $(userHeaderDetailsHTML);
            $("#profile-header-container").append($userHeaderDetails);
            if (response.profile.followby) {
                var userDetailsHTML = '<div class="user">' +
                    '<header>' +
                    '<img src="client/public/images/golang-icon2.jpeg" class="user-icon">' +
                    '<div class="user-title">' +
                    '<p>@' + response.profile.username + '</p>' +
                    '</div>' +
                    '<div class="user-bio">' +
                    '<p>' + response.profile.bio + '</p>' +
                    '</div>' +
                    '</header>' +
                    '<main>' +
                    '<div class="user-followby">' +
                    '<p>' + response.profile.followbycount + '</p>' +
                    '<p id="followers-name">Followers</p>' +
                    '</div>' +
                    '<div class="user-followto">' +
                    '<p>' + response.profile.followtocount + ' </p>' +
                    '<p id="following-name">Following</p>' +
                    '</div>' +
                    '</main>' +
                    '<footer>' +
                    '<div class="create-btn">' +
                    '<button id="following-btn">Following</button>' +
                    '<button id="follow-btn" style="display: none;">Follow</button>' +
                    '</div>' +
                    '</footer>' +
                    '</div>';
                var $userDetails = $(userDetailsHTML);
                $(document).ready(function () {
                    $("#following-btn").click(function () {
                        $("#following-btn").prop("disabled", true);
                        var pathParts = window.location.pathname.split("/");
                        var user_follow_to = pathParts[pathParts.length - 1];
                        $.ajax({
                            type: "POST",
                            url: "/unfollow",
                            data: {
                                username: user_follow_to,
                            },
                            success: function (response) {
                                $("#following-btn").hide();
                                $("#follow-btn").show();

                                $("#following-btn").text("Follow").prop("disabled", false);
                                console.log("Unfollowed successfully:", response);
                            },
                            error: function () {
                                $("#follow-btn").prop("disabled", false);
                                console.log("Error following user");
                            }
                        });
                    });
                });
                $("#user-profile-container").append($userDetails);
            } else {
                var userDetailsHTML = '<div class="user">' +
                    '<header>' +
                    '<img src="client/public/images/golang-icon2.jpeg" class="user-icon">' +
                    '<div class="user-title">' +
                    '<p>@' + response.profile.username + '</p>' +
                    '</div>' +
                    '<div class="user-bio">' +
                    '<p>' + response.profile.bio + '</p>' +
                    '</div>' +
                    '</header>' +
                    '<main>' +
                    '<div class="user-followby">' +
                    '<p>' + response.profile.followbycount + '</p>' +
                    '<p id="followers-name">Followers</p>' +
                    '</div>' +
                    '<div class="user-followto">' +
                    '<p>' + response.profile.followtocount + ' </p>' +
                    '<p id="following-name">Following</p>' +
                    '</div>' +
                    '</main>' +
                    '<footer>' +
                    '<div class="create-btn">' +
                    '<button id="follow-btn">Follow</button>' +
                    '<button id="following-btn" style="display: none;">Following</button>' +
                    '</div>' +
                    '</footer>' +
                    '</div>';
                var $userDetails = $(userDetailsHTML);
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
                                $("#follow-btn").hide();
                                $("#following-btn").show();

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
                $("#user-profile-container").append($userDetails);
            }
            $("#posts-container").empty();
            response.posts.forEach(function (post) {
                var postHTML = '<div class="post">' +
                    '<header>' +
                    '<img src="client/public/images/golang-icon2.jpeg" class="profile-icon">' +
                    '<div class="post-title">' +
                    '<div class="user-name-post">' +
                    '<p class="name-user' + post.postID + '">' + post.createdbyname + '</p>' +
                    '</div>' +
                    '<div class="user-username">' +
                    '<p class="username-user' + post.postID + '">@' + post.createdby + '</p>' +
                    '</div>' +
                    '</div>' +
                    '</header>' +
                    '<main>' +
                    '<div class="post-content">' +
                    '<p>' + post.content + '</p>' +
                    '</div>' +
                    '<div class="post-links">' +
                    '<img src="" alt="Comentário" class="comment-button" data-post-id="' + post.postID + '">' +
                    '<img src="" alt="Comment">' +
                    '<img src="" alt="Curtir" class="like-button" data-post-id="' + post.postID + '">' +
                    '</div>' +
                    '</main>' +
                    '<footer>' +
                    '</footer>' +
                    '</div>';
                var $post = $(postHTML);
                $post.find(".like-button").on("click", function () {
                    var postID = $(this).data("post-id");
                    // Alterar a imagem para a nova imagem de curtida
                    $(this).attr("src", "");
                });
                $("#posts-container").prepend($post);
            });
        },
        error: function (xhr, status, error) {
            console.error(error);
        }
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

    $("#follow-btn").click(followUser(username))


});