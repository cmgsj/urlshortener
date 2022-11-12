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
	"urlshortener/pkg/grpc/interceptor"
	"urlshortener/pkg/proto/apipb"
	"urlshortener/pkg/proto/urlspb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

type urlServer struct {
	urlspb.UnimplementedUrlsServer
	db *sql.DB
}

var (
	port    = flag.Int("port", 8081, "the port to serve on")
	urls_db = flag.String("urls_db", "urls.db", "the urls db")
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

	var server urlspb.UrlsServer = &urlServer{
		db: initSqliteDB(*urls_db),
	}

	urlspb.RegisterUrlsServer(grpcServer, server)

	log.Printf("Starting urls_service at %v\n", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (server *urlServer) GetUrl(ctx context.Context, req *urlspb.GetUrlRequest) (*urlspb.GetUrlResponse, error) {
	urlEntity, err := getUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, err
	}
	return &urlspb.GetUrlResponse{UrlId: urlEntity.UrlId, RedirectUrl: urlEntity.RedirectUrl}, nil
}

func (server *urlServer) CreateUrl(ctx context.Context, req *urlspb.CreateUrlRequest) (*urlspb.CreateUrlResponse, error) {
	urlId, err := generateID()
	if err != nil {
		return nil, err
	}
	err = createUrl(server.db, ctx, urlId, req.GetRedirectUrl())
	if err != nil {
		return nil, err
	}
	return &urlspb.CreateUrlResponse{UrlId: urlId}, nil
}

func (server *urlServer) UpdateUrl(ctx context.Context, req *urlspb.UpdateUrlRequest) (*apipb.NoContent, error) {
	err := updateUrl(server.db, ctx, req.GetUrlId(), req.GetRedirectUrl())
	if err != nil {
		return nil, err
	}
	return &apipb.NoContent{}, nil
}

func (server *urlServer) DeleteUrl(ctx context.Context, req *urlspb.DeleteUrlRequest) (*apipb.NoContent, error) {
	err := deleteUrl(server.db, ctx, req.GetUrlId())
	if err != nil {
		return nil, err
	}
	return &apipb.NoContent{}, nil
}

func (server *urlServer) Ping(ctx context.Context, req *apipb.PingRequest) (*apipb.PingResponse, error) {
	return &apipb.PingResponse{Message: "pong"}, nil
}

func generateID() (string, error) {
	var data [6]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data[:]), nil
}
