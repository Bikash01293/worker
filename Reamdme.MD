1. To build the application , run the following command:

    go build *.go

2. Hence to run the worker pool, run the command as:

    go run .

Now, the code snippet can be followed to have a better glimse:

WorkData - to hold the work request

Dispather - this struct holds the semaphore to restrict the number of goroutines, jobBuffer channel to insert or release the job request, worker interface to perfom the work on the work function and waitgroup to simulate the work of go routines.

Printer - method is used to finally print the work request.


    
