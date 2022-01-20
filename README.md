# Covid Data Analysis

## Applying work-stealing scheduling using a DEQueue

This command line application analysis a given data set of covid cases in the Chicago area and returns the number of cases, tests and deaths for a provided zipcode, month and year.
It does so by applying parallel computation, and in specific a work-stealing scheduling between spawned threads.

## How to use the program

### 1) Retrieve the data

Because of the size of the data set, download a zip folder here:

https://uchicago.box.com/s/6quo5pf75riwv6va6356g3yrgolmw4az

and place only the covid\_\*.csv files inside the "program/covid/data" directory.

### 2) Run the program

The command line application takes four arguments:<br/>
threads = the number of threads (i.e., goroutines to spawn)<br/>
zipcode = a possible Chicago zipcode<br/>
month = the month to display for that zipcode<br/>
year = the year to display for that zipcode<br/>

Within the "covid" directory, run
go run bounded.go worker.go covid.go threads zipcode month year
e.g. go run bounded.go worker.go covid.go 2 60603 5 2020
