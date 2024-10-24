//go:generate go run gen/gorm_gen/main.go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config .oapi-config.yaml spec/schema/@typespec/openapi3/openapi.yaml
//go:generate go run -mod=mod github.com/google/wire/cmd/wire ./...

package main

import (
	"crypto/rand"
	"fmt"
)

func main() {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	fmt.Println(result)
}
