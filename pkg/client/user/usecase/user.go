package usecase

import (
	"context"
	"errors"

	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/common"
	"github.com/AGODOVALOV/grader/pkg/token"
)

type UserService struct {
	repo  *repo.Repo
	token token.Maker
}

func NewUserService(r *repo.Repo, token token.Maker) *UserService {
	return &UserService{
		repo:  r,
		token: token,
	}
}

func (s *UserService) CreateUser(ctx context.Context, login, name, password string) error {

	_, err := s.repo.Queries.CreateUser(ctx, repo.CreateUserParams{
		Login:    login,
		Name:     name,
		Password: password,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) CheckUserLogin(ctx context.Context, login string, password string) (int64, error) {
	user, err := s.repo.Queries.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return 0, common.ErrRecordNotFound
		}
		return 0, err
	}

	err = common.CheckPassword(password, user.Password)
	if err != nil {
		return 0, common.ErrIncorrectPassword
	}

	return user.ID, nil

}

func (s *UserService) CheckUserIsAdmin(ctx context.Context, login string) (bool, error) {
	isAdmin, err := s.repo.Queries.IsAdmin(ctx, login)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return isAdmin, nil
}

func (s *UserService) GetNewToken(userID int64, login string) (string, *token.Payload, error) {
	jwtToken, payload, err := s.token.CreateToken(userID, login)
	if err != nil {
		return "", nil, err
	}
	return jwtToken, payload, nil
}
