package main

import (
	"os"

	"github.com/open-component-model/service-model/plugins/serviceplugin/cmds/services"
	"ocm.software/ocm/api/ocm/plugin/ppi"
	"ocm.software/ocm/api/ocm/plugin/ppi/clicmd"
	"ocm.software/ocm/api/ocm/plugin/ppi/cmds"
	// enable mandelsoft plugin logging configuration.
	"github.com/open-component-model/service-model/api/version"
	_ "ocm.software/ocm/api/ocm/plugin/ppi/logging"
)

func main() {
	p := ppi.NewPlugin("serviceplugin", version.Get().String())

	p.SetShort("Plugin to handle Service Models")
	p.SetLong("The plugin offers the basic evaluation f service model resources.")

	cmd, err := clicmd.NewCLICommand(services.New(), clicmd.WithCLIConfig(), clicmd.WithObjectType("services"), clicmd.WithVerb("get"))
	if err != nil {
		os.Exit(1)
	}
	p.RegisterCommand(cmd)
	p.ForwardLogging()

	err = cmds.NewPluginCommand(p).Execute(os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}
