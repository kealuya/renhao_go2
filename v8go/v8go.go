package v8go

import (
	"fmt"
	"io/ioutil"
	"net/http"
	v8 "rogchap.com/v8go"
)

func V8go() {
	iso := v8.NewIsolate() // creates a new JavaScript VM
	defer iso.Dispose()
	ctx := v8.NewContext(iso) // 绑定global全局环境，并生成上下文
	defer ctx.Close()
	// 为运行上下文ctx做准备
	// 定义一个普通函数======================================================
	myFunctionTemplate := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		//fmt.Printf("info.Args()::%v \n", info.Args())
		//fmt.Printf("info.Context()::%#v \n", info.Context())
		message := info.Args()[0].String()           // 获取 console.log 的第一个参数
		fmt.Println("myFunctionCall::", message[:2]) // 将日志信息打印到 Go 的标准输出中
		re, _ := v8.NewValue(iso, "返回结果11")
		return re
	})
	// 定义一个异步函数======================================================
	fetchFn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		args := info.Args()
		url := args[0].String()

		resolver, _ := v8.NewPromiseResolver(info.Context())

		go func() {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			val, _ := v8.NewValue(iso, string(body))
			resolver.Resolve(val)
		}()
		return resolver.GetPromise().Value
	})

	// 将自定义函数绑定到全局对象
	/*
		本身还支持ctx := v8.NewContext(iso,objectTemplate) 这样会把objectTemplate放到全局ctx里，可以使用
		myFunctionCall() 直接调用。下面这种写法可以把objectTemplate绑定到一个全局对象中。方便管理和调用。
		且，还可以替换掉原有js中对象，比如global.Set("console", instance),就可以替换掉console.log方法。
	*/
	global := ctx.Global()
	objectTemplate := v8.NewObjectTemplate(iso)
	_ = objectTemplate.Set("myFunctionCall", myFunctionTemplate)
	_ = objectTemplate.Set("fetch", fetchFn, v8.ReadOnly)
	instance, err := objectTemplate.NewInstance(ctx)
	if err != nil {
		panic(err)
	}

	err = global.Set("my", instance)
	if err != nil {
		panic(err)
	}
	//======================================================

	// ES6 示例脚本
	script := `
		// 使用 let 和 const
		let greeting = "Hello";

		const add = (a, b) => a + b;
		const result = add(param1, 4);
		const returnObj = {greeting,result};

		// 我自己的函数，可以与外界交互
		let callResult = my.myFunctionCall(result);

		// 异步函数
		let f;
		(async () => {
			f = await my.fetch('https://www.baidu.com')
			my.myFunctionCall(f);
    	})();
	`
	// 设定参数
	err1 := ctx.Global().Set("param1", int32(123))
	if err1 != nil {
		panic(err1)
	}

	// 执行脚本  origin字段用于错误调试，比如 es6_example.js:1: Uncaught Error: Something went wrong
	val, err2 := ctx.RunScript(script, "es6_example.js")
	if err2 != nil {
		panic(err2)
	}
	prom, _ := val.AsPromise()
	// wait for the promise to resolve
	for prom.State() == v8.Pending {
		continue
	}

	// 获取脚本内数据
	v2, _ := ctx.RunScript("returnObj", "version.js")
	v2Json, err := v2.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("  returnObj: %s \n", v2Json)
	// 获取脚本内数据
	v3, _ := ctx.RunScript("f", "version.js")
	v3Json, err := v3.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("  f: %s \n", v3Json[:10])
	// 获取脚本内数据
	v4, _ := ctx.RunScript("callResult", "version.js")
	v4Json, err := v4.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("  callResult: %s \n", v4Json)

	//vm := ctx.Isolate() // get the Isolate from the context
	//vm.TerminateExecution()

}
