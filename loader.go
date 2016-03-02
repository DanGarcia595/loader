package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"net/http"
	"time"
)

var ports [5000]bool

func main() {
//	Handle requests for the icon
	http.HandleFunc("/favicon.ico", handlerICon)

//	Loader "main" function. Executed when someone connects to port 80
	http.HandleFunc("/", serve)
        http.ListenAndServe(":80", nil)
}

func serve(w http.ResponseWriter, r *http.Request) {
//	Use the global boolean array to determine the free port number
	port := 1024
	for i:=0; ports[i] == true; i++ {
		port = 1024 + i + 1
	}

//	And reserve it
	ports[port - 1024] = true

//	Create the string for redirect
	redirectString := fmt.Sprintf("http://localhost:%d", port)

//	The goroutine is important as the cmd.Wait() will block execution if not performed concurrently
	go runApp(port)

//	Wait for app to start note: could use mutex instead
	time.Sleep(time.Second * 1)

//	Redirect to the app http: 307 (Temporary Redirect)
	http.Redirect(w,r,redirectString, 307)
}

func handlerICon(w http.ResponseWriter, r *http.Request) {}

func runApp(port int) {
//	Command to open app
	cmdString := fmt.Sprintf("shiny::runApp('/home/dan/srt.core', port = %d)", port)

//	Load command to be executed
	cmd := exec.Command("Rscript", "-e" ,cmdString)

//	Setup connect to stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	fmt.Println("Opening on port", port)

//	Start commmand and make sure it runs fine
	err := cmd.Start()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		//Dont hog the port if you can't start
		ports[port - 1024] = false
		return
	}

//	Wait for command to finish before freeing the port in the boolean array
	if cmd.Wait() == nil {
		fmt.Println("Free on port: ", port)
		ports[port - 1024] = false
	}

}

