/*
**TODO**
* Process data for merge cells
 */

package googlesheets

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func listSpreadsheetWithPath(ctx context.Context, p *plugin.Plugin, sheetName string) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		// Load spreadsheet data
		spreadsheetData, err := getSpreadsheetData(ctx, p, []string{sheetName})
		if err != nil {
			return nil, err
		}

		// No table
		if len(spreadsheetData) == 0 {
			return nil, nil
		}
		data := spreadsheetData[0]

		// Fetch spreadsheet header
		var spreadsheetHeaders []string
		for idx, i := range data.Values[0] {
			if len(i.(string)) == 0 {
				columnName := intToLetters(idx + 1) // since index in for is zero-based
				spreadsheetHeaders = append(spreadsheetHeaders, columnName)
			} else {
				spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
			}
		}

		// Remove the header row from the spreadsheet data
		cellData := append(data.Values[:0], data.Values[0+1:]...)

		// Iterate the spreadsheet rows
		for _, i := range cellData {
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
