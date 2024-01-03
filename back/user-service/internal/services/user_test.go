package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/0x726f6f6b6965/my-blog/db/pkg/repository"
	pbUser "github.com/0x726f6f6b6965/my-blog/protos/user/v1"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/client"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/utils"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	user                 = &userService{}
	expire time.Duration = time.Second
	db     *sql.DB
	rds    *redis.Client
	mock   sqlmock.Sqlmock
	rmock  redismock.ClientMock
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
func setup() {
	db, mock, _ = sqlmock.New()
	rds, rmock = redismock.NewClientMock()
	user.db = db
	user.rds = rds
	user.expire = expire
	boil.SetDB(db)
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	user = nil
	db.Close()
	rds.Close()
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestCreateUser(t *testing.T) {
	req := &pbUser.CreateUserRequest{
		Username: "abc",
		Email:    "a123@gmail.com",
		Password: "123",
	}
	// mock not exist
	query := repository.TUsers(qm.Select(repository.TUserColumns.Email), qm.Where(
		fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.Email)).Query
	rawsql, args := queries.BuildQuery(query)
	rawsql = regexp.QuoteMeta(rawsql)
	// OneG
	rawsql = rawsql[0:len(rawsql)-1] + " LIMIT 1;"
	mock.ExpectQuery(rawsql).WithArgs(args[0].(string)).WillReturnError(sql.ErrNoRows)

	salt := utils.CreateNewSalt()
	user.testSalt = &salt
	pwd, _ := salt.SaltInput(req.Password)
	// mock insert
	insert := "INSERT INTO \"t_user\" (\"username\",\"email\",\"salt\",\"password\") VALUES ($1,$2,$3,$4) RETURNING \"id\",\"create_time\",\"update_time\""
	insert, args = queries.BuildQuery(queries.Raw(insert, req.Username, req.Email, salt.SaltString, pwd))
	insert = regexp.QuoteMeta(insert)
	rows := mock.NewRows([]string{"id", "create_time", "update_time"}).AddRow(uuid.New().String(), time.Now(), time.Now())
	mock.ExpectQuery(insert).WithArgs(args[0].(string), args[1].(string), args[2].(string), args[3].(string)).WillReturnRows(rows)

	// mock sercet
	rmock.ExpectGet(client.Secret).SetVal("test-secret")

	resp, err := user.CreateUser(context.Background(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestGetTokenFromRedis(t *testing.T) {
	req := &pbUser.GetTokenRequest{
		Email:    "a123@gmail.com",
		Password: "123",
	}
	// mock token
	rmock.ExpectGet(fmt.Sprintf(client.UserToken, req.Email)).SetVal("test-token")
	resp, err := user.GetToken(context.Background(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestGetToken(t *testing.T) {
	req := &pbUser.GetTokenRequest{
		Email:    "a123@gmail.com",
		Password: "123",
	}

	salt := utils.CreateNewSalt()
	pwd, _ := salt.SaltInput(req.Password)

	// mock not exist
	mod := []qm.QueryMod{
		qm.Select(repository.TUserColumns.Email, repository.TUserColumns.Salt, repository.TUserColumns.Password),
		qm.Where(fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.GetEmail()),
	}
	query := repository.TUsers(mod...).Query
	rawsql, args := queries.BuildQuery(query)
	rawsql = regexp.QuoteMeta(rawsql)
	// OneG
	rawsql = rawsql[0:len(rawsql)-1] + " LIMIT 1;"
	rows := mock.NewRows([]string{repository.TUserColumns.Email, repository.TUserColumns.Salt, repository.TUserColumns.Password}).
		AddRow(req.GetEmail(), salt.SaltString, pwd)
	mock.ExpectQuery(rawsql).WithArgs(args[0].(string)).WillReturnRows(rows)

	// mock sercet
	rmock.ExpectGet(client.Secret).SetVal("test-secret")

	resp, err := user.GetToken(context.Background(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestUpdateToken(t *testing.T) {
	req := &pbUser.UpdateTokenRequest{
		Email: "a123@gmail.com",
	}
	// mock token
	rmock.ExpectGet(fmt.Sprintf(client.UserToken, req.Email)).SetVal("test-token")
	// mock sercet
	rmock.ExpectGet(client.Secret).SetVal("test-secret")
	resp, err := user.UpdateToken(context.Background(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Token)
}
