package app

import (
	"errors"
	"strings"

	"github.com/pamugk/social-nyetwork-server/internal/db"
	"github.com/pamugk/social-nyetwork-server/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CreateUser(data *domain.UserData, password string) (newId int64, err error) {
	sanitizeUserData(data)
	if err = validate.Struct(data); err != nil {
		return -1, err
	}
	if password == "" {
		return -1, errors.New("Password required")
	}
	tx, err := db.StartTransaction()
	if err != nil {
		return -1, err
	}
	defer db.FinishTransaction(tx, err)
	userExists, err := db.UserExistsByLogin(data.Login)
	if err != nil {
		return -1, err
	}
	if userExists {
		return -1, errors.New("Login is already taken")
	}
	newId, err = db.CreateUser(data, hashPassword(password))
	return newId, err
}

func GetUser(id int64) (*domain.User, error) {
	return db.GetUser(id)
}

func SearchUsers(loginPart string, page int32, limit int32) ([]domain.ShortUser, int64, error) {
	return db.SearchUsers(loginPart, page, limit)
}

func UpdateUser(id int64, data *domain.UserData) (err error) {
	sanitizeUserData(data)
	if err = validate.Struct(data); err != nil {
		return err
	}
	tx, err := db.StartTransaction()
	if err != nil {
		return err
	}
	defer db.FinishTransaction(tx, err)
	userExists, err := db.OtherUserExistsByLogin(id, data.Login)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("Login is already taken")
	}
	return db.UpdateUser(id, data)
}

func ChangePassword(id int64, newPassword string) (err error) {
	if newPassword == "" {
		return errors.New("Password required")
	}
	return db.ChangePassword(id, hashPassword(newPassword))
}

func DeleteUser(id int64) error {
	return db.DeleteUser(id)
}

func sanitizeUserData(data *domain.UserData) {
	data.Login = strings.TrimSpace(data.Login)
	data.PreferredLocale = strings.TrimSpace(data.PreferredLocale)
	data.Country = strings.TrimSpace(data.Country)
	data.Name = strings.TrimSpace(data.Name)
	data.Surname = strings.TrimSpace(data.Surname)
	sanitizeStringPointer(&data.Patronymic)
	sanitizeStringPointer(&data.About)
	sanitizeStringPointer(&data.Phone)
	sanitizeStringPointer(&data.Email)
}
