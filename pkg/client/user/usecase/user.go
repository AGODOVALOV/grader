package usecase

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/AGODOVALOV/grader/pkg/client/user/dto"
	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/common"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	repo     *repo.Repo
	fStorage *s3.FileStorage
	token    token.Maker
}

func NewUserService(r *repo.Repo, fStorage *s3.FileStorage, token token.Maker) *UserService {
	return &UserService{
		repo:     r,
		fStorage: fStorage,
		token:    token,
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

func (s *UserService) CreateNewReview(
	ctx context.Context,
	userID int64,
	taskNum int,
	objectName string,
	file multipart.File,
	size int64) (int64, error) {

	_, err := s.fStorage.UploadFile(
		ctx,
		file,
		size,
		objectName)
	if err != nil {
		return 0, err
	}

	reviewNew, err := s.repo.Queries.CreateReview(ctx, repo.CreateReviewParams{
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
		Task: int32(taskNum),
		Fileid: pgtype.Text{
			String: objectName,
			Valid:  true,
		},
	})

	if err != nil {
		return 0, err
	}

	return reviewNew.ID, nil

}

func (s *UserService) GetReviewsByUserID(ctx context.Context, userID int64) (*dto.AccountPageData, error) {
	data, err := s.repo.Queries.GetReviewsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, common.ErrRecordNotFound
	}

	header := data[0]

	result := dto.AccountPageData{
		ID:    int(header.ID),
		Login: header.Login,
		Name:  header.Name,
	}

	var msg string

	for _, v := range data {
		if v.Reviewid.Int64 == 0 && len(data) == 1 {
			msg = "You don't have any submitted review assignments yet."
		} else {
			msg = fmt.Sprintf("message for task - %s", v.Taskname.String)
		}

		result.Tasks = append(result.Tasks, dto.TaskData{
			ID:        int(v.Reviewid.Int64),
			Title:     v.Taskname.String,
			Status:    string(v.Status.ReviewStatus),
			Message:   msg, //fmt.Sprintf("message for task - %s", v.Taskname.String),
			UpdatedAt: v.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return &result, nil
}
