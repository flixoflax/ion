package mosaic

import (
	"github.com/pulumi/pulumi/sdk/v3/go/common/apitype"
	"github.com/sst/ion/cmd/sst/cli"
	"github.com/sst/ion/cmd/sst/mosaic/aws"
	"github.com/sst/ion/cmd/sst/mosaic/server"
	"github.com/sst/ion/cmd/sst/mosaic/ui"
	"github.com/sst/ion/pkg/project"
)

func CmdMosaicDeploy(c *cli.Cli) error {
	evts, err := server.Stream(c.Context, "http://localhost:13557",
		project.StackCommandEvent{},
		project.ConcurrentUpdateEvent{},
		project.StackCommandEvent{},
		project.BuildFailedEvent{},
		apitype.ResourcePreEvent{},
		apitype.ResOpFailedEvent{},
		apitype.ResOutputsEvent{},
		apitype.DiagnosticEvent{},
		project.CompleteEvent{},
		aws.FunctionInvokedEvent{},
		aws.FunctionResponseEvent{},
		aws.FunctionErrorEvent{},
		aws.FunctionLogEvent{},
		aws.FunctionBuildEvent{},
	)
	if err != nil {
		return err
	}

	u := ui.New(c.Context, ui.ProgressModeDev)
	u.Header("dev", "app", "foo")
	for {
		select {
		case <-c.Context.Done():
			u.Destroy()
			return nil
		case evt, ok := <-evts:
			if !ok {
				c.Cancel()
				return nil
			}
			u.Event(evt)
		}
	}
}
