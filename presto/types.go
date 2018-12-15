package presto

import (
	"net/http"
	js "encoding/json"
)

type JsonObject = map[string]interface{}

//Request : The request struct
type Request struct {
	Body    JsonObject
	Headers JsonObject
	Query   JsonObject
}

//Response : The response object
type Response struct {
	status  int
	headers map[string]string
	r       http.ResponseWriter
}

//Js : Return json response
func (r Response) Js(statusCode int) int {
	r.status = statusCode
	return r.status
}

func (res Response) Json(response JsonObject) bool {
	writer := res.r
	jsonResponse, err := js.Marshal(response)
	if err != nil {
		panic(err)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonResponse)
	return true
}
