let greeting = "Hello";

const add = (a, b) => a + b
const result = add(param1, 4)

// 我自己的函数，可以与外界交互
let fn = async () => {
    myFunctionCall(result)
}

const returnObj = {greeting, result}