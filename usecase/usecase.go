package usecase

import (
	"net/http"
	"log"
	"encoding/json"

	"fmt"
	"github.com/oranmoshai/go-mcp-streamable-sse-minimal/entity"
)

type Usecase struct {
}

func New() *Usecase {
	return &Usecase{}
}

func (u *Usecase) EventStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// You might want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	log.Printf("Client connected: %s", r.RemoteAddr)

	if r.Method == http.MethodPost || r.Method == http.MethodGet {

		decoder := json.NewDecoder(r.Body)

		var request entity.MCPRequest
		if err := decoder.Decode(&request); err != nil {
			return
		}
		switch request.Method {
		case "initialize":
			response(w, flusher, request,
				json.RawMessage(`{"protocolVersion":"2025-03-26","capabilities":{"logging":{},"prompts":{"listChanged":true},"resources":{"subscribe":true,"listChanged":true},"tools":{"listChanged":true}},"serverInfo":{"name":"ExampleServer","version":"1.0.0"},"instructions":"Optional instructions for the client"}`))
		case "tools/list":
			response(w, flusher, request,
				json.RawMessage(`{"tools":[{"name":"get_weather","description":"Get current weather information for a location","inputSchema":{"type":"object","properties":{"location":{"type":"string","description":"City name or zip code"}},"required":["location"]}}],"nextCursor":"next-page-cursor"}`))
		case "tools/call":
			response(w, flusher, request,
				json.RawMessage(`{"content":[{"type":"text","text":"Current weather in New York:\nTemperature: 72°F\nConditions: Partly cloudy"}],"isError":false}`))
		default:
			sendSSEError(w, flusher, request, "unimplemented method")
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

	flusher.Flush()
}

func response(w http.ResponseWriter, flusher http.Flusher, request entity.MCPRequest, result json.RawMessage) {

	for k, v := range request.Params {
		fmt.Printf("Request Param [%s]: %+v\n", k, v)
	}

	progressResponse := entity.MCPResponse{
		JsonRpc: "2.0",
		ID:      request.ID,
		Result:  result,
	}
	sendSSE(w, flusher, progressResponse)
}

func sendSSE(w http.ResponseWriter, flusher http.Flusher, data interface{}) {

	if data != nil {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			dataJSON = []byte(`{"error":"Failed to marshal data"}`)
		}
		fmt.Fprintf(w, "data: %s\n\n", dataJSON)
		fmt.Printf("data: %s\n", dataJSON)
	}

	flusher.Flush()
}

func sendSSEError(w http.ResponseWriter, flusher http.Flusher, request entity.MCPRequest, message string) {
	response := entity.MCPResponse{
		JsonRpc: "2.0",
		ID:      request.ID,
		Error: &entity.MCPError{
			Code:    1,
			Message: message,
		},
	}
	sendSSE(w, flusher, response)
}
