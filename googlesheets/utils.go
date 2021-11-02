package googlesheets

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns the ID of the current working spreadsheet
func getSpreadsheetID(_ context.Context, p *plugin.Plugin) string {
	spreadsheetCacheKey := "googlesheets.spreadsheetID"
	if cachedData, ok := p.ConnectionManager.Cache.Get(spreadsheetCacheKey); ok {
		return cachedData.(string)
	}

	googlesheetConfig := GetConfig(p.Connection)

	var spreadsheetID string
	if googlesheetConfig.SpreadsheetId != nil {
		spreadsheetID = *googlesheetConfig.SpreadsheetId
	}

	p.ConnectionManager.Cache.Set(spreadsheetCacheKey, spreadsheetID)

	return spreadsheetID
}
