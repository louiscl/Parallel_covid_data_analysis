package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

type StealingWorker struct {
	queues *[]DEQueue
	ctx    Context
	idx    int
}

func NewStealingWorker(assignedQueue int, ctx *Context, queues *[]DEQueue) *StealingWorker {

	newWorker := &StealingWorker{ctx: *ctx, queues: queues, idx: assignedQueue}

	return newWorker

}

func analyzeSingleCsv(i int, zipcode_input int, month_input int, year_input int, analysis_hash *Analysis_hash, flag *int32) {
	file_path := "./data/covid_" + strconv.Itoa(i+1) + ".csv"
	csv_file, _ := os.Open(file_path)

	r := csv.NewReader(csv_file)

	// iterate through the lines of the csv file
	for {
		line, err := r.Read()
		if err == io.EOF {
			// No more lines to read
			break
		} else if err != nil {
			log.Fatal(err)
		}

		zipcode, _ := strconv.Atoi(line[zipcodeCol])
		start_date := line[weekStart]
		start_month, _ := strconv.Atoi(start_date[:2])
		start_year, _ := strconv.Atoi(start_date[6:10])

		if zipcode != zipcode_input {
			continue
		} else if start_month != month_input {
			continue
		} else if start_year != year_input {
			continue
		} else if line[casesWeek] == "" || line[testsWeek] == "" || line[deathsWeek] == "" {
			continue
		}

		// critical section:
		for !atomic.CompareAndSwapInt32(flag, 0, 1) {
		}

		date, _ := strconv.Atoi(start_date[3:5])
		cases, _ := strconv.Atoi(line[casesWeek])
		tests, _ := strconv.Atoi(line[testsWeek])
		deaths, _ := strconv.Atoi(line[deathsWeek])

		if analysis_hash.analysis_map[date] != 1 {
			analysis_hash.analysis_map[date] = 1
			analysis_hash.total_cases += cases
			analysis_hash.total_tests += tests
			analysis_hash.total_deaths += deaths
		}
		atomic.StoreInt32(flag, 0)
	}
}

func createList(min, max int) []int {
	arr := make([]int, max-min+1)
	for i := range arr {
		arr[i] = min + i
	}
	return arr
}

func deleteElement(arr []int, i int) []int {
	arr[i] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}

//StealingWorker method
func (worker *StealingWorker) Run(wg *sync.WaitGroup, flag *int32, analysis_hash *Analysis_hash) {

	dq := (*worker.queues)[worker.idx]

	task := dq.PopBottom()

	for task != nil {
		fileIdx := task.ReturnFileNum()
		analyzeSingleCsv(fileIdx, worker.ctx.zipcode, worker.ctx.month, worker.ctx.year, analysis_hash, flag)
		task = dq.PopBottom()
	}

	runtime.Gosched()

	lenQ := len(*worker.queues)

	// Steal a task:
	indexSlice := createList(0, lenQ-1)
	indexSlice = deleteElement(indexSlice, worker.idx)

	for len(indexSlice) != 0 {
		randomIndex := rand.Intn(len(indexSlice))
		targetIndex := indexSlice[randomIndex]

		stolenTask := (*worker.queues)[targetIndex].PopTop()

		if stolenTask != nil {
			analyzeSingleCsv(stolenTask.ReturnFileNum(), worker.ctx.zipcode, worker.ctx.month, worker.ctx.year, analysis_hash, flag)
		}
		for stolenTask != nil {
			stolenTask = (*worker.queues)[targetIndex].PopTop()
			if stolenTask != nil {
				analyzeSingleCsv(stolenTask.ReturnFileNum(), worker.ctx.zipcode, worker.ctx.month, worker.ctx.year, analysis_hash, flag)
			}
		}

		indexSlice = deleteElement(indexSlice, randomIndex)

	}
	worker.Exit(wg)
}

func (worker *StealingWorker) Exit(wg *sync.WaitGroup) {
	wg.Done()
}
