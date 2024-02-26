function setupWebSocket() {
    var pathParts = window.location.pathname.split("/");
    var username = pathParts[pathParts.length - 1];
    var wsURL = "ws://" + window.location.host + "/websocket/" + username;
    var ws = new WebSocket(wsURL);

    ws.onopen = function () {
        console.log("WebSocket connection established.");
    };

    ws.onmessage = function (event) {


        loadMessages();
    };

    ws.onclose = function () {
        console.log("WebSocket connection closed. Reconnecting...");
        setTimeout(setupWebSocket, 1000);
        loadMessages();
    };
}

var autoScroll = true;

function loadMessages() {
    var pathParts = window.location.pathname.split("/");
    var username = pathParts[pathParts.length - 1];
    fetch("/chat/" + username, {
        method: "POST"
    })
        .then(response => response.json())
        .then(data => {
            $("#messages-container").empty();

            data.messages.forEach(post => {
                var imageData = post.iconbase64;
                var imageUrl = 'data:image/jpeg;base64,' + imageData;

                var postHTML = `
        <div class="${post.messagesession ? 'message-session' : 'message-to'}">
            <header>
                <div class="message-box">
                    <img src="${imageUrl}" class="profile-icon">
                    <div class="post-content">
                        <p>${post.content}</p>
                    </div>
                </div>
            </header>
            <div class="user-name">
                <p class="name-user${post.postid}">${post.createdbyname}</p>
            </div>
        </div>`;
                var $post = $(postHTML);
                var username = post.createdby;
                $post.find(".user-name").on("click", function () {
                    window.location.href = "/" + username;
                });
                $("#messages-container").prepend($post);
            });

            if (autoScroll) {
                window.scrollTo(0, document.body.scrollHeight);
            } else {
                document.documentElement.scrollTop = scrollOffset;
            }

        })
        .catch(error => {
            console.error('Erro ao carregar posts:', error);
        });
}

function handleCreatePostFormSubmit(event) {
    event.preventDefault();
    var pathParts = window.location.pathname.split("/");
    var username = pathParts[pathParts.length - 1];

    var formData = $(this).serialize();
    fetch("/create-message/" + username, {
        method: "POST",
        body: formData,
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        }
    })
        .then(response => response.json())
        .then(() => {
            $("#message-form")[0].reset();
            $("#create-message-form-container").hide();

            loadMessages();
        })
        .catch(error => {
            console.error(error);
        });
}

document.addEventListener("DOMContentLoaded", function () {
    loadMessages();
    setupWebSocket();
    $("#message-form").submit(handleCreatePostFormSubmit);
});

window.addEventListener("scroll", function () {
    if (window.scrollY < document.body.scrollHeight - window.innerHeight) {
        autoScroll = false;
    } else {
        autoScroll = true;
    }
});
