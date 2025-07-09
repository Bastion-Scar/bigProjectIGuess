package models

import (
	"gorm.io/gorm"
	"time"
)

var LogChan = make(chan OkLogs, 100)

func SendLogs(db *gorm.DB) {
	ticker := time.NewTicker(3 * time.Second)
	batch := make([]OkLogs, 0, 10)

	for {
		select {
		case logs := <-LogChan:
			batch = append(batch, logs)
			if len(batch) >= 10 {
				flushBatch(db, batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flushBatch(db, batch)
				batch = batch[:0]
			}
		}
	}
}

func flushBatch(db *gorm.DB, batch []OkLogs) {
	db.Create(&batch)
}
