package cmd

import (
	"fmt"
	"os"
	"os/exec"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config/cmd"
	gostore "github.com/micro/go-micro/v2/store"
	"github.com/micro/micro/v2/plugin"
	"github.com/micro/micro/v2/plugin/build"
	"github.com/micro/micro/v2/server"
	"github.com/micro/micro/v2/service"

	// clients
	"github.com/micro/micro/v2/client/api"
	"github.com/micro/micro/v2/client/bot"
	"github.com/micro/micro/v2/client/cli"
	"github.com/micro/micro/v2/client/cli/new"
	"github.com/micro/micro/v2/client/cli/util"
	"github.com/micro/micro/v2/client/proxy"
	"github.com/micro/micro/v2/client/web"

	// services
	"github.com/micro/micro/v2/service/auth"
	"github.com/micro/micro/v2/service/broker"
	"github.com/micro/micro/v2/service/config"
	"github.com/micro/micro/v2/service/debug"
	"github.com/micro/micro/v2/service/health"
	"github.com/micro/micro/v2/service/network"
	"github.com/micro/micro/v2/service/registry"
	"github.com/micro/micro/v2/service/router"
	"github.com/micro/micro/v2/service/runtime"
	"github.com/micro/micro/v2/service/store"
	"github.com/micro/micro/v2/service/tunnel"

	// internals
	inauth "github.com/micro/micro/v2/internal/auth"
	"github.com/micro/micro/v2/internal/helper"
	"github.com/micro/micro/v2/internal/platform"
	_ "github.com/micro/micro/v2/internal/plugins"
	"github.com/micro/micro/v2/internal/update"
	_ "github.com/micro/micro/v2/internal/usage"
)

var (
	GitCommit string
	GitTag    string
	BuildDate string

	name        = "micro"
	description = "A microservice runtime"
	version     = "latest"
)

func init() {
	// setup the build plugin
	plugin.Register(build.Flags())

	// set platform build date
	platform.Version = BuildDate
}

func setup(app *ccli.App) {
	app.Flags = append(app.Flags,
		&ccli.BoolFlag{
			Name:  "local",
			Usage: "Enable local only development: Defaults to true.",
		},
		&ccli.BoolFlag{
			Name:    "enable_acme",
			Usage:   "Enables ACME support via Let's Encrypt. ACME hosts should also be specified.",
			EnvVars: []string{"MICRO_ENABLE_ACME"},
		},
		&ccli.StringFlag{
			Name:    "acme_hosts",
			Usage:   "Comma separated list of hostnames to manage ACME certs for",
			EnvVars: []string{"MICRO_ACME_HOSTS"},
		},
		&ccli.StringFlag{
			Name:    "acme_provider",
			Usage:   "The provider that will be used to communicate with Let's Encrypt. Valid options: autocert, certmagic",
			EnvVars: []string{"MICRO_ACME_PROVIDER"},
		},
		&ccli.BoolFlag{
			Name:    "enable_tls",
			Usage:   "Enable TLS support. Expects cert and key file to be specified",
			EnvVars: []string{"MICRO_ENABLE_TLS"},
		},
		&ccli.StringFlag{
			Name:    "tls_cert_file",
			Usage:   "Path to the TLS Certificate file",
			EnvVars: []string{"MICRO_TLS_CERT_FILE"},
		},
		&ccli.StringFlag{
			Name:    "tls_key_file",
			Usage:   "Path to the TLS Key file",
			EnvVars: []string{"MICRO_TLS_KEY_FILE"},
		},
		&ccli.StringFlag{
			Name:    "tls_client_ca_file",
			Usage:   "Path to the TLS CA file to verify clients against",
			EnvVars: []string{"MICRO_TLS_CLIENT_CA_FILE"},
		},
		&ccli.StringFlag{
			Name:    "api_address",
			Usage:   "Set the api address e.g 0.0.0.0:8080",
			EnvVars: []string{"MICRO_API_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "namespace",
			Usage:   "Set the micro service namespace",
			EnvVars: []string{"MICRO_NAMESPACE"},
			Value:   "micro",
		},
		&ccli.StringFlag{
			Name:    "proxy_address",
			Usage:   "Proxy requests via the HTTP address specified",
			EnvVars: []string{"MICRO_PROXY_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "web_address",
			Usage:   "Set the web UI address e.g 0.0.0.0:8082",
			EnvVars: []string{"MICRO_WEB_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "network",
			Usage:   "Set the micro network name: local, go.micro",
			EnvVars: []string{"MICRO_NETWORK"},
		},
		&ccli.StringFlag{
			Name:    "network_address",
			Usage:   "Set the micro network address e.g. :9093",
			EnvVars: []string{"MICRO_NETWORK_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "router_address",
			Usage:   "Set the micro router address e.g. :8084",
			EnvVars: []string{"MICRO_ROUTER_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "gateway_address",
			Usage:   "Set the micro default gateway address e.g. :9094",
			EnvVars: []string{"MICRO_GATEWAY_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "tunnel_address",
			Usage:   "Set the micro tunnel address e.g. :8083",
			EnvVars: []string{"MICRO_TUNNEL_ADDRESS"},
		},
		&ccli.StringFlag{
			Name:    "api_handler",
			Usage:   "Specify the request handler to be used for mapping HTTP requests to services; {api, proxy, rpc}",
			EnvVars: []string{"MICRO_API_HANDLER"},
		},
		&ccli.StringFlag{
			Name:    "api_namespace",
			Usage:   "Set the namespace used by the API e.g. com.example.api",
			EnvVars: []string{"MICRO_API_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "web_namespace",
			Usage:   "Set the namespace used by the Web proxy e.g. com.example.web",
			EnvVars: []string{"MICRO_WEB_NAMESPACE"},
		},
		&ccli.StringFlag{
			Name:    "web_url",
			Usage:   "Set the host used for the web dashboard e.g web.example.com",
			EnvVars: []string{"MICRO_WEB_HOST"},
		},
		&ccli.BoolFlag{
			Name:    "enable_stats",
			Usage:   "Enable stats",
			EnvVars: []string{"MICRO_ENABLE_STATS"},
		},
		&ccli.BoolFlag{
			Name:    "auto_update",
			Usage:   "Enable automatic updates",
			EnvVars: []string{"MICRO_AUTO_UPDATE"},
		},
		&ccli.StringFlag{
			Name:    "update_url",
			Usage:   "Set the url to retrieve system updates from",
			EnvVars: []string{"MICRO_UPDATE_URL"},
			Value:   update.DefaultURL,
		},
		&ccli.BoolFlag{
			Name:    "report_usage",
			Usage:   "Report usage statistics",
			EnvVars: []string{"MICRO_REPORT_USAGE"},
			Value:   true,
		},
		&ccli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Override environment",
			EnvVars: []string{"MICRO_ENV"},
		},
	)

	plugins := plugin.Plugins()

	for _, p := range plugins {
		if flags := p.Flags(); len(flags) > 0 {
			app.Flags = append(app.Flags, flags...)
		}

		if cmds := p.Commands(); len(cmds) > 0 {
			app.Commands = append(app.Commands, cmds...)
		}
	}

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {

		if len(ctx.String("api_handler")) > 0 {
			api.Handler = ctx.String("api_handler")
		}
		if len(ctx.String("api_address")) > 0 {
			api.Address = ctx.String("api_address")
		}
		if len(ctx.String("proxy_address")) > 0 {
			proxy.Address = ctx.String("proxy_address")
		}
		if len(ctx.String("web_address")) > 0 {
			web.Address = ctx.String("web_address")
		}
		if len(ctx.String("network_address")) > 0 {
			network.Address = ctx.String("network_address")
		}
		if len(ctx.String("router_address")) > 0 {
			router.Address = ctx.String("router_address")
		}
		if len(ctx.String("tunnel_address")) > 0 {
			tunnel.Address = ctx.String("tunnel_address")
		}
		if len(ctx.String("api_namespace")) > 0 {
			api.Namespace = ctx.String("api_namespace")
		}
		if len(ctx.String("web_namespace")) > 0 {
			web.Namespace = ctx.String("web_namespace")
		}
		if len(ctx.String("web_host")) > 0 {
			web.Host = ctx.String("web_host")
		}

		for _, p := range plugins {
			if err := p.Init(ctx); err != nil {
				return err
			}
		}

		util.SetupCommand(ctx)
		// now do previous before
		if err := before(ctx); err != nil {
			// DO NOT return this error otherwise the action will fail
			// and help will be printed.
			fmt.Println(err)
			os.Exit(1)
			return err
		}

		var opts []gostore.Option

		// the database is not overriden by flag then set it
		if len(ctx.String("store_database")) == 0 {
			opts = append(opts, gostore.Database(cmd.App().Name))
		}

		// if the table is not overriden by flag then set it
		if len(ctx.String("store_table")) == 0 {
			table := cmd.App().Name

			// if an arg is specified use that as the name
			// so each service has its own table preconfigured
			if name := ctx.Args().First(); len(name) > 0 {
				table = name
			}

			opts = append(opts, gostore.Table(table))
		}

		// TODO: move this entire initialisation elsewhere
		// maybe in service.Run so all things are configured
		if len(opts) > 0 {
			(*cmd.DefaultCmd.Options().Store).Init(opts...)
		}

		// add the system rules if we're using the JWT implementation
		// which doesn't have access to the rules in the auth service
		if (*cmd.DefaultCmd.Options().Auth).String() == "jwt" {
			for _, rule := range inauth.SystemRules {
				if err := (*cmd.DefaultCmd.Options().Auth).Grant(rule); err != nil {
					return err
				}
			}
		}

		return nil
	}
}

func buildVersion() string {
	microVersion := version

	if GitTag != "" {
		microVersion = GitTag
	}

	if GitCommit != "" {
		microVersion += fmt.Sprintf("-%s", GitCommit)
	}

	if BuildDate != "" {
		microVersion += fmt.Sprintf("-%s", BuildDate)
	}

	return microVersion
}

// Init initialised the command line
func Init(options ...micro.Option) {
	Setup(cmd.App(), options...)

	cmd.Init(
		cmd.Name(name),
		cmd.Description(description),
		cmd.Version(buildVersion()),
	)
}

// Setup sets up a cli.App
func Setup(app *ccli.App, options ...micro.Option) {
	// Add the various commands
	app.Commands = append(app.Commands, api.Commands(options...)...)
	app.Commands = append(app.Commands, auth.Commands()...)
	app.Commands = append(app.Commands, bot.Commands()...)
	app.Commands = append(app.Commands, cli.Commands()...)
	app.Commands = append(app.Commands, broker.Commands(options...)...)
	app.Commands = append(app.Commands, health.Commands(options...)...)
	app.Commands = append(app.Commands, proxy.Commands(options...)...)
	app.Commands = append(app.Commands, router.Commands(options...)...)
	app.Commands = append(app.Commands, tunnel.Commands(options...)...)
	app.Commands = append(app.Commands, network.Commands(options...)...)
	app.Commands = append(app.Commands, registry.Commands(options...)...)
	app.Commands = append(app.Commands, runtime.Commands(options...)...)
	app.Commands = append(app.Commands, debug.Commands(options...)...)
	app.Commands = append(app.Commands, server.Commands(options...)...)
	app.Commands = append(app.Commands, service.Commands(options...)...)
	app.Commands = append(app.Commands, store.Commands(options...)...)
	app.Commands = append(app.Commands, new.Commands()...)
	app.Commands = append(app.Commands, build.Commands()...)
	app.Commands = append(app.Commands, web.Commands(options...)...)
	app.Commands = append(app.Commands, config.Commands(options...)...)

	// add the init command for our internal operator
	app.Commands = append(app.Commands, &ccli.Command{
		Name:  "init",
		Usage: "Run the micro operator",
		Action: func(c *ccli.Context) error {
			platform.Init(c)
			return nil
		},
		Flags: []ccli.Flag{},
	})

	// boot micro runtime
	app.Action = func(c *ccli.Context) error {
		if c.Args().Len() > 0 {
			command := c.Args().First()

			v, err := exec.LookPath(command)
			if err != nil {
				fmt.Println(helper.UnexpectedCommand(c))
				os.Exit(1)
			}

			// execute the command
			ce := exec.Command(v, c.Args().Slice()[1:]...)
			ce.Stdout = os.Stdout
			ce.Stderr = os.Stderr
			return ce.Run()
		}
		fmt.Println(helper.MissingCommand(c))
		os.Exit(1)
		return nil
	}

	setup(app)
}
