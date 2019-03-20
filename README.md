<a href="https://unlicense.org/">![License](https://img.shields.io/github/license/lu4p/ToRat_client.svg)</a>
## Build Client for Windows
You can build the client on Linux by executing the following
```
cd ~/go/src/github.com/lu4p/ToRat_client
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H windowsgui"
```
## Build Client for Linux
On Linux run 
```
cd ~/go/src/github.com/lu4p/ToRat_client
go build -ldflags "-s -w"
```

### [README](https://github.com/lu4p/ToRat/blob/master/README.md)
