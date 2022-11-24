package urls_test

import (
	"context"
	"net"
	"testing"

	"urlshortener/pkg/proto/urlspb"
	"urlshortener/pkg/services/urls"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var logger = zap.Must(zap.NewDevelopment())

func TestUrlsService(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer()
	urlspb.RegisterUrlsServiceServer(srv, urls.NewService())
	go func() {
		if err := srv.Serve(lis); err != nil {
			logger.Fatal("Server exited with error:", zap.Error(err))
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return lis.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("Failed to dial bufnet:", zap.Error(err))
	}
	defer conn.Close()

	redirectUrl := "https://www.example.com"
	client := urlspb.NewUrlsServiceClient(conn)
	createRes, err := client.CreateUrl(context.Background(), &urlspb.CreateUrlRequest{RedirectUrl: redirectUrl})
	if err != nil {
		t.Fatal("Failed to create url:", err)
	}
	t.Log("Created url:", createRes.GetUrlId())

	getRes, err := client.GetUrl(context.Background(), &urlspb.GetUrlRequest{UrlId: createRes.GetUrlId()})
	if err != nil {
		t.Fatal("Failed to get url:", err)
	}
	t.Log("Got url:", getRes.GetUrl())

	if getRes.GetUrl().GetRedirectUrl() != redirectUrl {
		t.Fatal("Got wrong url:", getRes.Url)
	}
}
