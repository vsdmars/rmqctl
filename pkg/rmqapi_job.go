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

	return createQueue(conn, &data)
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

	return createExchange(conn, &data)
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

	return createBind(conn, &data)
}

func createUserJob(ctx *cli.Context) error {
	data := createUserType{}

	if err := validateCreateUser(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return createUser(conn, &data)
}

func createVhostJob(ctx *cli.Context) error {
	data := createVhostType{}

	if err := validateCreateVhost(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return createVhost(conn, &data)
}

func listQueueJob(ctx *cli.Context) error {
	data := listQueueType{}

	if err := validateListQueue(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listQueue(conn, &data)
}

func listExchangeJob(ctx *cli.Context) error {
	data := listExchangeType{}

	if err := validateListExchange(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listExchange(conn, &data)
}

func listBindJob(ctx *cli.Context) error {
	data := listBindType{}

	if err := validateListBind(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listBind(conn, &data)
}

func listVhostJob(ctx *cli.Context) error {
	data := listVhostType{}

	if err := validateListVhost(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listVhost(conn, &data)
}

func listNodeJob(ctx *cli.Context) error {
	data := listNodeType{}

	if err := validateListNode(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listNode(conn, &data)
}

func listPolicyJob(ctx *cli.Context) error {
	data := listPolicyType{}

	if err := validateListPolicy(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listPolicy(conn, &data)
}

func listUserJob(ctx *cli.Context) error {
	data := listUserType{}

	if err := validateListUser(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return listUser(conn, &data)
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

	return deleteQueue(conn, &data)
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

	return deleteExchange(conn, &data)
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

	return deleteBind(conn, &data)
}

func deletePolicyJob(ctx *cli.Context) error {
	data := deletePolicyType{}

	if err := validateDeletePolicy(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return deletePolicy(conn, &data)
}

func deleteUserJob(ctx *cli.Context) error {
	data := deleteUserType{}

	if err := validateDeleteUser(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return deleteUser(conn, &data)
}

func deleteVhostJob(ctx *cli.Context) error {
	data := deleteVhostType{}

	if err := validateDeleteVhost(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return deleteVhost(conn, &data)
}

func updateUserJob(ctx *cli.Context) error {
	data := updateUserType{}

	if err := validateUpdateUser(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return updateUser(conn, &data)
}

func updateVhostJob(ctx *cli.Context) error {
	data := updateVhostType{}

	if err := validateUpdateVhost(ctx, &data); err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return updateVhost(conn, &data)
}
