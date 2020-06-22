package version

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var Version = "1.0.0"

func Run(cli *cli.Context) error {
	fmt.Println(Version)
	return nil
}
