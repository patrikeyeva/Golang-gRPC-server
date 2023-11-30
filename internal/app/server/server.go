package server

import (
	"context"
	"encoding/json"
	"errors"
	"homework8/internal/pkg/repository"
	pb "homework8/pkg/grpcserver"
	"homework8/pkg/logger"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ServerImplement struct {
	pb.UnimplementedDbHandlerServer
	ArticleRepo repository.ArticlesRepo
	CommentRepo repository.CommentsRepo
	Sender      *KafkaSender
}

func (server *ServerImplement) SendToKafka(ctx context.Context, path string, method string, reqBody []byte) error {
	message, err := NewServerMessage(path, method, reqBody)
	if err != nil {
		return err
	}
	messageToKafka, err := BuildMessage(message, server.Sender.topic)
	if err != nil {
		return err
	}
	err = server.Sender.SendMessage(ctx, messageToKafka)
	if err != nil {
		return err
	}
	return nil
}

func (server *ServerImplement) CreateArticle(ctx context.Context, request *pb.ArticleRequest) (*pb.ArticleResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "app: CreateArticle")
	defer span.Finish()

	sendKafkaSpan := opentracing.StartSpan("app: CreateArticle-sendkafka", opentracing.ChildOf(span.Context()))
	go func(ctx context.Context, span opentracing.Span) {
		defer span.Finish()

		body, errJson := json.Marshal(&request)
		if errJson != nil {
			logger.Errorf(ctx, errJson.Error())
		}
		err := server.SendToKafka(ctx, "article", "POST", body)
		if err != nil {
			logger.Errorf(ctx, "Error when trying to send a request to kafka", err)
		}
	}(ctx, sendKafkaSpan)

	article := &repository.Article{
		Name:   request.Name,
		Rating: request.Rating,
	}

	article, err := server.ArticleRepo.Add(ctx, article)
	if err != nil {
		return &pb.ArticleResponse{}, status.Error(codes.Internal, err.Error())
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		return &pb.ArticleResponse{}, status.Error(codes.Internal, err.Error())
	}
	return &pb.ArticleResponse{Answer: string(articleJson)}, nil
}

func (server *ServerImplement) GetArticle(ctx context.Context, request *pb.ArticleID) (*pb.GetArticleResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "app: GetArticle")
	defer span.Finish()

	sendKafkaSpan := opentracing.StartSpan("app: GetArticle-sendkafka", opentracing.ChildOf(span.Context()))
	go func(ctx context.Context, span opentracing.Span) {
		defer span.Finish()

		err := server.SendToKafka(ctx, "article", "GET", nil)
		if err != nil {
			logger.Errorf(ctx, "Error when trying to send a request to kafka", err)
		}
	}(ctx, sendKafkaSpan)

	//Get article data
	span, ctx = opentracing.StartSpanFromContext(ctx, "app: GetArticle-article")
	article, err := server.ArticleRepo.GetByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return &pb.GetArticleResponse{}, status.Error(codes.NotFound, err.Error())
		}
		return &pb.GetArticleResponse{}, status.Error(codes.Internal, err.Error())
	}
	articleJson, err := json.Marshal(&article)
	if err != nil {
		return &pb.GetArticleResponse{}, status.Error(codes.Internal, err.Error())
	}
	span.Finish()

	//Get comments for article
	span, ctx = opentracing.StartSpanFromContext(ctx, "app: GetArticle-comments")
	comments, err := server.CommentRepo.GetCommentsForArticle(ctx, request.Id)
	if err != nil {
		return &pb.GetArticleResponse{}, status.Error(codes.Internal, err.Error())
	}
	var commentsData []byte

	for idx, comment := range comments {
		commentJson, err := json.Marshal(&comment)
		if err != nil {
			return &pb.GetArticleResponse{}, status.Error(codes.Internal, err.Error())
		}
		if idx > 0 {
			commentsData = append(commentsData, []byte(",  ")...)
		}
		commentsData = append(commentsData, commentJson...)
	}
	span.Finish()

	return &pb.GetArticleResponse{Answer: string(articleJson), Comments: string(commentsData)}, nil

}

func (server *ServerImplement) DeleteArticle(ctx context.Context, request *pb.ArticleID) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "app: DeleteArticle")
	defer span.Finish()

	sendKafkaSpan := opentracing.StartSpan("app: DeleteArticle-sendkafka", opentracing.ChildOf(span.Context()))
	go func(ctx context.Context, span opentracing.Span) {
		defer span.Finish()
		err := server.SendToKafka(ctx, "article", "DELETE", nil)
		if err != nil {
			logger.Errorf(ctx, "Error when trying to send a request to kafka", err)
		}
	}(ctx, sendKafkaSpan)

	err := server.ArticleRepo.DeleteByID(ctx, request.Id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		} else {
			return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (server *ServerImplement) UpdateArticle(ctx context.Context, request *pb.UpdateRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "app: UpdateArticle")
	defer span.Finish()

	sendKafkaSpan := opentracing.StartSpan("app: UpdateArticle-sendkafka", opentracing.ChildOf(span.Context()))
	go func(ctx context.Context, span opentracing.Span) {
		defer span.Finish()
		body, errJson := json.Marshal(&request)
		if errJson != nil {
			logger.Errorf(ctx, errJson.Error())
		}
		err := server.SendToKafka(ctx, "article", "PUT", body)
		if err != nil {
			logger.Errorf(ctx, "Error when trying to send a request to kafka", err)
		}
	}(ctx, sendKafkaSpan)

	articleRepo := &repository.Article{
		ID:     request.Id,
		Name:   request.Name,
		Rating: request.Rating,
	}

	if err := server.ArticleRepo.Update(ctx, articleRepo); err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return &emptypb.Empty{}, status.Error(codes.NotFound, err.Error())
		} else {
			return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
		}
	}
	return &emptypb.Empty{}, nil
}

func (server *ServerImplement) CreateComment(ctx context.Context, request *pb.CommentRequest) (*pb.CommentResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "app: CreateComment")
	defer span.Finish()

	sendKafkaSpan := opentracing.StartSpan("app: CreateComment-sendkafka", opentracing.ChildOf(span.Context()))
	go func(ctx context.Context, span opentracing.Span) {
		defer span.Finish()
		body, errJson := json.Marshal(&request)
		if errJson != nil {
			logger.Errorf(ctx, errJson.Error())
		}
		err := server.SendToKafka(ctx, "comment", "POST", body)
		if err != nil {
			logger.Errorf(ctx, "Error when trying to send a request to kafka", err)
		}
	}(ctx, sendKafkaSpan)

	commentRepo := &repository.Comment{
		ArticleID: request.ArticleId,
		Text:      request.Text,
	}

	commentRepo, err := server.CommentRepo.AddComment(ctx, commentRepo)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return &pb.CommentResponse{}, status.Error(codes.NotFound, err.Error())
		}
		return &pb.CommentResponse{}, status.Error(codes.Internal, err.Error())

	}
	commentJson, err := json.Marshal(commentRepo)
	if err != nil {
		return &pb.CommentResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &pb.CommentResponse{Answer: string(commentJson)}, nil
}
