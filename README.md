
<h1 align="center">Zeniire</h1>

> A high performance server application for keeping track of a coin wallet balance. REST and gRPC APIs included. Written in Go utilizing Protocol Buffers.
> The name, **Zeniire**, comes from the old Japanese word (銭入れ) for "purse".
> Also this is a practice project.

### **[📰 CHANGELOG](docs/CHANGELOG.md)** | **[❤ CONTRIBUTING](docs/CONTRIBUTING.md)**

## 📌 Features
- Create new records of wallet balance either through gRPC or REST API.
- Read the records of the wallet by specifing a datetime range.
- Automatic deployment with docker-compose
- Database migrations
- TLS for gRPC

## 📝 TODO

- Proper user authentication system (e.g. with JWT)
- More testing
- Automatic benchmarks

## 💨 Benchmarks

### Instructions

- Install go-wrk `go install github.com/tsliwowicz/go-wrk@latest`
- Fill the database with data (e.g. with using [test_data.sql](docs/test_data.sql))
- Run the service
- Run the commands in the results section

### Results

- AMD Ryzen 9 6900HS
- 32GB 4800MHz
- Western Digital SN735 NVMe PCIe Gen3 x4
- 1 000 000 records in the PSQL DB
- 100 concurrent connections (-d 100)

**Command** `go-wrk -c 100 -d 10 http://localhost:3333/?limit=10`

```
Running 10s test @ http://localhost:3333/?limit=10
  100 goroutine(s) running concurrently
157462 requests in 9.948994505s, 31.23MB read
Requests/sec:           15826.93
Transfer/sec:           3.14MB
Avg Req Time:           6.318346ms
Fastest Request:        0s
Slowest Request:        258.9788ms
Number of Errors:       0
```

**Command** `go-wrk -c 100 -d 10 http://localhost:3333/?limit=10000`

```
Running 10s test @ http://localhost:3333/?limit=10000
  100 goroutine(s) running concurrently
156117 requests in 10.089704292s, 30.97MB read
Requests/sec:           15472.90
Transfer/sec:           3.07MB
Avg Req Time:           6.462911ms
Fastest Request:        0s
Slowest Request:        299.4244ms
Number of Errors:       0
```
