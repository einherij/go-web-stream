package main

// First Arg - time between frames
// Second Arg - server port

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/einherij/go-web-stream/auth"
	"github.com/einherij/go-web-stream/conmanager"
	"github.com/einherij/go-web-stream/controller"
	"github.com/einherij/go-web-stream/streamvideo"
	// "os"
	// "image/color"
)

const (
	// time to store cookie
	storeCookieTime = time.Minute * 60
)

var (
	// Default time between frames and server port
	timeBetweenFrames = 100 * time.Millisecond
	serverPort        = 9000
)

// Page contains page structure to parse by page template
type Page struct {
	Title string
	Body  template.HTML
}

// loadPage loades page template
func loadPage(title string) *Page {
	filename := "./templates/" + title + ".html"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: template.HTML(body)}
}

// keepAliveHandler is handler to perform keep alive requests from client
func keepAliveHandler(w http.ResponseWriter, r *http.Request) {
	// get cookie from request
	conID, err := r.Cookie("ConnectionID")
	if err != nil {
		log.Fatal("!There is no cookie in the request, error:", err)
	}
	// send keep alive for current connection id
	currentConnectionID := conID.Value
	conmanager.KeepConnectionAlive(currentConnectionID)
}

// videoFeedHandler is handler to perform requests for video stream
func videoFeedHandler(w http.ResponseWriter, r *http.Request) {
	// create new connection for connection manager
	connect := conmanager.NewConnection()
	// strat timeout for current connection
	go connect.StartTimer()
	// time to live for cookies
	expiration := time.Now().Add(storeCookieTime)
	// set connection id as cookie
	cookie := http.Cookie{Name: "ConnectionID", Value: strconv.Itoa(connect.ConnectionID), Expires: expiration}
	http.SetCookie(w, &cookie)
	// Set header for multipart video stream object
	w.Header().Set("Content-Type", "multipart/x-mixed-replace;boundary=frame")
	// Reading frames loop
	for {
		// stop sending frames to disconnected users
		if connect.MustDisconnect == true {
			fmt.Println("End connection id:", connect.ConnectionID)
			break
		}
		// get frame from frames channel
		currentFrame := <-streamvideo.Frames
		// write frame to web page in multipart object as image/jpeg
		w.Write(currentFrame)
	}
}

// mainHandler is handler to return main page of the server
func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		if auth.LoginUser(r.FormValue("Login"), r.FormValue("Password")) {
			// send html document to web page
			p := loadPage("main_page")
			t, err := template.ParseFiles("./templates/page.html")
			if err != nil {
				fmt.Println("Can't load template", p.Title, ", error:", err)
			}
			t.Execute(w, p)
		} else {
			p := loadPage("error_page")
			t, err := template.ParseFiles("./templates/page.html")
			if err != nil {
				fmt.Println("Can't load template", p.Title, ", error:", err)
			}
			t.Execute(w, p)
		}
	}
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	// send html document to web page
	p := loadPage("register_page")
	t, err := template.ParseFiles("./templates/page.html")
	if err != nil {
		fmt.Println("Can't load template", p.Title, ", error:", err)
	}
	t.Execute(w, p)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	// Register new user
	if r.Method == "POST" {
		r.ParseForm()
		p := auth.NewPerson(r.FormValue("Login"), r.FormValue("Name"), r.FormValue("Email"), r.FormValue("Password"))
		p.Save()
	}
	// send html document to web page
	p := loadPage("auth_page")
	t, err := template.ParseFiles("./templates/page.html")
	if err != nil {
		fmt.Println("Can't load template", p.Title, ", error:", err)
	}
	t.Execute(w, p)
}

func controlHandler(w http.ResponseWriter, r *http.Request) {
	command := r.URL.Path[len("/run/"):]
	// fmt.Println("Command number is:", command)
	com, _ := strconv.Atoi(command)
	controller.RunRobot(com)
}

func main() {

	// To know how it goes
	fmt.Println("Begin...")

	// Args parsing
	if len(os.Args) > 1 {
		intTimeBetweenFrames, _ := strconv.Atoi(os.Args[1])
		timeBetweenFrames = time.Duration(intTimeBetweenFrames) * time.Millisecond
	}
	if len(os.Args) > 2 {
		serverPort, _ = strconv.Atoi(os.Args[2])
	}

	// Start RPIO
	controller.StartRPIO()
	defer controller.StopRPIO()

	// Read video from camera
	go streamvideo.StartReadingVideo(timeBetweenFrames)

	// add static files directory
	staticFilesHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)

	http.HandleFunc("/", signInHandler)
	http.HandleFunc("/video_feed", videoFeedHandler)
	http.HandleFunc("/keep_alive", keepAliveHandler)
	http.HandleFunc("/register", registrationHandler)
	http.HandleFunc("/main", mainHandler)
	http.Handle("/data/", staticFilesHandler)
	http.HandleFunc("/run/", controlHandler)

	fmt.Println("Start server on", serverPort, "port and time between frames", timeBetweenFrames)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(serverPort), nil))
}
