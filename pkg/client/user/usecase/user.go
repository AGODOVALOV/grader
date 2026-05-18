// Package usecase provides user-related business logic.
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/parser"
	"io"
	"math"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"time"

	"go/token"

	"github.com/AGODOVALOV/grader/pkg/client/user/dto"
	"github.com/AGODOVALOV/grader/pkg/client/user/repo"
	"github.com/AGODOVALOV/grader/pkg/common"
	dtograder "github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/queue/config"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	jwttoken "github.com/AGODOVALOV/grader/pkg/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/streadway/amqp"
)

const (
	maxUploadFileSize = 1 << 20 // 1 MB
	goLanguage        = "go"
	// other like python and ...
)

var (
	ErrInvalidUploadFileName = errors.New("file name must be main.go")
	ErrUploadFileTooLarge    = errors.New("file is too large")
	ErrEmptyUploadFile       = errors.New("file is empty")
	ErrInvalidGoFile         = errors.New("file must contain valid Go code")
	ErrInvalidGoPackage      = errors.New("file must use package main")
)

// UserService provides methods for user-related operations.
type UserService struct {
	repo          *repo.Repo
	fStorage      *s3.FileStorage
	token         jwttoken.Maker
	tokenCallBack jwttoken.Maker
}

// NewUserService creates a new UserService instance.
func NewUserService(r *repo.Repo, fStorage *s3.FileStorage, tkn jwttoken.Maker, tknCallBack jwttoken.Maker) *UserService {
	return &UserService{
		repo:          r,
		fStorage:      fStorage,
		token:         tkn,
		tokenCallBack: tknCallBack,
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
func (s *UserService) GetNewToken(userID int64, login string) (string, *jwttoken.Payload, error) {
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
) (int64, error) {
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

// UploadFileToReviewS3 creates a new review assignment for a user.
func (s *UserService) UploadFileToReviewS3(
	ctx context.Context,
	objectName string,
	file multipart.File,
	size int64,
	eventID *uuid.UUID,
) error {
	_, err := s.fStorage.UploadFile(
		ctx,
		file,
		size,
		objectName,
		eventID)
	if err != nil {
		return err
	}

	return nil
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
			msg = fmt.Sprintf("result out for task - %v\n", v.ResultOut.String) + " " + v.LastError.String
		}

		result.Tasks = append(result.Tasks, dto.TaskData{
			ID:        int(v.Reviewid.Int64),
			Title:     v.Taskname.String,
			Status:    string(v.Status.ReviewStatus),
			Message:   msg,
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
			ReviewID:  v.Reviewid.Int64,
			UserID:    v.ID,
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

func (s *UserService) ProduceMessages(ctx context.Context, rCh *amqp.Channel, cfg *config.QueueMsgChannel) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tx, err := s.repo.DB.Pool.Begin(ctx)
		if err != nil {
			continue
		}

		qtx := s.repo.Queries.WithTx(tx)
		tasks, err := qtx.GetOutboxReviewsBatch(ctx)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logger.Z(ctx).Error(ctx, "rollback", err.Error())
			}
			continue
		}

		if len(tasks) == 0 {
			err = tx.Commit(ctx)
			if err != nil {
				logger.Z(ctx).Error(ctx, "commit", err.Error())
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(300 * time.Millisecond):
			}
			continue
		}

		ids := make([]int64, 0, len(tasks))
		for _, t := range tasks {
			ids = append(ids, t.ID)
		}

		err = qtx.MarkOutboxReviewsProcessingMany(ctx, ids)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				logger.Z(ctx).Error(ctx, "rollback", err.Error())
			}
			continue
		}

		err = tx.Commit(ctx)

		if err != nil {
			continue
		}

		for _, r := range tasks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			err := publishToRabbit(ctx, r, rCh, cfg)
			if err == nil {
				err = s.repo.Queries.MarkOutboxReviewProcessingOne(ctx, r.ID)
				if err != nil {
					logger.Z(ctx).Error(ctx, "Mark Outbox Review Processing One", err.Error())
				}
				err = s.repo.Queries.UpdateReviewStatusByID(ctx, repo.UpdateReviewStatusByIDParams{
					Status: repo.NullReviewStatus{
						ReviewStatus: "processing",
						Valid:        true,
					},
					ID: r.Reviewid.Int64,
				})
				if err != nil {
					logger.Z(ctx).Error(ctx, "Update Review Status By ID", err.Error())
				}
				continue
			}

			delay := time.Duration(math.Pow(2, float64(r.Attempts.Int32))) * time.Second
			if r.Attempts.Int32+1 >= r.MaxAttempts {
				_ = s.repo.Queries.MarkOutboxReviewFailedFinal(
					ctx,
					repo.MarkOutboxReviewFailedFinalParams{
						ID:        r.ID,
						LastError: pgtype.Text{String: err.Error(), Valid: true},
					},
				)
			} else {
				_ = s.repo.Queries.MarkOutboxReviewRetry(ctx, repo.MarkOutboxReviewRetryParams{
					ID:        r.ID,
					Column2:   int32(delay.Seconds()),
					LastError: pgtype.Text{String: err.Error(), Valid: true},
				})
			}
			logger.Z(ctx).Error(ctx, "publish msg", err.Error())
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(300 * time.Millisecond):
		}
	}
}

func publishToRabbit(ctx context.Context, review repo.OutboxReview, rCh *amqp.Channel, cfg *config.QueueMsgChannel) error {
	err := rCh.Publish(
		"",
		cfg.Name,
		false,
		false,
		amqp.Publishing{
			MessageId:    review.EventID.String(),
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         review.Payload,
		})
	if err != nil {
		logger.Z(ctx).Error(ctx, "publish msg", err.Error())
		return err
	}

	logger.Z(ctx).Debug(ctx, "publish message", "successful", map[string]string{
		"channel": cfg.Name,
		"msgID":   review.EventID.String(),
		"payload": string(review.Payload),
	})
	return nil
}

func (s *UserService) CreateAndOutboxReviewTx(
	ctx context.Context,
	userID int64,
	taskNum int64,
	filename string,
	eventID *uuid.UUID,
) error {
	tx, err := s.repo.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	qtx := s.repo.Queries.WithTx(tx)

	newReview, err := qtx.CreateReview(ctx, repo.CreateReviewParams{
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
		Task: int32(taskNum),
		Fileid: pgtype.Text{
			String: filename,
			Valid:  true,
		},
	})
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			logger.Z(ctx).Error(ctx, "rollback", err.Error())
			return err
		}
		return err
	}

	payloadGrader := dtograder.GraderPayload{
		UserID:   strconv.FormatInt(userID, 10),
		TaskID:   strconv.FormatInt(taskNum, 10),
		ReviewID: strconv.FormatInt(newReview.ID, 10),
		EventID:  eventID.String(),
		FileIDs: []dtograder.File{
			{"label_" + filename,
				filename},
		},
		ContainerName: "default",
	}

	jsonBytes, err := json.Marshal(payloadGrader)
	if err != nil {
		logger.Z(ctx).Error(ctx, "create outbox event", err.Error())
		err := tx.Rollback(ctx)
		if err != nil {
			logger.Z(ctx).Error(ctx, "rollback", err.Error())
			return err
		}
		return err
	}

	err = qtx.CreateOutboxReview(ctx, repo.CreateOutboxReviewParams{
		EventID: pgtype.UUID{
			Bytes: *eventID,
			Valid: true,
		},
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
		Reviewid: pgtype.Int8{
			Int64: newReview.ID,
			Valid: true,
		},
		Payload: jsonBytes,
	})
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			logger.Z(ctx).Error(ctx, "rollback", err.Error())
			return err
		}
		return err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) VerifyTokenCallBackToken(rawToken string) (*jwttoken.Payload, error) {
	payload, err := s.tokenCallBack.VerifyToken(rawToken)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func (s *UserService) ProcessGraderCallback(ctx context.Context, payload *dtograder.GraderPayloadCallback) error {
	var reviewStatusNew string

	tx, err := s.repo.DB.Pool.Begin(ctx)
	if err != nil {
		return err
	}

	if payload.Passed {
		reviewStatusNew = "done"
	} else {
		reviewStatusNew = "failed"
	}

	qtx := s.repo.Queries.WithTx(tx)

	reviewID, err := strconv.ParseInt(payload.ReviewID, 10, 64)
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(payload.UserID, 10, 64)
	if err != nil {
		return err
	}

	taskID, err := strconv.ParseInt(payload.TaskID, 10, 64)

	err = qtx.MarkOutboxReviewStatus(ctx, repo.MarkOutboxReviewStatusParams{
		Status: repo.OutboxStatus(reviewStatusNew),
		LastError: pgtype.Text{
			String: payload.ErrorText,
			Valid:  true,
		},
		ResultOut: pgtype.Text{
			String: payload.OutputMessage,
			Valid:  true,
		},
		Reviewid: pgtype.Int8{
			Int64: reviewID,
			Valid: true,
		},
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
	})
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = qtx.MarkReviewStatus(ctx, repo.MarkReviewStatusParams{
		Status: repo.NullReviewStatus{
			ReviewStatus: repo.ReviewStatus(reviewStatusNew),
			Valid:        true,
		},
		ID: reviewID,
		Userid: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
		Task: int32(taskID),
	})
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) ValidateUploadFile(ctx context.Context, id int64, file multipart.File, header *multipart.FileHeader) error {
	task, err := s.repo.Queries.GetTaskByID(ctx, int32(id))
	if err != nil {
		return err
	}

	if filepath.Base(header.Filename) != task.TargetFileName.String {
		return ErrInvalidUploadFileName
	}

	if task.TargetFileValidation.Bool {
		return validateFileByConditions(task.TargetFileValidationLanguage.String, file, header)
	}

	return nil
}

func validateFileByConditions(lang string, file multipart.File, header *multipart.FileHeader) error {
	if header.Size <= 0 {
		return ErrEmptyUploadFile
	}

	if header.Size > maxUploadFileSize {
		return ErrUploadFileTooLarge
	}

	switch lang {
	case goLanguage:
		return validateGoSourceFile(file, header)
		// case otherLanguage:
		// 	return validateOtherLanguageSourceFile(file)
	}

	return nil
}

func validateGoSourceFile(file multipart.File, header *multipart.FileHeader) error {
	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, header.Filename, content, parser.AllErrors)
	if err != nil {
		return ErrInvalidGoFile
	}

	if parsedFile.Name == nil || parsedFile.Name.Name != "main" {
		return ErrInvalidGoPackage
	}

	return nil
}
