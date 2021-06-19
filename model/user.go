package model

import (
	"context"

	"fmt"
	"time"
)

type User struct {
	ID        string    `json:"id"  db:"id"`
	Name      string    `json:"name"  db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type UsersMe struct {
	ID            string   `json:"id"  db:"id"`
	Name          string   `json:"name"  db:"name"`
	FavoriteFiles []string `json:"favoriteFiles" db:"sound_id"`
}

func GetUsers(ctx context.Context) ([]*User, error) {
	users := []*User{}
	err := db.SelectContext(ctx, &users, "SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("Failed to get users: %w", err)
	}

	return users, nil
}

func CreateUser(ctx context.Context, user *User) error {
	_, err := db.ExecContext(ctx, "INSERT IGNORE INTO users (id, name) VALUES (?, ?) ", user.ID, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(ctx context.Context, accessToken string, userID string) (*User, error) {
	var user User
	err := db.GetContext(ctx, &user, "SELECT id, name FROM users WHERE id = ? LIMIT 1", userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	return &user, nil
}

func GetUsersMe(ctx context.Context, accessToken string) (*UsersMe, error) {

	var usersMe UsersMe
	err1 := db.GetContext(ctx, &usersMe, "SELECT id, name FROM users")
	if err1 != nil {
		return nil, fmt.Errorf("Failed to get usersMeId: %w", err1)
	}

	var usersMeFavorites []string
	err2 := db.SelectContext(ctx, &usersMeFavorites, "SELECT sound_id FROM favorites")
	if usersMeFavorites == nil {
		usersMeFavorites = []string{}
	}

	if err2 != nil {
		return nil, fmt.Errorf("Failed to get usersMeFavorites: %w", err2)
	}
	usersMe.FavoriteFiles = usersMeFavorites
	return &usersMe, nil
}
