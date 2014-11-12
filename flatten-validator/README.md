This is a simple validator for checking that flattened and unflattened URLs perform in the expected ways.

Currently it repeats queries and compares timing and any errors incurred.

Build with

```
go build main.go
```

And run with

```
./flatten-validator
```

TODOs:
* Make environment it runs against configurable (currently staging)
* Parallelise GETs (currently the whole thing takes > 10 minutes to run)
* Improve the data validation and reporting
