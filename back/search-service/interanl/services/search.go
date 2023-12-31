package services

import (
	"context"
	"slices"

	"github.com/0x726f6f6b6965/my-blog/lib/checker"
	"github.com/0x726f6f6b6965/my-blog/lib/grpc"
	pbSearch "github.com/0x726f6f6b6965/my-blog/protos/search/v1"
	"github.com/0x726f6f6b6965/my-blog/search-service/interanl/helper"
	"github.com/0x726f6f6b6965/my-blog/search-service/interanl/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type searchService struct {
	pbSearch.UnimplementedSearchServiceServer
	cache *utils.WordDictionary
	index map[string][]string
}

// AddIndex implements v1.SearchServiceServer.
func (s *searchService) AddIndex(ctx context.Context, req *pbSearch.AddIndexRequest) (*emptypb.Empty, error) {
	if checker.IsEmpty(req.Index) {
		return &emptypb.Empty{}, grpc.BadRequestErr("index is required", "index", "index is empty")
	}
	s.cache.InsertWord(req.GetIndex())
	// indexer
	for _, token := range utils.GetTokens(req.GetIndex()) {
		if ids, ok := s.index[token]; ok && slices.Contains(ids, req.Id) {
			continue
		}
		s.index[token] = append(s.index[token], req.Id)
	}
	return &emptypb.Empty{}, nil
}

// Search implements v1.SearchServiceServer.
func (s *searchService) Search(ctx context.Context, req *pbSearch.SearchRequest) (*pbSearch.SearchResponse, error) {
	if checker.IsEmpty(req.Query) {
		return nil, grpc.BadRequestErr("query is required", "query", "query is empty")
	}
	// indexer
	var result []string
	for _, token := range utils.GetTokens(req.GetQuery()) {
		if ids, ok := s.index[token]; ok {
			if result == nil {
				result = ids
			} else {
				result = helper.Intersect(result, ids)
			}
		} else {
			// token doesn't exist
			continue
		}
	}

	return &pbSearch.SearchResponse{Ids: result}, nil
}

func (s *searchService) AutoComplete(ctx context.Context, req *pbSearch.AutoCompleteReuqest) (*pbSearch.AutoCompleteResponse, error) {
	if checker.IsEmpty(req.Words) {
		return nil, grpc.BadRequestErr("words is required", "words", "words is empty")
	}
	return &pbSearch.AutoCompleteResponse{Match: s.cache.SearchWord(req.Words)}, nil
}

func NewSearchService() pbSearch.SearchServiceServer {
	return &searchService{
		cache: utils.NewWordDictionary(),
		index: map[string][]string{},
	}
}
