# worker
In this worker job, work is performed through using buffered channel, waitgroup, context and semaphore.

1. To build the application , run the following command:

    go build *.go

2. Hence to run the worker pool, run the command as:

Now, the code snippet can be followed to have a better glimse:

WorkData - to hold the work request

Dispather - this struct holds the semaphore to restrict the number of goroutines, jobBuffer channel to insert or release the job request, worker interface to perfom the work on the work function and waitgroup to simulate the work of go routines.

Printer - method is used to finally print the work request.


Send the work data to be executed by the worker pool:

![Screenshot from 2022-04-15 00-22-09](https://user-images.githubusercontent.com/48493235/163457814-04badea6-f0ca-480b-84c6-2ff70ab77cbb.png)

Finally you will se the work done on the terminal:

![Screenshot from 2022-04-15 00-22-43](https://user-images.githubusercontent.com/48493235/163457639-04e3c9ef-3d89-4c04-8ea2-8f45152aa05f.png)


