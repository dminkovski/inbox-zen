package main

import (
	"context"
	"fmt"
	"inbox-zen/pkg/controller"
	"inbox-zen/pkg/jobs"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	backgroundContext := context.Background()

	// setting up the scheduler
	job := jobs.NewSummarizeNewlettersJob()
	err := job.StartScheduler(backgroundContext)
	if err != nil {
		fmt.Println("failed to start scheduler")
		os.Exit(1)
	}

	// setting up the server
	httpPathPrefix := "/inbox-zen"
	port := ":8080"
	router := gin.Default()
	inboxZenRouter := router.Group(httpPathPrefix)

	// Creating the summaryController for the APIs consumed by the frontend
	summaryController := controller.NewSummaryController()
	discountsController := controller.NewDiscountsController()

	// Setting the routes for the APIs consumed by the frontend
	inboxZenRouter.GET(summaryController.GetSummaryRoute, summaryController.GetSummary)
	inboxZenRouter.GET(discountsController.GetDiscountsRoute, summaryController.GetDiscounts)

	router.Run(port)
}
