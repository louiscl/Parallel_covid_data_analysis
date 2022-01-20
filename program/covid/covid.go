// See README.md
package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Analysis_hash struct {
	analysis_map map[int]int
	total_cases  int
	total_tests  int
	total_deaths int
}

// provided variables
const zipcodeCol = 0
const weekStart = 2
const casesWeek = 4
const testsWeek = 8
const deathsWeek = 14

type Context struct {
	wg      sync.WaitGroup
	zipcode int
	month   int
	year    int
}

func main() {
	const usage = "Usage: covid threads zipcode month year\n" +
		"    threads = the number of threads (i.e., goroutines to spawn)\n" +
		"    zipcode = a possible Chicago zipcode\n" +
		"    month = the month to display for that zipcode \n" +
		"    year  = the year to display for that zipcode \n"

	if len(os.Args) != 5 {
		fmt.Println(usage)
		os.Exit(0)
	}

	// Command line arguments
	cmdLineArgs := os.Args[1:]
	var threads, _ = strconv.Atoi(cmdLineArgs[0])
	var zipcode, _ = strconv.Atoi(cmdLineArgs[1])
	var month, _ = strconv.Atoi(cmdLineArgs[2])
	var year, _ = strconv.Atoi(cmdLineArgs[3])

	var wg sync.WaitGroup
	var atomicFlag int32

	context := Context{wg: wg, zipcode: zipcode, month: month, year: year}

	// Static initial allocation, with  dynamic stealing:
	var interval = 500

	// run sequentially
	if threads == 0 {
		threads = 1
	}
	var remainder = interval % threads
	var ops_per_thread = (interval - remainder) / threads

	var start_point = 1

	// critical variable
	analysis_hash := Analysis_hash{}
	analysis_hash.analysis_map = make(map[int]int)

	wg.Add(threads)

	var queues []DEQueue
	var workers []StealingWorker

	// Generate & fill DEqueues, generate StealingWorkers, assign DEqueues
	for i := 1; i <= threads; i++ {
		if i == threads {
			ops_per_thread = ops_per_thread + remainder
		}
		indQueueInt := i - 1
		individualDequeue := NewBoundedDEQueue()

		endPoint := start_point + ops_per_thread - 1
		// fill runnables into queue:
		for i := start_point; i < endPoint; i++ {
			newRun := NewRunnable(i)
			individualDequeue.PushBottom(newRun)
		}

		queues = append(queues, individualDequeue)

		individualWorker := NewStealingWorker(indQueueInt, &context, &queues)
		workers = append(workers, *individualWorker)

		start_point += ops_per_thread
	}

	for i := 0; i < threads; i++ {
		go workers[i].Run(&wg, &atomicFlag, &analysis_hash)
	}

	wg.Wait()

	fmt.Println(strconv.Itoa(analysis_hash.total_cases) + "," + strconv.Itoa(analysis_hash.total_tests) + "," + strconv.Itoa(analysis_hash.total_deaths))
}
