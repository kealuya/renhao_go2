package goutil

import (
	"fmt"
	"github.com/gookit/goutil/timex"
)

func GoUtil() {
	now := timex.Now()
	fmt.Println(now.String())

	fmt.Println(now.Format("2006-01-02 15:04:05.000000"))

	tx, err := timex.FromString("2022-04-20 19:40:34")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tx.String())
}
