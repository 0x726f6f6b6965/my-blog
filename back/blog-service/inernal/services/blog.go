package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/0x726f6f6b6965/my-blog/blog-service/inernal/helper"
	"github.com/0x726f6f6b6965/my-blog/blog-service/inernal/utils"
	"github.com/0x726f6f6b6965/my-blog/db/pkg/repository"
	"github.com/0x726f6f6b6965/my-blog/lib/checker"
	"github.com/0x726f6f6b6965/my-blog/lib/grpc"
	pbBlog "github.com/0x726f6f6b6965/my-blog/protos/blog/v1"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blogService struct {
	pbBlog.UnimplementedBlogServiceServer
	db *sql.DB
}

// CreateBlog: create blog
func (s *blogService) CreateBlog(ctx context.Context, req *pbBlog.CreateBlogRequest) (*pbBlog.Blog, error) {
	if checker.IsEmpty(req.Title) {
		return &pbBlog.Blog{}, grpc.BadRequestErr("title is required", "title", "title is empty")
	}
	if checker.IsEmpty(req.Author) || !checker.ValidateEmail(req.Author) {
		return &pbBlog.Blog{}, grpc.BadRequestErr("author is required", "author", "author is empty")
	}
	mod := []qm.QueryMod{qm.Select(repository.TUserColumns.Email), qm.Where(
		fmt.Sprintf("%s = ?", repository.TUserColumns.Email), req.Author)}
	_, err := repository.TUsers(mod...).OneG(ctx)
	if err != nil {
		return &pbBlog.Blog{}, grpc.InvalidErr("author is invalid", "author", "author is invalid")
	}
	data := &repository.TBlog{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}
	err = data.InsertG(ctx, boil.Infer())
	if err != nil {
		return &pbBlog.Blog{}, grpc.InternalErr("create data error, try again later")
	}
	resp := helper.TBlogToPb(data)
	return resp, nil
}

// DeleteBlog: delete specific blog
func (*blogService) DeleteBlog(ctx context.Context, req *pbBlog.DeleteBlogRequest) (*emptypb.Empty, error) {
	if checker.IsEmpty(req.Id) {
		return &emptypb.Empty{}, grpc.BadRequestErr("id is required", "id", "id is empty")
	}
	data, err := repository.FindTBlogG(ctx, req.Id)
	if err != nil || req.Author != data.Author {
		return &emptypb.Empty{}, grpc.InvalidErr("id is invalid", "id", "id is invalid")
	}
	_, err = data.DeleteG(ctx)
	if err != nil {
		return &emptypb.Empty{}, grpc.InternalErr("delete data error, try again later")
	}
	return &emptypb.Empty{}, nil
}

// EditBlog: edit exist blog
func (*blogService) EditBlog(ctx context.Context, req *pbBlog.EditBlogRequest) (*pbBlog.Blog, error) {
	if checker.IsEmpty(req.Id) {
		return &pbBlog.Blog{}, grpc.BadRequestErr("id is required", "id", "id is empty")
	}
	data, err := repository.FindTBlogG(ctx, req.Id)
	if err != nil || data.Author != req.Blog.Author {
		return &pbBlog.Blog{}, grpc.InvalidErr("id is invalid", "id", "id is invalid")
	}
	update := []string{}

	for _, key := range req.UpdateMask.GetPaths() {
		switch key {
		case "blog.title":
			data.Title = req.Blog.Title
			update = append(update, repository.TBlogColumns.Title)
		case "blog.content":
			data.Content = req.Blog.Content
			update = append(update, repository.TBlogColumns.Content)
		}
	}
	_, err = data.UpdateG(ctx, boil.Whitelist(update...))
	if err != nil {
		return &pbBlog.Blog{}, grpc.InternalErr("update data error, try again later")
	}
	resp := helper.TBlogToPb(data)
	return resp, nil
}

// GetBlog: get blog information
func (*blogService) GetBlog(ctx context.Context, req *pbBlog.GetBlogRequest) (*pbBlog.Blog, error) {
	if checker.IsEmpty(req.Id) {
		return &pbBlog.Blog{}, grpc.BadRequestErr("id is required", "id", "id is empty")
	}
	data, err := repository.FindTBlogG(ctx, req.Id)
	if err != nil {
		return &pbBlog.Blog{}, grpc.InvalidErr("id is invalid", "id", "id is invalid")
	}
	resp := helper.TBlogToPb(data)
	return resp, nil

}

// GetBlogList: get blog list
func (s *blogService) GetBlogList(ctx context.Context, req *pbBlog.GetBlogListRequest) (*pbBlog.GetBlogListResponse, error) {
	var (
		mod   = []qm.QueryMod{}
		token = &utils.PageToken{}
		size  = 25
	)
	if !checker.IsEmpty(req.PageToken) {
		utils.DecodePageTokenStruct(req.PageToken, token)
		if len(token.Authors) <= 0 {
			mod = append(mod, qm.Where("create_time <= ? AND id > ?", token.CreateTime, token.Id))
		} else {
			mod = append(mod, qm.Where("author IN (?) AND create_time <= ? AND id > ?", strings.Join(token.Authors, ","), token.CreateTime, token.Id))
		}
		if token.Size != 0 {
			size = token.Size
		}
	} else {
		if len(req.Authors) > 0 {
			mod = append(mod, qm.Where("author IN (?)", strings.Join(req.Authors, ",")))
		}

		if req.PageSize != 0 {
			size = int(req.PageSize)
		}
	}

	mod = append(mod, qm.OrderBy(fmt.Sprintf("%s %s, %s %s",
		repository.TBlogColumns.CreateTime, "DESC",
		repository.TBlogColumns.ID, "DESC")), qm.Limit(size))
	blogs, err := repository.TBlogs(mod...).AllG(ctx)
	if err != nil {
		return &pbBlog.GetBlogListResponse{}, grpc.InternalErr("please try again later")
	}

	if len(blogs) > 0 {
		resp := &pbBlog.GetBlogListResponse{}
		last := blogs[len(blogs)-1]
		token.Id = last.ID
		if len(req.Authors) > 0 {
			token.Authors = req.Authors
		}
		token.CreateTime = last.CreateTime
		token.Size = size
		resp.NextToken = token.String()
		for _, blog := range blogs {
			resp.Blog = append(resp.Blog, helper.TBlogToPb(blog))
		}
		return resp, nil
	}
	return &pbBlog.GetBlogListResponse{}, nil
}

func NewBlogService(db *sql.DB) pbBlog.BlogServiceServer {
	ser := &blogService{db: db}
	boil.SetDB(ser.db)
	return ser
}
