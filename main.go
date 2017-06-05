package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fc := Explore("pictures", os.Args[1])

	jfc, _ := json.Marshal(fc)

	fmt.Println(string(jfc))
}
