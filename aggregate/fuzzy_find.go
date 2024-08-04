package aggregate

import (
	"fmt"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// FuzzyFind Fuzzy Search
func FuzzyFind() {
	words := []string{"中华人民共和国", "中国", "中心国", "中华人民"}
	fmt.Println(fuzzy.Match("中国", "中国"))
	fmt.Println(fuzzy.Find("中国", words))
	fmt.Println(fuzzy.RankFind("中国", words))
	fmt.Println(fuzzy.RankMatch("中国", "中华人民共和国"))
}
