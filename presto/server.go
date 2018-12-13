package presto

import (
	"net/http"
)

type handler struct {
	path     string
	function func(Request, Response)
}

//Server : The server struct lol
type Server struct {
	port       string
	middleware []map[string]func(Request, Response)
}

func handleRequest(w http.ResponseWriter, r *http.Request, s *Server) {
	middleware := s.middleware
	path := r.RequestURI
	// gotta build the Request and Response to pass lol
	req := Request{}
	res := Response{}
	req.Body = make(map[string]interface{})
	req.Body["test"] = true
	var handler func(Request, Response)
	// this is incredibly slow and I have to say I am embarrassed
	// but hey this is just to see if it all works
	for i := 0; i < len(middleware); i++ {
		current := middleware[i]
		for key := range current {
			if key == path || key == "*" {
				handler = current[key]
				break;
			}
		}
		if handler != nil {
			break;
		}
	}
	if handler != nil {
		handler(req, res)
	}
}

func (s *Server) Get(path string, handler func(Request, Response)) {
	temp := make(map[string]func(Request, Response))
	temp[path] = handler
	s.middleware = append(s.middleware, temp)
}

func (s *Server) Use(handler func(Request, Response)) {
	temp := make(map[string]func(Request, Response))
	temp["*"] = handler
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
	http.ListenAndServe(":" + port, nil)
}
