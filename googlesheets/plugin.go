package googlesheets

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// Plugin creates this (googlesheets) plugin
func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-googlesheets",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromGo().NullIfZero(),
		SchemaMode:       plugin.SchemaModeDynamic,
		TableMapFunc:     PluginTables,
	}
	return p
}

// Map of spreadsheet headers along with the sheet name
var googleSpreadsheetHeadersMap = map[string][]string{}

func PluginTables(ctx context.Context, d *plugin.TableMapData) (map[string]*plugin.Table, error) {
	// Initialize tables
	tables := map[string]*plugin.Table{}

	/* Static tables */
	tables["googlesheets_cell"] = tableGoogleSheetsCell(ctx)
	tables["googlesheets_sheet"] = tableGoogleSheetsSheet(ctx)
	tables["googlesheets_spreadsheet"] = tableGoogleSheetsSpreadsheet(ctx)

	/* Dynamic tables */

	// Get the list of sheets to be retrieved from the spreadsheet
	googleSheetsConfig := GetConfig(d.Connection)

	// Get the headers along with sheet name
	availableSheets, err := getSpreadsheets(ctx, pluginData.Table.Plugin)
	if err != nil {
		return tables, nil
	}

	// Return all valid sheets
	var validSheets []string
	for _, i := range googleSheetsConfig.Sheets {
		if helpers.StringSliceContains(availableSheets, i) {
			validSheets = append(validSheets, i)
		}
	}

	// Get spreadsheet details
	spreadsheetData, err := getSpreadsheetHeaders(ctx, pluginData.Table.Plugin, validSheets)
	if err != nil {
		return tables, nil
	}

	// Create tablemap for all the available sheets
	for _, sheetName := range googleSheetsConfig.Sheets {
		for _, data := range spreadsheetData {
			// Return if empty sheet
			if len(data.Values) == 0 {
				continue
			}

			// Return if first row is empty
			if len(data.Values[0]) == 0 {
				continue
			}

			// Return if A1 cell is empty
			if len(data.Values[0][0].(string)) == 0 {
				continue
			}

			// API wraps the range inside quotes if sheet name contains whitespaces
			// For example, if the sheet name is 'Sheet 1', then range comes as "'Sheet 1'!A1:Z1"
			str := strings.Split(data.Range, "!")[0]
			if len(str) > 0 && str[0] == '\'' {
				str = str[1 : len(str)-1]
			}

			var spreadsheetHeaders []string
			if sheetName == str {
				mergeCellInfo, _ := getMergeCells(ctx, pluginData.Table.Plugin, sheetName)
				maxColsLength := getMaxLength(data.Values)
				for idx, i := range data.Values[0] {
					mergeRow, mergeColumn, _, parentColumn := findMergeCells(mergeCellInfo, int64(1), int64(idx+1))
					if mergeRow != nil && mergeColumn != nil { // Merge cell
						parentData := data.Values[0][*parentColumn-1]
						spreadsheetHeaders[len(spreadsheetHeaders)-1] = fmt.Sprintf("%s [%s]", spreadsheetHeaders[len(spreadsheetHeaders)-1], intToLetters(idx))
						spreadsheetHeaders = append(spreadsheetHeaders, fmt.Sprintf("%s [%s]", parentData.(string), intToLetters(idx+1)))
					} else if len(i.(string)) == 0 {
						columnName := intToLetters(idx + 1) // since index in for is zero-based
						spreadsheetHeaders = append(spreadsheetHeaders, columnName)
					} else {
						if helpers.StringSliceContains(spreadsheetHeaders, i.(string)) {
							spreadsheetHeaders = append(spreadsheetHeaders, fmt.Sprintf("%s [%s]", i.(string), intToLetters(idx+1)))
						} else {
							spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
						}
					}
				}

				/*
					* Case:
					  | Col A |       |       |
					  | ----- | ----- | ----- |
					  | val A | val B | val C |
					* Expected output
					  | Col A | B     | C     |
					  | ----- | ----- | ----- |
					  | val A | val B | val C |
				*/
				if len(data.Values[0]) < maxColsLength {
					for i := len(data.Values[0]); i < maxColsLength; i++ {
						columnName := intToLetters(i + 1)
						spreadsheetHeaders = append(spreadsheetHeaders, columnName)
					}
				}
				googleSpreadsheetHeadersMap[sheetName] = spreadsheetHeaders

				// Create columns
				cols := []*plugin.Column{}
				for idx, j := range spreadsheetHeaders {
					cols = append(cols, &plugin.Column{Name: j, Type: proto.ColumnType_STRING, Transform: transform.FromField(j), Description: fmt.Sprintf("Field %d.", idx)})
				}

				// Create table definition
				tables[sheetName] = &plugin.Table{
					Name:        sheetName,
					Description: fmt.Sprintf("Retrieves data from %s.", sheetName),
					List: &plugin.ListConfig{
						Hydrate: listSpreadsheetWithPath(ctx, pluginData.Table.Plugin, sheetName),
					},
					Columns: cols,
				}
			}
		}
	}

	return tables, nil
}
