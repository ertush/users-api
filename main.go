package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ertush/collectibles"
	"github.com/muesli/cache2go"
)

var cache = cache2go.Cache("data-store")

type user interface{}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

// RegexpHandler is not exported
type RegexpHandler struct {
	routes []*route
}

// Handler is not exported
func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

// HandleFunc is not exported
func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

// ServeHTTP is not exported
func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func listUsersHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Add("Host", "localhost:8080")

		n := cache.Count()

		payloads := make([]*cache2go.CacheItem, n)

		for i := 0; i < n; i++ {
			var r int
			r = i + 1
			payloads[i], _ = cache.Value(r)

		}

		allUsers := make([]user, n)

		for user := 0; user < n; user++ {
			allUsers[user] = payloads[user].Data()

		}

		var uID int
		var i string
		var b bool

		if b = strings.Contains(req.RequestURI, ":"); b {
			i = strings.Split(req.RequestURI, ":")[1]
			uID, _ = strconv.Atoi(i)

		} else {
			w.WriteHeader(http.StatusOK)
			str := fmt.Sprintf("%s", allUsers)
			s := strings.ReplaceAll(str, "} {", "},{")
			fmt.Fprintf(w, "%v", s)

		}

		if uID > 0 && uID <= n {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "%+v", allUsers[int(uID)-1])
		} else {
			if b {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "{\"status\": 404,\"message\": \"resource not found\"}")
			}

		}
	}
}

// Sanitize Json
func cx(b byte) bool {
	retVal := false
	if b != byte('[') && b != byte(']') {
		retVal = true
	}

	return retVal

}

func createUsersHandler(w http.ResponseWriter, req *http.Request) {

	var b *collectibles.Byt
	buffer := make([]byte, req.ContentLength)

	var userID int

	if req.Method == http.MethodPost {

		if _, err := req.Body.Read(buffer); err.Error() != "EOF" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("[ERROR]: ", err)
		}

		userID += cache.Count()
		userID++

		stringBuffer := string(buffer)

		if strings.Contains(stringBuffer, "[") {
			if b.Include([]byte(stringBuffer), '[') {
				json := strings.Split(string(b.Filter([]byte(stringBuffer), cx)), "},")

				for index, item := range json {

					index += 1

					if !strings.Contains(item, "}") {
						item += "}"
					}

					if strings.Contains(item, "name") && strings.Contains(item, "gender") && strings.Contains(item, "email") {
						cache.Add(index, 0, item)
						// w.WriteHeader(http.StatusOK)
						w.Header().Set("Content-Type", "application/json")
						w.Header().Add("Host", "localhost:8080")
						fmt.Fprintf(w, "{\"Success\": \"true\",\"Message\": \"User Created with id %d\",\"Content-Size\": %d}", index, req.ContentLength)
					} else {
						w.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(w, "{\"Success\": \"false\",\"Message\": \"Expected at least fields - name, gender and email\",\"Content-Size\": %d}", req.ContentLength)
					}

				}

			}
		} else {

			if strings.Contains(stringBuffer, "name") && strings.Contains(stringBuffer, "gender") && strings.Contains(stringBuffer, "email") {
				cache.Add(userID, 0, stringBuffer)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Add("Host", "localhost:8080")

				fmt.Fprintf(w, "{\"Success\": \"true\",\"Message\": \"User Created with id %d\",\"Content-Size\": %d}", userID, req.ContentLength)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "{\"Success\": \"false\",\"Message\": \"Expected at least fields - name, gender and email\",\"Content-Size\": %d}", req.ContentLength)
			}

		}

	}
}

func deleteUsersHandler(w http.ResponseWriter, req *http.Request) {

	buffer := make([]byte, req.ContentLength)

	if req.Method == http.MethodDelete {
		if _, err := req.Body.Read(buffer); err.Error() != "EOF" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("[ERROR]: ", err)

			reqBuffer := string(buffer)
			fmt.Println(reqBuffer)

			var uID int
			var i string
			var b bool

			if b = strings.Contains(req.RequestURI, ":"); b {
				i = strings.Split(req.RequestURI, ":")[1]
				uID, _ = strconv.Atoi(i)

				if _, err := cache.Delete(uID); len(err.Error()) > 0 {
					fmt.Fprintf(w, "{\"Success\": \"false\",\"Message\": \"Unable to delete user with id [%v]\",\"Content-Size\": %d}", err.Error(), req.ContentLength)
				}
				w.WriteHeader(http.StatusNoContent)
				fmt.Fprintf(w, "{\"Success\": \"true\",\"Message\": \"User with id %d has been deleted\",\"Content-Size\": %d}", uID, req.ContentLength)

			} else {
				w.WriteHeader(http.StatusNoContent)

				fmt.Fprintf(w, "{\"Success\": \"false\",\"Message\": \"Please specify a user to delete with an id\",\"Content-Size\": %d}", req.ContentLength)
			}

		}
	}

}
func main() {

	handler := &RegexpHandler{}

	reg1, _ := regexp.Compile("/list-users:\\d")
	handler.HandleFunc(reg1, listUsersHandler)

	reg2, _ := regexp.Compile("/list-users")
	handler.HandleFunc(reg2, listUsersHandler)

	reg3, _ := regexp.Compile("/create-user")
	handler.HandleFunc(reg3, createUsersHandler)

	reg4, _ := regexp.Compile("/delete-user:\\d")
	handler.HandleFunc(reg4, deleteUsersHandler)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 Mib
	}

	log.Println("[+] Server is listening on localhost:8080")
	log.Fatal(s.ListenAndServe())

}
