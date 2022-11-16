/*******
* @Author:qingmeng
* @Description:
* @File:user
* @Date2022/7/14
 */

package dao

import (
	"user-center/model"
)

type UserDao struct {
}

func (d *UserDao) CreateUser(user *model.User) error {
	tx := db.Create(&user)
	return tx.Error
}

func (d *UserDao) GetUserById(id int) (user model.User, err error) {
	tx := db.First(&user, id)
	return user, tx.Error
}

func (d *UserDao) GetUsers() (userList []*model.User, err error) {
	tx := db.Find(&userList)
	return userList, tx.Error
}

func (d *UserDao) DeleteUserById(id int) error {
	return db.Delete(&model.User{}, id).Error
}

func (d *UserDao) SaveUser(user model.User) error {
	return db.Save(&user).Error
}

func (d *UserDao) DisableUserById(id int) (user model.User, err error) {
	user.ID = uint(id)
	tx := db.Model(&user).Update("state", "1")
	return user, tx.Error
}

func (d *UserDao) EnableUserById(id int) (user model.User, err error) {
	user.ID = uint(id)
	tx := db.Model(&user).Update("state", "0")
	return user, tx.Error
}

func (d *UserDao) GetUserByUsername(username string) (user model.User, err error) {
	tx := db.Where("username=?", username).First(&user)
	return user, tx.Error
}
