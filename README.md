Worker Pool

In this worker job, work is performed through using buffered channel, waitgroup, context and semaphore.

To build the application , run the following command:

go build *.go

Hence to run the worker pool, run the command as:

go run .

![Screenshot from 2022-04-30 09-53-38](https://user-images.githubusercontent.com/48493235/166091080-99576f8b-0022-447b-9057-e67077dd89c5.png)

Further it will ask to input the number of workers which will be performing the task.

Finally the output will be:

![Screenshot from 2022-04-30 09-54-20](https://user-images.githubusercontent.com/48493235/166091136-0fe71a1e-540b-4a89-ad9a-911dc8d46c52.png)

Now, the code snippet can be followed to have a better glimse:

WorkData - to hold the work request

Dispather - this struct holds the semaphore to restrict the number of goroutines, jobBuffer channel to insert or release the job request, worker interface to perfom the work on the work function and waitgroup to simulate the work of go routines.

