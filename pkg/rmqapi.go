package pkg

import (
	"fmt"
	"net/url"
	"time"

	rh "github.com/michaelklishin/rabbit-hole"
	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

// ---- Job ----

func createQueueHaJob(data *createQueueType) error {
	conn, err := connectAPI(&data.amqpConnectionType)
	if err != nil {
		return err
	}

	err = createQueueHA(conn, data)
	if err != nil {
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

// ---- End Job ----

func connectAPI(amqpconn *amqpConnectionType) (*rh.Client, error) {
	apiURL := url.URL{Scheme: "http",
		Host: fmt.Sprintf("%s:%d",
			amqpconn.Host,
			amqpconn.APIPort)}

	logger.Debug("api connecti URL",
		zap.String("api", apiURL.String()))

	conn, err := rh.NewClient(apiURL.String(),
		amqpconn.Username,
		amqpconn.Password)
	if err != nil {
		logger.Debug("Opening API Connection failed.",
			zap.String("error", err.Error()))

		return nil, cli.NewExitError(err.Error(), 1)
	}

	// http client connection timeout in 3 seconds
	conn.SetTimeout(3 * time.Second)

	return conn, nil
}

func createQueueHA(conn *rh.Client, data *createQueueType) error {

	policyDefFunc := func() rh.PolicyDefinition {
		if data.HaMode == "all" {
			return rh.PolicyDefinition{
				"ha-mode":      data.HaMode,
				"ha-sync-mode": data.HaSyncMode,
			}
		}

		return rh.PolicyDefinition{
			"ha-mode":      data.HaMode,
			"ha-params":    data.HaParam,
			"ha-sync-mode": data.HaSyncMode,
		}
	}

	policy := rh.Policy{
		Vhost:      data.Vhost,
		Pattern:    data.QueueName,
		ApplyTo:    "queues", // queues, exchanges, all
		Definition: policyDefFunc()}

	res, err := conn.PutPolicy(
		data.Vhost,
		fmt.Sprintf("%s_%s", data.QueueName, "HA"),
		policy)

	logger.Debug("HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("queue created; however, set HA policy failed", 1)
	}

	if err != nil {
		logger.Debug("set queue HA policy failed",
			zap.String("queue", data.QueueName),
			zap.String("HA mode", data.HaMode),
			zap.String("HA param", data.HaParam),
			zap.String("HA sync mode", data.HaSyncMode),
			zap.String("http response", fmt.Sprintf("%v", res)),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createUser(conn *rh.Client, data *createUserType) error {
	// check user existence.
	userInfo, err := conn.GetUser(data.RmqUsername)
	if err == nil {
		logger.Debug("user exists",
			zap.String("user", userInfo.Name),
			zap.String("tag", userInfo.Tags))

		return cli.NewExitError("user already exists", 1)
	}

	setting := rh.UserSettings{
		Name:     data.RmqUsername,
		Tags:     data.Tag,
		Password: data.RmqPassword}

	res, err := conn.PutUser(data.RmqUsername, setting)

	logger.Debug("HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("create user failed", 1)
	}

	if err != nil {
		logger.Debug("create user failed",
			zap.String("user", data.RmqUsername),
			zap.String("tag", data.Tag),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func createVhost(conn *rh.Client, data *createVhostType) error {
	// check vhost existence.
	vhostInfo, err := conn.GetVhost(data.VhostName)
	if err == nil {
		logger.Debug("vhost exists",
			zap.String("vhost", vhostInfo.Name),
			zap.Bool("tracing", vhostInfo.Tracing))

		return cli.NewExitError("vhost already exists", 1)
	}

	setting := rh.VhostSettings{
		Tracing: data.Tracing}

	res, err := conn.PutVhost(data.VhostName, setting)

	logger.Debug("HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("create vhost failed", 1)
	}

	if err != nil {
		logger.Debug("create vhost failed",
			zap.String("vhost", data.VhostName),
			zap.Bool("tracing", data.Tracing),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func listQueue(conn *rh.Client, data *listQueueType) error {
	if data.QueueName != "" {
		queueInfo, err := conn.GetQueue(data.Vhost, data.QueueName)
		if err != nil {
			logger.Debug("no such queue",
				zap.String("vhost", data.Vhost),
				zap.String("queue", data.QueueName))

			return cli.NewExitError(err.Error(), 1)
		}

		return detailedQueueInfoF(data.Formatter, queueInfo)
	}

	queueInfos, err := conn.ListQueuesIn(data.Vhost)
	if err != nil {
		logger.Debug("list queues failed",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return detailedQueueInfoSliceF(data.Formatter, queueInfos)
}

func listExchange(conn *rh.Client, data *listExchangeType) error {
	if data.ExchangeName != "" {
		exgInfo, err := conn.GetExchange(data.Vhost, data.ExchangeName)
		if err != nil {
			logger.Debug("no such exchange",
				zap.String("vhost", data.Vhost),
				zap.String("exchange", data.ExchangeName))

			return cli.NewExitError(err.Error(), 1)
		}

		return detailedExchangeInfoF(data.Formatter, exgInfo)
	}

	exgInfos, err := conn.ListExchangesIn(data.Vhost)
	if err != nil {
		logger.Debug("list exchange failed",
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return detailedExchangeInfoSliceF(data.Formatter, exgInfos)
}

func listBind(conn *rh.Client, data *listBindType) error {
	var binfos []rh.BindingInfo
	var err error

	if data.All {
		binfos, err = conn.ListBindings()
		if err != nil {
			logger.Debug("list bindings from all vhosts failed")

			return cli.NewExitError(err.Error(), 1)
		}

	} else {
		binfos, err = conn.ListBindingsIn(data.Vhost)
		if err != nil {
			logger.Debug("list bindings from vhost failed",
				zap.String("vhost", data.Vhost),
				zap.String("error", err.Error()))

			return cli.NewExitError(err.Error(), 1)
		}
	}

	return bindingInfoF(data.Formatter, binfos)
}

func listVhost(conn *rh.Client, data *listVhostType) error {
	if data.VhostName != "" {
		vhostInfo, err := conn.GetVhost(data.VhostName)
		if err != nil {
			logger.Debug("no such vhost",
				zap.String("vhost", data.VhostName))

			return cli.NewExitError(err.Error(), 1)
		}

		return vhostInfoF(data.Formatter, vhostInfo)
	}

	vhostInfos, err := conn.ListVhosts()
	if err != nil {
		logger.Debug("list vhosts failed")

		return cli.NewExitError(err.Error(), 1)
	}

	return vhostInfoSliceF(data.Formatter, vhostInfos)
}

func listNode(conn *rh.Client, data *listNodeType) error {
	if data.NodeName != "" {
		nodeInfo, err := conn.GetNode(data.NodeName)
		if err != nil {
			logger.Debug("no such node",
				zap.String("node", data.NodeName))

			return cli.NewExitError(err.Error(), 1)
		}

		return nodeInfoF(data.Formatter, nodeInfo)
	}

	nodeInfos, err := conn.ListNodes()
	if err != nil {
		logger.Debug("list nodes failed")

		return cli.NewExitError(err.Error(), 1)
	}

	return nodeInfoSliceF(data.Formatter, nodeInfos)
}

func listPolicy(conn *rh.Client, data *listPolicyType) error {
	if data.PolicyName != "" {
		policyInfo, err := conn.GetPolicy(data.Vhost, data.PolicyName)
		if err != nil {
			logger.Debug("no such policy",
				zap.String("policy", data.PolicyName))

			return cli.NewExitError(err.Error(), 1)
		}

		return policyInfoF(data.Formatter, policyInfo)
	}

	if data.All {
		policyInfos, err := conn.ListPolicies()
		if err != nil {
			logger.Debug("list all vhosts' policy failed",
				zap.String("err", err.Error()))

			return cli.NewExitError(err.Error(), 1)
		}

		return policyInfoSliceF(data.Formatter, policyInfos)
	}

	policyInfos, err := conn.ListPoliciesIn(data.Vhost)
	if err != nil {
		logger.Debug("list vhost's policy failed",
			zap.String("vhost", data.Vhost),
			zap.String("err", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return policyInfoSliceF(data.Formatter, policyInfos)
}

func listUser(conn *rh.Client, data *listUserType) error {
	if data.RmqUsername != "" {
		userInfo, err := conn.GetUser(data.RmqUsername)
		if err != nil {
			logger.Debug("no such user",
				zap.String("user", data.RmqUsername))

			return cli.NewExitError(err.Error(), 1)
		}

		return userInfoF(data.Formatter, userInfo)
	}

	userInfos, err := conn.ListUsers()
	if err != nil {
		logger.Debug("list users failed",
			zap.String("err", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return userInfoSliceF(data.Formatter, userInfos)
}

func deletePolicy(conn *rh.Client, data *deletePolicyType) error {
	res, err := conn.DeletePolicy(data.Vhost, data.PolicyName)
	if err != nil {
		logger.Debug("delete policy failed",
			zap.String("policy", data.PolicyName),
			zap.String("http response", fmt.Sprintf("%v", res)),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteUser(conn *rh.Client, data *deleteUserType) error {
	res, err := conn.DeleteUser(data.RmqUsername)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("delete user failed", 1)
	}

	if err != nil {
		logger.Debug("delete user failed",
			zap.String("user", data.RmqUsername),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func deleteVhost(conn *rh.Client, data *deleteVhostType) error {
	res, err := conn.DeleteVhost(data.VhostName)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("delete vhost failed", 1)
	}

	if err != nil {
		logger.Debug("delete vhost failed",
			zap.String("vhost", data.VhostName),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func updateUser(conn *rh.Client, data *updateUserType) error {
	// check user existence.
	_, err := conn.GetUser(data.RmqUsername)
	if err != nil {
		logger.Debug("user doesn't exist",
			zap.String("user", data.RmqUsername),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	setting := rh.UserSettings{
		Name:     data.RmqUsername,
		Tags:     data.Tag,
		Password: data.RmqPassword}

	res, err := conn.PutUser(data.RmqUsername, setting)

	logger.Debug("HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("update user failed", 1)
	}

	if err != nil {
		logger.Debug("update user failed",
			zap.String("user", data.RmqUsername),
			zap.String("tag", data.Tag),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}

func updateVhost(conn *rh.Client, data *updateVhostType) error {
	// check vhost existence.
	_, err := conn.GetVhost(data.VhostName)
	if err != nil {
		logger.Debug("vhost doesn't exist",
			zap.String("vhost", data.VhostName),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	setting := rh.VhostSettings{
		Tracing: data.Tracing}

	res, err := conn.PutVhost(data.VhostName, setting)

	logger.Debug("HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError("create vhost failed", 1)
	}

	if err != nil {
		logger.Debug("create vhost failed",
			zap.String("vhost", data.VhostName),
			zap.Bool("tracing", data.Tracing),
			zap.String("error", err.Error()))

		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
