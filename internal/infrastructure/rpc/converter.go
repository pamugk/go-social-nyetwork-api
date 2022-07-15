package rpc

import (
	"time"

	"github.com/pamugk/social-nyetwork-server/internal/domain"
	date "google.golang.org/genproto/googleapis/type/date"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func wrapStringPointer(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func convertFromDate(in *date.Date) time.Time {
	return time.Date(int(in.Year), time.Month(in.Month), int(in.Day), 0, 0, 0, 0, time.Now().Location())
}

func convertToDate(out time.Time) *date.Date {
	return &date.Date{Year: int32(out.Year()), Month: int32(out.Month()), Day: int32(out.Day())}
}

func convertFromGender(in Gender) domain.Gender {
	switch in {
	case Gender_MALE:
		{
			return domain.MALE
		}
	case Gender_FEMALE:
		{
			return domain.FEMALE
		}
	}
	return domain.MALE
}

func convertToGender(out domain.Gender) Gender {
	switch out {
	case domain.MALE:
		{
			return Gender_MALE
		}
	case domain.FEMALE:
		{
			return Gender_FEMALE
		}
	}
	return Gender_MALE
}

func convertFromUserData(in *UserData) *domain.UserData {
	if in == nil {
		return &domain.UserData{}
	}
	return &domain.UserData{
		Login:           in.Login,
		PreferredLocale: in.PreferredLocale,
		Country:         in.Country,
		Name:            in.Name,
		Surname:         in.Surname,
		Patronymic:      &in.Patronymic,
		Phone:           &in.Phone,
		Email:           &in.Email,
		About:           &in.About,
		Birthday:        convertFromDate(in.Birthday),
		Gender:          convertFromGender(in.Gender),
	}
}

func convertToUser(out *domain.User) *GetUserResponse {
	if out == nil {
		return &GetUserResponse{}
	}
	return &GetUserResponse{
		Id:         out.Id,
		Login:      out.Login,
		Created:    timestamppb.New(out.Created),
		Name:       out.Name,
		Surname:    out.Surname,
		Patronymic: wrapStringPointer(out.Patronymic),
		About:      wrapStringPointer(out.About),
		Phone:      wrapStringPointer(out.Phone),
		Email:      wrapStringPointer(out.Email),
		Birthday:   convertToDate(out.Birthday),
		Gender:     convertToGender(out.Gender),
	}
}

func convertToShortUser(out *domain.ShortUser) *ShortUser {
	if out == nil {
		return nil
	}
	return &ShortUser{
		Id:         out.Id,
		Login:      out.Login,
		Name:       out.Name,
		Surname:    out.Surname,
		Patronymic: wrapStringPointer(out.Patronymic),
		About:      wrapStringPointer(out.About),
	}
}

func convertToSearchUserResponse(page int32, limit int32, total int64, out []domain.ShortUser) *SearchUsersResponse {
	items := make([]*ShortUser, len(out))
	for i, s := range out {
		items[i] = convertToShortUser(&s)
	}
	return &SearchUsersResponse{PageNumber: page, PageSize: limit, Total: total, Page: items}
}
