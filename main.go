package main

import (
	"awesomeProject10/goSquare"
	"awesomeProject10/models"
	"awesomeProject10/zapLogger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"sync"
)

func main() {

	var errChan = make(chan error)
	var wg sync.WaitGroup

	zapLogger.Init()
	logger := zapLogger.Log

	db := models.InitDb()
	go models.SendLogs(db)

	r := gin.New()
	r.Use(zapLogger.CustomLogger())
	r.POST("/calculate", func(c *gin.Context) {
		jobs := make(chan int, 10)
		square := make(chan int, 10)

		wg.Add(1)

		go goSquare.GetSquare(jobs, square, &wg)
		for i := 0; i <= 10; i++ {
			var sq []int
			sq = append(sq, <-square)
			c.JSON(200, gin.H{
				"square": sq,
			})
		}
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Сервер работает")
		if err := r.Run(":8080"); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := <-errChan; err != nil {
			logger.Fatal("Не удалось запустить сервер", zap.Error(err))
		}
	}()
	wg.Wait()
}
