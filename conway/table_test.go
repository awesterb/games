package conway

import (
	"fmt"
)

func ExampleSerial() {
	t := NewTable()
	for i := 0; i < 10; i++ {
		fmt.Println(t.Describe(Serial(i)))
	}
	// Output:
	// {|}
	// {#0|}
	// {#0|#0}
	// {|#0}
	// {#1|}
	// {#1|#1}
	// {|#1}
	// {#0|#1}
	// {#1|#0}
	// {|#2}
}
