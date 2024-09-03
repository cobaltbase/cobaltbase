package main

import (
	"github.com/cobaltbase/cobaltbase"
)

func main() {
	cb := cobaltbase.New()
	cb.Run(":3000")
}
