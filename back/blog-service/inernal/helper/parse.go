package helper

import (
	"github.com/0x726f6f6b6965/my-blog/db/pkg/repository"
	pbBlog "github.com/0x726f6f6b6965/my-blog/protos/blog/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TBlogToPb(data *repository.TBlog) *pbBlog.Blog {
	blog := &pbBlog.Blog{
		Id:         data.ID,
		Title:      data.Title,
		Content:    data.Content,
		Author:     data.Author,
		CreateTime: timestamppb.New(data.CreateTime),
		UpdateTime: timestamppb.New(data.UpdateTime),
	}
	return blog
}
