/*******
* @Author:qingmeng
* @Description:
* @File:token
* @Date2022/7/16
 */

package model

import "time"
import "github.com/dgrijalva/jwt-go"

type MyClaims struct {
	UserId   int
	UserName string
	Password string
	Time     time.Time
	jwt.StandardClaims
}

//第三方登陆的token
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
