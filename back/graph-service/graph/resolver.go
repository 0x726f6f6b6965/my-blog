package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	pbBlog "github.com/0x726f6f6b6965/my-blog/protos/blog/v1"
	pbSearch "github.com/0x726f6f6b6965/my-blog/protos/search/v1"
	pbUser "github.com/0x726f6f6b6965/my-blog/protos/user/v1"
	"go.uber.org/zap"
)

type Resolver struct {
	BlogService   pbBlog.BlogServiceClient
	SearchService pbSearch.SearchServiceClient
	UserService   pbUser.UserServiceClient
	Log           *zap.Logger
}
