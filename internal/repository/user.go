package repository

import (
	"context"
	"crud/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	Create(user *model.UserModel) (*model.UserModel, error)
	GetById(id int) (*model.UserModel, error)
	Update(user *model.UserModel) (*model.UserModel, error)
	Delete(id int) (*model.UserModel, error)
	GetAll(offset, limit int) ([]*model.UserModel, error)
}

type UserRepository struct {
	dbPool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) IUserRepository {
	return &UserRepository{dbPool: pool}
}

func (repository *UserRepository) Create(user *model.UserModel) (*model.UserModel, error) {
	row := repository.dbPool.QueryRow(context.Background(),
		"INSERT INTO users(name, email, age) values($1, $2, $3) RETURNING id",
		user.Name, user.Email, user.Age)
	response := *user
	err := row.Scan(&response.ID)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repository *UserRepository) GetById(id int) (*model.UserModel, error) {
	row := repository.dbPool.QueryRow(context.Background(), "SELECT id, name, email, age FROM users WHERE id = $1",
		id)
	response := model.UserModel{}
	err := row.Scan(&response.ID, &response.Name, &response.Email, &response.Age)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repository *UserRepository) Update(user *model.UserModel) (*model.UserModel, error) {
	row := repository.dbPool.QueryRow(context.Background(),
		"UPDATE users SET name = $1, age = $2 WHERE id = $3 RETURNING id, name, email, age",
		user.Name, user.Age, user.ID)
	response := model.UserModel{}
	err := row.Scan(&response.ID, &response.Name, &response.Email, &response.Age)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repository *UserRepository) Delete(id int) (*model.UserModel, error) {
	row := repository.dbPool.QueryRow(context.Background(),
		"DELETE FROM users WHERE id = $1 RETURNING id, name, email, age", id)
	response := model.UserModel{}
	err := row.Scan(&response.ID, &response.Name, &response.Email, &response.Age)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repository *UserRepository) GetAll(offset, limit int) ([]*model.UserModel, error) {
	rows, err := repository.dbPool.Query(context.Background(), "SELECT * FROM users LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	users := make([]*model.UserModel, 0)
	for rows.Next() {
		user := &model.UserModel{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
