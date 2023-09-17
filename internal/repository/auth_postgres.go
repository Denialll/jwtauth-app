package repository

import (
	"context"
	"fmt"
	"github.com/Denialll/jwtauth-app/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)

	return user, err
}

//func (r *AuthPostgres) SetSession(ctx context.Context, studentID int, session models.Session) error {
//	fmt.Println(studentID)
//	query := fmt.Sprintf("UPDATE %s SET refresh_token = $1 AND expires_at = $2 WHERE id = $3", usersTable)
//
//	_, err := r.db.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, studentID)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (r *AuthPostgres) SetSession(ctx context.Context, studentID int, session models.Session) error {
	fmt.Println(studentID)
	query := fmt.Sprintf("UPDATE %s SET refresh_token = $1, expires_at = $2 WHERE id = $3", usersTable)
	_, err := r.db.ExecContext(ctx, query, session.RefreshToken, session.ExpiresAt, studentID)
	if err != nil {
		return err
	}
	return nil
}

//func (r *AuthPostgres) SetSession(ctx context.Context, studentID int, session models.Session) error {
//	query := fmt.Sprintf("UPDATE %s SET refresh_token = $1 AND expires_at = $2 WHERE id = $3", usersTable)
//	err := r.db.Get(userId)
//
//	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
//
//	return fmt.Errorf("Произошла ошибка")
//}

//func (r *StudentsRepo) SetSession(ctx context.Context, studentID primitive.ObjectID, session domain.Session) error {
//	_, err := r.db.UpdateOne(ctx, bson.M{"_id": studentID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
//
//	return err
//}
