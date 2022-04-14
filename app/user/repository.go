package user

import (
	"context"
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Save(ctx context.Context, user User) (User, error)
	FindByID(ctx context.Context, userID string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
}

type repository struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) Repository {
	return &repository{DB}
}

const (
	layoutDateTime = "2006-01-02 15:04:05"
)

func (r *repository) Save(ctx context.Context, user User) (User, error) {
	sqlQuery := "INSERT INTO users (id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	_, err := r.DB.ExecContext(ctx, sqlQuery, user.ID,
		user.Name,
		user.Occupation,
		user.Email,
		user.PasswordHash,
		user.AvatarFileName,
		user.Role,
		time.Now().Format(layoutDateTime),
		time.Now().Format(layoutDateTime))

	if err != nil {
		log.Println(err.Error())
		return user, err
	}

	log.Info("Success insert new user!")
	return user, nil
}

func (r *repository) FindByID(ctx context.Context, userID string) (User, error) {
	user := User{}
	var createdAt, updatedAt string

	sqlQuery := "SELECT id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at FROM users WHERE id = $1"

	rows, err := r.DB.QueryContext(ctx, sqlQuery, userID)
	if err != nil {
		return user, err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Occupation,
			&user.Email,
			&user.PasswordHash,
			&user.AvatarFileName,
			&user.Role,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return user, err
		}
	}

	if createdAt != "" || updatedAt != "" {
		if user.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			log.Fatal(err)
		}

		if user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			log.Fatal(err)
		}
	}

	log.Info(user)
	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	var createdAt, updatedAt string

	sqlQuery := "SELECT id, name, occupation, email, password_hash, avatar_file_name, role, created_at, updated_at FROM users WHERE email = $1"

	rows, err := r.DB.QueryContext(ctx, sqlQuery, email)
	if err != nil {
		return user, err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Occupation,
			&user.Email,
			&user.PasswordHash,
			&user.AvatarFileName,
			&user.Role,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return user, err
		}
	}

	if createdAt != "" || updatedAt != "" {
		if user.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			log.Fatal(err)
		}

		if user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			log.Fatal(err)
		}
	}

	log.Info(user)
	return user, nil
}
