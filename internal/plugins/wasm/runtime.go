package wasm

import (
	"context"
	"fmt"
	"os"

	"github.com/tetratelabs/wazero"
)

type PluginRuntime struct {
	runtime wazero.Runtime
	modules map[string]wazero.CompiledModule
}

func NewPluginRuntime() *PluginRuntime {
	ctx := context.Background()
	runtime := wazero.NewRuntime(ctx)

	return &PluginRuntime{
		runtime: runtime,
		modules: make(map[string]wazero.CompiledModule),
	}
}

func (pr *PluginRuntime) LoadPlugin(name, wasmPath string) error {
	wasmBytes, err := os.ReadFile(wasmPath)
	if err != nil {
		return fmt.Errorf("failed to read WASM file: %w", err)
	}

	ctx := context.Background()
	compiled, err := pr.runtime.CompileModule(ctx, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to compile WASM module: %w", err)
	}

	pr.modules[name] = compiled
	return nil
}

func (pr *PluginRuntime) ExecuteFunction(pluginName, functionName string, args []uint64) ([]uint64, error) {
	module, exists := pr.modules[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not loaded", pluginName)
	}

	ctx := context.Background()
	config := wazero.NewModuleConfig().WithStdout(os.Stdout).WithStderr(os.Stderr)
	mod, err := pr.runtime.InstantiateModule(ctx, module, config)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %w", err)
	}
	defer mod.Close(ctx)

	fn := mod.ExportedFunction(functionName)
	if fn == nil {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	result, err := fn.Call(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to call function: %w", err)
	}

	return result, nil
}

func (pr *PluginRuntime) Close() error {
	ctx := context.Background()
	return pr.runtime.Close(ctx)
}
