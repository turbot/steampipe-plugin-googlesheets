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
	Hyperlink   string
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
				{
					Name:    "ranges",
					Require: plugin.Optional,
				},
				{
					Name:    "column_name",
					Require: plugin.Optional,
				},
				{
					Name:    "row_name",
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
				Name:        "hyperlink",
				Description: "A hyperlink this cell points to, if any. If the cell contains multiple hyperlinks, this field will be empty.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ranges",
				Description: "The ranges to retrieve from the spreadsheet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("ranges"),
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

	resp := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Fields(googleapi.Field("sheets(properties.title,data(rowData,startColumn,startRow),merges)"))

	// Additional filters
	quals := d.KeyColumnQuals
	if quals["sheet_name"] != nil {
		if quals["ranges"] != nil {
			ranges := fmt.Sprintf("%s!%s", quals["sheet_name"].GetStringValue(), quals["ranges"].GetStringValue())
			resp.Ranges(ranges)
		} else if quals["row_name"] != nil && quals["column_name"] != nil {
			ranges := fmt.Sprintf("%s!%s%d", quals["sheet_name"].GetStringValue(), quals["column_name"].GetStringValue(), quals["row_name"].GetInt64Value())
			resp.Ranges(ranges)
		} else {
			resp.Ranges(quals["sheet_name"].GetStringValue())
		}
	}
	data, err := resp.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	if data.Sheets != nil {
		for _, sheet := range data.Sheets {
			if sheet.Data != nil {
				for _, i := range sheet.Data {
					for rowCount, row := range i.RowData {
						rowCount = rowCount + int(i.StartRow)
						for colCount, value := range row.Values {
							colCount = colCount + int(i.StartColumn)
							mergeRow, mergeColumn, parentRow, parentColumn := findMergeCells(sheet.Merges, int64(rowCount+1), int64(colCount+1))
							if mergeRow != nil && mergeColumn != nil { // Merge cell
								parentData := i.RowData[*parentRow-1].Values[*parentColumn-1]
								rowInfo := getCellInfo(sheet.Properties.Title, rowCount, colCount, parentData)
								d.StreamListItem(ctx, rowInfo)
							} else if value.UserEnteredValue != nil && value.UserEnteredValue.FormulaValue != nil { // Image in cell
								rowInfo := getCellInfo(sheet.Properties.Title, rowCount, colCount, value)
								d.StreamListItem(ctx, rowInfo)
							} else if value.FormattedValue != "" {
								rowInfo := getCellInfo(sheet.Properties.Title, rowCount, colCount, value)
								d.StreamListItem(ctx, rowInfo)
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

func getCellInfo(sheetName string, rowCount int, colCount int, data *sheets.CellData) cellInfo {
	var formulaValue string
	if data.UserEnteredValue.FormulaValue != nil {
		formulaValue = *data.UserEnteredValue.FormulaValue
	}
	result := cellInfo{
		SheetName:   sheetName,
		ColumnName:  intToLetters(colCount + 1),
		RowName:     rowCount + 1,
		CellAddress: fmt.Sprintf("%s%d", intToLetters(colCount+1), rowCount+1),
		Value:       data.FormattedValue,
		Formula:     formulaValue,
		Note:        data.Note,
		Hyperlink:   data.Hyperlink,
	}

	return result
}
