# workday-go

Simple go application for monitoring worktime. 
The application is counting the time from the moment of starting the application and saving the elapsed time to dmp file. 
Saving of dmp file currently every minute. 

## Development

**The application is still under development. Some features may not work yet!!!**

## Getting started

Download the executable file for your operating OS. 
Create a shell/batch script for starting it and add it to te autostart of your OS. 

See Application usage for details:

```
workday-<VERSION>-windows-amd64.exe 

# or

workday-<VERSION>-linux-amd64.bin 

# Example: 

./workday-0.3-20240325-linux-amd64.bin


Workday - No params provided


Workday-go - Application monitors work time. 
 Every few minutes it save the current progress of work time counter.
Usage:
    workday.exe --daemon
      to start counting of work time 
    workday.exe --pause
      to pause counting of work time, for example: 
      when you make a break in work to go out for a walk
    workday.exe --resume
      to resume counting of work time 
    workday.exe --status
      to print current day status of work time 
    workday.exe --report
      to print all statuses of work time 

All aguments you provided: 
[./workday-0.3-20240325-linux-amd64.bin]


```

### Windows workday-START.bat

```
@echo  off

workday-<VERSION>-windows-amd64.exe --daemon

```


