package REST

import (
	"net/http"
	"log"
	"bytes"
	"os"
	"encoding/json"
)

type Response interface {
	WriteBody(response http.ResponseWriter, request *http.Request)
	GetStatus() int
	GetContentType() string
}

type baseResponse struct {
	Status int
	contentType string
}

func (br baseResponse) SetStatus(status int) {
	br.Status = status
	logger.Debug("HAHAHHAHAHHAHAHHHHHH :::: %v", br)
}

func (br baseResponse) GetStatus() int {
	return br.Status
}

func (br baseResponse) GetContentType() string {
	return br.contentType
}

/// JSON RESPONSE

type JsonResponse struct {
	baseResponse
	json map[string]interface{}
}

func NewJsonResponse(json map[string]interface{}) *JsonResponse {
	return &JsonResponse{baseResponse{200, "application/json; charset=UTF-8"}, json}
}

func (jr JsonResponse) Merge(other map[string]interface{}) {

}

func (jr JsonResponse) WriteBody(response http.ResponseWriter, request *http.Request) {
	r, err := json.Marshal(jr.json)
	if err != nil {
		GenerateErrorResponse(response, err.Error(), 500)
		return
	}

	response.Write(r)
	//response.WriteHeader(200)
	//response.Header().Add("Content-type", "application/json; charset=UTF-8")
}

// BYTE JSON

type JsonByteResponse struct {
	baseResponse
	Json []byte
}

func NewJsonByteResponse(json1 interface{}) *JsonByteResponse {
	b, err := json.Marshal(json1)
	if err != nil {
		logger.Warning("Failed: %v", err)
		b = nil
	}
	return &JsonByteResponse{baseResponse{200, "application/json; charset=UTF-8"}, b}
}

func (jr JsonByteResponse) WriteBody(response http.ResponseWriter, request *http.Request) {
	response.Write(jr.Json)
}


/// Buffered response

type BufferedResponse struct {
	baseResponse
	Bytes *bytes.Buffer
}

func NewBufferedResponse() *BufferedResponse {
	return &BufferedResponse{baseResponse{200, "text/html; charset=UTF-8"}, new(bytes.Buffer)}
}

func (br BufferedResponse) Append(content string) {
	br.Bytes.WriteString(content)
}

func (br BufferedResponse) WriteBody(response http.ResponseWriter, request *http.Request) {
	//br.Bytes.WriteTo(response)
	response.Write(br.Bytes.Bytes())
}

// Redirect

type RedirectResponse struct  {
	baseResponse
	url string
}

func NewRedirectResponse(url string) *RedirectResponse {
	return &RedirectResponse{baseResponse{200, "text/plain; charset=UTF-8"}, url}
}

func (rr RedirectResponse) WriteBody(response http.ResponseWriter, request *http.Request) {
	logger.Debug("REDIRECTNG------ %v", rr.url)
	http.Redirect(response, request, rr.url, 301)
}

// File response

type FileResponse struct  {
	baseResponse
	file string
}

func NewFileResponse(file string) *FileResponse {
	return &FileResponse{baseResponse{200, ""}, file}
}

func (fr FileResponse) WriteBody(response http.ResponseWriter, request *http.Request) {
	d, _ := os.Getwd()
	log.Printf("Serving %v from %v", fr.file, d)
	if _, err := os.Stat(fr.file); os.IsNotExist(err) {
		log.Printf("no such file or directory: %s", fr.file)
		return
	}

	http.ServeFile(response, request, fr.file)
}

/// Error response

func GenerateErrorResponse(response http.ResponseWriter, error string, status int) {
	response.Write([]byte(error))
	//response.WriteHeader(status)
}