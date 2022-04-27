// switched to new-branch
package work

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Emailst func()

type WorkData struct {
	emailreq Emailst
}

func (d *Dispatcher) Collector(emfunc func()) {

	work := WorkData{
		emailreq: emfunc,
	}

	d.Add(&work)
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
	time.Sleep(3 * time.Second)

}

func (d *Dispatcher) stop() {
	d.wg.Done()
}

func (d *Dispatcher) loop(ctx context.Context) {
Loop:
	for {
		select {
		case <-ctx.Done():
			// block until all the jobs finishes
			d.wg.Wait()
			break Loop
		case job := <-d.jobBuffer:
			// Increment the waitgroup
			d.wg.Add(1)
			// Incrementing a semaphore count
			d.sem <- struct{}{}
			fmt.Printf("\n\nPlease wait worker %d is perfoming the work...\n", len(d.sem))
			go func(job *WorkData) {
				defer d.wg.Done()
				// After the job finished, decremented a semaphore count
				defer func() { <-d.sem }()

				d.worker.Work(job)
				fmt.Printf("worker %d done the work and ready to accept a new work.\n", len(d.sem))
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
	fmt.Println("Wait for 2 second to complete the work")
	time.Sleep(2 * time.Second)
	j.emailreq()

}

func Workerpool() {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		// wait until receiving the signal
		<-sigCh
		cancel()
	}()

	maxNumWorkers := 0
	fmt.Print("Please enter the maximum number of workers and if you want to use the default value for maximum worker then please press 0: ")
	fmt.Scanf("%d", &maxNumWorkers)
	if maxNumWorkers == 0 {
		maxNumWorkers = runtime.GOMAXPROCS(1)
	}
	for i := 0; i < maxNumWorkers; i++ {
		fmt.Println("Starting worker", i+1)
		NewWorkerId(i + 1)
	}

	p := NewPrinter()
	d := NewDispatcher(p, maxNumWorkers, maxNumWorkers+1)
	d.Start(ctx)
	for i := 0; i < 10; i++ {
		d.Collector(Email)
	}

}
