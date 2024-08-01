package progress_go

import (
	"log"
	"time"

	"github.com/vardius/progress-go"
)

func ProgressGo() {
	bar := progress.New(0, 100)

	_, _ = bar.Start()
	defer func() {
		if _, err := bar.Stop(); err != nil {
			log.Printf("failed to finish progress: %v", err)
		}
	}()

	for i := 0; i < 100; i++ {
		_, _ = bar.Advance(1)
		time.Sleep(10 * time.Millisecond)
	}
}
