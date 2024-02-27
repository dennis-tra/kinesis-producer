package producer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func ExampleSimple() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := kinesis.NewFromConfig(cfg)
	pr := New(&Config{
		StreamName:   "test",
		BacklogCount: 2000,
		Client:       client,
	})

	pr.Start(context.TODO())

	// Handle failures
	go func() {
		for r := range pr.NotifyFailures() {
			// r contains `Data`, `PartitionKey` and `Error()`
			slog.Error("detected put failure", r.error)
		}
	}()

	go func() {
		for i := 0; i < 5000; i++ {
			err := pr.Put([]byte("foo"), "bar")
			if err != nil {
				slog.Error("error producing", err)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	pr.Stop()
}
