package command

import (
    "net/http"
    "strconv"
    "fmt"
    "io"
)

type PartialDataFromGetCommand struct {
    RemoteHost string
    numRetries int
    MaxNumRetries int
}

func (partialDataFromGetCommand *PartialDataFromGetCommand) Execute(responseWriter http.ResponseWriter, request *http.Request) (err error, handled bool) {
    handled = true

    request.URL.Host = partialDataFromGetCommand.RemoteHost
    request.RequestURI = ""

    if request.URL.Scheme == "" {
        request.URL.Scheme = "http"
    }

    httpClient := http.Client{}
    var response *http.Response
    response, err = httpClient.Do(request)

    if err != nil {
        fmt.Println(err)
        return
    }

    defer response.Body.Close()

    contentLengthHeaderId := "Content-Length"

    contentLen, err := strconv.Atoi(response.Header.Get(contentLengthHeaderId))

    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("=========>, Http request content length: %d\n", contentLen)

    var numBytesToTransfer int

    partialDataFromGetCommand.numRetries++

    if partialDataFromGetCommand.numRetries <= partialDataFromGetCommand.MaxNumRetries {
        numBytesToTransfer = contentLen / 4
    } else {
        numBytesToTransfer = contentLen
    }

    fmt.Printf("=========>, Content length we will return: %d\n", numBytesToTransfer)

    writeHeaders(responseWriter, response)
    responseWriter.Header().Set(contentLengthHeaderId, fmt.Sprintf("%d", contentLen))
    responseWriter.WriteHeader(response.StatusCode)

    readBuffer := make([]byte, 1024 * 1024)
    var totalNumBytesRead int
    var dataToTransfer []byte
    limitedReader := io.LimitedReader{R : response.Body, N : int64(numBytesToTransfer)}

    for totalNumBytesRead < numBytesToTransfer {
        var numBytesRead int
        numBytesRead, err = limitedReader.Read(readBuffer)

        fmt.Printf("=========>, Num bytes read: %d\n", numBytesRead)

        if err == nil || err == io.EOF {
            if numBytesRead > 0 {
                dataToTransfer = append(dataToTransfer, readBuffer[0:numBytesRead]...)
            }

            if err == io.EOF {
                break
            }
        } else {
            fmt.Println(err)
            return
        }

        totalNumBytesRead += numBytesRead
    }

    var numBytesWritten int

    numBytesWritten, err = responseWriter.Write(dataToTransfer)

    if err != nil {
        fmt.Println(err)
        return
    }

    if numBytesWritten != len(dataToTransfer) {
        fmt.Printf("numBytesWritten: %d, data length: %d\n", numBytesWritten, len(dataToTransfer))
    }

    fmt.Println("  Response header: ", response.Header)
    fmt.Println("  Response body: ", response.Body)

    return
}
