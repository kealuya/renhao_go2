package protobuf

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
)

/*
	需要运行 protoc --go_out=. person.proto 来自动生成person.pb.go文件

	protoc --go_out=. person.proto
*/

func Protobuf() {
	// 创建一个Person实例
	p := &Person{
		Name:        "Alice",
		Id:          123,
		Email:       "alice@example.com",
		Age:         30,
		Address:     "123 Main St, Springfield",
		PhoneNumber: 1234567890,
		IsActive:    true,
		Height:      5.6,
		Weight:      130.5,
		Job:         "Software Engineer",
		Company:     "Tech Corp",
	}

	//json测试
	b, err := json.Marshal(p)
	// 将字节流保存到文件
	err = ioutil.WriteFile("protobuf/person.json", b, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	// 序列化为二进制格式
	data, err := proto.Marshal(p)
	if err != nil {
		log.Fatal("Marshaling error: ", err)
	}
	// 将字节流保存到文件
	err = ioutil.WriteFile("protobuf/person.bin", data, 0644)
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
	// 反序列化二进制格式
	newPerson := &Person{}
	err = proto.Unmarshal(data, newPerson)
	if err != nil {
		log.Fatal("Unmarshaling error: ", err)
	}

	// 输出新创建的Person实例
	fmt.Println(newPerson)
}
