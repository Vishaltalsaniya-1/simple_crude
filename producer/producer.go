	package producer

	import (
		"context"
		"errors"
		cnf "simple_crude/config"

		"github.com/sirupsen/logrus"
		"github.com/surendratiwari3/paota/config"
		"github.com/surendratiwari3/paota/schema"
		"github.com/surendratiwari3/paota/workerpool"
	)

	type RMP struct {
		WorkerPool *workerpool.Pool
	}

	func NewProducer() ProducerInterface {
		return &RMP{}
	}

	func (rmp *RMP) Initialize() error {
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
			logrus.Fatalf("WorkerPool creation failed: %v", err)
			return err
		}

		rmp.WorkerPool = &workerPool

		if rmp.WorkerPool == nil {
			logrus.Fatal("WorkerPool is nil after initialization")
			return errors.New("failed to initialize worker pool")
		}

		logrus.Info("WorkerPool initialized successfully")
		return nil
	}

	func (rmp *RMP) Publish(Data []byte, ) error {
		consumerConfig := cnf.Consumerconfig

		if rmp.WorkerPool == nil {
			return errors.New("worker pool is not initialized")
		}

		task := &schema.Signature{
			Name: consumerConfig.QueueTaskName,
			Args: []schema.Arg{
				{
					Type:  "string",
					Value: string(Data),
				},
			},
			RetryCount:                  10,
			IgnoreWhenTaskNotRegistered: true,
		}

		logrus.Infof("Created task: %+v", task)

		state, err := (*rmp.WorkerPool).SendTaskWithContext(context.Background(), task)
		if err != nil {
			logrus.Error("Failed to publish task:", err)
			return err
		}

		logrus.Infof("Task State: %+v", state)
		logrus.Info("Task published successfully.")
		return nil
	}
