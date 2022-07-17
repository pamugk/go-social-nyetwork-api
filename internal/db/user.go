package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/pamugk/social-nyetwork-server/internal/domain"
)

func UserExistsByLogin(login string) (exists bool, err error) {
	result := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM \"user\" WHERE login = $1 AND active)", login)
	err = result.Scan(&exists)
	return exists, err
}

func CreateUser(data *domain.UserData, hashedPassword string) (int64, error) {
	var id int64
	result := pool.QueryRow(context.Background(),
		`INSERT INTO "user"(login, "password", preferred_locale, country, "name", surname, patronymic, birthday, "gender", phone, email, about)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`, data.Login, hashedPassword, data.PreferredLocale, data.Country, data.Name, data.Surname, data.Patronymic, data.Birthday, data.Gender, data.Phone, data.Email, data.About)
	if err := result.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func GetUser(id int64) (*domain.User, error) {
	user := &domain.User{}
	result := pool.QueryRow(context.Background(), "SELECT id, login, \"name\", surname, patronymic, about, created, phone, email, birthday, \"gender\" FROM \"user\" WHERE id = $1", id)

	err := result.Scan(&user.Id, &user.Login, &user.Name, &user.Surname, &user.Patronymic, &user.About, &user.Created, &user.Phone, &user.Email, &user.Birthday, &user.Gender)
	if err == pgx.ErrNoRows {
		return nil, errors.New("Not found")
	} else if err != nil {
		return user, err
	}
	return user, nil
}

func SearchUsers(loginPart string, page int32, limit int32) ([]domain.ShortUser, int64, error) {
	users := []domain.ShortUser{}
	var total int64
	var rows pgx.Rows
	noFilter := loginPart == ""

	tx, err := StartTransaction()
	
	if err != nil {
		return users, total, err
	}
	
	defer FinishTransaction(tx, err)
	
	if noFilter {
		rows, err = pool.Query(context.Background(), "SELECT id, login, \"name\", surname, patronymic, about FROM \"user\" WHERE active ORDER BY id LIMIT $1 OFFSET $2", limit, page*limit)
	} else {
		rows, err = pool.Query(context.Background(), "SELECT id, login, \"name\", surname, patronymic, about FROM \"user\" WHERE login ILIKE $1 AND active ORDER BY id LIMIT $2 OFFSET $3", loginPart+"%", limit, page*limit)
	}

	if err == pgx.ErrNoRows {
		err = nil
	} else if err == nil {
		var countResult pgx.Row
		if noFilter {
			countResult = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM \"user\" WHERE active")
		} else {
			countResult = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM \"user\" WHERE login ILIKE $1 AND active", loginPart+"%")
		}
		
		_ = countResult.Scan(&total)
		
		for rows.Next() {
			user := domain.ShortUser{}
			if rows.Scan(&user.Id, &user.Login, &user.Name, &user.Surname, &user.Patronymic, &user.About) == nil {
				users = append(users, user)
			}
		}
	}

	return users, total, err
}

func OtherUserExistsByLogin(id int64, login string) (exists bool, err error) {
	result := pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM \"user\" WHERE id <> $1 AND login = $2 AND active)", id, login)
	err = result.Scan(&exists)
	return exists, err
}

func UpdateUser(id int64, data *domain.UserData) error {
	outcome, err := pool.Exec(context.Background(),
		`UPDATE "user"
		SET login = $1, preferred_locale = $2, country = $3, "name" = $4, surname = $5, patronymic = $6, birthday = $7, "gender" = $8, phone = $9, email = $10, about = $11
		WHERE id = $12 AND active`,
		data.Login, data.PreferredLocale, data.Country, data.Name, data.Surname, data.Patronymic, data.Birthday, data.Gender, data.Phone, data.Email, data.About,
		id)
	if err != nil {
		return err
	} else if outcome.RowsAffected() == 0 {
		return errors.New("Not found")
	}
	return nil
}

func ChangePassword(id int64, newPassword string) error {
	outcome, err := pool.Exec(context.Background(), "UPDATE \"user\" SET \"password\" = $1 WHERE id = $2 AND active", newPassword, id)
	if err != nil {
		return err
	} else if outcome.RowsAffected() == 0 {
		return errors.New("Not found")
	}
	return nil
}

func DeleteUser(id int64) error {
	tag, err := pool.Exec(context.Background(), "UPDATE \"user\" SET active = FALSE WHERE id = $1 AND active", id)
	if err != nil {
		return err
	} else if tag.RowsAffected() == 0 {
		return errors.New("Not found")
	}
	return nil
}
