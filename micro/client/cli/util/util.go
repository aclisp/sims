// Package cliutil contains methods used across all cli commands
// @todo: get rid of os.Exits and use errors instread
package util

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	ccli "github.com/micro/cli/v2"
	"github.com/micro/micro/v2/internal/config"
	"github.com/micro/micro/v2/internal/platform"
	"github.com/micro/micro/v2/service/runtime/profile"
)

const (
	// EnvLocal is a builtin environment, it means services launched
	// with `micro run` will use default, zero dependency implementations for
	// interfaces, like mdns for registry.
	EnvLocal = "local"
	// EnvServer is a builtin environment, it represents your local `micro server`
	EnvServer = "server"
	// EnvPlatform is a builtin environment, the One True Micro Live(tm) environment.
	EnvPlatform = "platform"
)

const (
	// localProxyAddress is the default proxy address for environment local
	// local env does not use other services so talking about a proxy
	localProxyAddress = "none"
	// serverProxyAddress is the default proxy address for environment server
	serverProxyAddress = "127.0.0.1:8081"
	// platformProxyAddress is teh default proxy address for environment platform
	platformProxyAddress = "proxy.micro.mu"
)

var defaultEnvs = map[string]Env{
	EnvLocal: {
		Name:         EnvLocal,
		ProxyAddress: localProxyAddress,
	},
	EnvServer: {
		Name:         EnvServer,
		ProxyAddress: serverProxyAddress,
	},
	EnvPlatform: {
		Name:         EnvPlatform,
		ProxyAddress: platformProxyAddress,
	},
}

func isBuiltinService(command string) bool {
	for _, service := range platform.Services {
		if command == service {
			return true
		}
	}
	return false
}

// SetupCommand includes things that should run for each command.
func SetupCommand(ctx *ccli.Context) {
	switch ctx.Args().First() {
	case "new", "server", "help":
		return
	}
	if ctx.Args().Len() == 1 && isBuiltinService(ctx.Args().First()) {
		return
	}
	if ctx.Args().Len() >= 1 && ctx.Args().First() == "env" {
		return
	}

	toFlag := func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, "MICRO_", ""))
	}
	setFlags := func(envars []string) {
		for _, envar := range envars {
			// setting both env and flags here
			// as the proxy settings for example did not take effect
			// with only flags
			parts := strings.Split(envar, "=")
			key := toFlag(parts[0])
			os.Setenv(parts[0], parts[1])
			ctx.Set(key, parts[1])
		}
	}

	env := GetEnv(ctx)

	// if we're running a local environment return here
	if len(env.ProxyAddress) == 0 || env.Name == EnvLocal {
		return
	}

	switch env.Name {
	case EnvServer:
		setFlags(profile.ServerCLI())
	case EnvPlatform:
		setFlags(profile.PlatformCLI())
	default:
		// default case for ad hoc envs, see comments above about tests
		setFlags(profile.ServerCLI())
	}

	// Set the proxy
	setFlags([]string{"MICRO_PROXY=" + env.ProxyAddress})
}

type Env struct {
	Name         string
	ProxyAddress string
}

func AddEnv(env Env) {
	envs := getEnvs()
	envs[env.Name] = env
	setEnvs(envs)
}

func getEnvs() map[string]Env {
	envsJSON, err := config.Get("envs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	envs := map[string]Env{}
	if len(envsJSON) > 0 {
		err := json.Unmarshal([]byte(envsJSON), &envs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	for k, v := range defaultEnvs {
		envs[k] = v
	}
	return envs
}

func setEnvs(envs map[string]Env) {
	envsJSON, err := json.Marshal(envs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = config.Set(string(envsJSON), "envs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetEnv returns the current selected environment
// Does not take
func GetEnv(ctx *ccli.Context) Env {
	var envName string
	if len(ctx.String("env")) > 0 {
		envName = ctx.String("env")
	} else {
		env, err := config.Get("env")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if env == "" {
			env = EnvLocal
		}
		envName = env
	}

	return GetEnvByName(envName)
}

func GetEnvByName(env string) Env {
	envs := getEnvs()

	envir, ok := envs[env]
	if !ok {
		fmt.Println(fmt.Sprintf("Env \"%s\" not found. See `micro env` for available environments.", env))
		os.Exit(1)
	}

	if len(envir.ProxyAddress) == 0 {
		return envir
	}

	// default to :443
	if _, port, _ := net.SplitHostPort(envir.ProxyAddress); len(port) == 0 {
		envir.ProxyAddress = net.JoinHostPort(envir.ProxyAddress, "443")
	}

	return envir
}

func GetEnvs() []Env {
	envs := getEnvs()
	ret := []Env{defaultEnvs[EnvLocal], defaultEnvs[EnvServer], defaultEnvs[EnvPlatform]}
	nonDefaults := []Env{}
	for _, env := range envs {
		if _, isDefault := defaultEnvs[env.Name]; !isDefault {
			nonDefaults = append(nonDefaults, env)
		}
	}
	// @todo order nondefault envs alphabetically
	ret = append(ret, nonDefaults...)
	return ret
}

// SetEnv selects an environment to be used.
func SetEnv(envName string) {
	envs := getEnvs()
	_, ok := envs[envName]
	if !ok {
		fmt.Printf("Environment '%v' does not exist\n", envName)
		os.Exit(1)
	}
	config.Set(envName, "env")
}

func IsLocal(ctx *ccli.Context) bool {
	return GetEnv(ctx).Name == EnvLocal
}

func IsServer(ctx *ccli.Context) bool {
	return GetEnv(ctx).Name == EnvServer
}

func IsPlatform(ctx *ccli.Context) bool {
	return GetEnv(ctx).Name == EnvPlatform
}
