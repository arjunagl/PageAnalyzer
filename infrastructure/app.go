package infrastructure

import (
	"arjunagl/htmlAnalyzer/domain/port"
	"arjunagl/htmlAnalyzer/domain/service"
)

type Application struct {
	ContentDownloader port.ContentDownloader
	SchedulerService  service.Scheduler
	ContentReader     port.ContentReader
}

func NewApplication(cd port.ContentDownloader, sc *service.SchedulerService, cr port.ContentReader) *Application {
	return &Application{ContentDownloader: cd, SchedulerService: sc, ContentReader: cr}
}
