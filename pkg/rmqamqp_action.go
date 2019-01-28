package pkg

import (
	"net"
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

	logger.Debug("amqp connect URL", zap.String("amqp", amqpURL.String()))

	amqpURL.Password = amqpconn.Password

	// tcp connection timeout in 3 seconds.
	connection, err := amqp.DialConfig(amqpURL.String(),
		amqp.Config{
			Vhost: amqpconn.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
			Heartbeat: 10 * time.Second,
			Locale:    "en_US"})
	if err != nil {
		logger.Debug("Opening amqp connection failed.",
			zap.String("error", err.Error()))

		return nil, cli.NewExitError(err.Error(), 1)
	}

	return connection, nil
}

func publishMsg(conn *amqp.Connection, data *publishType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	if err = channel.Publish(
		data.ExchangeName,
		data.Key, // routing key
		data.Mandatory,
		data.Immediate,
		data.Message,
	); err != nil {
		logger.Debug("Publish failed, channel closed.",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

// consumeMsg consumes message from the cluster.
// it uses it's own formatter instead of generic ones due
// to looping.
func consumeMsg(conn *amqp.Connection, data *consumeType) error {
	channel, err := conn.Channel()
	if err != nil {
		logger.Debug("Opening channel failed.",
			zap.String("error", err.Error()))

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
