package gRpcUtils

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func fetchOAuthToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "my-access-token",
	}
}

func GetTokenSource() oauth.TokenSource {
	return oauth.TokenSource{TokenSource: oauth2.StaticTokenSource(fetchOAuthToken())}
}

func EnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "metadata 가 없습니다.")
	}

	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "잘못된 token")
	}
	return handler(ctx, req)
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")

	fmt.Println("token", token)
	return token == "my-access-token"
}
