package googlesheets

import (
	"context"
	"fmt"

	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCSV(ctx context.Context, p *plugin.Plugin) *plugin.Table {
	sheetRange := ctx.Value("range").(string)

	// Load spreadsheet data
	spreadsheetData, err := getSpreadsheetData(ctx, p, sheetRange)
	if err != nil {
		panic(err)
	}

	// Get headers
	var csvHeaders []string
	for _, i := range spreadsheetData.Values[0] {
		csvHeaders = append(csvHeaders, i.(string))
	}

	// Create columns
	cols := []*plugin.Column{}
	for idx, i := range csvHeaders {
		cols = append(cols, &plugin.Column{Name: i, Type: proto.ColumnType_STRING, Transform: transform.FromField(i), Description: fmt.Sprintf("Field %d.", idx)})
	}

	// Create table definition
	return &plugin.Table{
		Name:        sheetRange,
		Description: fmt.Sprintf("CSV file at %s", sheetRange),
		List: &plugin.ListConfig{
			Hydrate: listCSVWithPath(spreadsheetData),
		},
		Columns: cols,
	}
}

func listCSVWithPath(data *sheets.ValueRange) func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		// Fetch spreadsheet header
		var csvHeaders []string
		for _, i := range data.Values[0] {
			csvHeaders = append(csvHeaders, i.(string))
		}

		// Remove the header row from the spreadsheet data
		data := append(data.Values[:0], data.Values[0+1:]...)

		// Iterate the spreadsheet rows
		for _, i := range data {
			row := map[string]string{}
			for idx, j := range i {
				row[csvHeaders[idx]] = j.(string)
			}
			d.StreamListItem(ctx, row)
		}

		return nil, nil
	}
}
