package pkg

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	cli "gopkg.in/urfave/cli.v1"
)

func handleHTTPResponse(res *http.Response, action string, obj string) error {

	logger.Debug(
		"HTTP response",
		zap.String("response", fmt.Sprintf("%v", res)),
	)

	// 204 No Content, object existed.
	// 404 Not Found, no such resource.
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return cli.NewExitError(fmt.Sprintf("%s %s failed", action, obj), 1)
	}

	return nil
}
