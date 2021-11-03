package googlesheets

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func listSpreadsheetWithPath(ctx context.Context, p *plugin.Plugin, sheetName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		// Load spreadsheet data
		spreadsheetData, err := getSpreadsheetData(ctx, p, sheetName)
		if err != nil {
			return nil, err
		}

		// Return if no rows found
		if len(spreadsheetData.Values) == 0 {
			return nil, err
		}

		// Fetch spreadsheet header
		var spreadsheetHeaders []string
		for _, i := range spreadsheetData.Values[0] {
			header := i.(string)
			// If no value passed as header use ?column? as column name
			if len(header) == 0 {
				header = "?column?"
			}
			spreadsheetHeaders = append(spreadsheetHeaders, header)
		}

		// Remove the header row from the spreadsheet data
		data := append(spreadsheetData.Values[:0], spreadsheetData.Values[0+1:]...)

		// Iterate the spreadsheet rows
		for _, i := range data {
			row := map[string]string{}
			for idx, j := range i {
				row[spreadsheetHeaders[idx]] = j.(string)
			}
			d.StreamListItem(ctx, row)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}

		return nil, nil
	}
}
