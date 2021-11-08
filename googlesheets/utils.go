package googlesheets

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns the ID of the current working spreadsheet
func getSpreadsheetID(_ context.Context, p *plugin.Plugin) string {
	spreadsheetCacheKey := "googlesheets.spreadsheetID"
	if cachedData, ok := p.ConnectionManager.Cache.Get(spreadsheetCacheKey); ok {
		return cachedData.(string)
	}

	googlesheetConfig := GetConfig(p.Connection)

	var spreadsheetID string
	if googlesheetConfig.SpreadsheetId != nil {
		spreadsheetID = *googlesheetConfig.SpreadsheetId
	}

	p.ConnectionManager.Cache.Set(spreadsheetCacheKey, spreadsheetID)

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
		contents, err := ioutil.ReadFile(path)
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

func intToLetters(number int) (letters string){
	number--
	if firstLetter := number/26; firstLetter >0{
			letters += intToLetters(firstLetter)
			letters += string('A' + number%26)
	} else {
			letters += string('A' + number)
	}
			
	return 
}
