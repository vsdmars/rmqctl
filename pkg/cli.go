package pkg

import (
	"fmt"
	"os"

	"github.com/urfave/cli/altsrc"
	"go.uber.org/zap/zapcore"
	cli "gopkg.in/urfave/cli.v1"
)

func mainBefore(flags []cli.Flag) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		if ctx.Bool("debug") {
			logger.setLevel(zapcore.DebugLevel)
		}

		// check existence of the config yaml file under current
		// directory
		if _, err := os.Stat(configDir); err == nil {
			// load config file
			return altsrc.InitInputSourceWithContext(flags,
				altsrc.NewYamlSourceFromFlagFunc("config"))(ctx)
		}
		return nil
	}
}

func doneAfter(ctx *cli.Context) error {
	fmt.Println("done")
	return nil
}
