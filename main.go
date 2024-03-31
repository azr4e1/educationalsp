package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/azr4e1/educationalsp/analysis"
	"github.com/azr4e1/educationalsp/lsp"
	"github.com/azr4e1/educationalsp/rpc"
)

var Shutdown bool

func main() {
	logger := getLogger("/home/ld/Desktop/Projects/educationalsp/log.txt")
	logger.Println("Hey, I started!")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analysis.NewState()
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Got an error: %s", err)
		}

		handleMessage(logger, writer, state, method, contents)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
	logger.Printf("Received msg with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)
		// gotta reply to this
		msg := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, msg)

		logger.Println("Sent the reply")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		logger.Printf("Opened: %s %s",
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Text,
		)
		diagnostics := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		writeResponse(writer, lsp.PublishDiagnosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params: lsp.PublishDiagnosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})

	case "textDocument/didChange":
		var request lsp.DidChangeTextDocument
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		logger.Printf("Changed document: %s, Version: %d",
			request.Params.TextDocument.URI,
			request.Params.TextDocument.Version,
		)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.Update(request.Params.TextDocument.URI, change.Text)
			writeResponse(writer, lsp.PublishDiagnosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params: lsp.PublishDiagnosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/hover":
		var request lsp.HoverTextDocumentRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		response := state.Hover(request.ID, request.Param.TextDocument.URI, request.Param.Position)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionTextDocumentRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}

		response := state.Definition(request.ID, request.Param.TextDocument.URI, request.Param.Position)
		writeResponse(writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionTextDocumentRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
			return
		}

		// create a response
		response := state.CodeAction(request.ID, request.Params.TextDocument.URI)

		// write it back
		writeResponse(writer, response)
	case "textDocument/completion":
		var request lsp.CompletionTextDocumentRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}
		response := state.Completion(request.ID, request.Params.TextDocument.URI)

		writeResponse(writer, response)
	case "shutdown":
		var request lsp.ShutdownRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Hey, we couldn't parse this: %s", err)
		}
		response := lsp.ShutdownResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &request.ID,
			},
			Result: nil,
		}
		writeResponse(writer, response)
		Shutdown = true
	case "exit":
		if !Shutdown {
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("hey, you didn't give me a good file")
	}

	return log.New(logFile, "[educationalsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
