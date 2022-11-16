/*******
* @Author:qingmeng
* @Description:
* @File:authInterceptor
* @Date2022/7/15
 */

package api

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

//获取拦截器
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	log.Printf("RPC: %s,req:%v start time: %s, end time: %s, err: %v", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}


func getAuthInterceptor(auth func(ctx context.Context)error) grpc.UnaryServerInterceptor {
	authInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		//拦截普通方法请求，验证token
		err = auth(ctx)
		if err != nil {
			return
		}
		return handler(ctx, req)
	}
	return authInterceptor
}

//认证方法
func authToken(ctx context.Context) error {
	//拿到传输的用户名和密码
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}

	var user string
	var password string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}
	if user != "qi" || password != "qi123456" {
		return status.Errorf(codes.Unauthenticated, "token不合法")
	}
	return nil

}

