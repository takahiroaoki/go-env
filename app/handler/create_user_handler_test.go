package handler

import (
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/grpc-sample/app/entity"
	"github.com/takahiroaoki/grpc-sample/app/pb"
	"github.com/takahiroaoki/grpc-sample/app/testutil"
	"github.com/takahiroaoki/grpc-sample/app/testutil/mock"
	"github.com/takahiroaoki/grpc-sample/app/util"
	"gorm.io/gorm"
)

func Test_createUserHandlerImpl_execute(t *testing.T) {
	t.Parallel()

	db, _ := testutil.GetDatabase()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockCreateUserService(ctrl)

	type fields struct {
		db                *gorm.DB
		createUserService *mock.MockCreateUserService
	}
	type args struct {
		ctx context.Context
		req *pb.CreateUserRequest
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		mockFunc       func(mockRepository *mock.MockCreateUserService)
		expected       *pb.CreateUserResponse
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name: "Success",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "user@example.com",
				},
			},
			mockFunc: func(mockService *mock.MockCreateUserService) {
				mockService.EXPECT().CreateUser(gomock.Any(), entity.User{
					Email: "user@example.com",
				}).Return(&entity.User{
					ID:    1,
					Email: "user@example.com",
				}, nil)
			},
			expected: &pb.CreateUserResponse{
				Id: "1",
			},
			expectErr: false,
		},
		{
			name: "Error(validation)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "invalid value",
				},
			},
			mockFunc: func(mockService *mock.MockCreateUserService) {
				mockService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).MaxTimes(0)
			},
			expected:       nil,
			expectErr:      true,
			expectedErrMsg: "email: must be in a valid format.",
		},
		{
			name: "Error(createUserService.CreateUser)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "user@example.com",
				},
			},
			mockFunc: func(mockService *mock.MockCreateUserService) {
				mockService.EXPECT().CreateUser(gomock.Any(), entity.User{
					Email: "user@example.com",
				}).Return(nil, util.NewError("err"))
			},
			expected:       nil,
			expectErr:      true,
			expectedErrMsg: "err",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &createUserHandlerImpl{
				db:                tt.fields.db,
				createUserService: tt.fields.createUserService,
			}
			tt.mockFunc(tt.fields.createUserService)
			actual, err := h.execute(tt.args.ctx, tt.args.req)

			assert.Equal(t, tt.expected, actual)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_createUserHandlerImpl_validate(t *testing.T) {
	t.Parallel()

	db, _ := testutil.GetDatabase()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockCreateUserService(ctrl)

	type fields struct {
		db                *gorm.DB
		createUserService *mock.MockCreateUserService
	}
	type args struct {
		ctx context.Context
		req *pb.CreateUserRequest
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		expected       error
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name: "Success",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "user@example.com",
				},
			},
			expectErr: false,
		},
		{
			name: "Success(Email right boundary safe)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: strings.Repeat("a", 308) + "@example.com",
				},
			},
			expectErr: false,
		},
		{
			name: "Error(Email right boundary over)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: strings.Repeat("a", 309) + "@example.com",
				},
			},
			expectErr:      true,
			expectedErrMsg: "email: the length must be between 1 and 320.",
		},
		{
			name: "Error(Email is nil)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{},
			},
			expectErr:      true,
			expectedErrMsg: "email: cannot be blank.",
		},
		{
			name: "Error(Email is empty)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "",
				},
			},
			expectErr:      true,
			expectedErrMsg: "email: cannot be blank.",
		},
		{
			name: "Error(Email is in an invalid format)",
			fields: fields{
				db:                db,
				createUserService: mockService,
			},
			args: args{
				ctx: context.Background(),
				req: &pb.CreateUserRequest{
					Email: "invalid format",
				},
			},
			expectErr:      true,
			expectedErrMsg: "email: must be in a valid format.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := &createUserHandlerImpl{
				db:                tt.fields.db,
				createUserService: tt.fields.createUserService,
			}
			err := h.validate(tt.args.ctx, tt.args.req)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErrMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
