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
	v1 "github.com/0x726f6f6b6965/my-blog/protos/blog/v1"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var (
	blog v1.BlogServiceServer
	db   *sql.DB
	mock sqlmock.Sqlmock
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	db, mock, _ = sqlmock.New()
	blog = NewBlogService(db)
	boil.DebugMode = true
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	blog = nil
	db.Close()
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestCreateBlog(t *testing.T) {
	req := &v1.CreateBlogRequest{
		Title:   "abc",
		Content: "test-content",
		Author:  "abc@gmail.com",
	}
	mod := []qm.QueryMod{qm.Select(repository.TUserColumns.Email), qm.Where(
		fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.Author)}
	query := repository.TUsers(mod...).Query
	rawsql, args := queries.BuildQuery(query)
	rawsql = regexp.QuoteMeta(rawsql)
	// OneG
	rawsql = rawsql[0:len(rawsql)-1] + " LIMIT 1;"
	user := mock.NewRows([]string{repository.TUserColumns.Email}).AddRow(req.Author)
	mock.ExpectQuery(rawsql).WithArgs(args[0].(string)).WillReturnRows(user)

	insert := "INSERT INTO \"t_blog\" (\"title\",\"content\",\"author\") VALUES ($1,$2,$3) RETURNING \"id\",\"create_time\",\"update_time\""
	insert, args = queries.BuildQuery(queries.Raw(insert, req.Title, req.Content, req.Author))
	insert = regexp.QuoteMeta(insert)
	rows := mock.NewRows([]string{"id", "create_time", "update_time"}).AddRow(uuid.New().String(), time.Now(), time.Now())
	mock.ExpectQuery(insert).WithArgs(args[0].(string), args[1].(string), args[2].(string)).WillReturnRows(rows)

	resp, err := blog.CreateBlog(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Content, resp.Content)
	assert.Equal(t, req.Author, resp.Author)
}

func TestGetBlog(t *testing.T) {
	req := &v1.GetBlogRequest{Id: uuid.NewString()}

	rawsql := "select * from \"t_blog\" where \"id\"=$1"
	rawsql = regexp.QuoteMeta(rawsql)
	rows := mock.NewRows([]string{repository.TBlogColumns.ID,
		repository.TBlogColumns.Author, repository.TBlogColumns.Content, repository.TBlogColumns.Title,
		repository.TBlogColumns.CreateTime, repository.TBlogColumns.UpdateTime}).
		AddRow(req.GetId(), "abc@gmail.com", "testcontent", "title", time.Now(), time.Now())
	mock.ExpectQuery(rawsql).WithArgs(req.GetId()).WillReturnRows(rows)
	resp, err := blog.GetBlog(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, "title", resp.Title)
	assert.Equal(t, "testcontent", resp.Content)
	assert.Equal(t, "abc@gmail.com", resp.Author)
}

func TestDeleteBlog(t *testing.T) {
	req := &v1.DeleteBlogRequest{
		Id:     uuid.NewString(),
		Author: "abc@gmail.com",
	}
	rawsql := "select * from \"t_blog\" where \"id\"=$1"
	rawsql = regexp.QuoteMeta(rawsql)
	rows := mock.NewRows([]string{repository.TBlogColumns.ID,
		repository.TBlogColumns.Author, repository.TBlogColumns.Content, repository.TBlogColumns.Title,
		repository.TBlogColumns.CreateTime, repository.TBlogColumns.UpdateTime}).
		AddRow(req.GetId(), "abc@gmail.com", "testcontent", "title", time.Now(), time.Now())
	mock.ExpectQuery(rawsql).WithArgs(req.GetId()).WillReturnRows(rows)
	delSql := "DELETE FROM \"t_blog\" WHERE \"id\"=$1"
	delSql = regexp.QuoteMeta(delSql)
	mock.ExpectExec(delSql).WithArgs(req.GetId()).WillReturnResult(sqlmock.NewResult(1, 1))
	_, err := blog.DeleteBlog(context.Background(), req)
	assert.Nil(t, err)
}

func TestEditBlog(t *testing.T) {
	req := &v1.EditBlogRequest{
		Id: uuid.NewString(),
		Blog: &v1.Blog{
			Title:  "bbbb",
			Author: "abc@gmail.com",
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"blog.title"}},
	}
	rawsql := "select * from \"t_blog\" where \"id\"=$1"
	rawsql = regexp.QuoteMeta(rawsql)
	rows := mock.NewRows([]string{repository.TBlogColumns.ID,
		repository.TBlogColumns.Author, repository.TBlogColumns.Content, repository.TBlogColumns.Title,
		repository.TBlogColumns.CreateTime, repository.TBlogColumns.UpdateTime}).
		AddRow(req.GetId(), "abc@gmail.com", "testcontent", "title", time.Now(), time.Now())
	mock.ExpectQuery(rawsql).WithArgs(req.GetId()).WillReturnRows(rows)
	updateSQL := "UPDATE \"t_blog\" SET \"title\"=$1 WHERE \"id\"=$2"
	updateSQL = regexp.QuoteMeta(updateSQL)
	mock.ExpectExec(updateSQL).WithArgs(req.Blog.Title, req.GetId()).WillReturnResult(sqlmock.NewResult(1, 1))
	resp, err := blog.EditBlog(context.Background(), req)
	assert.Nil(t, err)
	assert.Equal(t, req.Blog.Title, resp.Title)
}
