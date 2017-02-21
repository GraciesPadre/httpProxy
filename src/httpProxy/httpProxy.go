package main

import (
	"os"
	"net/http"
	"net"
	"command"
	"fmt"
    "config"
)

var httpClientClientConnectionInfo config.Ds3HttpClientConnectionInfo
var httpListener net.Listener
var defaultCommand command.DefaultCommand

// TODO set the remote host and other configuration stuff from admin request

func main() {
	getDs3ConnectionInfoFromEnvironmentVars()

	defaultCommand = command.DefaultCommand{RemoteHost:httpClientClientConnectionInfo.RemoteHost}

	go listenOnAdminPort()

	listenOnProxyPort(func () {
		fmt.Printf("Listening on proxy port %s and admin port %s and forwarding to %s\n",
			httpClientClientConnectionInfo.ProxyPort,
			httpClientClientConnectionInfo.AdminPort,
            httpClientClientConnectionInfo.RemoteHost)
	})
}

func getDs3ConnectionInfoFromEnvironmentVars() {
	httpClientClientConnectionInfo.RemoteHost = getEnvironmentVar("DS3_TARGET_SYSTEM_DNS_NAME", "sm2u-11.eng.sldomain.com")
	httpClientClientConnectionInfo.ProxyPort = getEnvironmentVar("HTTP_PROXY_PORT", ":9080")
	httpClientClientConnectionInfo.AdminPort = getEnvironmentVar("HTTP_PROXY_ADMIN_PORT", ":9090")
}

func getEnvironmentVar(envVarToLookFor string, defaultIfVarNotSet string) string {
	result := os.Getenv(envVarToLookFor)

	if len(result) == 0 {
		return defaultIfVarNotSet
	}

	return result
}

func listenOnAdminPort() {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/close", finit)

	err := http.ListenAndServe(httpClientClientConnectionInfo.AdminPort, adminMux)

	if err != nil {
		fmt.Printf("Error listening on admin port: %s\n", httpClientClientConnectionInfo.AdminPort)
		panic(err)
	}
}

func listenOnProxyPort(onSuccess func()) {
	var err error

	httpListener, err = net.Listen("tcp", httpClientClientConnectionInfo.ProxyPort)

	if err != nil {
		fmt.Printf("Error listening on port: %s\n", httpClientClientConnectionInfo.ProxyPort)
		panic(err)
	}

	onSuccess()

	httpHandler := http.NewServeMux()
	httpHandler.HandleFunc("/", proxyHandler)
	http.Serve(httpListener, httpHandler)
}

var commandToRun command.Command

func proxyHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

    fmt.Println("Processing ", request.Method, ": ", request.RequestURI)
    fmt.Println("  Request header: ", request.Header)
    fmt.Println("  Request body: ", request.Body)

	commandToRun = config.CommandForUrlPrefix(request.RequestURI, request.Method, &httpClientClientConnectionInfo)

	if commandToRun != nil {
		err, handled := commandToRun.Execute(responseWriter, request)

		if(err != nil) {
			fmt.Println("Error: ", err, " processing http request: ", request);
		}

		if handled {
			return
		}
	}

	commandToRun = &defaultCommand

	commandToRun.Execute(responseWriter, request)
}

func finit (responseWriter http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	if request.Method == "PUT" {
		responseWriter.Write([]byte("Closing...\n"))
		fmt.Printf("Later dude\n")
		httpListener.Close()
	}
}

