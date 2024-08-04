package aggregate

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Jwt() {
	// 创建一个新的令牌对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  "John Doe",
		"admin": true,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // 1小时过期
	})

	// 使用秘密字符串签名令牌。该字符串需要确保保密，不要嵌入代码中。
	secretKey := []byte("6666")

	// 生成JWT字符串
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return
	}

	// 输出生成的JWT
	fmt.Println("Generated Token:", tokenString)

	// 解析和验证JWT
	token1, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证令牌使用的签名方法是否符合我们预期的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}

	if claims, ok := token1.Claims.(jwt.MapClaims); ok && token1.Valid {
		// JWT 验证成功，打印声明
		fmt.Println("Claims:", claims)
	} else {
		fmt.Println("Invalid token")
	}

}
