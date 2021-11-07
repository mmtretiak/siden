## Architecture
HTTP webservice with two endpoints
1. POST `/file` for saving file on disk
2. GET `/file?fileName=?` for getting "processed" file

For saving large files used multipart/form as Content-Type, size of part stored in RAM controlled by `maxMemory` parameter.
Another part of file stored on disk, so usage of RAM is small.

For reading large files used reading into buffer, means that we are reading not more than N bytes in once.
Since buffer size controlled as parameter, it's easy to config it according to memory specs of server. Tested that most expensive
operation from time perspective is syscall.Read, which is expected.
If we look into total allocated bytes per operation from benchmark, it will be big.
But idea is to build application which will not use huge amount of memory at once.

## Benchmarks
Benchmarks results:

totalAlloc - how many memory was used from start of bench

maxAlloc - pick of memory usage

`12GB` file with buffer `~100MB`: timeSpent - `~103s` maxAlloc - `~300MB` totalAlloc - `~38GB`

`1GB` file with buffer `~100MB`: timeSpent - `~7.6s` maxAlloc - `~300MB` totalAlloc - `~2.9GB`

## Improvements

1. Rework `/file?fileName=?` endpoint, so it will start background job, result and status of which will be available on `/job` endpoint.
This will improve architecture of system, so user won't wait for long response in case of large files, instead
he will ask `job/` endpoint time to time about job result and status.
2. Try performance of worker pool for `processing file`
3. Add integration tests
4. Depends on exact system usage, we could store processed file with unprocessed, so it will speed up response of `/file?fileName=?` 
endpoint, but will require more disk space
5. Unhandled case when `processed file` is bigger than 8GB, this is good case, but need to think how to handle it
because in this case we can't guarantee that user will receive file without duplications, at least from first point of view
