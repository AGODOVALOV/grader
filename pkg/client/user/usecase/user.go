// Package usecase provides user-related business logic.
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
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserService provides methods for user-related operations.
type UserService struct {
	repo     *repo.Repo
	fStorage *s3.FileStorage
	token    token.Maker
}

// NewUserService creates a new UserService instance.
func NewUserService(r *repo.Repo, fStorage *s3.FileStorage, tkn token.Maker) *UserService {
	return &UserService{
		repo:     r,
		fStorage: fStorage,
		token:    tkn,
	}
}

// CreateUser creates a new user.
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

// CheckUserLogin checks if a user exists and the provided password is correct.
func (s *UserService) CheckUserLogin(ctx context.Context, login, password string) (int64, error) {
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

// CheckUserIsAdminByLogin checks if a user exists and the provided password is correct.
func (s *UserService) CheckUserIsAdminByLogin(ctx context.Context, login string) (bool, error) {
	isAdmin, err := s.repo.Queries.IsAdminByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return isAdmin, nil
}

// CheckUserIsAdminByUserID checks if a user exists and the provided password is correct.
func (s *UserService) CheckUserIsAdminByUserID(ctx context.Context, id int64) (bool, error) {
	isAdmin, err := s.repo.Queries.IsAdminByUseID(ctx, id)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return isAdmin, nil
}

// GetNewToken generates a new JWT token for the user.
func (s *UserService) GetNewToken(userID int64, login string) (string, *token.Payload, error) {
	jwtToken, payload, err := s.token.CreateToken(userID, login)
	if err != nil {
		return "", nil, err
	}
	return jwtToken, payload, nil
}

// CreateNewReview creates a new review assignment for a user.
func (s *UserService) CreateNewReview(
	ctx context.Context,
	userID int64,
	taskNum int32,
	objectName string,
	file multipart.File,
	size int64,
) (int64, error) {
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
		Task: taskNum,
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

// GetReviewsByUserID returns a list of review assignments for a user.
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

	for i := range data {
		v := &data[i]
		if v.Reviewid.Int64 == 0 && len(data) == 1 {
			msg = "You don't have any submitted review assignments yet."
		} else {
			msg = fmt.Sprintf("message for task - %s", v.Taskname.String)
		}

		result.Tasks = append(result.Tasks, dto.TaskData{
			ID:        int(v.Reviewid.Int64),
			Title:     v.Taskname.String,
			Status:    string(v.Status.ReviewStatus),
			Message:   msg, // fmt.Sprintf("message for task - %s", v.Taskname.String),
			UpdatedAt: v.CreatedAt.Time.Format(time.RFC3339),
		})
	}

	return &result, nil
}

// GetReviews returns a list of all review assignments.
func (s *UserService) GetReviews(ctx context.Context) (*dto.AdminReviewsPageData, error) {
	data, err := s.repo.Queries.GetReviewsAll(ctx)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, common.ErrRecordNotFound
	}

	result := dto.AdminReviewsPageData{}

	for i := range data {
		v := &data[i]
		result.Reviews = append(result.Reviews, dto.AdminReviewData{
			ID:        v.ID,
			UserLogin: v.Login,
			FileName:  v.Name,
			TaskTitle: v.Taskname.String,
			Status:    string(v.Status.ReviewStatus),
			CreatedAt: v.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt: v.UpdatedAt.Time.Format(time.RFC3339),
		},
		)
	}

	return &result, nil
}

// UpdateReviewStatus updates the status of a review assignment.
func (s *UserService) UpdateReviewStatus(ctx context.Context, id int64, status string) error {
	return s.repo.Queries.UpdateReviewStatusByID(ctx, repo.UpdateReviewStatusByIDParams{
		Status: repo.NullReviewStatus{
			ReviewStatus: repo.ReviewStatus(status),
			Valid:        true,
		},
		ID: id,
	})
}

func (s *UserService) CreateOutboxReview(ctx context.Context,
	eventID uuid.UUID, userID int64, reviewID int64, payload []byte) error {
	return s.repo.Queries.CreateOutboxReview(ctx, repo.CreateOutboxReviewParams{
		EventID: pgtype.UUID{
			Bytes: eventID,
			Valid: true,
		},
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
		Reviewid: pgtype.Int8{
			Int64: reviewID,
			Valid: true,
		},
		Payload: payload,
	})
}
