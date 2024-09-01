package service

import (
	"database/sql"

	"social-network-server/internal/models"
	"social-network-server/pkg/repositories"
	"social-network-server/pkg/validators"
)

// SignupService handles the signup logic.
func SignupService(db *sql.DB, user models.User) (map[string]string, error) {
	resp := make(map[string]string)

	// Validate user
	if user.Username == "" || user.Name == "" || user.Email == "" || user.Password == "" || user.ConfirmPassword == "" {
		resp["missing"] = "Some values are missing!"
	}
	if len(user.Username) < 4 || len(user.Username) > 32 {
		resp["username"] = "Username should be between 4 and 32"
	}
	if len(user.Name) < 1 || len(user.Name) > 70 {
		resp["name"] = "Name should be between 1 and 70"
	}
	if len(user.Bio) > 150 {
		resp["bio"] = "Bio should be up to 150 characters"
	}
	if err := validators.ValidateFormatEmail(user.Email); err != nil {
		resp["email"] = "Invalid email format!"
	}
	existEmail, err := repositories.CheckEmailExistence(db, user.Email)
	if err != nil {
		return resp, err
	}
	if existEmail {
		resp["email"] = "Email already exists!"
	}
	if len(user.Password) < 8 || len(user.Password) > 16 {
		resp["password"] = "Passwords should be between 8 and 16 characters"
	}
	if user.Password != user.ConfirmPassword {
		resp["confirm_password"] = "Passwords don't match"
	}

	if len(resp) > 0 {
		return resp, nil
	}

	// Hash the password before saving
	hashedPassword, err := validators.Hash(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	err = repositories.CreateUser(db, user)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
