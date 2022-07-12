package schedules

import (
	"time"

	"github.com/go-co-op/gocron"
)

var (
	scheduler *gocron.Scheduler
)

func deleteTempImage() {

}

func InitSchedule() error {
	scheduler = gocron.NewScheduler(time.Local)
	return nil
}
