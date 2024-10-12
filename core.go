package main

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Operation ====================================================
// Operation 定义一个操作接口
type Operation interface {
	Execute(a, b int) int
}

// OperationFunc 使用一个类型，将函数适配成符合接口的方法
type OperationFunc func(a, b int) int

// Execute 让函数类型实现接口
func (f OperationFunc) Execute(a, b int) int {
	return f(a, b)
}

// OperationFunc2 使用一个类型，将函数适配成符合接口的方法
type OperationFunc2 func(a, b int) int

type OperationStruct struct{}

func (receiver OperationStruct) Execute(a, b int) int {
	return a * b
}

// Calculate1 计算函数，接受一个 Operation 接口
func Calculate1(a, b int, op Operation) int {
	return op.Execute(a, b)
}
func Calculate2(a, b int, op OperationFunc2) int {
	return op(a, b)
}
func Calculate3(a, b int, op OperationStruct) int {
	return op.Execute(a, b)
}

func CoreRun() {

	/*
		OperationFunc作为函数，实现了Operation接口，如果某个函数内需要该函数时，可以通过Calculate，传入该函数
		且这个函数可以由用户实现。
		如果这个功能用struct实现，那么每一个不同业务，都需要创建一个struct。不方便。
		http的HandleFunc采用的就是这种【由func实现接口】的形式实现的。
	*/

	// 1的场合，可以提前定义各种实现（因为已经实现了Operation接口）
	addition := OperationFunc(func(a, b int) int {
		return a + b
	})
	subtraction := OperationFunc(func(a, b int) int {
		return a - b
	})
	fmt.Println("Subtraction:", Calculate1(10, 5, subtraction), "Addition:", Calculate1(10, 5, addition)) // Output: Subtraction: 5

	// 2的场合，因为没有实现具体接口，所以每次传入参数时实现，虽然可以实现各种业务，但是没有接口控制，无法做到适配器模式
	fmt.Println("Addition:", Calculate2(10, 5, func(a, b int) int {
		return a + b
	}))
	// 3的场合，每次都需要创建一个struct，实现接口
	fmt.Println("MyImplement:", Calculate3(10, 5, OperationStruct{})) //

}

// Operation ====================================================

func LlmChat() {
	cmd := exec.Command("./app", "如何改进这些代码")
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}

func HttpService() {

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()

		fmt.Println(request.RequestURI)
		fmt.Println(request.URL)
		fmt.Println(request.RemoteAddr)

		_, err := writer.Write([]byte("hello world"))
		if err != nil {
			log.Panicln(err)
		}
	})

	logs.Info("it's work~!")

	err := http.ListenAndServe(":9009", nil)
	if err != nil {
		log.Panicln(err)
	}

}
