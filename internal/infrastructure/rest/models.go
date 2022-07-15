package rest

import (
	"time"

	"github.com/pamugk/social-nyetwork-server/internal/domain"
)

type errorResponse struct {
	Err string
}

type createUserRequest struct {
	UserData
	Password string
}

type getUserResponse struct {
	User
}

type updateUserRequest struct {
	UserData
}

type changePasswordRequest struct {
	NewPassword string
}

type searchResponse struct {
	PageNumber int32
	PageSize   int32
	Total      int64
}

type searchUsersResponse struct {
	searchResponse
	Page []domain.ShortUser
}

type User struct {
	domain.ShortUser
	Created      time.Time
	Phone, Email *string
	Birthday     string
	Gender       domain.Gender
}

type UserData struct {
	Login           string
	PreferredLocale string
	Country         string
	Name            string
	Surname         string
	Patronymic      *string
	Phone           *string
	Email           *string
	About           *string
	Birthday        string
	Gender          domain.Gender
}
