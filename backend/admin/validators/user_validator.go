package validators

import (
	"errors"
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	validRoles = map[string]bool{
		"admin": true,
		"user":  true,
	}
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

type ResetPasswordRequest struct{
	NewPassword string `json:"new_password"`
}

func (r *CreateUserRequest) Validate() error {
	if r.Username == "" {
		return errors.New("username is required")
	}
	if len(r.Username) < 3 || len(r.Username) > 50 {
		return errors.New("username must be between 3 and 50 characters")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(r.Username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	if r.Email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(strings.ToLower(r.Email)) {
		return errors.New("invalid email format")
	}

	if r.Password == "" {
		return errors.New("password is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	if r.FullName == "" {
		return errors.New("full_name is required")
	}
	if len(r.FullName) > 255 {
		return errors.New("full_name must be less than 255 characters")
	}

	return nil
}

func (r *UpdateUserRequest) Validate() error {
	if r.Email != nil {
		if *r.Email == "" {
			return errors.New("email cannot be empty")
		}
		if !emailRegex.MatchString(strings.ToLower(*r.Email)) {
			return errors.New("invalid email format")
		}
	}

	if r.FullName != nil {
		if *r.FullName == "" {
			return errors.New("full_name cannot be empty")
		}
		if len(*r.FullName) > 255 {
			return errors.New("full_name must be less than 255 characters")
		}
	}

	return nil
}

func (r *ResetPasswordRequest) Validate() error {
	if r.NewPassword == "" {
		return errors.New("new_password is required")
	}
	if len(r.NewPassword) < 6 {
		return errors.New("new_password must be at least 6 characters")
	}
	return nil
}
