package config

import (
	"command"
	"strings"
	"net/http"
)

type urlCommandEntry struct {
    url string
    httpVerb string
    commandToRun *command.PartialDataFromGetCommand
}

var urlCommandEntryTable []urlCommandEntry
var failFirstAttemptCommand = command.FailFirstAttemptCommand{}

func CommandForUrlPrefix(url string, httpVerb string, ds3HttpClientConnectionInfo *Ds3HttpClientConnectionInfo) command.Command {
    if strings.HasPrefix(url, "/Put_Job_Management_Test/lesmis-copies.txt") && httpVerb == http.MethodPut {
        return &failFirstAttemptCommand
    } else if strings.HasPrefix(url, "/Get_Job_Management_Test/lesmis-copies.txt") && httpVerb == http.MethodGet {
        const maxNunRetries int = 1
        return getPartialDataFromGetCommand(url, httpVerb, ds3HttpClientConnectionInfo, maxNunRetries)
    } else if strings.HasPrefix(url, "/Get_Job_Management_Test/GreatExpectations.txt") && httpVerb == http.MethodGet {
        const maxNunRetries int = 2
        return getPartialDataFromGetCommand(url, httpVerb, ds3HttpClientConnectionInfo, maxNunRetries)
    }

    return nil
}

func getPartialDataFromGetCommand(url string,
                                  httpVerb string,
                                  ds3HttpClientConnectionInfo *Ds3HttpClientConnectionInfo,
                                  maxNumRetries int) *command.PartialDataFromGetCommand {
    var partialDataFromGetCommand *command.PartialDataFromGetCommand

    partialDataFromGetCommand = findGetJobManagementUrlCommandEntry(url, httpVerb)

    if partialDataFromGetCommand == nil {
        partialDataFromGetCommand = &command.PartialDataFromGetCommand{RemoteHost:ds3HttpClientConnectionInfo.RemoteHost,
            MaxNumRetries: maxNumRetries}
        urlCommandEntryTable = append(urlCommandEntryTable, urlCommandEntry{url, httpVerb, partialDataFromGetCommand})
    }

    return partialDataFromGetCommand
}

func findGetJobManagementUrlCommandEntry(url string, httpVerb string) *command.PartialDataFromGetCommand {
    for _, tableEntry := range urlCommandEntryTable {
        if strings.HasPrefix(url, tableEntry.url) && httpVerb == tableEntry.httpVerb {
            return tableEntry.commandToRun
        }
    }

    return nil
}


