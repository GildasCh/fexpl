package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "scan":
		scan()
	default:
		usage()
	}

}

func usage() {
	fmt.Println("Bad command")
}

func scan() {
	if len(os.Args) < 5 {
		usage()
		return
	}

	fc := Explore(os.Args[2], os.Args[3])

	fout, _ := os.Create(os.Args[4])

	enc := gob.NewEncoder(fout)
	enc.Encode(fc)

	// jfc, _ := json.Marshal(fc)

	// fmt.Println(string(jfc))
}
