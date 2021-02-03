package conway

import (
	"fmt"
)

func ExampleSimplify() {
	fmt.Println(Describe(
		Simplify(Add(Int(13), Int(-12)))))
	// Output:  {{|}|}
}
