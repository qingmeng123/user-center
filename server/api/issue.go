/*******
* @Author:qingmeng
* @Description:证书
* @File:issue
* @Date2022/7/15
 */

package api

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

func getCreds() credentials.TransportCredentials {
	//双向证书
	//从证书相关文件中读取和解析信息，得到证书公钥，私钥对
	cert, err := tls.LoadX509KeyPair("server/pbfile/cert/server.pem", "server/pbfile/cert/server.key")
	if err != nil {
		log.Fatal("证书读取错误", err)
	}
	//创建一个新的，空的CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("server/pbfile/cert/ca.crt")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}

	//尝试解析所传入的PEM编码的证书，如果解析成功会将其加到CertPool中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	//构建基于TLS的TransportCredentials选项
	creds := credentials.NewTLS(&tls.Config{
		//设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		//要求必须校验客户端的证书。根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		//设置证书的集合，校验方式使用ClientAuth中设定的模式
		ClientCAs: certPool,
	})

	return creds
}
