# Go stack #
The goal of this small project is to benchmark the performance of a R/W lock based stack implementation with an optimistic concurrency control implementation.

## Lock description ##
The locking stack uses a classic R/W lock that allows multiple readers and has exclusive writers. additionally, incomming readers are block when there is a waiting writer.

## OCC description ##
The optimistic concurrency control stack uses atomic loads and compare and swap instructions in a retry loop to perform concurrent operations. This method relies on the golang garbage collector to not reuse memory while it is still referenced from another stack, thus avoiding the ABA problem.

## Testing methods ##
To save on time, the benchmark results in [results.txt](result.txt ) were run on my local laptop and did not include a hot start period to just capture steady state throughput. Additionally I did not run the benchmark multiple times for statistical significance.

## Results ##
That being said the test show that OCC out performs locking for this algorithm. 

In this chart plotting 5 second throughput as a function of the degrees of parallelsim using a 10% write workload, OCC outperforms but both implementations slow down as DOP increases.

![alt text](dop.JPG?raw=true)


In this second chart of 5 second throughput agaisnt the write workload percentage at 4 degrees of parallelism, both perform comparable at 0% write. Once the stack approach starts seeing write workload it's performance tanks compared to OCC. Finally the two converge when the write percent of workload goes to 100%
![alt text](write.JPG?raw=true)

## Reproducing results ##
To launch the full benchmark  the command, ```go run driver/driver.go -fullBenchmark```