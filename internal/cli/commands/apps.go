package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewApps(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "apps", Short: "Zen apps"}
	cmd.AddCommand(
		simpleList(g, "list", "List Zen apps", "/firewall/apps"),
		endpointCommand(g, endpointCommandConfig{Use: "create", Short: "Create app", Method: http.MethodPost, Path: staticPath("/firewall/apps")}),
		endpointCommand(g, endpointCommandConfig{Use: "get <app-id>", Short: "Get app", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s")}),
		endpointCommand(g, endpointCommandConfig{Use: "update <app-id>", Short: "Update app", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s")}),
		endpointCommand(g, endpointCommandConfig{Use: "delete <app-id>", Short: "Delete app", Method: http.MethodDelete, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s"), Confirm: true}),
		endpointCommand(g, endpointCommandConfig{Use: "rotate-token <app-id>", Short: "Rotate app token", Method: http.MethodPost, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/token")}),
		endpointCommand(g, endpointCommandConfig{Use: "blocking <service-id>", Short: "Update blocking mode", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/blocking")}),
		endpointCommand(g, endpointCommandConfig{Use: "update-user <app-id> <user-id>", Short: "Update Zen user", Method: http.MethodPut, Args: cobra.ExactArgs(2), Path: func(args []string) string { return "/firewall/" + args[0] + "/users/" + args[1] }}),
		endpointCommand(g, endpointCommandConfig{Use: "ip-blocklist <app-id>", Short: "Update IP blocklist", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/ip-blocklist")}),
		endpointCommand(g, endpointCommandConfig{Use: "bot-lists <app-id>", Short: "Get bot lists", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/bot-lists")}),
		endpointCommand(g, endpointCommandConfig{Use: "update-bot-lists <app-id>", Short: "Update bot lists", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/bot-lists")}),
		endpointCommand(g, endpointCommandConfig{Use: "ip-lists <app-id>", Short: "Get threat lists", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/ip-lists")}),
		endpointCommand(g, endpointCommandConfig{Use: "update-ip-lists <app-id>", Short: "Update threat lists", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/ip-lists")}),
		endpointCommand(g, endpointCommandConfig{Use: "countries <app-id>", Short: "Get countries", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/countries")}),
		endpointCommand(g, endpointCommandConfig{Use: "update-countries <app-id>", Short: "Update countries", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/countries")}),
		endpointCommand(g, endpointCommandConfig{Use: "users <app-id>", Short: "List Zen users", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/firewall/apps/%s/users")}),
		endpointCommand(g, endpointCommandConfig{Use: "event <app-id> <event-id>", Short: "Get event", Method: http.MethodGet, Args: cobra.ExactArgs(2), Path: func(args []string) string { return "/firewall/apps/" + args[0] + "/events/" + args[1] }}),
	)
	return cmd
}
