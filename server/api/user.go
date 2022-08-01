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
	"user-center/server/model"
	"user-center/server/pbfile/pb"
	"user-center/server/service"
	"user-center/server/tool"
)

func UserServer(tcpPort string) {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}
	creds := getCreds()
	//添加拦截器
	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(unaryInterceptor))
	pb.RegisterUserServiceServer(s, &server{})

	//开始处理
	if err = s.Serve(lis); err != nil {
		log.Fatalln("failed to server:", err)
	}

}


type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginReq) (resp *pb.LoginResp,err error) {
	//必须先手动为嵌套的结构体指针分配指向，即初始化一个内存，返回的是指针类型时也得先分配，他只是声明了结构体的指针，不然爆空指针异常
	resp=new(pb.LoginResp)
	resp.Resp=new(pb.Resp)

	us := service.UserService{}
	ok, err := us.IsExistUsername(req.Username)
	if !ok {
		resp.Resp.Status = false
		resp.Resp.Data = tool.USERNAMEERROR
		return resp, nil
	}
	user, err := us.GetUserByUsername(req.Username)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return
	}

	//验证密码
	if !tool.CheckPassword(user.Password, req.Password) {
		resp.Resp.Status = false
		resp.Resp.Data = tool.PASSWORERROR
		return resp, nil
	}

	resp.Resp.Status = true
	resp.Resp.Data = tool.SUCCESS
	return
}

func (s *server) Register(ctx context.Context, req *pb.RegisterReq) (resp *pb.RegisterResp,err error) {
	resp=new(pb.RegisterResp)
	resp.Resp=new(pb.Resp)

	us := service.UserService{}

	//用户名验证
	ok, err := us.IsExistUsername(req.User.Username)
	if err != nil {
		log.Println("is exist user err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return
	}
	//验证敏感词和sql注入
	if ok || tool.CheckIfSensitive(req.User.Username) {
		resp.Resp.Status = false
		resp.Resp.Data = tool.USERNAMEERROR
		return resp, nil
	}

	//密码处理
	level := tool.CheckPasswordLever(req.User.Password)
	if level < 1 {
		resp.Resp.Status = false
		resp.Resp.Data = tool.PASSWORERROR
		return resp, nil
	}
	pwd, err := tool.AddSalt(req.User.Password)
	if err != nil {
		log.Println("add salt err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
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
		GroupId: int(req.User.GroupId),
	}

	err = us.CreateUser(&user)
	if err != nil {
		log.Println("create user err:", err)
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return
	}
	resp.Resp.Status = true
	resp.Resp.Data = tool.SUCCESS
	return resp, nil
}


func (s * server) ChangePwd(ctx context.Context, req *pb.ChangePwdReq) (*pb.ChangePwdResp, error) {
	resp:=new(pb.ChangePwdResp)
	resp.Resp=new(pb.Resp)
	us:=service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.ID))
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	if !ok{
		resp.Resp.Status=false
		resp.Resp.Data=tool.USERIDERROR
		return resp,nil
	}
	user, err := us.GetUserById(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return	resp,err
	}

	//验证密码
	if !tool.CheckPassword(user.Password, req.OldPassword) {
		resp.Resp.Status = false
		resp.Resp.Data = tool.PASSWORERROR
		return resp, nil
	}

	//更改密码
	pwd, err := tool.AddSalt(req.NewPassword)
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	user.Password=pwd
	_, err = us.UpdateUser(user)
	if err != nil {
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp, err
	}
	resp.Resp.Data=tool.SUCCESS
	resp.Resp.Status=true
	return resp,nil
}

func (s * server) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserResp, error) {
	resp:=new(pb.UpdateUserResp)
	resp.Resp=new(pb.Resp)
	resp.User=new(pb.User)
	us:=service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.User.ID))
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	if !ok{
		resp.Resp.Status=false
		resp.Resp.Data=tool.USERIDERROR
		return resp,nil
	}
	user, err := us.GetUserById(int(req.User.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return	resp,err
	}

	//更改用户名
	if req.User.Username!=""{
		//用户名验证
		ok, err = us.IsExistUsername(req.User.Username)
		if err != nil {
			log.Println("is exist user err:", err)
			resp.Resp.Status = false
			resp.Resp.Data = tool.SERVERERROR
			return	resp,err
		}
		//验证敏感词和sql注入
		if ok || tool.CheckIfSensitive(req.User.Username) {
			resp.Resp.Status = false
			resp.Resp.Data = tool.USERNAMEERROR
			return resp, nil
		}
		user.Username=req.User.Username
	}

	//更改成员组
	if int(req.User.GroupId)!=user.GroupId{
		user.GroupId=int(req.User.GroupId)
	}

	//更改昵称
	if req.User.Name!=""{
		user.Name=req.User.Name
	}

	//更改邮箱
	if req.User.Email!=""{
		user.Email=req.User.Email
	}

	//更改电话
	if req.User.Phone!=""{
		user.Phone=req.User.Phone
	}

	//更改性别
	if int(req.User.Gender)!=user.Gender{
		user.Gender=int(req.User.Gender)
	}

	//更改状态
	if int(req.User.State)!=user.State{
		user.State=int(req.User.State)
	}

	res, err := us.UpdateUser(user)
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	resp.User=&pb.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId: int32(res.GroupId),
	}
	resp.Resp.Status=true
	resp.Resp.Data=tool.SUCCESS
	return resp,nil
}

func (s * server) GetUserById(ctx context.Context, req *pb.GetUserByIdReq) (*pb.GetUserByIdResp, error) {
	resp:=new(pb.GetUserByIdResp)
	resp.Resp=new(pb.Resp)
	resp.User=new(pb.User)
	us:=service.UserService{}

	//验证用户
	ok, err := us.IsExistID(int(req.ID))
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	if !ok{
		resp.Resp.Status=false
		resp.Resp.Data=tool.USERIDERROR
		return resp,nil
	}
	res, err := us.GetUserById(int(req.ID))
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return	resp,err
	}

	resp.Resp.Status=true
	resp.Resp.Data=tool.SUCCESS
	resp.User=&pb.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId: int32(res.GroupId),
	}
	return resp,nil
}


func (s * server) GetUserByUserName(ctx context.Context, req *pb.GetUserByUserNameReq) (*pb.GetUserByUserNameResp, error) {
	resp:=new(pb.GetUserByUserNameResp)
	resp.Resp=new(pb.Resp)
	resp.User=new(pb.User)
	us:=service.UserService{}

	//验证用户
	ok, err := us.IsExistUsername(req.UserName)
	if err!=nil{
		resp.Resp.Status=false
		resp.Resp.Data=tool.SERVERERROR
		return resp,err
	}
	if !ok{
		resp.Resp.Status=false
		resp.Resp.Data=tool.USERIDERROR
		return resp,nil
	}
	res, err := us.GetUserByUsername(req.UserName)
	if err != nil {
		resp.Resp.Status = false
		resp.Resp.Data = tool.SERVERERROR
		return	resp,err
	}

	resp.Resp.Status=true
	resp.Resp.Data=tool.SUCCESS
	resp.User=&pb.User{
		ID:       int32(res.ID),
		Username: res.Username,
		Gender:   int32(res.Gender),
		Name:     res.Name,
		Phone:    res.Phone,
		Email:    res.Email,
		State:    int32(res.State),
		GroupId: int32(res.GroupId),
	}
	return resp,nil
}

