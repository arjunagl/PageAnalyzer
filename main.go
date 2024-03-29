package main

import (
	"arjunagl/htmlAnalyzer/domain/service"
	"arjunagl/htmlAnalyzer/infrastructure"
	"arjunagl/htmlAnalyzer/infrastructure/adapter"
	"arjunagl/htmlAnalyzer/infrastructure/adapter/controller"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func buildApp() *infrastructure.Application {

	// initialize the application
	wcd := &adapter.WebContentDownloader{}
	tas := &service.TitleAnalyzeService{}
	vas := &service.VersionAnalyzeService{}
	has := &service.HeadingAnalyzeService{}
	las := service.NewLinkAnalyzeService(wcd)
	lgas := &service.LoginAnalyzeService{}
	cr := adapter.NewGoqueryAdapter()
	as := service.NewAnalyzerService([]*service.ItemAnalyzeService{
		{Title: "Title", ContentAnalyzer: tas},
		{Title: "Version", ContentAnalyzer: vas},
		{Title: "Link", ContentAnalyzer: las},
		{Title: "Login", ContentAnalyzer: lgas},
		{Title: "Heading", ContentAnalyzer: has}}, cr, wcd)
	ss := service.NewSchedulerService(as)

	app := infrastructure.NewApplication(wcd, ss, cr)
	return app
}

func startHttpServer() {

	// initialize the application
	app := buildApp()
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})

	ca := controller.NewAnalyzeController(app)
	r.HandleFunc("/analyze", ca.Analyze).Methods(http.MethodPost)
	r.HandleFunc("/analyze", ca.GetResults).Methods(http.MethodGet)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	httpServer := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:80", // Listen on all IP addresses on the machine
	}

	// Start the server
	log.Println("Starting server on port 80")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func main() {
	startHttpServer()
}
