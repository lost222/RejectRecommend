package gRPC_client

import (
	"context"
	"fmt"
	pb "ginrss/gRPC_Client/pb"
	"ginrss/middleware"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"time"
)

var serviceAddress = "127.0.0.1:1234"

func GrpcTokenGenerate(claim middleware.MyClaims) (string, error)  {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	bookClient := pb.NewTokenServiceClient(conn)

	req := &pb.UserClaim{
		Name: claim.Username,
		NotBefore : claim.NotBefore,
		ExpiresAt: claim.ExpiresAt,
		Issuer: claim.Issuer,
	}
	reply, err := bookClient.CreateToken(context.Background(), req)

	if err != nil{
		return "", err
	}

	return reply.Token , err
}


func GrpcTokenParser(tokenString string) (*middleware.MyClaims, error) {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	bookClient := pb.NewTokenServiceClient(conn)
	req := &pb.Token{
		Token: tokenString,
	}
	UserClaim, err := bookClient.ParserToken(context.Background(), req)

	if err != nil {
		return &middleware.MyClaims{}, err
	}
	ans := middleware.MyClaims{
		Username: UserClaim.Name,
		StandardClaims: jwt.StandardClaims{
			Audience:  UserClaim.Audience,
			ExpiresAt: UserClaim.ExpiresAt,
			Id:        UserClaim.Id,
			IssuedAt:  UserClaim.IssuedAt,
			Issuer:    UserClaim.Issuer,
			NotBefore: UserClaim.NotBefore,
			Subject:   UserClaim.Subject,
		},
	}
	return &ans , err
}


//一个claim变成token又解码回来，所带数据应该前后一致
func testGrpc() bool {
	var claim = middleware.MyClaims{
		Username: "test6",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 100,
			ExpiresAt: time.Now().Unix() + 7200,
			Issuer:    "GinRss",
		},
	}

	tokenStr, err := GrpcTokenGenerate(claim)

	if err != nil{
		panic(err)
	}

	fmt.Println("token:")
	fmt.Println(tokenStr)

	var ansClaim *middleware.MyClaims

	ansClaim , err = GrpcTokenParser(tokenStr)
	if err != nil{
		panic(err)
	}
	return claim == *ansClaim
}

