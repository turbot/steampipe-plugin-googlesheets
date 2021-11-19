package main

import (
	"github.com/turbot/steampipe-plugin-googlesheets/googlesheets"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		PluginFunc: googlesheets.Plugin})
}
