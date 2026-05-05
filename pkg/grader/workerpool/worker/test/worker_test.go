package test_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/AGODOVALOV/grader/pkg/config"
	"github.com/AGODOVALOV/grader/pkg/dto"
	"github.com/AGODOVALOV/grader/pkg/grader/workerpool/worker"
	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/storage/s3"
	"github.com/stretchr/testify/require"
)

func TestWorker_DoJob(t *testing.T) {

	// create and read config
	appCfg, err := config.GetApplicationConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// create logger
	z, err := logger.NewAppLogger(appCfg.GetConfig().Log)
	if err != nil {
		fmt.Println(err)
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
			name: "good test",
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
							FileName: "test.file",
						},
					},
					ContainerName: appCfg.GetConfig().FileStorage.Bucket,
				},
			},
			wantErr: false,
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
			w := worker.NewWorker(tt.fields.fStorage)
			err := w.DoJob(tt.args.ctx, tt.args.payload)
			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
