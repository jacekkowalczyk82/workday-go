package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
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

var CONST_8H_SECONDS int64 = 60 * 60 * 8

var debugLog *log.Logger
var infoLog *log.Logger
var errorLog *log.Logger

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err == nil {

		// fmt.Println("File " + fileName + " exists.")
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		// file does *not* exist
		// fmt.Println("File " + fileName + " does NOT exist.")
		return false

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		errorLog.Println("Errors wile checking if " + fileName + " exists.")
		errorLog.Println(err)

		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		panic(err)

	}
}

func GetFilesInDir(dirPath string) []fs.DirEntry {
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirEntries {
		debugLog.Println(f.Name())
	}
	return dirEntries
}

func configureLogs(date_time string) (*log.Logger, *log.Logger, *log.Logger) {
	// configure logs
	logFile, err := openLogFile("./logs/", "./logs/workday-"+date_time+".log")
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	// log.Println("Log file created")

	debugLog = log.New(logFile, "[debug]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	// debugLog.Println("this is debug")

	infoLog = log.New(mw, "[info]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	// infoLog.Println("this is info")

	errorLog = log.New(mw, "[error]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	// errorLog.Println("this is error")

	return debugLog, infoLog, errorLog

}

func openLogFile(dirPath string, path string) (*os.File, error) {
	os.MkdirAll(dirPath, 0755)

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
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
	debugLog.Println("")
	debugLog.Println("Saving dump file "+filePath, currentTime.Format(CONST_TIME_FORMAT))

	dumpFile, err := os.Create(filePath)
	check(err)

	//It’s idiomatic to defer a Close immediately after opening a file.
	defer dumpFile.Close()

	bufferedWriter := bufio.NewWriter(dumpFile)
	writtenBytes, err := bufferedWriter.WriteString(fmt.Sprintf("%v", worktime))
	check(err)
	bufferedWriter.Flush()
	debugLog.Printf("wrote %d bytes\n", writtenBytes)
	debugLog.Println("")
}

func ParseInt64(numberString string) int64 {
	number, err := strconv.ParseInt(numberString, 10, 64)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("DEBUG::Parsed number, %v with type %s!\n", number, reflect.TypeOf(number))

	return number
}

func GetHumanReadableTime(secondsTime int64) string {
	fullMinutes := secondsTime / 60
	seconds := secondsTime % 60

	hours := fullMinutes / 60
	minutes := fullMinutes % 60

	return fmt.Sprintf("%v hours %v minutes %v seconds", hours, minutes, seconds)
}

func GetHumanReadableHours(secondsTime int64) string {
	fullMinutes := secondsTime / 60

	hours := fullMinutes / 60

	return fmt.Sprintf("%v hours", hours)
}

func ReadWorkTimeFromDumpFile(filePath string) int64 {
	//currentTime := time.Now()
	//fmt.Println("DEBUG::Reading dump file "+filePath, currentTime.Format(CONST_TIME_FORMAT))

	data, err := os.ReadFile(filePath)
	check(err)
	debugLog.Println("read ", string(data))
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
	var fromPauseToResume string = ""

	var totalWorkTimeSeconds int64 = 0
	var dumpFileWorkTimeSeconds int64 = 0
	// read it only once from dump file, do NOT modify this VARiABLE ever

	var startPausedTimeSeconds int64 = 0
	var totalPausedTimeSeconds int64 = 0
	var currentPausedTimeSeconds int64 = 0

	startTime := time.Now()
	startTimeUnix := startTime.Unix()
	lastDumpTimeUnix := startTimeUnix

	// configure logs
	debugLog, infoLog, errorLog := configureLogs(startTime.Format(CONST_DATE_FORMAT))

	dumpFilePath := CONST_WORKDAY_RECORDS_DIR_PATH + "/" + CONST_WORKDAY_RECORDS_FILE_PREFIX + startTime.Format(CONST_DATE_FORMAT) + ".dmp"

	if len(os.Args) > 1 {
		appCommandParam := os.Args[1]
		if appCommandParam == "--daemon" {

			infoLog.Println("Counting Workday -- daemon STARTED " + startTime.String())
			// fmt.Println("Counting Workday -- daemon STARTED " + startTime.GoString())

			//do code for counting work time

			// check for existing dump file and load it and set values for totalWorkTimeSeconds

			infoLog.Println("Workday time - COUNTING")

			if FileExists(dumpFilePath) {
				dumpFileWorkTimeSeconds = ReadWorkTimeFromDumpFile(dumpFilePath)
				debugLog.Println("Loaded elapsed time: " + GetHumanReadableTime(dumpFileWorkTimeSeconds))
				// read it only once from dump file, do NOT modify this VARiABLE ever
				infoLog.Println("")
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
					if fromPauseToResume == "" || fromPauseToResume == "pause->resume" {
						fromPauseToResume = "resume->pause"
						fmt.Println("") // for correct printing of status line
						infoLog.Println("Workday time counting - PAUSED\n")
						debugLog.Println("Workday time - PAUSED")
						// } else if fromPauseToResume == "pause->resume" {
						// 	//"pause->resume"
						// 	fromPauseToResume = "resume->pause"
						// 	fmt.Println("") // for correct printing of status line
						// 	infoLog.Println("Workday time counting - PAUSED\n")
						// 	debugLog.Println("Workday time - PAUSED")
					} else {
						// if "resume->pause" do nothing

					}

					pausedTimeModulo60 := currentPausedTimeSeconds % int64(CONST_DUMP_PERIOD_SECONDS)
					if startPausedTimeSeconds == 0 {
						if (pausedTimeModulo60 >= 0) && (pausedTimeModulo60 < 2) {

							//paused = true
							debugLog.Println("Workday time counting - PAUSE detected, marking pause time")

						}

						startPausedTimeSeconds = currentTimeUnixSeconds
						currentPausedTimeSeconds = 0
					} else {
						currentPausedTimeSeconds = currentTimeUnixSeconds - startPausedTimeSeconds

					}
					// status line
					fmt.Print("\rWorkday time counting - PAUSED for this pause: ", currentPausedTimeSeconds,
						" seconds ", GetHumanReadableTime(currentPausedTimeSeconds),
						" in total: ", totalPausedTimeSeconds, " seconds ", GetHumanReadableTime(totalPausedTimeSeconds))

				} else {
					if fromPauseToResume == "resume->pause" {
						fmt.Println("") // for correct printing of status line
						infoLog.Println("Workday time - Resumed from pause - COUNTING\n")

						fromPauseToResume = "pause->resume"
					}
					// else {
					// 	//empty

					// }

					// count total pause,
					totalPausedTimeSeconds = totalPausedTimeSeconds + currentPausedTimeSeconds
					//if counting then reset startPausedTimeSeconds
					startPausedTimeSeconds = 0
					currentPausedTimeSeconds = 0

					totalWorkTimeSeconds = dumpFileWorkTimeSeconds + currentTimeUnixSeconds - totalPausedTimeSeconds - startTimeUnix

					fmt.Print("\rWorkday time counting - total work time: ", totalWorkTimeSeconds,
						" seconds = "+GetHumanReadableTime(totalWorkTimeSeconds),
						", paused: ", totalPausedTimeSeconds, " seconds                                              ")

					if elapsedFromLastDumpTimeSeconds >= int64(CONST_DUMP_PERIOD_SECONDS) {
						debugLog.Println("")
						debugLog.Println("Workday time - COUNTING, 5 minutes passed")
						debugLog.Println("currentTime.date format : ", currentTime.Format(CONST_DATE_FORMAT))
						debugLog.Println("currentTime.Format : ", currentTime.Format(CONST_TIME_FORMAT))

						debugLog.Println("currentTime.RFC3339 : ", currentTime.Format(time.RFC3339))
						debugLog.Println("currentTime.UNIX Epoch Seconds: ", currentTimeUnixSeconds)

						//do dump file of total worktime for the given date
						//make dir workday_records
						os.MkdirAll(CONST_WORKDAY_RECORDS_DIR_PATH, 0755)

						SaveWorkTimetoDumpFile(totalWorkTimeSeconds, dumpFilePath)
						debugLog.Println("Saved elapsed time: " + GetHumanReadableTime(totalWorkTimeSeconds))
						debugLog.Println("")

						lastDumpTimeUnix = currentTimeUnixSeconds

					}

				}

				time.Sleep(time.Second)

			}

		} else if appCommandParam == "--pause" {
			infoLog.Println("")
			infoLog.Println("Pausing Workday time counting")
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
			debugLog.Printf("wrote %d bytes\n", writtenBytes)
			bufferedWriter.Flush()

			infoLog.Println("Pausing Workday time counting - DONE")
			infoLog.Println("")

		} else if appCommandParam == "--resume" {
			infoLog.Println("")
			infoLog.Println("Resuming Workday time counting")
			//do resume code
			//remove a file workday-pause.txt
			// the main process should start counting time

			err := os.Remove("workday-pause.txt") //remove the file
			if err != nil {
				errorLog.Println("Error: ", err) //print the error if file is not removed
				errorLog.Println("Resuming Workday time counting - FAILED")
			} else {
				infoLog.Println("Successfully deleted file: ", "workday-pause.txt") //print success if file is removed
				infoLog.Println("Resuming Workday time counting - DONE")

			}
			infoLog.Println("")

		} else if appCommandParam == "--status" {
			infoLog.Println("")
			infoLog.Println("Status of current Workday ")
			//print current work day hours, minutes

			if FileExists(dumpFilePath) {
				dumpFileWorkTimeSeconds = ReadWorkTimeFromDumpFile(dumpFilePath)

				// read it only once from dump file, do NOT modify this VARiABLE ever
				infoLog.Println("")
				infoLog.Println("Workday time ", startTime.Format(CONST_DATE_FORMAT), " - ", dumpFileWorkTimeSeconds,
					" seconds = ", GetHumanReadableTime(dumpFileWorkTimeSeconds))

			} else {
				infoLog.Println("No Workday records for Today: ", startTime.Format(CONST_DATE_FORMAT))
			}

			// TO BE implemented
			fmt.Println("")

		} else if appCommandParam == "--report" {
			infoLog.Println("")
			infoLog.Println("Report of ALL Workdays , use command line grep for filtering per month")
			//print current work day hours, minutes
			reportFiles := GetFilesInDir(CONST_WORKDAY_RECORDS_DIR_PATH)
			var workdayRecordsTotalSecondsPerDay map[string]int64
			workdayRecordsTotalSecondsPerDay = make(map[string]int64)

			var workdayRecordsTotalSecondsPerMonth map[string]int64
			workdayRecordsTotalSecondsPerMonth = make(map[string]int64)

			var numberOfWorkDaysPerMonth map[string]int64
			numberOfWorkDaysPerMonth = make(map[string]int64)

			for _, f := range reportFiles {
				debugLog.Println(f.Name())

				dumpFilePath := CONST_WORKDAY_RECORDS_DIR_PATH + "/" + f.Name()
				dumpFileWorkTimeSeconds = ReadWorkTimeFromDumpFile(dumpFilePath)

				debugLog.Println(GetHumanReadableTime(dumpFileWorkTimeSeconds))

				fileNameElements := strings.Split(f.Name(), "_")
				debugLog.Println(fileNameElements[1])
				fileNameElements2 := strings.Split(fileNameElements[1], ".")
				reportDate := fileNameElements2[0]
				debugLog.Println(reportDate)
				dateElements := strings.Split(reportDate, "-")
				yearMonthPart := dateElements[0] + "-" + dateElements[1]

				workdayRecordsTotalSecondsPerDay[reportDate] = dumpFileWorkTimeSeconds

				sumWorkTimeSecondsPerMonth, sumFound := workdayRecordsTotalSecondsPerMonth[yearMonthPart]
				if sumFound {
					//already added
					sumWorkTimeSecondsPerMonth = sumWorkTimeSecondsPerMonth + dumpFileWorkTimeSeconds

				} else {
					sumWorkTimeSecondsPerMonth = 0 + dumpFileWorkTimeSeconds
				}
				workdayRecordsTotalSecondsPerMonth[yearMonthPart] = sumWorkTimeSecondsPerMonth

				numberOfWorkDays, numberFound := numberOfWorkDaysPerMonth[yearMonthPart]
				if numberFound {
					//already added
					numberOfWorkDays = numberOfWorkDays + 1

				} else {
					// 	//we must add it
					numberOfWorkDays = 1
				}
				numberOfWorkDaysPerMonth[yearMonthPart] = numberOfWorkDays

			} //for report files

			//count avaerage, per Month
			for yearMonthKey, totalSeconds := range workdayRecordsTotalSecondsPerMonth {

				numberOfWorkDays, numberFound2 := numberOfWorkDaysPerMonth[yearMonthKey]
				debugLog.Println("yearMonthKey", yearMonthKey)
				debugLog.Println("totalSeconds", totalSeconds)
				debugLog.Println("numberOfWorkDays", numberOfWorkDays)
				if numberFound2 {
					averagePerDay := totalSeconds / numberOfWorkDays
					expectedWorkTime := numberOfWorkDays * CONST_8H_SECONDS
					infoLog.Println("Month:", yearMonthKey, ", Total work seconds:", totalSeconds,
						"=", GetHumanReadableTime(totalSeconds), "/ Expected Time", expectedWorkTime, "=",
						GetHumanReadableHours(expectedWorkTime),
						", Average per day:", averagePerDay, "=", GetHumanReadableTime(averagePerDay))
					for dateKey, dateSeconds := range workdayRecordsTotalSecondsPerDay {
						if strings.HasPrefix(dateKey, yearMonthKey) {
							infoLog.Println("    ", dateKey, "->", GetHumanReadableTime(dateSeconds))
						}
					}

				}

			}

			// TO BE implemented
			infoLog.Println("")
		} else {
			errorLog.Println("\n\n")
			errorLog.Println("Workday - INVALID params provided")
			ShowUsage()
		}
		infoLog.Println(appCommandParam, "Exit !!!")

	} else {
		errorLog.Println("\n\nWorkday - No params provided")
		ShowUsage()

	}

}
