package oteleport

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/fujiwara/ssm-lookup/ssm"
	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	"github.com/samber/oops"
)

func MakeVM(ctx context.Context) (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, oops.Wrapf(err, "failed to load default aws config")
	}
	cache := &sync.Map{}
	ssmlookup := ssm.New(awsCfg, cache)
	for _, nf := range ssmlookup.JsonnetNativeFuncs(ctx) {
		vm.NativeFunction(nf)
	}
	for _, nf := range NativeFunctions {
		vm.NativeFunction(nf)
	}
	return vm, nil
}

var NativeFunctions = []*jsonnet.NativeFunction{
	MastEnvNativeFunction,
	EnvNativeFunction,
	JsonescapeNativeFunction,
}

var MastEnvNativeFunction = &jsonnet.NativeFunction{
	Name:   "must_env",
	Params: []ast.Identifier{"name"},
	Func: func(args []interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, oops.Errorf("must_env: invalid arguments length expected 1 got %d", len(args))
		}
		key, ok := args[0].(string)
		if !ok {
			return nil, oops.Errorf("must_env: invalid arguments, expected string got %T", args[0])
		}
		val, ok := os.LookupEnv(key)
		if !ok {
			return nil, oops.Errorf("must_env: %s not set", key)
		}
		return val, nil
	},
}
var EnvNativeFunction = &jsonnet.NativeFunction{
	Name:   "env",
	Params: []ast.Identifier{"name", "default"},
	Func: func(args []interface{}) (interface{}, error) {
		if len(args) != 2 {
			return nil, oops.Errorf("env: invalid arguments length expected 2 got %d", len(args))
		}
		key, ok := args[0].(string)
		if !ok {
			return nil, oops.Errorf("env: invalid 1st arguments, expected string got %T", args[0])
		}
		val := os.Getenv(key)
		if val == "" {
			return args[1], nil
		}
		return val, nil
	},
}

var JsonescapeNativeFunction = &jsonnet.NativeFunction{
	Name:   "json_escape",
	Params: []ast.Identifier{"str"},
	Func: func(args []interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, oops.Errorf("jsonescape: invalid arguments length expected 1 got %d", len(args))
		}
		str, ok := args[0].(string)
		if !ok {
			return nil, oops.Errorf("jsonescape: invalid arguments, expected string got %T", args[0])
		}
		bs, err := json.Marshal(str)
		if err != nil {
			return nil, oops.Wrapf(err, "jsonescape")
		}
		return string(bs), nil
	},
}
