package hypnus

import (
	"fmt"
	"testing"
)

func TestJoinPaths(t *testing.T) {
	path := JoinPaths("/", "/b/")
	fmt.Println(path)
}
