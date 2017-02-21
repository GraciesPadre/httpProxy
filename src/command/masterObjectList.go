package command

import (
    "encoding/xml"
    "fmt"
    "strings"
)

type Node struct {
    EndPoint string `xml:"EndPoint,attr"`
    HttpPort string `xml:"HttpPort,attr"`
    HttpsPort string `xml:"HttpsPort,attr"`
    Id string `xml:"Id,attr"`
}

type Nodes struct {
    Node []*Node
}

type Object struct {
    Bucket string `xml:"Bucket,attr"`
    Id string `xml:"Id,attr"`
    InCache string `xml:"InCache,attr"`
    Latest string `xml:"Latest,attr"`
    Length string `xml:"Length,attr"`
    Name string `xml:"Name,attr"`
    Offset string `xml:"Offset,attr"`
    Version string `xml:"Version,attr"`
}

type Objects struct {
    ChunkId string `xml:"ChunkId,attr"`
    ChunkNumber string `xml:"ChunkNumber,attr"`
    NodeId string `xml:"NodeId,attr"`
    Object []Object
}

type MasterObjectList struct {
    Aggregating string `xml:"Aggregating,attr"`
    BucketName string `xml:"BucketName,attr"`
    CachedSizeInBytes string `xml:"CachedSizeInBytes,attr"`
    ChunkClientProcessingOrderGuarantee string `xml:"ChunkClientProcessingOrderGuarantee,attr"`
    CompletedSizeInBytes string `xml:"CompletedSizeInBytes,attr"`
    EntirelyInCache string `xml:"EntirelyInCache,attr"`
    JobId string `xml:"JobId,attr"`
    Naked string `xml:"Naked,attr"`
    Name string `xml:"Name,attr"`
    OriginalSizeInBytes string `xml:"OriginalSizeInBytes,attr"`
    Priority string `xml:"Priority,attr"`
    RequestType string `xml:"RequestType,attr"`
    StartDate string `xml:"StartDate,attr"`
    Status string `xml:"Status,attr"`
    UserId string `xml:"UserId,attr"`
    UserName string `xml:"UserName,attr"`
    Nodes Nodes
    Objects Objects
}


func NewMasterObjectList(masterObjectList []byte) *MasterObjectList {
    newMasterObjectList := new(MasterObjectList)

    if err := xml.Unmarshal(masterObjectList, newMasterObjectList); err != nil {
        fmt.Println("Error deserializing MasterObjectList xml: ", err)
        return nil
    }

    return newMasterObjectList
}

func (masterObjectList *MasterObjectList) SetNodesEndpoint(endPoint string) {
    if endPoint == "" {
        return
    }

    for _, node := range masterObjectList.Nodes.Node {
        node.EndPoint = endPoint
    }
}

func (masterObjectList *MasterObjectList) SetNodesHttpPort(httpPort string) {
    if httpPort == "" {
        return
    }

    for _, node := range masterObjectList.Nodes.Node {
        node.HttpPort = httpPort
    }
}

func IsMasterObjectListXml(xmlString string) bool {
    return strings.HasPrefix(xmlString, "<MasterObjectList")
}
