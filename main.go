package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type WorkData struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	WorkDays int    `json:"workdays"`
	Salary   int    `json:"salary"`
}

func (d *Dispatcher) Collector(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	decoder := json.NewDecoder(r.Body)
	var wd *WorkData
	err := decoder.Decode(&wd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(wd)

	d.Add(wd)
}

type Worker interface {
	Work(j *WorkData)
}

type Dispatcher struct {
	sem       chan struct{} // semaphore
	jobBuffer chan *WorkData
	worker    Worker
	wg        sync.WaitGroup
}

type WorkerId struct {
	Id int
}

func NewWorkerId(id int) *WorkerId {
	return &WorkerId{
		Id: id,
	}
}

// NewDispatcher will create a new instance of job dispatcher.
// maxWorkers means the maximum number of goroutines that can work concurrently.
// buffers means the maximum size of the queue.
func NewDispatcher(p *Printer, maxWorkers int, buffers int) *Dispatcher {
	return &Dispatcher{

		// Restrict the number of goroutine using buffered channel (as counting semaphor)
		sem:       make(chan struct{}, maxWorkers),
		jobBuffer: make(chan *WorkData, buffers),
		worker:    p,
	}
}

// Start starts a dispatcher.
// This dispatcher will stops when it receive a value from `ctx.Done`.
func (d *Dispatcher) Start(ctx context.Context) {
	// fmt.Printf("worker %d revieved the work request:", d.id)
	d.wg.Add(1)
	go d.loop(ctx)
}

// Wait blocks until the dispatcher stops.
func (d *Dispatcher) Wait() {
	d.wg.Wait()
}

// Add enqueues a job into the queue.
// If the number of enqueued jobs has already reached to the maximum size,
// this will block until the other job has finish and the queue has space to accept a new job.
func (d *Dispatcher) Add(job *WorkData) {
	d.jobBuffer <- job
}

func (d *Dispatcher) stop() {
	d.wg.Done()
}

func (d *Dispatcher) loop(ctx context.Context) {
	var wg sync.WaitGroup
Loop:
	for {
		select {
		case <-ctx.Done():
			// block until all the jobs finishes
			wg.Wait()
			break Loop
		case job := <-d.jobBuffer:
			// Increment the waitgroup
			wg.Add(1)
			// Incrementing a semaphore count
			d.sem <- struct{}{}
			fmt.Printf("\n\nPlease wait worker %d is perfoming the work...\n\n", len(d.sem))
			go func(job *WorkData) {
				defer wg.Done()
				// After the job finished, decremented a semaphore count
				defer func() { <-d.sem }()
				time.Sleep(time.Second * 2)
				d.worker.Work(job)
				fmt.Printf("\n\nworker %d done the work\n\n", len(d.sem))
			}(job)
		}
	}
	d.stop()
}

// Printer is a dummy worker that just prints received URL.
type Printer struct{}

func NewPrinter() *Printer {
	return &Printer{}
}

// Work waits for a few seconds and print a received URL.
func (p *Printer) Work(j *WorkData) {
	fmt.Printf("Name: %s\nCity: %s\nWorkdays: %d\nSalary: %d\n", j.Name, j.City, j.WorkDays, j.Salary)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		// wait until receiving the signal
		<-sigCh
		cancel()
	}()

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Printf("Please enter the maximum number of workers to perfom the work request: \n")
	// maxNumWorkers, _ := reader.ReadString('\n')
	// fmt.Println(maxNumWorkers)
	// maxNumWorkersCs, err := strconv.Atoi(maxNumWorkers)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// fmt.Println(maxNumWorkersCs)

	maxNumWorkers := 0
	fmt.Print("Please enter the maximum number of workers and if you want to use the default value for maximum worker then please press 0: ")
	fmt.Scanf("%d", &maxNumWorkers)
	if maxNumWorkers == 0{
		maxNumWorkers = 4
	}
	for i := 0; i < maxNumWorkers; i++ {
		fmt.Println("Starting worker", i+1)
		NewWorkerId(i + 1)
	}

	p := NewPrinter()
	d := NewDispatcher(p, maxNumWorkers, maxNumWorkers+1)
	d.Start(ctx)
	http.HandleFunc("/collector", d.Collector)
	//starting the web server to listen for the request.
	fmt.Println("server is starting...")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Println(err.Error())
	}
}
