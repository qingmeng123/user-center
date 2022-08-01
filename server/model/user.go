/*******
* @Author:qingmeng
* @Description:
* @File:user
* @Date2022/7/14
 */

package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string ` gorm:"column:username;size:20;NOT NULL;unique;primaryKey" form:"username"`
	Password string `gorm:"column:password;NOT NULL" form:"password"`
	Gender   int    `gorm:"column:gender;default:0;type:tinyint;NOT NULL" form:"gender"` //0为男，1为女
	Name     string `gorm:"column:name;size:30;default:'null'" form:"name"`              //昵称
	Phone    string `gorm:"column:phone;size:20" form:"phone"`
	Email    string `gorm:"column:email;size:30" form:"email"`
	State    int    `gorm:"column:state;default:0;type:tinyint;NOT NULL" form:"state"`       //(0为有效用户，1为无效)
	GroupId  int    `gorm:"column:group_id;default:0;type:tinyint;NOT NULL" form:"group_id"` //成员组id,1为超级管理员，0为普通用户
}
