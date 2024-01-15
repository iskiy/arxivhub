package postgresql

import (
	"arxivhub/internal/models"
	"arxivhub/internal/repository"
	sqlc2 "arxivhub/internal/repository/postgresql/sqlc"
	"context"
	"database/sql"
)

type UserPSQLRepository struct {
	queries *sqlc2.Queries
}

var _ repository.UserRepository = (*UserPSQLRepository)(nil)

func NewUserPSQLRepository(db *sql.DB) *UserPSQLRepository {
	return &UserPSQLRepository{queries: sqlc2.New(db)}
}

func (r *UserPSQLRepository) CreateUser(ctx context.Context, arg models.RegisterUserRequest) (models.User, error) {
	params := sqlc2.CreateUserParams{
		Username:       arg.Username,
		Email:          arg.Email,
		HashedPassword: arg.Password,
	}
	createdUser, err := r.queries.CreateUser(ctx, params)

	return models.User(createdUser), err
}

func (r *UserPSQLRepository) GetUser(ctx context.Context, id int64) (models.User, error) {
	user, err := r.queries.GetUser(ctx, id)

	return models.User(user), err
}

func (r *UserPSQLRepository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)

	return models.User(user), err
}

func (r *UserPSQLRepository) UpdateUserEmail(ctx context.Context, arg models.UpdateUserEmailParams) (models.User, error) {
	user, err := r.queries.UpdateUserEmail(ctx, sqlc2.UpdateUserEmailParams(arg))

	return models.User(user), err
}
