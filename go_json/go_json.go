package go_json

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

func GoJson() {
	jsonString := `
{
  "name": {"first": "Tom", "last": "Anderson"},
  "age":37,
  "children": ["Sara","Alex","Jack"],
  "fav.movie": "Deer Hunter",
  "friends": [
    {"first": "Dale", "last": "Murphy", "age": 44, "nets": ["ig", "fb", "tw"]},
    {"first": "Roger", "last": "Craig", "age": 68, "nets": ["fb", "tw"]},
    {"first": "Jane", "last": "Murphy", "age": 47, "nets": ["ig", "tw"]}
  ]
}`

	gjson.AddModifier("case", func(json, arg string) string {

		fmt.Println("handle::", json)

		if arg == "upper" {
			return strings.ToUpper(json)
		}
		if arg == "lower" {
			return strings.ToLower(json)
		}
		return json
	})
	result := gjson.Parse(jsonString)

	m := result.Map()
	fmt.Printf("%+v\n", m)
	fmt.Println("========")
	fmt.Printf("%+v\n", gjson.Get(jsonString, "friends.#.first|@case"))

	fmt.Println("========")
	gjson.ForEachLine(jsonString, func(line gjson.Result) bool {
		fmt.Println(line.Value())
		return true
	})
	fmt.Println("========")
	customMethodDemo()
}

/*
 * ======================   实现自定义解析json ======================
 */
type Order struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	CreatedTime time.Time `json:"created_time"`
}

const layout = "2006-01-02 15:04:05"

// MarshalJSON 为 Order 类型实现自定义的 MarshalJSON 方法
func (o *Order) MarshalJSON() ([]byte, error) {
	type TempOrder Order // 定义与 Order 字段一致的新类型
	return json.Marshal(struct {
		CreatedTime string `json:"created_time"`
		*TempOrder         // 避免直接嵌套 Order 进入死循环
	}{
		CreatedTime: o.CreatedTime.Format(layout),
		TempOrder:   (*TempOrder)(o),
	})
}

// UnmarshalJSON 为 Order 类型实现自定义的 UnmarshalJSON 方法
func (o *Order) UnmarshalJSON(data []byte) error {
	type TempOrder Order // 定义与 Order 字段一致的新类型
	ot := struct {
		CreatedTime string `json:"created_time"`
		*TempOrder         // 避免直接嵌套 Order 进入死循环
	}{
		TempOrder: (*TempOrder)(o),
	}
	if err := json.Unmarshal(data, &ot); err != nil {
		return err
	}
	var err error
	o.CreatedTime, err = time.Parse(layout, ot.CreatedTime)
	if err != nil {
		return err
	}
	return nil
}

// 自定义序列化方法
func customMethodDemo() {
	o1 := Order{
		ID:          123456,
		Title:       "《七米的Go学习笔记》",
		CreatedTime: time.Now(),
	}
	// 通过自定义的 MarshalJSON 方法实现 struct -> json string
	b, err := json.Marshal(&o1)
	if err != nil {
		fmt.Printf("json.Marshal o1 failed, err:%v\n", err)
		return
	}
	fmt.Printf("str:%s\n", b)
	// 通过自定义的 UnmarshalJSON 方法实现 json string -> struct
	jsonStr := `{"created_time":"2020-04-05 10:18:20","id":123456,"title":"《七米的Go学习笔记》"}`
	var o2 Order
	if err := json.Unmarshal([]byte(jsonStr), &o2); err != nil {
		fmt.Printf("json.Unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Printf("o2:%#v\n", o2)
}
