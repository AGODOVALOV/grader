package test_test

import (
	"context"
	"testing"

	"github.com/AGODOVALOV/grader/pkg/config"
	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/grader/client"
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool/worker"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/AGODOVALOV/grader/pkg/token"
	tokenconfig "github.com/AGODOVALOV/grader/pkg/token/config"
	"github.com/stretchr/testify/require"
)

func TestWorker_DoJob(t *testing.T) {
	// create and read config
	appCfg, err := config.GetApplicationConfig()
	if err != nil {
		t.Error(err)
		return
	}

	tokenMaker, err := token.NewJWTMaker((*tokenconfig.Config)(&appCfg.GetConfig().Grader.Callback.JWT))
	if err != nil {
		t.Error(err)
		return
	}

	callbackClient := client.NewClient(&appCfg.GetConfig().Grader.Callback, tokenMaker)

	// create logger
	z, err := logger.NewAppLogger(appCfg.GetConfig().Log)
	if err != nil {
		t.Error(err)
		return
	}

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	// ctx add logger
	ctx = logger.CtxWithLogger(ctx, z)
	defer cancel()

	// init file storage
	fStorage, err := s3.NewFileStorage(ctx, &appCfg.GetConfig().FileStorage)
	if err != nil {
		z.Error(ctx, "init file storage", err.Error())
		return
	}

	type fields struct {
		fStorage *s3.FileStorage
	}

	type args struct {
		ctx     context.Context
		payload *dto.GraderPayload
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "good test with HW1",
			fields: fields{
				fStorage: fStorage,
			},
			args: args{
				ctx: ctx,
				payload: &dto.GraderPayload{
					TaskID:   "1",
					UserID:   "9",
					ReviewID: "20",
					FileIDs: []dto.File{
						{
							Label:    "label_",
							FileName: "review_9_1_main.go",
						},
					},
					ContainerName: appCfg.GetConfig().FileStorage.Bucket,
				},
			},
			wantErr: false,
		},
		{
			name: "bad test with HW2",
			fields: fields{
				fStorage: fStorage,
			},
			args: args{
				ctx: ctx,
				payload: &dto.GraderPayload{
					TaskID:   "2",
					UserID:   "9",
					ReviewID: "20",
					FileIDs: []dto.File{
						{
							Label:    "label_",
							FileName: "review_9_2_main.go",
						},
					},
					ContainerName: appCfg.GetConfig().FileStorage.Bucket,
				},
			},
			wantErr: true,
		},
		{
			name: "file is not exist test",
			fields: fields{
				fStorage: fStorage,
			},
			args: args{
				ctx: ctx,
				payload: &dto.GraderPayload{
					UserID:   "9",
					ReviewID: "20",
					FileIDs: []dto.File{
						{
							Label:    "label_",
							FileName: "badfile.file",
						},
					},
					ContainerName: appCfg.GetConfig().FileStorage.Bucket,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := worker.NewWorker(tt.fields.fStorage, callbackClient)
			err := w.DoJob(tt.args.ctx, tt.args.payload)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
