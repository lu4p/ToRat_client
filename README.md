## Build Client for Windows
You can build the client on Linux by executing the following in the ```ToRAT_client``` directory.
```
cd ~/go/src/github.com/lu4p/ToRat_client
env GOOS=windows GOARCH=amd64 go build ldflags "-s -w -H windowsgui"
```
