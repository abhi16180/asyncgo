# asyncgo

Asyncgo is zero-dependency asynchronous task executor written in pure go, that prioritises speed and ease of use.

###  Features
- Asynchronous Task Execution: Submit tasks to execute asynchronously and retrieve results.
- No Manual Goroutine Management: Abstracts away the complexity of managing goroutines, and simplifying the code.
- Worker Pool Management: Asyncgo carefully handles worker pool creation & task execution.
- Graceful Shutdown: Guarantees all existing tasks are completed before shutting down the workers
- Task Cancellation: Support for terminating workers.

### Usecases

- Asynchronous HTTP Requests for Microservices
- Background Job Execution
- Infinite concurrent pollling with worker pool (receiving messages from AWS SQS or similar services)


### Documentation