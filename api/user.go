/*******
* @Author:qingmeng
* @Description:
* @File:server
* @Date2022/7/15
 */

package api

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-center/model"
	pb2 "user-center/pbfile/pb"
	"user-center/service"
	tool2 "user-center/tool"
)

func UserServer(tcpPort string) {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}
	creds := getCreds()
	//添加拦截器
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor))
	pb2.RegisterUserServiceServer(s, &server{})

	//开始处理
	if err = s.Serve(lis); err != nil {
		log.Fatalln("failed to server:", err)
	}

}

type server struct {
	pb2.UnimplementedUserServiceServer
}

func (s *server) Login(ctx context.Context, req *pb2.LoginReq) (resp *pb2.LoginResp, err error) {
	//必须先手动为嵌套的结构体指针分配指向，即初始化一个内存，返回的是指针类型时也得先分配，他只是声明了结构体的指针，不然爆空指针异常
	resp = new(pb2.LoginResp)
	resp.Resp = new(pb2.Resp)

	us := service.UserService{}
	ok, err := us.IsExistUsername(req.Username)
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERNAMEERROR
		return resp, nil
	}
	user, err := us.GetUserByUsername(req.Username)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return
	}

	//验证密码
	if !tool2.CheckPassword(user.Password, req.Password) {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.PASSWORERROR
		return resp, nil
	}

	resp.Resp.Status = true
	resp.Resp.Data = tool2.SUCCESS
	return
}

func (s *server) Register(ctx context.Context, req *pb2.RegisterReq) (resp *pb2.RegisterResp, err error) {
	resp = new(pb2.RegisterResp)
	resp.Resp = new(pb2.Resp)

	us := service.UserService{}

	//用户名验证
	ok, err := us.IsExistUsername(req.User.Username)
	if err != nil {
		log.Println("is exist user err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return
	}
	//验证敏感词和sql注入
	if ok || tool2.CheckIfSensitive(req.User.Username) {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERNAMEERROR
		return resp, nil
	}

	//密码处理
	level := tool2.CheckPasswordLever(req.User.Password)
	if level < 1 {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.PASSWORERROR
		return resp, nil
	}
	pwd, err := tool2.AddSalt(req.User.Password)
	if err != nil {
		log.Println("add salt err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return
	}

	//创建实体
	user := model.User{
		Username: req.User.Username,
		Password: pwd,
		Gender:   int(req.User.Gender),
		Name:     req.User.Name,
		Phone:    req.User.Phone,
		Email:    req.User.Email,
		GroupId:  int(req.User.GroupId),
	}

	err = us.CreateUser(&user)
	if err != nil {
		log.Println("create user err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return
	}
	resp.Resp.Status = true
	resp.Resp.Data = tool2.SUCCESS
	return resp, nil
}

func (s *server) ChangePwd(ctx context.Context, req *pb2.ChangePwdReq) (*pb2.ChangePwdResp, error) {
	resp := new(pb2.ChangePwdResp)
	resp.Resp = new(pb2.Resp)
	us := service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERIDERROR
		return resp, nil
	}
	user, err := us.GetUserById(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}

	//验证密码
	if !tool2.CheckPassword(user.Password, req.OldPassword) {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.PASSWORERROR
		return resp, nil
	}

	//更改密码
	pwd, err := tool2.AddSalt(req.NewPassword)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	user.Password = pwd
	_, err = us.UpdateUser(user)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	resp.Resp.Data = tool2.SUCCESS
	resp.Resp.Status = true
	return resp, nil
}

func (s *server) UpdateUser(ctx context.Context, req *pb2.UpdateUserReq) (*pb2.UpdateUserResp, error) {
	resp := new(pb2.UpdateUserResp)
	resp.Resp = new(pb2.Resp)
	resp.User = new(pb2.User)
	us := service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.User.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERIDERROR
		return resp, nil
	}
	user, err := us.GetUserById(int(req.User.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}

	//更改用户名
	if req.User.Username != "" {
		//用户名验证
		ok, err = us.IsExistUsername(req.User.Username)
		if err != nil {
			log.Println("is exist user err:", err)
			resp.Resp.Status = false
			resp.Resp.Data = tool2.SERVERERROR
			return resp, err
		}
		//验证敏感词和sql注入
		if ok || tool2.CheckIfSensitive(req.User.Username) {
			resp.Resp.Status = false
			resp.Resp.Data = tool2.USERNAMEERROR
			return resp, nil
		}
		user.Username = req.User.Username
	}

	//更改成员组
	if int(req.User.GroupId) != user.GroupId {
		user.GroupId = int(req.User.GroupId)
	}

	//更改昵称
	if req.User.Name != "" {
		user.Name = req.User.Name
	}

	//更改邮箱
	if req.User.Email != "" {
		user.Email = req.User.Email
	}

	//更改电话
	if req.User.Phone != "" {
		user.Phone = req.User.Phone
	}

	//更改性别
	if int(req.User.Gender) != user.Gender {
		user.Gender = int(req.User.Gender)
	}

	//更改状态
	if int(req.User.State) != user.State {
		user.State = int(req.User.State)
	}

	res, err := us.UpdateUser(user)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	resp.User = &pb2.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId:  int32(res.GroupId),
	}
	resp.Resp.Status = true
	resp.Resp.Data = tool2.SUCCESS
	return resp, nil
}

func (s *server) GetUserById(ctx context.Context, req *pb2.GetUserByIdReq) (*pb2.GetUserByIdResp, error) {
	resp := new(pb2.GetUserByIdResp)
	resp.Resp = new(pb2.Resp)
	resp.User = new(pb2.User)
	us := service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERIDERROR
		return resp, nil
	}
	res, err := us.GetUserById(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}

	resp.Resp.Status = true
	resp.Resp.Data = tool2.SUCCESS
	resp.User = &pb2.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId:  int32(res.GroupId),
	}
	return resp, nil
}

func (s *server) GetUserByUserName(ctx context.Context, req *pb2.GetUserByUserNameReq) (*pb2.GetUserByUserNameResp, error) {
	resp := new(pb2.GetUserByUserNameResp)
	resp.Resp = new(pb2.Resp)
	resp.User = new(pb2.User)
	us := service.UserService{}

	//验证用户
	ok, err := us.IsExistUsername(req.UserName)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.USERIDERROR
		return resp, nil
	}
	res, err := us.GetUserByUsername(req.UserName)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool2.SERVERERROR
		return resp, err
	}

	resp.Resp.Status = true
	resp.Resp.Data = tool2.SUCCESS
	resp.User = &pb2.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId:  int32(res.GroupId),
	}
	return resp, nil
}
