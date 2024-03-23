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

func showUsage() {
	fmt.Println("\n")
	fmt.Println("Workday-go - Application monitors work time. \n " +
		"Every few minutes it save the current progress of work time counter.")

	fmt.Println("Usage:")
	fmt.Println("    workday.exe --daemon\n" +
		"      to start counting of work time ")
	fmt.Println("    workday.exe --pause\n" +
		"      to pause counting of work time, for example: \n" +
		"      when you make a break in work to go out for a walk")
	fmt.Println("    workday.exe --resume\n" +
		"      to resume counting of work time ")
	fmt.Println("    workday.exe --status\n" +
		"      to print current day status of work time ")
	fmt.Println("    workday.exe --report\n" +
		"      to print all statuses of work time ")

	fmt.Println("\nAll aguments you provided: ")
	fmt.Println(os.Args)

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
		appCommandParam := os.Args[1]
		if appCommandParam == "--daemon" {
			fmt.Println("Counting Workday -- daemon STARTED")
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

		} else if appCommandParam == "--pause" {
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

		} else if appCommandParam == "--resume" {
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

		} else if appCommandParam == "--status" {
			fmt.Println("Status of current Workday ")
			//print current work day hours, minutes
		} else if appCommandParam == "--report" {
			fmt.Println("Report of ALL Workdays , use command line grep for filtering per month")
			//print current work day hours, minutes
		} else {
			fmt.Println("\n\nWorkday - INVALID params provided")
			showUsage()
		}
		fmt.Println("\n\n\n Exit !!!")

	} else {
		fmt.Println("\n\nWorkday - No params provided")
		//do code show usage
		showUsage()

	}

}
