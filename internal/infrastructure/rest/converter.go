package rest

import (
	"time"

	"github.com/pamugk/social-nyetwork-server/internal/domain"
)

const (
	ISO_DATE_FORMAT = "2006-01-02"
)

func convertFromUserData(in *UserData) (*domain.UserData, error) {
	birthday, err := time.Parse(ISO_DATE_FORMAT, in.Birthday)
	return &domain.UserData{
		Login:           in.Login,
		PreferredLocale: in.PreferredLocale,
		Country:         in.Country,
		Name:            in.Name,
		Surname:         in.Surname,
		Patronymic:      in.Patronymic,
		Phone:           in.Phone,
		Email:           in.Email,
		About:           in.About,
		Birthday:        birthday,
		Gender:          in.Gender,
	}, err
}

func convertToUser(out *domain.User) *User {
	return &User{
		out.ShortUser,
		out.Created,
		out.Phone,
		out.Email,
		out.Birthday.Format(ISO_DATE_FORMAT),
		out.Gender,
	}
}
