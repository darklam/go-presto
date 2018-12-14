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
	handle   func(Request, Response)
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

func jsonBodyParser(r *http.Request) (json, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	result := make(json)

	err = js.Unmarshal(data, &result)

	return result, err
}

//Server : The server struct lol
type Server struct {
	port       string
	middleware []handler
}

func handleRequest(w http.ResponseWriter, r *http.Request, s *Server) {
	middleware := s.middleware
	path := r.RequestURI
	method := r.Method
	// gotta build the Request and Response to pass lol
	req := Request{}
	res := Response{}
	body, err := jsonBodyParser(r)
	if err != nil {
		panic(err)
	}
	req.Body = body
	// this is incredibly slow and I have to say I am embarrassed
	// but hey this is just to see if it all works
	for i := 0; i < len(middleware); i++ {
		current := middleware[i]
		if handlerMatches(current, method, path) {
			current.handle(req, res)
			break;
		}
	}
}

func addHandler(s *Server, path string, handlerFunction func(Request, Response), method string) {
	temp := handler{}
	temp.handle = handlerFunction
	temp.method = method
	temp.path = path
	s.middleware = append(s.middleware, temp)
}

func (s *Server) Get(path string, handlerFunction func(Request, Response)) {
	addHandler(s, path, handlerFunction, "GET")
}

func (s *Server) Post(path string, handlerFunction func(Request, Response)) {
	addHandler(s, path, handlerFunction, "POST")
}

func (s *Server) Use(handlerFunction func(Request, Response), path string) {
	temp := handler{}
	temp.method = "*"
	temp.path = "*"
	if path != "" {
		temp.path = path
	}
	s.middleware = append(s.middleware, temp)
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
