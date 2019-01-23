package pkg

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	rh "github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

const padding = 1

// Darn... Miss C++ generics...

func jsonFormat(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logger.Debug("format output failed")

		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println(string(b))
	return nil
}

func detailedQueueInfoF(format string, v *rh.DetailedQueueInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Durable",
		"AutoDelete",
		"MasterNode",
		"Status",
		"Consumers",
		"Policy",
		"Messages")
	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		v.Name,
		v.Vhost,
		v.Durable,
		v.AutoDelete,
		v.Node,
		v.Status,
		v.Consumers,
		v.Policy,
		v.Messages)

	w.Flush()
	return nil
}

func detailedQueueInfoSliceF(format string, v []rh.QueueInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Durable",
		"AutoDelete",
		"MasterNode",
		"Status",
		"Consumers",
		"Policy",
		"Messages")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
			sv.Name,
			sv.Vhost,
			sv.Durable,
			sv.AutoDelete,
			sv.Node,
			sv.Status,
			sv.Consumers,
			sv.Policy,
			sv.Messages)
	}

	w.Flush()
	return nil
}

func detailedExchangeInfoF(format string, v *rh.DetailedExchangeInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Type",
		"Durable",
		"AutoDelete")
	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
		v.Name,
		v.Vhost,
		v.Type,
		v.Durable,
		v.AutoDelete)

	w.Flush()
	return nil
}

func detailedExchangeInfoSliceF(format string, v []rh.ExchangeInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Type",
		"Durable",
		"AutoDelete")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
			sv.Name,
			sv.Vhost,
			sv.Type,
			sv.Durable,
			sv.AutoDelete)
	}

	w.Flush()
	return nil
}

func bindingInfoF(format string, v []rh.BindingInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Source",
		"Destination",
		"Vhost",
		"Key",
		"DestinationType")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\n",
			sv.Source,
			sv.Destination,
			sv.Vhost,
			sv.RoutingKey,
			sv.DestinationType)
	}

	w.Flush()
	return nil
}

func vhostInfoF(format string, v *rh.VhostInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		"Name",
		"Tracing",
		"Messages")
	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		v.Name,
		v.Tracing,
		v.Messages)

	w.Flush()
	return nil
}

func vhostInfoSliceF(format string, v []rh.VhostInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		"Name",
		"Tracing",
		"Messages")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
			sv.Name,
			sv.Tracing,
			sv.Messages)
	}

	w.Flush()
	return nil
}

func nodeInfoF(format string, v *rh.NodeInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		"Name",
		"NodeType",
		"IsRunning")
	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		v.Name,
		v.NodeType,
		v.IsRunning)

	w.Flush()
	return nil
}

func nodeInfoSliceF(format string, v []rh.NodeInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
	fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
		"Name",
		"NodeType",
		"IsRunning")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\n",
			sv.Name,
			sv.NodeType,
			sv.IsRunning)
	}

	w.Flush()
	return nil
}

func policyInfoF(format string, v *rh.Policy) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Pattern",
		"Priority",
		"ApplyTo",
		"Definition")
	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		v.Name,
		v.Vhost,
		v.Pattern,
		v.Priority,
		v.ApplyTo,
		v.Definition)

	w.Flush()
	return nil
}

func policyInfoSliceF(format string, v []rh.Policy) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
	fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
		"Name",
		"Vhost",
		"Pattern",
		"Priority",
		"ApplyTo",
		"Definition")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\t|%v\t|%v\t|%v\t|%v\n",
			sv.Name,
			sv.Vhost,
			sv.Pattern,
			sv.Priority,
			sv.ApplyTo,
			sv.Definition)
	}

	w.Flush()
	return nil
}

func userInfoF(format string, v *rh.UserInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\n",
		"Name",
		"Tag")
	fmt.Fprintf(w, "|%v\t|%v\n",
		v.Name,
		v.Tags)

	w.Flush()
	return nil
}

func userInfoSliceF(format string, v []rh.UserInfo) error {

	if format == "json" {
		return jsonFormat(v)
	}

	w := tabwriter.NewWriter(
		os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)

	fmt.Fprintf(w, "|%v\t|%v\n",
		"Name",
		"Tag")

	for _, sv := range v {
		fmt.Fprintf(w, "|%v\t|%v\n",
			sv.Name,
			sv.Tags)
	}

	w.Flush()
	return nil
}

func daemonConsumeF(channel *amqp.Channel, data *consumeType) error {

	F := func() error {
		var w *tabwriter.Writer

		if data.Formatter == "plain" {
			w = tabwriter.NewWriter(
				os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
			defer w.Flush()

			fmt.Fprintf(w, "|%v\n", "Message")
		}

		deliveries, err := channel.Consume(
			data.QueueName,
			"", // ConsumerTag is unnecessary in rmqctl's usecase.
			data.AutoAck,
			false, // exclusive is unnecessary in rmqctl's usecase.
			false, // nolocal is not supported by rabbitmq.
			data.NoWait,
			data.Args,
		)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		for d := range deliveries {
			ackFunction(&d, data)

			if data.Formatter == "plain" {
				fmt.Fprintf(w, "|%v\n", string(d.Body))
			} else {
				b, err := json.MarshalIndent(d, "", "  ")
				if err != nil {
					logger.Debug("json formatting failed",
						zap.String("delivery", string(d.Body)))

					logger.Error(
						"json formatting failed, fallback to plain format")

					fmt.Println(string(d.Body))
					continue
				}

				fmt.Println(string(b))
			}
		}

		return nil
	}

	return F()
}

func noneDaemonConsumeF(channel *amqp.Channel, data *consumeType) error {

	F := func() error {
		var w *tabwriter.Writer

		if data.Formatter == "plain" {
			w = tabwriter.NewWriter(
				os.Stdout, 0, 0, padding, ' ', tabwriter.StripEscape)
			defer w.Flush()

			fmt.Fprintf(w, "|%v\n", "Message")
		}

		for {
			d, ok, err := channel.Get(
				data.QueueName,
				data.AutoAck,
			)

			if ok {
				ackFunction(&d, data)

				if data.Formatter == "plain" {
					fmt.Fprintf(w, "|%v\n", string(d.Body))
				} else {
					b, err := json.MarshalIndent(d, "", "  ")
					if err != nil {
						logger.Debug("json formatting failed",
							zap.String("delivery", string(d.Body)))

						logger.Error(
							"json formatting failed, fallback to plain format")

						fmt.Println(string(d.Body))
						continue
					}

					fmt.Println(string(b))
				}
			} else {
				if err != nil {
					logger.Debug("consume message failed",
						zap.String("error", err.Error()))

					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			}
		}
	}

	return F()
}
