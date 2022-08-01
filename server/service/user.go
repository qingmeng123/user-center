/*******
* @Author:qingmeng
* @Description:
* @File:user
* @Date2022/7/14
 */

package service

import (
	"gorm.io/gorm"
	"user-center/server/dao"
	"user-center/server/model"
)

type UserService struct {
}

//创建user
func (s *UserService) CreateUser(user *model.User) error {
	us := dao.UserDao{}
	return us.CreateUser(user)
}

//通过id获取user
func (s *UserService) GetUserById(id int) (model.User, error) {
	us := dao.UserDao{}
	return us.GetUserById(id)
}

//通过username获取user
func (s *UserService) GetUserByUsername(username string) (model.User, error) {
	us := dao.UserDao{}
	return us.GetUserByUsername(username)
}

//获取所有user
func (s *UserService) GetUsers() ([]*model.User, error) {
	us := dao.UserDao{}
	return us.GetUsers()
}

//删除user
func (s *UserService) DeleteUserById(id int) error {
	us := dao.UserDao{}
	//先寻找是否存在该用户
	_, err := us.GetUserById(id)
	if err != nil {
		return err
	}
	return us.DeleteUserById(id)
}

//更新user
func (s *UserService) UpdateUser(user model.User) (model.User, error) {
	us := dao.UserDao{}
	ok, err := s.IsExistID(int(user.ID))
	if !ok {
		return model.User{}, err
	}
	err = us.SaveUser(user)
	return user, err
}

//判断Id是否存在
func (s *UserService) IsExistID(id int) (bool, error) {
	d := dao.UserDao{}
	_, err := d.GetUserById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsExistUsername 判断用户名是否存在
func (s *UserService) IsExistUsername(username string) (bool, error) {
	d := dao.UserDao{}
	_, err := d.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//使用户无效，1为有效，2无效
func (s *UserService) DisableUserById(id int) (model.User, error) {
	us := dao.UserDao{}
	_, err := us.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}
	return us.DisableUserById(id)
}

func (s *UserService) EnableUserById(id int) (model.User, error) {
	us := dao.UserDao{}
	_, err := us.GetUserById(id)
	if err != nil {
		return model.User{}, err
	}
	return us.EnableUserById(id)
}
