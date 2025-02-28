package cmd

import (
	ext "github.com/linuxsuren/api-testing/pkg/extension"
	"github.com/linuxsuren/atest-ext-store-redis/pkg"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the root Command
func NewRootCommand() (c *cobra.Command) {
	opt := &options{
		Extension: ext.NewExtension("redis", "store", 7073),
	}
	c = &cobra.Command{
		Use:   opt.GetFullName(),
		Short: "A store extension for redis",
		RunE:  opt.runE,
	}
	opt.AddFlags(c.Flags())
	return
}

type options struct {
	*ext.Extension
}

func (o *options) runE(c *cobra.Command, _ []string) (err error) {
	remoteServer := pkg.NewRemoteServer()
	err = ext.CreateRunner(o.Extension, c, remoteServer)
	return
}
