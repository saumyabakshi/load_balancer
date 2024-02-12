package main

import (
	"fmt"
	"log"
	"net/http"
	"io"
)

func main() {
	fmt.Println("Starting server at port 80\n")

	backendURLS := []string{"http://localhost:8081/backend", "http://localhost:8082/backend", "http://localhost:8083/backend"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		backendURL := chooseServer(backendURLS)
		//create a new request
		req, err := http.NewRequest(r.Method, backendURL, r.Body)
		if err != nil {
			log.Fatal(err)
		}
		//copy the headers from the original request
		log.Printf("Header from the original request: %s\n", r.Header)
		// for k,v := range r.Header {
		// 	req.Header.Set(k,v)
		// }
		req.Header = r.Header

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		// defer resp.Body.Close()
		// _, err = io.Copy(w, resp.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		fmt.Printf("Recieved request from %s\n", r.RemoteAddr)
		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		fmt.Printf("Host: %s\n", r.Host)
		fmt.Printf("User-Agent: %s\n", r.UserAgent())
		fmt.Printf("Accept: %s\n", r.Header.Get("Accept"))
		fmt.Printf("Response from the server: %s %s\n", resp.Proto, resp.Status)
		io.WriteString(w, "Hello from the backend server\n")
	})

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}

}

func chooseServer(servers []string) string {
	//  round robin
	
