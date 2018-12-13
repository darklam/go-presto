package presto

type json = map[string]interface{}

//Request : The request struct
type Request struct {
	Body    json
	Headers json
	Query   json
}

//Response : The response object
type Response struct {
	status int
}

//Js : Return json response
func (r Response) Js(statusCode int) int {
	r.status = statusCode
	return r.status
}
