package cmd

import (
	"github.com/spf13/cobra"

	"github.com/oasisprotocol/oasis-core-ledger/common"
)

// InitVersions sets a custom version template for the given cobra command.
func InitVersions(cmd *cobra.Command) {
	cobra.AddTemplateFunc("additionalVersions", func() interface{} { return common.Versions })

	cmd.SetVersionTemplate(`Software version: {{.Version}}
{{- with additionalVersions }}
Go toolchain version: {{ .Toolchain }}
{{ end -}}
`)
}
