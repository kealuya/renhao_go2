package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Perform simple calculations",
	Long:  `Perform simple arithmetic calculations like addition.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("This command requires exactly two arguments.")
			return
		}

		num1, err1 := strconv.Atoi(args[0])
		if err1 != nil {
			fmt.Printf("Error converting %s to number: %v\n", args[0], err1)
			return
		}

		num2, err2 := strconv.Atoi(args[1])
		if err2 != nil {
			fmt.Printf("Error converting %s to number: %v\n", args[1], err2)
			return
		}

		result := num1 + num2
		fmt.Printf("The sum of %d and %d is %d\n", num1, num2, result)
	},
}
