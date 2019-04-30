package xtime

import (
	"fmt"
	"testing"
	"time"
)

func TestUnmarshalText(t *testing.T) {

	var d Duration
	err := d.UnmarshalText([]byte("1s"))
	fmt.Println(err)
	fmt.Println(time.Duration(d))

	dd := 3 * time.Minute
	fmt.Println(dd)

}
