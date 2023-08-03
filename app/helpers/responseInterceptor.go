package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InterceptResponse represents the structure of the intercepted response.
type InterceptResponse struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Error   bool        `json:"error"`
	Message string      `json:"message"`
}

// ResponseInterceptor is the middleware function that formats the response.
func ResponseInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before calling the main handler, capture the response and its status code.
		captureWriter := &interceptResponseWriter{ResponseWriter: c.Writer}
		c.Writer = captureWriter
		c.Next()

		fmt.Println(captureWriter.responseData)
		// Build the formatted response.
		var recievedData interface{}
		if captureWriter.statusCode >= 400 {
			recievedData = nil
		} else {
			recievedData =  captureWriter.responseData
		}
		response := InterceptResponse{
			Data:    recievedData,
			Status:  captureWriter.statusCode,
			Error:   captureWriter.statusCode >= 400,
			Message: captureWriter.statusText,
		}

		// Set the formatted response.
		c.JSON(captureWriter.statusCode, response)
	}
}

// interceptResponseWriter is a custom response writer to capture response data and status code.
type interceptResponseWriter struct {
	gin.ResponseWriter
	responseData interface{}
	statusCode   int
	statusText   string
}

// Write method captures the response data.
func (w *interceptResponseWriter) Write(data []byte) (int, error) {
	if w.responseData == nil {
		w.responseData = make(map[string]interface{})
	}

	var raw interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return w.ResponseWriter.Write(data)
	}

	w.responseData = raw
	return w.ResponseWriter.Write(data)
}

// WriteHeader method captures the response status code and status text.
func (w *interceptResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.statusText = http.StatusText(statusCode)
	w.ResponseWriter.WriteHeader(statusCode)
}
