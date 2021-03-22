# go-speed
go-speed is a small wrapper for speedtest-cli to measure the speed of your internet connection every 30 minutes and show the result as a graph.

## Steps to get it working
1. Install speedtest-cli
1. Install go
1. Build binary by `go build main.go`
1. Run it by `./main` (or `.\main.exe` if you are on Windows machine)

This will also start a webserver distributing a plot that is available under `http://localhost:8070`.
So you can run this on your RasPi and access the results plot quite easily.
Detailed results are also stored in a SQLite-Database `data.db`.
