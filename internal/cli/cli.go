package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/jgfranco17/llm-workspace/internal/config"
	"github.com/jgfranco17/llm-workspace/toolbox/registry"
	"github.com/jgfranco17/llm-workspace/toolbox/tools"
)

type appState struct {
	registry *registry.Manager
	config   string
}

type Runner struct {
	state *appState
	cmd   *cobra.Command
}

func New() *Runner {
	state := &appState{}

	rootCmd := &cobra.Command{
		Use:           "tool-runner",
		Short:         "Ollama Agent Tools CLI",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			configPath := state.config
			if configPath == "" {
				configPath = defaultConfigPath()
			}

			file, err := os.Open(configPath)
			if err != nil {
				return fmt.Errorf("open config: %w", err)
			}
			defer file.Close()

			cfg, err := config.Read(file)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			registry := registry.New(nil)
			if err := registry.Add(cfg.Tools...); err != nil {
				return fmt.Errorf("register tools: %w", err)
			}

			state.registry = registry
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&state.config, "config", "c", "tools.json", "Path to tools configuration file")

	rootCmd.AddCommand(newListCommand(state))
	rootCmd.AddCommand(newSchemaCommand(state))
	rootCmd.AddCommand(newExecuteCommand(state))
	rootCmd.AddCommand(newInfoCommand(state))

	return &Runner{
		state: state,
		cmd:   rootCmd,
	}
}

func (r *Runner) Execute() error {
	return r.cmd.Execute()
}

func newListCommand(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			tools := state.registry.ListTools()
			if len(tools) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No tools registered")
				return nil
			}

			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 2, 4, 2, ' ', 0)
			fmt.Fprintln(writer, "NAME\tDESCRIPTION\tPARAMETERS")
			for _, tool := range tools {
				params := formatParams(tool.Parameters)
				fmt.Fprintf(writer, "%s\t%s\t%s\n", tool.Name, tool.Description, params)
			}
			return writer.Flush()
		},
	}
}

func newSchemaCommand(state *appState) *cobra.Command {
	var output string

	cmd := &cobra.Command{
		Use:   "schema",
		Short: "Export tool schemas in Ollama format",
		RunE: func(cmd *cobra.Command, args []string) error {
			schemas := state.registry.AsOllamaSchema()
			payload, err := json.MarshalIndent(schemas, "", "  ")
			if err != nil {
				return fmt.Errorf("serialize schema: %w", err)
			}

			if output != "" {
				if err := os.WriteFile(output, payload, 0o644); err != nil {
					return fmt.Errorf("write schema: %w", err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Schemas exported to %s\n", output)
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(payload))
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file path")
	return cmd
}

func newExecuteCommand(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "execute <tool_name> [key=value ...]",
		Short: "Execute a tool with given parameters",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]
			params, err := parseParams(args[1:])
			if err != nil {
				return err
			}

			result, err := state.registry.Execute(cmd.Context(), toolName, params)
			if err != nil {
				return err
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Result:")
			fmt.Fprintln(cmd.OutOrStdout(), result)
			return nil
		},
	}
}

func newInfoCommand(state *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "info <tool_name>",
		Short: "Show detailed information about a tool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			toolName := args[0]
			schema, _, ok := state.registry.GetHandler(toolName)
			if !ok {
				return fmt.Errorf("tool '%s' not found", toolName)
			}

			fmt.Fprintln(cmd.OutOrStdout(), schema.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "\nDescription: %s\n", schema.Description)

			if len(schema.Parameters) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "\nParameters: none")
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "\nParameters:")
			for _, name := range sortedKeys(schema.Parameters) {
				param := schema.Parameters[name]
				required := ""
				if param.Required {
					required = "*"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  - %s%s: %s\n", name, required, param.Description)
				fmt.Fprintf(cmd.OutOrStdout(), "    Type: %s\n", param.Type)
				if len(param.Enum) > 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "    Options: %s\n", strings.Join(param.Enum, ", "))
				}
			}
			return nil
		},
	}
}

func parseParams(params []string) (tools.Parameters, error) {
	result := make(tools.Parameters)
	for _, param := range params {
		parts := strings.SplitN(param, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid parameter format '%s'. Use key=value", param)
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func formatParams(parameters map[string]config.ToolParameter) string {
	if len(parameters) == 0 {
		return "none"
	}

	parts := make([]string, 0, len(parameters))
	for _, name := range sortedKeys(parameters) {
		param := parameters[name]
		suffix := ""
		if param.Required {
			suffix = "*"
		}
		parts = append(parts, fmt.Sprintf("%s%s", name, suffix))
	}

	return strings.Join(parts, ", ")
}

func sortedKeys(parameters map[string]config.ToolParameter) []string {
	keys := make([]string, 0, len(parameters))
	for key := range parameters {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func defaultConfigPath() string {
	candidates := []string{}
	if env := os.Getenv("TOOL_CONFIG_PATH"); env != "" {
		candidates = append(candidates, env)
	}

	candidates = append(candidates,
		config.DefaultConfigFile,
		filepath.Join("toolbox-go", "config", "tools.json"),
		filepath.Join("config", "tools.json"),
	)

	for _, candidate := range candidates {
		if fileExists(candidate) {
			return candidate
		}
	}

	return config.DefaultConfigFile
}

func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
