package googlesheets

import (
	"context"
	"fmt"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

type cellInfo = struct {
	Column    string
	Row       int
	Cell      string
	Value     string
	Formula   string
	Note      string
	Hyperlink string
	SheetName string
}

//// TABLE DEFINITION

func tableGoogleSheetsCell(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesheets_cell",
		Description: "Retrieve information of cells of a sheet in a spreadsheet.",
		List: &plugin.ListConfig{
			Hydrate: listGoogleSheetCells,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "sheet_name",
					Require: plugin.Optional,
				},
				{
					Name:    "range",
					Require: plugin.Optional,
				},
				{
					Name:    "cell",
					Require: plugin.Optional,
				},
				{
					Name:    "col",
					Require: plugin.Optional,
				},
				{
					Name:    "row",
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
				Name:        "col",
				Description: "The ID of the column.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Column"),
			},
			{
				Name:        "row",
				Description: "The index of the row.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "cell",
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
				Name:        "range",
				Description: "The ranges to retrieve from the spreadsheet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("range"),
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

func listGoogleSheetCells(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	opts, err := getSessionConfig(ctx, d.Table.Plugin)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("listGoogleSheetCells", "connection_error", err)
		return nil, err
	}

	// Get the ID of the spreadsheet
	spreadsheetID := getSpreadsheetID(ctx, d.Table.Plugin)

	resp := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Fields(googleapi.Field("sheets(properties.title,data(rowData,startColumn,startRow),merges)"))

	// Additional filters
	quals := d.KeyColumnQuals

	/*
		 * If `range` qual is defined, API will use that filter directly
		 * If `sheet_name` qual is defined, the following checks will be performed
		    - If `cell` qual is defined, API will use that value to build `ranges` filter (i.e. '[<sheet_name>!<cell>]')
				- If only `row` is defined, `ranges` filter value will be '[<sheet_name>!<row>:<row>]'
				- If only `col` is defined, `ranges` filter value will be '[<sheet_name>!<col>:<col>]'
				- If both `row` and `col` are defined, `ranges` filter value will be '[<sheet_name>!<col><row>]'
		 * If only `sheet_name` is defined, it will query the whole sheet
	*/

	// If `range` qual is defined, API will use that filter directly
	if quals["range"] != nil && quals["range"].GetStringValue() != "" {
		resp.Ranges(quals["range"].GetStringValue())
	} else if quals["sheet_name"] != nil {
		sheetName := quals["sheet_name"].GetStringValue()
		if quals["cell"] != nil { // only `cell` defined
			sheetRange := fmt.Sprintf("%s!%s", sheetName, quals["cell"].GetStringValue())
			resp.Ranges(sheetRange)
		} else if quals["row"] != nil && quals["col"] != nil { // both `row` and `col` defined
			row := quals["row"].GetInt64Value()
			col := quals["col"].GetStringValue()
			sheetRange := fmt.Sprintf("%s!%s%d", sheetName, col, row)
			resp.Ranges(sheetRange)
		} else if quals["row"] != nil { // only `row` defined
			row := quals["row"].GetInt64Value()
			sheetRange := fmt.Sprintf("%s!%d:%d", sheetName, row, row)
			resp.Ranges(sheetRange)
		} else if quals["col"] != nil { // only `col` defined
			col := quals["col"].GetStringValue()
			sheetRange := fmt.Sprintf("%s!%s:%s", sheetName, col, col)
			resp.Ranges(sheetRange)
		} else {
			ranges := getQualListValues(quals)
			resp.Ranges(ranges...)
		}
	}

	data, err := resp.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	/*
	   * JSON representation of response
	    {
	      "sheets": [
	        {
	          "data": [
	            {
	              "rowData": [
	                {
	                  "values": [{column properties...}]
	                }
	              ]
	            }
	          ]
	        }
	      ]
	    }
	*/
	if data.Sheets != nil {
		for _, sheet := range data.Sheets {
			if sheet.Data != nil {
				for _, i := range sheet.Data {
					for rowCount, row := range i.RowData {
						// If a range has been passed to query a particular range, `StartRow` will indicate the start row index(zero-based)
						rowCount = rowCount + int(i.StartRow)
						for colCount, value := range row.Values {
							var rowInfo cellInfo

							// If a range has been passed to query a particular range, `StartColumn` will indicate the start column index(zero-based)
							colCount = colCount + int(i.StartColumn)
							mergeRow, mergeColumn, parentRow, parentColumn := findMergeCells(sheet.Merges, int64(rowCount+1), int64(colCount+1))
							if mergeRow != nil && mergeColumn != nil { // Merge cell
								if len(i.RowData) > int(*parentRow) && i.RowData[*parentRow-1] != nil && i.RowData[*parentRow-1].Values[*parentColumn-1] != nil {
									parentData := i.RowData[*parentRow-1].Values[*parentColumn-1]
									rowInfo = getCellInfo(sheet.Properties.Title, rowCount, colCount, parentData)
								}
							} else if value.UserEnteredValue != nil && value.UserEnteredValue.FormulaValue != nil { // Image in cell
								rowInfo = getCellInfo(sheet.Properties.Title, rowCount, colCount, value)
							} else if value.FormattedValue != "" {
								rowInfo = getCellInfo(sheet.Properties.Title, rowCount, colCount, value)
							}

							if rowInfo.Value != "" {
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

// findMergeCells identifies the merge cells and returns the merge cell along with its parent cell details
func findMergeCells(mergeInfo []*sheets.GridRange, currentRow int64, currentColumn int64) (*int64, *int64, *int64, *int64) {
	for _, mergeData := range mergeInfo {
		/*
		 * Calculate difference between startRow, endRow index; and startColumn, endColumn index
		 * If rowDiff is >1, it is vertically merged; and
		 * If colDiff is >1, it is horizontally merged
		 * else multiple rows are merged (e.g. A2,B2,A3,B3 are merged together)
		 * This function will identify the parent of the merge cell, and returns the parent cell data; and
		 * the merge cell will use the value of their corresponding parent
		 */
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
	if data.UserEnteredValue != nil && data.UserEnteredValue.FormulaValue != nil {
		formulaValue = *data.UserEnteredValue.FormulaValue
	}
	result := cellInfo{
		SheetName: sheetName,
		Column:    intToLetters(colCount + 1),
		Row:       rowCount + 1,
		Cell:      fmt.Sprintf("%s%d", intToLetters(colCount+1), rowCount+1),
		Value:     data.FormattedValue,
		Formula:   formulaValue,
		Note:      data.Note,
		Hyperlink: data.Hyperlink,
	}

	return result
}
