package go_house

import (
	"fmt"
	"github.com/gohouse/t"
)

func GoHouse() {

	a1 := 123
	a2 := "321.34"

	fmt.Println(t.New(a1).String())
	fmt.Println(t.New(a2).Float64())
}
