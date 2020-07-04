# On Windows, send a ctrl-break to a process

This sample repo is a standalone version of the test in the signal package of golang.

#### Instructions:

```
go get github.com/iwdgo/sigint-windows
cd <path to sigint-windows>
go run ctrlbreak.go
cat ctrlbreak.log # type ctrlbreak.log
```

It is reference for SO questions:

https://stackoverflow.com/questions/45309984/signal-other-than-sigkill-not-terminating-process-on-windows
https://stackoverflow.com/questions/55092139/gracefully-terminate-a-process-on-windows
