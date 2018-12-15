package presto

import (
	"net/http"
	"fmt"
	js "encoding/json"
	"io/ioutil"
)

type handler struct {
	path     string
	method   string
	handle   func(Request, Response) bool
}

func handlerMatches(current handler, method string, path string) bool {
	if current.method == "*" && current.path == "*" {
		return true
	}

	if current.method != method && current.method != "*" {
		return false
	}

	handlerPath := current.path
	if handlerPath == path {
		return true
	}
	return false
}

func jsonBodyParser(r *http.Request) (JsonObject, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if len(data) == 0 {
		return nil, nil
	}

	result := make(JsonObject)

	err = js.Unmarshal(data, &result)

	return result, err
}

//Server : The server struct lol
type Server struct {
	port       string
	middleware []handler
}

func initVars(w http.ResponseWriter, r *http.Request) (Request, Response){
	// gotta build the Request and Response to pass lol
	req := Request{}
	res := Response{}
	body, err := jsonBodyParser(r)
	if err != nil {
		println("okay here")
		fmt.Println(err)
		panic(err)
	}
	if body != nil {
		req.Body = body
	} else {
		req.Body = make(JsonObject)
	}
	res.r = w;

	return req, res
}

func handleRequest(w http.ResponseWriter, r *http.Request, s *Server) {
	middleware := s.middleware
	path := r.RequestURI
	method := r.Method
	var requestDone bool
	req, res := initVars(w, r)
	// this is incredibly slow and I have to say I am embarrassed
	// but hey this is just to see if it all works
	for i := 0; i < len(middleware); i++ {
		current := middleware[i]
		if handlerMatches(current, method, path) {
			requestDone = current.handle(req, res)
			fmt.Println(requestDone)
			if (requestDone) {
				break
			}
		}
	}
}

func addHandler(s *Server, path string, handlerFunction func(Request, Response) bool, method string) {
	temp := handler{}
	temp.handle = handlerFunction
	temp.method = method
	temp.path = path
	if method == "" {
		temp.method = "*"
	}
	if path == "" {
		temp.path = "*"
	}
	s.middleware = append(s.middleware, temp)
}

func (s *Server) Get(path string, handlerFunction func(Request, Response) bool) {
	addHandler(s, path, handlerFunction, "GET")
}

func (s *Server) Post(path string, handlerFunction func(Request, Response) bool) {
	addHandler(s, path, handlerFunction, "POST")
}

func (s *Server) Use(handlerFunction func(Request, Response) bool, path string) {
	addHandler(s, path, handlerFunction, "*")
}

func (s *Server) Start(port string) {
	if port == "" {
		port = "1234"
	}
	s.port = port
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		println("something")
		handleRequest(w, r, s)
	})
	error := http.ListenAndServe(":" + port, nil)
	if error != nil {
		fmt.Println(error)
	}
}
