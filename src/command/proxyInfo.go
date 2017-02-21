package command

import (
    "os"
    "strings"
)

type ProxyInfo struct {
    ipAddress string
    httpPort string
}

func (proxyInfo *ProxyInfo) IpAddress() string {
    if proxyInfo.ipAddress == "" {
        populateFields(proxyInfo)
    }

    return proxyInfo.ipAddress
}

func populateFields(proxyInfo *ProxyInfo)  {
    if ds3EndPoint := os.Getenv("DS3_ENDPOINT"); ds3EndPoint != "" {
        if ds3EndPointParts := strings.Split(ds3EndPoint, ":"); len(ds3EndPointParts) >= 3 {
            proxyInfo.httpPort = ds3EndPointParts[2]

            ipAddressParts := strings.Split(ds3EndPointParts[1], "//")
            if len(ipAddressParts) >= 2 {
                proxyInfo.ipAddress = ipAddressParts[1]
            }
        }
    }
}

func (proxyInfo *ProxyInfo) HttpPort() string {
    if proxyInfo.httpPort == "" {
        populateFields(proxyInfo)
    }

    return proxyInfo.httpPort
}




