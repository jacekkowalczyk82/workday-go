package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"time"
)

/*
https://gobyexample.com/command-line-arguments
*/

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {

		fmt.Println("File " + fileName + " exists.")
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		// file does *not* exist
		fmt.Println("File " + fileName + " does NOT exist.")
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		fmt.Println("Errors wile checking if " + fileName + " exists.")
		fmt.Println(err)

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		panic(err)

	}
}

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

			pauseFile, err := os.Create("workday-pause.txt")
			check(err)

			//It’s idiomatic to defer a Close immediately after opening a file.
			defer pauseFile.Close()

			bufferedWriter := bufio.NewWriter(pauseFile)
			writtenBytes, err := bufferedWriter.WriteString("PAUSED\n")
			check(err)
			fmt.Printf("wrote %d bytes\n", writtenBytes)
			bufferedWriter.Flush()

			fmt.Println("Pausing Workday time counting - DONE")

		} else if name == "--resume" {
			fmt.Println("Resuming Workday time counting")
			//do resume code
			//remove a file workday-pause.txt
			// the main process should start counting time

			err := os.Remove("workday-pause.txt") //remove the file
			if err != nil {
				fmt.Println("Error: ", err) //print the error if file is not removed
				fmt.Println("Resuming Workday time counting - FAILED")
			} else {
				fmt.Println("Successfully deleted file: ", "workday-pause.txt") //print success if file is removed
				fmt.Println("Resuming Workday time counting - DONE")
			}

		}
		fmt.Println("\n\n\n Exit !!!")

	} else {
		fmt.Println("Counting Workday ")
		//do code for counting work time
		for true {

			time.Sleep(time.Second)

			var paused bool = fileExists("workday-pause.txt")
			if paused {

				fmt.Println("Workday time counting - PAUSED")

			} else {
				fmt.Println("Workday time - COUNTING")

			}

		}

	}

}
