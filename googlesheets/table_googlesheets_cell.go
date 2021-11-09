package googlesheets

import (
	"context"
	"fmt"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type cellInfo = struct {
	ColumnName  string
	RowName     int
	CellAddress string
	Value       string
	Formula     string
	Note        string
	SheetName   string
}

//// TABLE DEFINITION

func tableGooglesheetsCell(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesheets_cell",
		Description: "Retrieve information of cells of a sheet in a spreadsheet.",
		List: &plugin.ListConfig{
			Hydrate: listGooglesheetCells,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "sheet_name",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "sheet_name",
				Description: "The name of the sheet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "column_name",
				Description: "The ID of the column.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "row_name",
				Description: "The index of the row.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cell_address",
				Description: "The address of a cell.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "value",
				Description: "The value of a cell.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "formula",
				Description: "The formula configured for a cell.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "note",
				Description: "A user defined note on a cell.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "spreadsheet_id",
				Description: "The ID of the spreadsheet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     spreadsheetID,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// LIST FUNCTION

func listGooglesheetCells(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	opts, err := getSessionConfig(ctx, d.Table.Plugin)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetData", "connection_error", err)
		return nil, err
	}

	// Get the ID of the spreadsheet
	spreadsheetID := getSpreadsheetID(ctx, d.Table.Plugin)

	resp := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Fields(googleapi.Field("sheets(properties.title,data(rowData),merges)"))

	// Additional filters
	if d.KeyColumnQuals["sheet_name"] != nil {
		resp.Ranges(d.KeyColumnQuals["sheet_name"].GetStringValue())
	}
	data, err := resp.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	if data.Sheets != nil {
		for _, sheet := range data.Sheets {
			if sheet.Data != nil {
				for _, i := range sheet.Data {
					for row_count, row := range i.RowData {
						for col_count, value := range row.Values {
							mergeRow, mergeColumn, parentRow, parentColumn := findMergeCells(sheet.Merges, int64(row_count+1), int64(col_count+1))
							if mergeRow != nil && mergeColumn != nil {
								parentData := i.RowData[*parentRow-1].Values[*parentColumn-1]
								var formulaValue string
								if parentData.UserEnteredValue.FormulaValue != nil {
									formulaValue = *parentData.UserEnteredValue.FormulaValue
								}
								d.StreamListItem(ctx, cellInfo{
									SheetName:   sheet.Properties.Title,
									ColumnName:  intToLetters(col_count + 1),
									RowName:     row_count + 1,
									CellAddress: fmt.Sprintf("%s%d", intToLetters(col_count+1), row_count+1),
									Value:       parentData.FormattedValue,
									Formula:     formulaValue,
									Note:        value.Note,
								})
							}
							if value.UserEnteredValue != nil {
								var formulaValue string
								if value.UserEnteredValue.FormulaValue != nil {
									formulaValue = *value.UserEnteredValue.FormulaValue
								}
								d.StreamListItem(ctx, cellInfo{
									SheetName:   sheet.Properties.Title,
									ColumnName:  intToLetters(col_count + 1),
									RowName:     row_count + 1,
									CellAddress: fmt.Sprintf("%s%d", intToLetters(col_count+1), row_count+1),
									Value:       value.FormattedValue,
									Formula:     formulaValue,
									Note:        value.Note,
								})
							}
						}
					}
				}
			}
		}
	}

	return nil, nil
}

func findMergeCells(mergeInfo []*sheets.GridRange, currentRow int64, currentColumn int64) (*int64, *int64, *int64, *int64) {
	for _, mergeData := range mergeInfo {
		rowDiff := mergeData.EndRowIndex - mergeData.StartRowIndex
		colDiff := mergeData.EndColumnIndex - mergeData.StartColumnIndex

		if rowDiff > 1 && colDiff == 1 { // vertically merged
			if currentRow >= mergeData.StartRowIndex+2 && currentRow <= mergeData.EndRowIndex && currentColumn == mergeData.StartColumnIndex+1 {
				parentRow := mergeData.StartRowIndex + 1
				parentColumn := mergeData.StartColumnIndex + 1
				return &currentRow, &currentColumn, &parentRow, &parentColumn
			}
		} else if colDiff > 1 && rowDiff == 1 { // horizontally merged
			if currentColumn >= mergeData.StartColumnIndex+2 && currentColumn <= mergeData.EndColumnIndex && currentRow == mergeData.StartRowIndex+1 {
				parentRow := mergeData.StartRowIndex + 1
				parentColumn := mergeData.StartColumnIndex + 1
				return &currentRow, &currentColumn, &parentRow, &parentColumn
			}
		} else { // Mixed
			if currentRow >= mergeData.StartRowIndex+1 && currentRow <= mergeData.EndRowIndex && currentColumn >= mergeData.StartColumnIndex+1 && currentColumn <= mergeData.EndColumnIndex {
				parentRow := mergeData.StartRowIndex + 1
				parentColumn := mergeData.StartColumnIndex + 1
				return &currentRow, &currentColumn, &parentRow, &parentColumn
			}
		}
	}
	return nil, nil, nil, nil
}
