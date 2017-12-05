package command

import (
    "net"
    "net/http"
    "fmt"
    "bufio"
    "strings"
    "strconv"
    "io"
)

type PartialDataFromGetCommand struct {
    RemoteHost string
    numRetries int
    MaxNumRetries int
}

func (partialDataFromGetCommand *PartialDataFromGetCommand) Execute(responseWriter http.ResponseWriter, request *http.Request) (err error, handled bool) {
    handled = true

    fmt.Println("\n=========> PartialDataFromGetCommand")

    socket, err := net.Dial("tcp", partialDataFromGetCommand.RemoteHost + ":80")

    if err != nil {
        fmt.Println(err)
        return
    }

    defer socket.Close()

    writeHttpRequest(&socket, request)

    fmt.Println("Reading result...")

    bufferedSocket := bufio.NewReader(socket)

    bodyLength, httpResult := readAndSetHttpHeaders(bufferedSocket, responseWriter)
    responseWriter.WriteHeader(httpResult)

    var numBytesToTransfer int

    partialDataFromGetCommand.numRetries++

    if partialDataFromGetCommand.numRetries <= partialDataFromGetCommand.MaxNumRetries {
        numBytesToTransfer = int(float32(bodyLength) * 0.8)
    } else {
        numBytesToTransfer = bodyLength
    }

    fmt.Printf("Content length we will return: %d\n", numBytesToTransfer)

    dataToTransfer := readHttpResponseBody(bufferedSocket, numBytesToTransfer)

    numBytesWritten, err := responseWriter.Write(dataToTransfer)

    if err != nil {
        fmt.Println(err)
        return
    }

    if numBytesWritten < numBytesToTransfer {
        fmt.Println("numBytesWritten: ", numBytesWritten, ", numBytesToTransfer: ", numBytesToTransfer)
    }

    fmt.Print("=========> PartialDataFromGetCommand\n\n")

    return
}

func writeHttpRequest(socket *net.Conn, request *http.Request) {
    (*socket).Write([]byte(fmt.Sprintf("GET %s HTTP/1.1\r\n", request.URL)))
    fmt.Printf("GET %s HTTP/1.1\r\n", request.URL)

    for key, value := range request.Header {
        (*socket).Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value[0])))
        fmt.Printf("%s: %s\r\n", key, value[0])
    }

    (*socket).Write([]byte(fmt.Sprintf("Host: %s\r\n\r\n", request.Host)))
    fmt.Println(fmt.Sprintf("Host: %s\r\n\r\n", request.Host))
}

func readAndSetHttpHeaders(bufferedSocket *bufio.Reader, responseWriter http.ResponseWriter) (bodyLength int, httpResult int) {
    bodyLength = 0
    httpResult = 500

    for {
        rawLine, _, err := bufferedSocket.ReadLine()

        if err != nil {
            fmt.Println(err)
            break
        }

        headerLine := string(rawLine)

        setHeaderFromServerReply(responseWriter, headerLine)

        bodyLengthFromHeader := contentLength(headerLine)
        if bodyLengthFromHeader > 0 {
            bodyLength = bodyLengthFromHeader
        }

        httpResultFromHeader := httpResultCode(headerLine)
        if httpResultFromHeader > 0 {
            httpResult = httpResultFromHeader
        }

        if stringIsBlank(headerLine) {
            break
        }
    }

    return
}

func setHeaderFromServerReply(responseWriter http.ResponseWriter, aLine string) {
    headerFields := strings.Split(aLine, ":")

    if len(headerFields) > 1 {
        responseWriter.Header().Set(strings.TrimSpace(headerFields[0]), strings.TrimSpace(headerFields[1]))
    }
}

func contentLength(aString string) int {
    if ! strings.HasPrefix(aString, "Content-Length:") {
        return 0
    }

    stringFields := strings.Split(aString, ":")

    if len(stringFields) < 2 {
        return 0
    }

    contentLength, err := strconv.Atoi(strings.TrimSpace(stringFields[1]))

    fmt.Println("Content-Length", contentLength)

    if err != nil {
        fmt.Println("Error getting content length: ", err)
        contentLength = 0
    }

    return contentLength
}

func httpResultCode(aString string) int {
    if ! strings.HasPrefix(aString, "HTTP") {
        return 0
    }

    stringFields := strings.Split(aString, " ")

    if len(stringFields) < 2 {
        return 0
    }

    httpResultCode, err := strconv.Atoi(strings.TrimSpace(stringFields[1]))

    fmt.Println("Http result: ", httpResultCode)

    if err != nil {
        fmt.Println("err: ", err)
        httpResultCode = 0
    }

    return httpResultCode
}

func stringIsBlank(aString string) bool {
    return len(strings.Trim(aString, "\r\n")) == 0
}

func readHttpResponseBody(bufferedSocket *bufio.Reader, numBytesToRead int) (dataRead []byte) {
    readBuffer := make([]byte, 1024 * 1024)
    var totalNumBytesRead int
    limitedReader := io.LimitedReader{R : bufferedSocket, N : int64(numBytesToRead)}

    for totalNumBytesRead < numBytesToRead {
        var numBytesRead int
        numBytesRead, err := limitedReader.Read(readBuffer)

        if err == nil || err == io.EOF {
            if numBytesRead > 0 {
                dataRead = append(dataRead, readBuffer[0:numBytesRead]...)
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

    return
}
