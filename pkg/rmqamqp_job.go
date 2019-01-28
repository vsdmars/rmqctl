package pkg

import (
	cli "gopkg.in/urfave/cli.v1"
)

// miss C++ template...
// waiting for generic in Golang2 :-)

func publishJob(ctx *cli.Context) error {
	data := publishType{}

	if err := validatePublish(ctx, &data); err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return publishMsg(conn, &data)
}

func consumeJob(ctx *cli.Context) error {
	data := consumeType{}

	if err := validateConsume(ctx, &data); err != nil {
		return err
	}

	conn, err := connect(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return consumeMsg(conn, &data)
}
