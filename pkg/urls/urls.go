package urls

import (
	context "context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"urlshortener/pkg/api"

	"urlshortener/pkg/interceptor"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

type urlServer struct {
	UnimplementedUrlsServer
	db *sql.DB
}

var (
	port   = flag.Int("url_port", 50051, "the port to serve on")
	url_db = flag.String("url_db", "url.db", "the url db")
)

func RunService() {

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcInterceptor := interceptor.NewGrpcInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor.UnaryLogger),
		grpc.StreamInterceptor(grpcInterceptor.StreamLogger))

	server := &urlServer{
		db: initSqliteDB(*url_db),
	}

	RegisterUrlsServer(grpcServer, server)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (server *urlServer) GetUrl(ctx context.Context, req *GetUrlRequest) (*GetUrlResponse, error) {
	urlEntity, err := getUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, err
	}
	return &GetUrlResponse{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}, nil
}

func (server *urlServer) CreateUrl(ctx context.Context, req *CreateUrlRequest) (*CreateUrlResponse, error) {
	urlId, err := generateID()
	if err != nil {
		return nil, err
	}
	err = createUrl(server.db, ctx, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, err
	}
	return &CreateUrlResponse{UrlId: urlId}, nil
}

func (server *urlServer) UpdateUrl(ctx context.Context, req *UpdateUrlRequest) (*api.NoContent, error) {
	err := updateUrl(server.db, ctx, req.GetUrlId(), req.GetRedirectUrl())
	if err != nil {
		return nil, err
	}
	return &api.NoContent{}, nil
}

func (server *urlServer) DeleteUrl(ctx context.Context, req *DeleteUrlRequest) (*api.NoContent, error) {
	err := deleteUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, err
	}
	return &api.NoContent{}, nil
}

func (server *urlServer) Ping(ctx context.Context, req *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Message: "pong"}, nil
}

func generateID() (string, error) {
	var data [6]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}
