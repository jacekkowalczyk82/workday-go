package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

/*
https://gobyexample.com/command-line-arguments
*/

var CONST_TIME_FORMAT = "2006-01-02 15:04:05"
var CONST_DATE_FORMAT = "2006-01-02"
var CONST_WORKDAY_RECORDS_DIR_PATH = "workday_records"
var CONST_WORKDAY_RECORDS_FILE_PREFIX = "worktime_"

// var CONST_DUMP_PERIOD_SECONDS = 300
var CONST_DUMP_PERIOD_SECONDS = 60 // every minute

// var CONST_DUMP_PERIOD_SECONDS = 30 // for testing

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {

		fmt.Println("File " + fileName + " exists.")
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		// file does *not* exist
		// fmt.Println("File " + fileName + " does NOT exist.")
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		fmt.Println("Errors wile checking if " + fileName + " exists.")
		fmt.Println(err)

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		panic(err)

	}
}

func ShowUsage() {
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

func SaveWorkTimetoDumpFile(worktime int64, filePath string) {
	currentTime := time.Now()
	fmt.Println("Saving dump file "+filePath, currentTime.Format(CONST_TIME_FORMAT))

	dumpFile, err := os.Create(filePath)
	check(err)

	//It’s idiomatic to defer a Close immediately after opening a file.
	defer dumpFile.Close()

	bufferedWriter := bufio.NewWriter(dumpFile)
	writtenBytes, err := bufferedWriter.WriteString(fmt.Sprintf("%v", worktime))
	check(err)
	bufferedWriter.Flush()
	fmt.Printf("wrote %d bytes\n", writtenBytes)
}

func ParseInt64(numberString string) int64 {
	number, err := strconv.ParseInt(numberString, 10, 64)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed number, %v with type %s!\n", number, reflect.TypeOf(number))

	return number
}

func GetHumanReadableTime(secondsTime int64) string {
	fullMinutes := secondsTime / 60
	seconds := secondsTime % 60

	hours := fullMinutes / 60
	minutes := fullMinutes % 60

	return fmt.Sprintf("%v hours %v minutes %v seconds", hours, minutes, seconds)
}

func ReadWorkTimeFromDumpFile(filePath string) int64 {
	currentTime := time.Now()
	fmt.Println("Reading dump file "+filePath, currentTime.Format(CONST_TIME_FORMAT))

	data, err := os.ReadFile(filePath)
	check(err)
	fmt.Println(string(data))
	return ParseInt64(string(data))
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
	var paused bool = false
	var totalWorkTimeSeconds int64 = 0
	var dumpFileWorkTimeSeconds int64 = 0
	// read it only once from dump file, do NOT modify this VARiABLE ever

	var startPausedTimeSeconds int64 = 0
	var totalPausedTimeSeconds int64 = 0
	var currentPausedTimeSeconds int64 = 0

	if len(os.Args) > 1 {
		appCommandParam := os.Args[1]
		if appCommandParam == "--daemon" {

			startTime := time.Now()
			startTimeUnix := startTime.Unix()
			lastDumpTimeUnix := startTimeUnix
			fmt.Println("Counting Workday -- daemon STARTED " + startTime.String())
			// fmt.Println("Counting Workday -- daemon STARTED " + startTime.GoString())

			dumpFilePath := CONST_WORKDAY_RECORDS_DIR_PATH + "/" + CONST_WORKDAY_RECORDS_FILE_PREFIX + startTime.Format(CONST_DATE_FORMAT) + ".dmp"
			//do code for counting work time

			// check for existing dump file and load it and set values for totalWorkTimeSeconds

			fmt.Println("Workday time - COUNTING")

			if FileExists(dumpFilePath) {
				dumpFileWorkTimeSeconds = ReadWorkTimeFromDumpFile(dumpFilePath)
				fmt.Println("Loaded elapsed time: " + GetHumanReadableTime(dumpFileWorkTimeSeconds))
				// read it only once from dump file, do NOT modify this VARiABLE ever
				fmt.Println("")
			}

			for true {

				currentTime := time.Now()

				currentTimeUnixSeconds := currentTime.Unix()
				// fmt.Println("currentTime.UNIX Epoch Seconds: ", currentTimeUnixSeconds)

				// y1, m1, d1 := currentTime.Date()
				// h1, min1, s1 := currentTime.Clock()

				// fmt.Printf("%d-%2d-%d %d:%d:%d\n",
				// 	currentTime.Year(),
				// 	currentTime.Month(),
				// 	currentTime.Day(),
				// 	currentTime.Hour(),
				// 	currentTime.Hour(),
				// 	currentTime.Second())

				elapsedFromLastDumpTimeSeconds := currentTimeUnixSeconds - lastDumpTimeUnix
				paused = FileExists("workday-pause.txt")
				if paused {
					fmt.Println("")
					fmt.Println("Workday time counting - PAUSED")
					//paused = true
					if startPausedTimeSeconds == 0 {
						fmt.Println("DEBUG::Workday time counting - PAUSE detected, marking pause time")
						startPausedTimeSeconds = currentTimeUnixSeconds
						currentPausedTimeSeconds = 0
					} else {
						currentPausedTimeSeconds = currentTimeUnixSeconds - startPausedTimeSeconds

					}
					fmt.Println("")
					fmt.Println("Workday time counting - PAUSED for total: ", totalPausedTimeSeconds, " seconds ", GetHumanReadableTime(totalPausedTimeSeconds))
					fmt.Println("Workday time counting - PAUSED for this pause: ", currentPausedTimeSeconds, " seconds ", GetHumanReadableTime(currentPausedTimeSeconds))
					fmt.Println("")
				} else {

					// count total pause,
					totalPausedTimeSeconds = totalPausedTimeSeconds + currentPausedTimeSeconds
					//if counting then reset startPausedTimeSeconds
					startPausedTimeSeconds = 0
					currentPausedTimeSeconds = 0

					totalWorkTimeSeconds = dumpFileWorkTimeSeconds + currentTimeUnixSeconds - totalPausedTimeSeconds - startTimeUnix
					fmt.Print("\rWorkday time counting - total work time: ", totalWorkTimeSeconds, " seconds = "+GetHumanReadableTime(totalWorkTimeSeconds)+", paused: ", totalPausedTimeSeconds, " seconds       ")

					if elapsedFromLastDumpTimeSeconds > int64(CONST_DUMP_PERIOD_SECONDS) {
						fmt.Println("")
						fmt.Println("DEBUG::Workday time - COUNTING, 5 minutes passed")
						fmt.Println("DEBUG::currentTime.date format : ", currentTime.Format(CONST_DATE_FORMAT))
						fmt.Println("DEBUG::currentTime.Format : ", currentTime.Format(CONST_TIME_FORMAT))

						fmt.Println("DEBUG::currentTime.RFC3339 : ", currentTime.Format(time.RFC3339))
						fmt.Println("DEBUG::currentTime.UNIX Epoch Seconds: ", currentTimeUnixSeconds)

						//do dump file of total worktime for the given date
						//make dir workday_records
						os.MkdirAll(CONST_WORKDAY_RECORDS_DIR_PATH, 0755)

						SaveWorkTimetoDumpFile(totalWorkTimeSeconds, dumpFilePath)
						fmt.Println("Saved elapsed time: " + GetHumanReadableTime(totalWorkTimeSeconds))
						fmt.Println("")

						lastDumpTimeUnix = currentTimeUnixSeconds

					}

				}

				time.Sleep(time.Second)

			}

		} else if appCommandParam == "--pause" {
			fmt.Println("")
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
			fmt.Println("")

		} else if appCommandParam == "--resume" {
			fmt.Println("")
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
			fmt.Println("")

		} else if appCommandParam == "--status" {
			fmt.Println("")
			fmt.Println("Status of current Workday ")
			//print current work day hours, minutes

			// TO BE implemented
			fmt.Println("")

		} else if appCommandParam == "--report" {
			fmt.Println("")
			fmt.Println("Report of ALL Workdays , use command line grep for filtering per month")
			//print current work day hours, minutes

			// TO BE implemented
			fmt.Println("")
		} else {
			fmt.Println("\n\nWorkday - INVALID params provided")
			ShowUsage()
		}
		fmt.Println("\n\n\n Exit !!!")

	} else {
		fmt.Println("\n\nWorkday - No params provided")
		ShowUsage()

	}

}
