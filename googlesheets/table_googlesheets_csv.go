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
	return func(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
		// Get spreadsheet details
		spreadsheetData, err := getSpreadsheetData(ctx, p, []string{sheetName})
		if err != nil {
			return nil, err
		}
		spreadsheetHeaders := googleSpreadsheetHeadersMap[sheetName]

		for _, sheet := range spreadsheetData {
			for _, i := range sheet.Data {
				for row_count, row := range i.RowData {
					// Skip first row, or header
					if row_count == 0 {
						continue
					}
					rowData := map[string]string{}
					for col_count, value := range row.Values {
						mergeRow, mergeColumn, parentRow, parentColumn := findMergeCells(sheet.Merges, int64(row_count+1), int64(col_count+1))
						if mergeRow != nil && mergeColumn != nil {
							parentData := i.RowData[*parentRow-1].Values[*parentColumn-1]
							rowData[spreadsheetHeaders[col_count]] = parentData.FormattedValue
						} else {
							rowData[spreadsheetHeaders[col_count]] = value.FormattedValue
						}
					}
					d.StreamListItem(ctx, rowData)
				}
			}
		}
		return nil, nil
	}
}
