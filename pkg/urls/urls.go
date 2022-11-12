package urls

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/mike9107/urlshortener/pkg/grpc/interceptor"
	"github.com/mike9107/urlshortener/pkg/protobuf/urlspb"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

var (
	port   = flag.Int("port", 8081, "the port to serve on")
	urlsDB = flag.String("urls_db", "urls.db", "the urls db")
)

func NewService() *urlServer {
	flag.Parse()

	server := &urlServer{
		db: initSqliteDB(*urlsDB),
	}
	return server
}

func (server *urlServer) Run() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcInterceptor := interceptor.NewGrpcInterceptor()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcInterceptor.UnaryLogger),
		grpc.StreamInterceptor(grpcInterceptor.StreamLogger))

	urlspb.RegisterUrlsServer(grpcServer, server)

	log.Printf("Starting urls_service at %v\n", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
