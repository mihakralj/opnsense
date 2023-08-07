package cmd

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/mihakralj/opnsense/internal"
	"github.com/spf13/cobra"
)

var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.`,
	Run: func(cmd *cobra.Command, args []string) {

		path := "actions"
		if len(args) >= 1 {
			trimmedArg := strings.Trim(args[0], "/")
			if trimmedArg != "" {
				path = trimmedArg
			}
			parts := strings.Split(path, "/")
			if parts[0] != "actions" {
				path = "actions/" + path
			}
		}
		internal.Checkos()
		bash := `echo "<actions>" && for file in /usr/local/opnsense/service/conf/actions.d/actions_*.conf; do service_name=$(basename "$file" | sed 's/actions_\(.*\).conf/\1/'); echo "  <${service_name}>"; awk 'function escape_xml(str) { gsub(/&/, "&amp;", str); gsub(/</, "&lt;", str); gsub(/>/, "&gt;", str); return str; } BEGIN {FS=":"; action = "";} /\[.*\]/ { if (action != "") {print "    </" action ">"} action = substr($0, 2, length($0) - 2); print "    <" action ">";} !/\[.*\]/ && NF > 1 { gsub(/^[ \t]+|[ \t]+$/, "", $2); value = escape_xml($2); print "      <" $1 ">" value "</" $1 ">";} END { if (action != "") {print "    </" action ">"} }' "$file"; echo "  </${service_name}>"; done && echo "</actions>"`
		config, err := internal.ExecuteCmd(bash, host)
		if err != nil {
			panic(err)
		}
		config_doc := etree.NewDocument()
		config_doc.ReadFromString(config)
		configtty := internal.ConfigToTTY(config_doc, path, depth)
		fmt.Println(configtty)
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)
	// Here you will define your flags and configuration settings.
}