package googlesheets

import (
	"context"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

// Returns the ID of the current working spreadsheet
func getSpreadsheetID(ctx context.Context, p *plugin.TableMapData) string {
	spreadsheetCacheKey := "googlesheets.spreadsheetID"
	if cachedData, ok := p.ConnectionCache.Get(ctx, spreadsheetCacheKey); ok {
		return cachedData.(string)
	}

	googleSheetsConfig := GetConfig(p.Connection)

	var spreadsheetID string
	if googleSheetsConfig.SpreadsheetId != nil {
		spreadsheetID = *googleSheetsConfig.SpreadsheetId
	}

	err := p.ConnectionCache.Set(ctx, spreadsheetCacheKey, spreadsheetID)

	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetID", "connection set error", err)
		return err.Error()
	}

	return spreadsheetID
}

// Returns the ID of the working spreadsheet for static tables
func getSpreadsheetIDStatic(ctx context.Context, p *plugin.QueryData) string {
	spreadsheetCacheKey := "googlesheets.spreadsheetID"
	if cachedData, ok := p.ConnectionCache.Get(ctx, spreadsheetCacheKey); ok {
		return cachedData.(string)
	}

	googleSheetsConfig := GetConfig(p.Connection)

	var spreadsheetID string
	if googleSheetsConfig.SpreadsheetId != nil {
		spreadsheetID = *googleSheetsConfig.SpreadsheetId
	}

	err := p.ConnectionCache.Set(ctx, spreadsheetCacheKey, spreadsheetID)

	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetID", "connection set error", err)
		return err.Error()
	}

	return spreadsheetID
}

// Returns the content of given file, or the inline JSON credential as it is
func pathOrContents(poc string) (string, error) {
	if len(poc) == 0 {
		return poc, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, err
		}
	}

	// Check for valid file path
	if _, err := os.Stat(path); err == nil {
		contents, err := os.ReadFile(path)
		if err != nil {
			return string(contents), err
		}
		return string(contents), nil
	}

	// Return error if content is a file path and the file doesn't exist
	if len(path) > 1 && (path[0] == '/' || path[0] == '\\') {
		return "", fmt.Errorf("%s: no such file or dir", path)
	}

	// Return the inline content
	return poc, nil
}

// Convert column index number to corresponding letter
// For example, 1:A, 2:B, 27:AA, 55:BC
func intToLetters(colIndex int) (letter string) {
	colIndex--
	if firstLetter := colIndex / 26; firstLetter > 0 {
		letter += intToLetters(firstLetter)
		letter += string(rune('A' + colIndex%26))
	} else {
		letter += string(rune('A' + colIndex))
	}

	return
}

// Return the maximum length of a column in a sheet
func getMaxLength(values [][]interface{}) int {
	var maxColsLength int
	for _, value := range values {
		if len(value) > maxColsLength {
			maxColsLength = len(value)
		}
	}
	return maxColsLength
}

func getQualListValues(quals map[string]*proto.QualValue) []string {
	if quals["sheet_name"].GetStringValue() != "" {
		return []string{quals["sheet_name"].GetStringValue()}
	} else if quals["sheet_name"].GetListValue() != nil {
		values := make([]string, 0)
		for _, value := range quals["sheet_name"].GetListValue().Values {
			values = append(values, value.GetStringValue())
		}
		return values
	}

	return nil
}
