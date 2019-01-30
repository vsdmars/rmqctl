package pkg

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

func connect(amqpconn *amqpConnectionType) (*amqp.Connection, error) {
	amqpURL := amqp.URI{Scheme: "amqp",
		Host:     amqpconn.Host,
		Username: amqpconn.Username,
		Password: "XXXXX",
		Port:     amqpconn.Port,
		Vhost:    amqpconn.Vhost,
	}

	logger.Debug(
		"amqp URL",
		zap.String("service", "amqp"),
		zap.String("URL", amqpURL.String()),
	)

	amqpURL.Password = amqpconn.Password

	logger.Debug(
		"timeout",
		zap.String("service", "amqp"),
		zap.String("timeout", (3*time.Second).String()),
	)

	// tcp connection timeout in 3 seconds.
	connection, err := amqp.DialConfig(
		amqpURL.String(),
		amqp.Config{
			Vhost: amqpconn.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
			Heartbeat: 10 * time.Second,
			Locale:    "en_US"},
	)
	if err != nil {
		logger.Debug(
			"amqp create connection failed",
			zap.String("service", "amqp"),
		)

		return nil, cli.NewExitError(err.Error(), 1)
	}

	return connection, nil
}

func burstPublish(
	gcount int,
	channel *amqp.Channel,
	consume <-chan struct{},
	payload amqp.Publishing,
	data *publishType) chan string {

	status := make(chan string)

	go func() {
		defer func() {
			close(status)
		}()

		for range consume {
			channel.Publish(
				data.ExchangeName,
				data.Key, // routing key
				data.Mandatory,
				data.Immediate,
				payload,
			)
		}
		status <- fmt.Sprintf("publish goroutine number: %v done.", gcount)
	}()

	return status
}

func publishMsg(conn *amqp.Connection, data *publishType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug(
			"amqp create channel failed",
			zap.String("service", "amqp"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	gcount := int(float64(runtime.NumCPU()) * float64(0.9))

	logger.Debug(
		"publish using number of goroutines",
		zap.String("service", "amqp"),
		zap.Int("goroutines", gcount),
	)

	consumeChannel := make(chan struct{}, 100)
	var gatherStatus []chan string

	payload := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(fmt.Sprintf("%s", data.Message)),
		DeliveryMode: data.Mode,
		MessageId:    "!42!",
	}

	// starts goroutines
	for i := 0; i <= gcount; i++ {
		gatherStatus = append(
			gatherStatus,
			burstPublish(
				i,
				channel,
				consumeChannel,
				payload,
				data),
		)
	}

	// burst publish
	for i := 0; i < data.Burst; i++ {
		consumeChannel <- struct{}{}
	}

	close(consumeChannel)

	// reference:
	// ContextType: https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Complete_list_of_MIME_types

	for _, s := range gatherStatus {
		for m := range s {
			logger.Debug(
				"publish goroutine done",
				zap.String("message", m),
			)
		}
	}

	logger.Debug(
		"publish done",
		zap.String("service", "amqp"),
	)

	return nil
}

// consumeMsg consumes message from the cluster.
// it uses it's own formatter instead of generic ones due
// to looping.
func consumeMsg(conn *amqp.Connection, data *consumeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug(
			"amqp create channel failed",
			zap.String("service", "amqp"),
		)

		return cli.NewExitError(err.Error(), 1)
	}

	if data.Daemon {
		return daemonConsumeF(channel, data)
	}

	return noneDaemonConsumeF(channel, data)
}

func ackFunction(d *amqp.Delivery, data *consumeType) {
	// Only acknowledge if no autoack.
	if !data.AutoAck {
		switch data.AckType {
		case "ack":
			d.Ack(true)
		case "nack":
			d.Nack(true, false)
		case "reject":
			d.Reject(false)
		}
	}
}
