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

	err := validateCreateUser(ctx, &data)
	if err != nil {
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

	err := validateCreateVhost(ctx, &data)
	if err != nil {
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

	err := validateListQueue(ctx, &data)
	if err != nil {
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

	err := validateListExchange(ctx, &data)
	if err != nil {
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

	err := validateListBind(ctx, &data)
	if err != nil {
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

	err := validateListVhost(ctx, &data)
	if err != nil {
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

	err := validateListNode(ctx, &data)
	if err != nil {
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

	err := validateListPolicy(ctx, &data)
	if err != nil {
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

	err := validateListUser(ctx, &data)
	if err != nil {
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

	err := validateDeletePolicy(ctx, &data)
	if err != nil {
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

	err := validateDeleteUser(ctx, &data)
	if err != nil {
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

	err := validateDeleteVhost(ctx, &data)
	if err != nil {
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

	err := validateUpdateUser(ctx, &data)
	if err != nil {
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

	err := validateUpdateVhost(ctx, &data)
	if err != nil {
		return err
	}

	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	return updateVhost(conn, &data)
}
