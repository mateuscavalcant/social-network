package models

type User struct {
	ID              int    `json:"id"`
	Username        string `json:"username" binding:"required, min=4,max=32"`
	Name            string `json:"name" binding:"required, min=1,max=70"`
	Icon            []byte `json:"icon"`
	Bio             string `json:"bio" binding:"required, max=70"`
	Email           string `json:"email" binding:"required, email"`
	Password        string `json:"password" binding:"required, min=8, max=16"`
	ConfirmPassword string `json:"cpassword" binding:"required"`
}

type UserLogin struct {
	Credential string `json:"credential"`
	Password   string
}

type UserFollow struct {
	FollowID int `json:"follow-id"`
	FollowBy int `json:"follow-by"`
	FolloTo  int `json:"follow-to"`
}

type UserPost struct {
	PostID     int    `json:"post-id"`
	PostUserID int    `json:"post-user-id"`
	UserID     int    `json:"user-id"`
	Content    string `json:"content"`
	Icon       []byte `json:"icon"`
	IconBase64 string `json:"iconbase64"`
	CreatedBy  string `json:"createdby"`
	Name       string `json:"createdbyname"`
}

type UserProfile struct {
	ID            int    `json:"user-id"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Icon          []byte `json:"icon"`
	IconBase64    string `json:"iconbase64"`
	Bio           string `json:"bio"`
	Posts         int    `json:"countposts"`
	FollowBy      bool   `json:"followby"`
	FollowByCount int    `json:"followbycount"`
	FollowToCount int    `json:"followtocount"`
}


type UserMessage struct {
	MessageSession bool `json:"messagesession"`
	MessageID     int    `json:"post-id"`
	MessageUserID int    `json:"post-user-id"`
	UserID     int    `json:"user-id"`
	Content    string `json:"content"`
	Icon       []byte `json:"icon"`
	IconBase64 string `json:"iconbase64"`
	CreatedBy  string `json:"createdby"`
	Name       string `json:"createdbyname"`
	MessageBy int `json:"message-by"`
	MessageTo  int `json:"message-to"`
}