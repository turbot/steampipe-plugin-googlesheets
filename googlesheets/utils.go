package googlesheets

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns the ID of the current working spreadsheet
func getSpreadsheetID(_ context.Context, p *plugin.Plugin) string {
	googlesheetConfig := GetConfig(p.Connection)

	var spreadsheetID string
	if googlesheetConfig.SpreadsheetId != nil {
		spreadsheetID = *googlesheetConfig.SpreadsheetId
	}

	return spreadsheetID
}
