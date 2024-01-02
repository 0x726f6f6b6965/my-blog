package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/0x726f6f6b6965/my-blog/db/pkg/repository"
	"github.com/0x726f6f6b6965/my-blog/lib/checker"
	"github.com/0x726f6f6b6965/my-blog/lib/grpc"
	pbUser "github.com/0x726f6f6b6965/my-blog/protos/user/v1"
	"github.com/redis/go-redis/v9"

	"github.com/0x726f6f6b6965/my-blog/user-service/internal/client"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/helper"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/utils"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type userService struct {
	pbUser.UnimplementedUserServiceServer
	expire time.Duration
	rds    *redis.Client
	db     *sql.DB
}

// CreateUser: create a new user account.
func (s *userService) CreateUser(ctx context.Context, req *pbUser.CreateUserRequest) (*pbUser.Token, error) {
	if checker.IsEmpty(req.GetUsername()) {
		return &pbUser.Token{}, grpc.BadRequestErr("username is required", "username", "username is empty")
	}
	if checker.IsEmpty(req.GetEmail()) {
		return &pbUser.Token{}, grpc.BadRequestErr("email is required", "email", "email is empty")
	}
	if checker.ValidateEmail(req.GetEmail()) {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	if checker.IsEmpty(req.GetPassword()) {
		return &pbUser.Token{}, grpc.BadRequestErr("password is required", "password", "password is empty")
	}

	token, _ := client.GetToken(ctx, req.GetEmail(), s.rds)
	if !checker.IsEmpty(token) {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	exist, err := repository.TUsers(qm.Where(
		fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.Email)).Exists(ctx, s.db)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("database error, try again later")
	}
	if exist {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	salt := utils.CreateNewSalt()
	pwd, err := salt.SaltInput(req.Password)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("salt error, try again later")
	}
	data := repository.TUser{
		Username: req.Username,
		Email:    req.Email,
		Password: pwd,
		Salt:     salt.SaltString,
	}
	err = data.InsertG(ctx, boil.Infer())
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create data error, try again later")
	}
	secret, err := client.GetSecret(ctx, s.rds)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	token, err = helper.GenerateJWT(req.Email, secret, s.expire)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	client.SetToken(ctx, req.GetEmail(), token,
		helper.GeneralDuration(s.expire, 3, 10, time.Minute), s.rds)
	return &pbUser.Token{Token: token}, nil
}

// GetToken: get a validate token based on user information.
func (s *userService) GetToken(ctx context.Context, req *pbUser.GetTokenRequest) (*pbUser.Token, error) {
	if checker.IsEmpty(req.GetEmail()) {
		return &pbUser.Token{}, grpc.BadRequestErr("email is required", "email", "email is empty")
	}
	if checker.ValidateEmail(req.GetEmail()) {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	if checker.IsEmpty(req.GetPassword()) {
		return &pbUser.Token{}, grpc.BadRequestErr("password is required", "password", "password is empty")
	}
	token, _ := client.GetToken(ctx, req.GetEmail(), s.rds)
	if !checker.IsEmpty(token) {
		return &pbUser.Token{Token: token}, nil
	}
	mod := []qm.QueryMod{
		qm.Select(repository.TUserColumns.Email, repository.TUserColumns.Salt, repository.TUserColumns.Password),
		qm.Where(fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.GetEmail()),
	}
	info, err := repository.TUsers(mod...).OneG(ctx)
	if err != nil {
		return &pbUser.Token{}, nil
	}
	salt, err := utils.CreateSaltByString(info.Salt)
	if err != nil {
		return &pbUser.Token{}, nil
	}
	if pwd, _ := salt.SaltInput(req.Password); pwd != info.Password {
		return &pbUser.Token{}, nil
	}
	secret, err := client.GetSecret(ctx, s.rds)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	token, err = helper.GenerateJWT(req.Email, secret, s.expire)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	client.SetToken(ctx, req.GetEmail(), token,
		helper.GeneralDuration(s.expire, 3, 10, time.Minute), s.rds)
	return &pbUser.Token{Token: token}, nil
}

// UpdateToken: update a used verified token to extend its expiration.
func (s *userService) UpdateToken(ctx context.Context, req *pbUser.UpdateTokenRequest) (*pbUser.Token, error) {
	if checker.IsEmpty(req.GetEmail()) {
		return &pbUser.Token{}, grpc.BadRequestErr("email is required", "email", "email is empty")
	}
	if checker.ValidateEmail(req.GetEmail()) {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	token, err := client.GetToken(ctx, req.GetEmail(), s.rds)
	if err != nil {
		return &pbUser.Token{}, grpc.InvalidErr("email invalid", "email", "email invalid")
	}
	if checker.IsEmpty(token) {
		return &pbUser.Token{}, grpc.BadRequestErr("please login", "email", "email invalid")
	}
	secret, err := client.GetSecret(ctx, s.rds)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	token, err = helper.GenerateJWT(req.Email, secret, s.expire)
	if err != nil {
		return &pbUser.Token{}, grpc.InternalErr("create token error, try again later")
	}
	client.SetToken(ctx, req.GetEmail(), token,
		helper.GeneralDuration(s.expire, 3, 10, time.Minute), s.rds)
	return &pbUser.Token{Token: token}, nil
}

func NewUserService(expire time.Duration, db *sql.DB, rds *redis.Client) pbUser.UserServiceServer {
	return &userService{
		expire: expire,
		db:     db,
		rds:    rds,
	}
}
