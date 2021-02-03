package conway

import (
	"fmt"
)

func ExampleInt() {
	fmt.Println(Below(Int(-2), Int(2)))
	fmt.Println(Describe(Int(2)))
	fmt.Println(Gift(Int(3), Int(3)))
	fmt.Println(Gift(Int(2), Int(3)))
	fmt.Println(Gift(Int(3), Int(2)))
	// Output: true
	// {{{|}|}|}
	// false
	// true
	// false
}
