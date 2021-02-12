package conway

import (
	"fmt"
)

func ExampleSimplify() {
	fmt.Println(Describe(
		Simplify(Add(Int(3), Int(-2)))))
	// Output:  {{|}|}
}
