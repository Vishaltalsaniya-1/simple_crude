package main

import (
	"log"
	"os"
	"os/signal"
	cnf "simple_crude/config"
	"simple_crude/consumer"
	"simple_crude/controller"
	"simple_crude/db"
	"simple_crude/manager"
	"simple_crude/producer"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartConsumerService() error {
	consumerInstance := consumer.NewConsumer()
	if err := consumerInstance.Initialize(); err != nil {
		log.Fatalf("Failed to start consumer service: %v", err)
		return err
	}
	log.Println("Consumer service started successfully")
	return nil
}

func StartProducerService() (*producer.RMP, error) {
	producerInstance := producer.NewProducer().(*producer.RMP)
	if err := producerInstance.Initialize(); err != nil {
		log.Fatalf("Failed to start producer service: %v", err)
		return nil, err
	}
	log.Println("Producer service started successfully")
	return producerInstance, nil
}

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("Database connection failed: %v", err)

	}
	cnf.LoadConsumer()

	go StartConsumerService()
	go StartProducerService()

	log.Println("Application Started Successfully")

	userManager := manager.NewUserManager()

	userController := controller.NewUserController(userManager)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/students", userController.CreateUser)
	e.GET("/students", userController.Getall)
	e.DELETE("/students/:id", userController.DeleteUser)
	e.PUT("/students/:id", userController.UpdateUser)

	go func() {
		log.Println("Server started on port 8082")
		if err := e.Start(":8082"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	log.Println(" Shutting down application gracefully...")
}

// rmqProducer := producer.NewProducer()
// producerService := producer.NewProducerService(rmqProducer)

// if err := producerService.Initialize(); err != nil {
// 	log.Fatalf("Failed to initialize producer service: %v", err)
// }

// consumerService := consumer.NewConsumer()
// go func() {
// 	if err := consumerService.StartListening("UserCreatedTask"); err != nil {
// 		log.Fatalf("Consumer failed to start: %v", err)
// 	}
// }()
