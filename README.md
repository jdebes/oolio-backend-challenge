# Shopping Cart

Solution to the Backend Challenge.

## Setup

1. Copy 3 couponbase gzip files to `data/`
2. Run the application with the `data` argument for promo code validation first (explanation on how this works below). This takes about 2-3 minutes.

## Run

Pre-requisite Promo Code Run (do this before running server):
```bash
go run main.go data
```

Start the API Server:
```bash
go run main.go server
```

Run the tests:
```bash
go test ./...
```

## Notes about API implementation
- Used Hashmaps to simulate database to save time
- Did not use OpenAPI generator to generate boilerplate from spec to save time on dealing with templating and other quirks.
- Ideally constants such as api key secret etc should be passed as environment variables. I have skipped this to save time.
- Using API spec linked in the advanced challenge readme [here](https://orderfoodonline.deno.dev/public/openapi.yaml). There are inconsistencies between this one and the one in the repo.
  - Did not add `images` to the order response as this is not in the linked spec.
  - Did not add `couponCode` to the order response. This is existing in the demo server response but not linked spec.
  - `price` is represented as float in all specs. I configured `decimcal.Decimal` to do this as well. But this adds risk of clients handling this incorrectly. String representation might be better.

## Promo Code Validation Approach

Some observations and assumptions:
- Out of ~313 million promo codes there are less than 10 valid ones. Indicating the data is almost entirely noise.
- A promo code is valid if it appears in at least two files, suggesting that each file represents a distinct batch rather than an ongoing stream of codes.
- There is no indication that the system must handle real-time updates (i.e., new promo codes being continuously added).
- The promo codes are too large to store in memory (taking into account overhead).

### Solution

Based off the above assumptions, a way to solve this problem is to find all valid codes up front (since there are so few) and let the API search only valid ones. This can be done with a combination of external sort and k-way merges. 

#### How this works
Finding duplicates across the 3 files becomes easy then they are sorted. External sort allows us to sort very large files (larger what can fit into memory) efficiently.

__External Sort__:
- Read chunks (N lines) of the file into memory.
- Sort these chunks in Memory
- Write sorted chunks back to disk as smaller temp files.
- Do repeated K-way merges on temp files

__K-way Merge__:
- Read one line at a time from each sorted file.
- Use a priority queue to keep track of the smallest line (in terms of sort order) across all files.
- Write the smallest line to an output file
- When a line is removed from PQ, read next line from same file.

The implementation can be found under [`data/valid_promos.go`](data/valid_promos.go) and does the following:
1. Use `github.com/lanrat/extsort` to sort each file in 1 million line chunks
   - Sorting of each file is done concurrently
   - While sorting remove any duplicates within the same file
2. Perform a K-way merge to find duplicates across sorted files
   - When we encounter duplicates, write them to an output file
3. Produce `couponbase_validpromos` which the API can use for fast lookups.

#### Advantages and Alternatives

Advantages:
- Run time to find all valid promos is 2-3 minutes. Which is ideally suited for an assessment style challenge.
- Very little memory overhead
- Efficient lookup since we have eliminated noise
- Don't have to spend time on setting up and inserting into a database for the assessment

Alternatives:

This solution would not be suitable if promo codes are being created continuously, and we frequently have to recompute valid ones. 

In this case an alternative approach could be to insert all promo codes into a low overhead key-value DB such as BadgerDB, which offers high write throughput and fast lookups based on promo code.