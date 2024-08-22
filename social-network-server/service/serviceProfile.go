package service

import (
	"social-network-server/database"
	"social-network-server/pkg/models"
)

type UserServiceProfile struct {
	userRepo *database.UserRepository
}

func NewUserServiceProfile() *UserServiceProfile {
	return &UserServiceProfile{
		userRepo: database.NewUserRepository(),
	}
}

func (us *UserServiceProfile) GetUserProfileAndPosts(username string, currentUserID int) (*models.UserProfile, []models.UserPost, bool, *models.UserIconResponse, error) {
	// Obter o ID do usuário alvo
	targetUserID, err := us.userRepo.GetUserIDByUsername(username)
	if err != nil {
		return nil, nil, false, nil, err
	}

	// Obter perfil do usuário
	profile, err := us.userRepo.GetUserProfile(targetUserID)
	if err != nil {
		return nil, nil, false, nil, err
	}

	// Obter posts do usuário
	posts, err := us.userRepo.GetUserPosts(targetUserID)
	if err != nil {
		return nil, nil, false, nil, err
	}

	profile.Posts = len(posts)

	// Verificar se o usuário atual segue o usuário alvo
	profile.FollowBy, err = us.userRepo.IsFollowing(currentUserID, targetUserID)
	if err != nil {
		return nil, nil, false, nil, err
	}

	// Verificar se o usuário alvo segue o usuário atual
	profile.FollowTo, err = us.userRepo.IsFollowing(targetUserID, currentUserID)
	if err != nil {
		return nil, nil, false, nil, err
	}

	// Verificar se o perfil visualizado é do usuário atual
	currentUsername, err := us.userRepo.GetUsernameByID(currentUserID)
	if err != nil {
		return nil, nil, false, nil, err
	}

	isCurrentUser := username == currentUsername

	// Obter ícone e nome do chat partner
	userService := UserService{}
	chatPartner, err := userService.GetUserIcon(currentUserID)

	if err != nil {
		return nil, nil, false, nil, err
	}

	return profile, posts, isCurrentUser, chatPartner, nil
}
