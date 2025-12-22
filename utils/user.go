package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword 加密用户密码（生成bcrypt哈希值）
// password：用户输入的明文密码
// 返回：加密后的哈希字符串，错误信息
func EncryptPassword(password string) (string, error) {
	// 1. 生成密码哈希
	// bcrypt.GenerateFromPassword参数：
	// - []byte(password)：明文密码转字节切片
	// - bcrypt.DefaultCost：计算成本（默认10，范围4-31，值越大计算越慢，安全性越高）
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// 2. 哈希值转字符串返回（可直接存储到数据库）
	return string(hashBytes), nil
}

// VerifyPassword 验证密码是否正确
// password：用户输入的明文密码
// hash：数据库中存储的密码哈希值
// 返回：是否验证通过，错误信息
func VerifyPassword(password, hash string) (bool, error) {
	// 1. 对比明文密码与哈希值
	// bcrypt.CompareHashAndPassword会自动提取哈希值中的盐值和成本，进行一致性验证
	fmt.Println(password, hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// 2. 无错误表示验证通过
		return false, err
	}

	return true, nil
}

func MkJWT(userID uint) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "GAPI APP",
		Subject:   fmt.Sprintf("%d", userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})
	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte("AllYourBase"))
}
