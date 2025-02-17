package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	cnf "simple_crude/config"

	"github.com/sirupsen/logrus"
	"github.com/surendratiwari3/paota/config"
	"github.com/surendratiwari3/paota/schema"
	"github.com/surendratiwari3/paota/workerpool"
)

type Consumer struct {
	WorkerPool *workerpool.Pool
}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Initialize() error {
	consumerConfig := cnf.Consumerconfig
	cnf := config.Config{
		Broker:        "amqp",
		TaskQueueName: consumerConfig.QueueTaskName,
		AMQP: &config.AMQPConfig{
			Url:                consumerConfig.Url,
			Exchange:           consumerConfig.Exchange,
			ExchangeType:       "direct",
			BindingKey:         consumerConfig.BindingKeyName,
			PrefetchCount:      consumerConfig.PrefetchCount,
			ConnectionPoolSize: consumerConfig.ConnectionPoolSize,
			DelayedQueue:       consumerConfig.DelayedQueueName,
		},
	}

	workerPool, err := workerpool.NewWorkerPoolWithConfig(context.Background(), 10, "testcourses", cnf)
	if err != nil {
		logrus.Errorf("WorkerPool creation failed: %v", err)
		return err
	}

	c.WorkerPool = &workerPool
	if c.WorkerPool == nil {
		logrus.Error("WorkerPool is nil after initialization")
		return errors.New("failed to initialize worker pool")
	}

	regTasks := map[string]interface{}{
		consumerConfig.QueueTaskName: c.ProcessTask,
	}

	for key, value := range regTasks {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
	// fmt.Printf("  Rregtaske %+v",regTasks)
	log.Printf("Consumer RegisterTasks.......")

	if err := workerPool.RegisterTasks(regTasks); err != nil {
		logrus.Errorf("Error registering tasks: %v", err)
		return err
	}
	log.Printf("Consumer Registering.......\n")

	if err := workerPool.Start(); err != nil {
		logrus.Errorf("Error starting worker: %v", err)
		return err
	}

	logrus.Info("WorkerPool consumer initialized and started successfully")
	return nil
}

func (c *Consumer) ProcessTask(arg *schema.Signature) error {
	log.Printf("Consumer Processtask.......\n")

	if len(arg.Args) == 0 {
		logrus.Warn("No arguments received in the task")
		return nil
	}

	for _, argItem := range arg.Args {
		logrus.Infof("Processing task - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

		argStr, ok := argItem.Value.(string)
		if !ok {
			logrus.Error("Received argument is not a string")
			return errors.New("invalid argument type")
		}

		var userData map[string]interface{}
		if err := json.Unmarshal([]byte(argStr), &userData); err != nil {
			logrus.Errorf("Error decoding task data: %v", err)
			return err
		}

		logrus.Infof("Task successfully processed: %+v", userData)
	}
	return nil
}

// func (c *Consumer) Print(taskData string) error {
// 	logrus.Infof("Received Task Data: %s", taskData)

// 	return nil
// }
// func (c *Consumer) Consume(data []byte, taskName string) error {
// 	if c.WorkerPool == nil {
// 		logrus.Error("Worker pool is not initialized")
// 		return errors.New("worker pool is not initialized")
// 	}

// 	task := &schema.Signature{
// 		Name: taskName,
// 		Args: []schema.Arg{
// 			{
// 				Name:  taskName,
// 				Type:  "string",
// 				Value: string(data),
// 			},
// 		},
// 		RetryCount:                  10,
// 		RoutingKey:                  cnf.Consumerconfig.BindingKeyName,
// 		IgnoreWhenTaskNotRegistered: true,
// 	}

// 	state, err := (*c.WorkerPool).SendTaskWithContext(context.Background(), task)
// 	if err != nil {
// 		logrus.Errorf("Error consuming task: %v", err)
// 		return err
// 	}

// 	logrus.Infof("Task sent successfully. State: %+v", state)
// 	return nil
// }

// func (c *Consumer) Print(arg *schema.Signature) error {
// 	mongoCollection := db.GetDB().Database("Courses").Collection("students")

// 	if len(arg.Args) == 0 {
// 		logrus.Warn("No arguments received in the task")
// 		return nil
// 	}

// 	for _, argItem := range arg.Args {
// 		logrus.Infof("Processing task - Arg Type: %s, Arg Value: %v", argItem.Type, argItem.Value)

// 		argStr, ok := argItem.Value.(string)
// 		if !ok {
// 			logrus.Error("Received argument is not a string")
// 			return errors.New("invalid argument type")
// 		}

// 		var movie models.Std
// 		if err := json.Unmarshal([]byte(argStr), &movie); err != nil {
// 			logrus.Errorf("Error decoding movie data: %v", err)
// 			return err
// 		}
// 		if mongoCollection == nil {
// 			logrus.Error("MongoDB collection is nil")
// 			return errors.New("database connection not established")
// 		}

// 		_, err := mongoCollection.InsertOne(context.Background(), movie)
// 		if err != nil {
// 			logrus.Errorf("Database insertion failed: %v", err)
// 			return err
// 		}

// 		logrus.Info("Movie successfully inserted into the database.")
// 	}

// 	return nil
// }
