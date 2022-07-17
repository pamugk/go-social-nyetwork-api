package domain

import "time"

type Gender string

const (
	MALE   Gender = "MALE"
	FEMALE        = "FEMALE"
)

type UserData struct {
	Login           string    `validate:"required,excludesall= ,max=100"`
	PreferredLocale string    `validate:"required,bcp47_language_tag,max=35"`
	Country         string    `validate:"required,iso3166_1_alpha3"`
	Name            string    `validate:"required,max=150"`
	Surname         string    `validate:"required,max=200"`
	Patronymic      *string   `validate:"omitempty,max=175"`
	Phone           *string   `validate:"omitempty,e164"`
	Email           *string   `validate:"omitempty,email,max=400"`
	About           *string   `validate:"omitempty,max=300"`
	Birthday        time.Time `validate:"required,not_future"`
	Gender          Gender    `validate:"required"`
}

type ShortUser struct {
	Id            int64
	Login         string
	Name, Surname string
	Patronymic    *string
	About         *string
}

type User struct {
	ShortUser
	Created      time.Time
	Phone, Email *string
	Birthday     time.Time
	Gender       Gender
}
