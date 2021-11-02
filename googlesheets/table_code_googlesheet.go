package googlesheets

import (
	"context"
	"fmt"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableSpreadsheets(ctx context.Context, p *plugin.Plugin) *plugin.Table {
	sheetName := ctx.Value("sheet").(string)

	// Load spreadsheet data
	spreadsheetData, err := getSpreadsheetData(ctx, p, sheetName + "!1:1")
	if err != nil {
		panic(err)
	}

	// Get headers
	var spreadsheetHeaders []string

	// Return if no rows found
	if len(spreadsheetData.Values) == 0 {
		return nil
	}

	// Extract spreadsheet headers
	for _, i := range spreadsheetData.Values[0] {
		spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
	}

	// Create columns
	cols := []*plugin.Column{}
	for idx, i := range spreadsheetHeaders {
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(i), Description: fmt.Sprintf("Field %d.", idx)})
	}

	// Create table definition
	return &plugin.Table{
		Name:        sheetName,
		Description: fmt.Sprintf("Retrieves data from %s.", sheetName),
		List: &plugin.ListConfig{
			Hydrate: listSpreadsheetWithPath(ctx, p, sheetName),
		},
		Columns: cols,
	}
}

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
			spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
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
