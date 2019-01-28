package pkg

import (
	cli "gopkg.in/urfave/cli.v1"
)

func createQueueJob(ctx *cli.Context) error {
	data := createQueueType{}

	if err := validateCreateQueue(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = createQueue(conn, &data); err != nil {
		return err
	}

	return nil
}

func createExchangeJob(ctx *cli.Context) error {
	data := createExchangeType{}

	if err := validateCreateExchange(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = createExchange(conn, &data); err != nil {
		return err
	}

	return nil
}

func createBindJob(ctx *cli.Context) error {
	data := createBindType{}

	if err := validateCreateBind(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = createBind(conn, &data); err != nil {
		return err
	}

	return nil
}

func createUserJob(ctx *cli.Context) error {
	data := createUserType{}

	err := validateCreateUser(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createUser(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func createVhostJob(ctx *cli.Context) error {
	data := createVhostType{}

	err := validateCreateVhost(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createVhost(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listQueueJob(ctx *cli.Context) error {
	data := listQueueType{}

	err := validateListQueue(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listQueue(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listExchangeJob(ctx *cli.Context) error {
	data := listExchangeType{}

	err := validateListExchange(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listExchange(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listBindJob(ctx *cli.Context) error {
	data := listBindType{}

	err := validateListBind(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listBind(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listVhostJob(ctx *cli.Context) error {
	data := listVhostType{}

	err := validateListVhost(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listVhost(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listNodeJob(ctx *cli.Context) error {
	data := listNodeType{}

	err := validateListNode(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listNode(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listPolicyJob(ctx *cli.Context) error {
	data := listPolicyType{}

	err := validateListPolicy(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listPolicy(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func listUserJob(ctx *cli.Context) error {
	data := listUserType{}

	err := validateListUser(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = listUser(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteQueueJob(ctx *cli.Context) error {
	data := deleteQueueType{}

	if err := validateDeleteQueue(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = deleteQueue(conn, &data); err != nil {
		return err
	}

	return nil
}

func deleteExchangeJob(ctx *cli.Context) error {
	data := deleteExchangeType{}

	if err := validateDeleteExchange(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = deleteExchange(conn, &data); err != nil {
		return err
	}

	return nil
}

func deleteBindJob(ctx *cli.Context) error {
	data := deleteBindType{}

	if err := validateDeleteBind(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	if err = deleteBind(conn, &data); err != nil {
		return err
	}

	return nil
}

func deletePolicyJob(ctx *cli.Context) error {
	data := deletePolicyType{}

	err := validateDeletePolicy(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deletePolicy(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteUserJob(ctx *cli.Context) error {
	data := deleteUserType{}

	err := validateDeleteUser(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteUser(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func deleteVhostJob(ctx *cli.Context) error {
	data := deleteVhostType{}

	err := validateDeleteVhost(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = deleteVhost(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func updateUserJob(ctx *cli.Context) error {
	data := updateUserType{}

	err := validateUpdateUser(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = updateUser(conn, &data)
	if err != nil {
		return err
	}

	return nil
}

func updateVhostJob(ctx *cli.Context) error {
	data := updateVhostType{}

	err := validateUpdateVhost(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = updateVhost(conn, &data)
	if err != nil {
		return err
	}

	return nil
}
