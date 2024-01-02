package services

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/0x726f6f6b6965/my-blog/db/pkg/repository"
	"github.com/0x726f6f6b6965/my-blog/lib/config"
	pbSearch "github.com/0x726f6f6b6965/my-blog/protos/search/v1"
	"github.com/0x726f6f6b6965/my-blog/search-service/interanl/client"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
)

var (
	search pbSearch.SearchServiceServer
	add1   = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "What's the best thing about Switzerland?",
	}
	add2 = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "I don't know, but the flag is a big plus.",
	}
	add3 = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "I invented a new word!",
	}
	add4 = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "Plagiarism!",
	}
	add5 = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "Did you hear about the mathematician who's afraid of negative numbers?",
	}
	add6 = pbSearch.AddIndexRequest{
		Id:    uuid.New().String(),
		Index: "He'll stop at nothing to avoid them.",
	}
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
func setup() {
	search = NewSearchService()
	ctx := context.Background()
	search.AddIndex(ctx, &add1)
	search.AddIndex(ctx, &add2)
	search.AddIndex(ctx, &add3)
	search.AddIndex(ctx, &add4)
	search.AddIndex(ctx, &add5)
	search.AddIndex(ctx, &add6)
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	search = nil
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestSearchOneResult(t *testing.T) {
	ctx := context.Background()
	result, err := search.Search(ctx, &pbSearch.SearchRequest{Query: "big plus"})
	assert.Nil(t, err)
	assert.Contains(t, result.Ids, add2.Id)
	assert.Equal(t, len(result.Ids), 1)
}

func TestSearchMultiResult(t *testing.T) {
	ctx := context.Background()
	result, err := search.Search(ctx, &pbSearch.SearchRequest{Query: "about"})
	assert.Nil(t, err)
	assert.Contains(t, result.Ids, add1.Id)
	assert.Contains(t, result.Ids, add5.Id)
	assert.Equal(t, len(result.Ids), 2)
}

func TestAutoCompleteOne(t *testing.T) {
	ctx := context.Background()
	result, err := search.AutoComplete(ctx, &pbSearch.AutoCompleteRequest{Words: "What's"})
	assert.Nil(t, err)
	assert.Contains(t, result.Match, add1.Index)
	assert.Equal(t, len(result.Match), 1)
}
func TestAutoCompleteMulti(t *testing.T) {
	ctx := context.Background()
	result, err := search.AutoComplete(ctx, &pbSearch.AutoCompleteRequest{Words: "I"})
	assert.Nil(t, err)
	assert.Contains(t, result.Match, add2.Index)
	assert.Contains(t, result.Match, add3.Index)
	assert.Equal(t, len(result.Match), 2)
}

func TestRealDBLoad(t *testing.T) {
	ctx := context.Background()
	db, cleanup, _ := client.NewPostgres(&config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "docker",
		DBName:   "postgres",
		SSLmode:  "disable",
	})
	defer cleanup()
	logger, _ := zap.NewDevelopment()
	SetLoadStorageFunc(10, db, logger, search.AddIndex)
	boil.SetDB(db)
	testUser := repository.TUser{Username: "test-user",
		Email: "abc@gmail.com", Salt: uuid.NewString(), Password: uuid.NewString()}
	_ = testUser.InsertG(ctx, boil.Infer())
	defer testUser.DeleteG(ctx)
	testBlog := repository.TBlog{
		Title:   "Test blog title",
		Content: "Test content",
		Author:  testUser.Email,
	}
	_ = testBlog.InsertG(ctx, boil.Infer())
	defer testBlog.DeleteG(ctx)
	LoadStorageFunc()
	resp, err := search.Search(ctx, &pbSearch.SearchRequest{Query: testBlog.Title})
	assert.Nil(t, err)
	assert.Contains(t, resp.Ids, testBlog.ID)

}
