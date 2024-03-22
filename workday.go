package main

import (
	"fmt"
	"os"
)

/*
https://gobyexample.com/command-line-arguments
*/

func main() {

	//argsWithProgramName := os.Args
	//argsWithoutProgramName := os.Args[1:]

	//arg_3 := os.Args[3]

	//fmt.Println(argsWithProgramName)
	//fmt.Println(argsWithoutProgramName)
	//fmt.Println(arg_3)

	// otherArgs := os.Args[2:]
	// 	fmt.Println("\n\n\n Your other aguments: ")
	// 	fmt.Println(otherArgs)

	// fmt.Println("\n\n\n")
	// fmt.Println("\n\n\n Your all aguments: ")
	// fmt.Println(os.Args)
	// fmt.Println("\n\n\n")

	if len(os.Args) > 1 {
		name := os.Args[1]
		if name == "--pause" {
			fmt.Println("Pausing Workday time counting")
			//do pause code
			//create a file workday-pause.txt
			//when this file will be detected by the main process
			// we stop counting time
		} else if name == "--resume" {
			fmt.Println("Resuming Workday time counting")
			//do resume code
			//remove a file workday-pause.txt
			// the main process should start counting time
		}
		fmt.Println("\n\n\n Exit !!!")

	} else {
		fmt.Println("Counting Workday ")
		//do code for counting work time
	}

}
