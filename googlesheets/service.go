package googlesheets

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns all the cells of a sheet in given spreadsheet
func getSpreadsheetHeaders(ctx context.Context, d *plugin.Plugin, sheetNames []string) (map[string][]string, error) {
	// To get config arguments from plugin config file
	opts, err := getSessionConfig(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetHeaders", "connection_error", err)
		return nil, err
	}
	spreadsheetID := getSpreadsheetID(ctx, d)

	// Create table range to get the first row of every sheet
	var sheetRanges []string
	for _, i := range sheetNames {
		sheetRanges = append(sheetRanges, fmt.Sprintf("%s!1:1", i))
	}

	resp, err := svc.Spreadsheets.Values.BatchGet(spreadsheetID).ValueRenderOption("FORMATTED_VALUE").Ranges(sheetRanges...).Fields(googleapi.Field("valueRanges")).Context(ctx).Do()
	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetHeaders", "spreadsheet_batchget", err)
		return nil, err
	}

	// Create map of headers along with corresponding sheet name
	spreadsheetHeadersMap := make(map[string][]string)
	for _, i := range resp.ValueRanges {
		if len(i.Values) == 0 {
			continue
		}
		var spreadsheetHeaders []string
		str := strings.Split(i.Range, "!")[0]

		// API wraps the range inside quotes if sheet name contains whitespaces
		// For example, if the sheet name is 'Sheet 1', then range comes as "'Sheet 1'!A1:Z1"
		if len(str) > 0 && str[0] == '\'' {
			str = str[1 : len(str)-1]
		}

		for _, i := range i.Values[0] {
			spreadsheetHeaders = append(spreadsheetHeaders, i.(string))
		}
		spreadsheetHeadersMap[str] = spreadsheetHeaders
	}

	return spreadsheetHeadersMap, nil
}

// Returns all the cells of a sheet in given spreadsheet
func getSpreadsheetData(ctx context.Context, d *plugin.Plugin, sheetRange string) (*sheets.ValueRange, error) {
	// To get config arguments from plugin config file
	opts, err := getSessionConfig(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheetData", "connection_error", err)
		return nil, err
	}

	spreadsheetID := getSpreadsheetID(ctx, d)

	resp, err := svc.Spreadsheets.Values.Get(spreadsheetID, sheetRange).ValueRenderOption("FORMATTED_VALUE").Fields(googleapi.Field("values")).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Returns all the spreadsheets at the given ID
func getSpreadsheets(ctx context.Context, d *plugin.Plugin) ([]string, error) {
	// To get config arguments from plugin config file
	opts, err := getSessionConfig(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("getSpreadsheets", "connection_error", err)
		return nil, err
	}

	// Read spreadsheetID from config
	spreadsheetID := getSpreadsheetID(ctx, d)

	// Get the spreadsheets
	resp, err := svc.Spreadsheets.Get(spreadsheetID).Fields(googleapi.Field("sheets(properties.title)")).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var spreadsheetList []string
	for _, sheet := range resp.Sheets {
		spreadsheetList = append(spreadsheetList, sheet.Properties.Title)
	}

	return spreadsheetList, nil
}

func getSessionConfig(ctx context.Context, d *plugin.Plugin) ([]option.ClientOption, error) {
	opts := []option.ClientOption{}

	// Get credential file path, and user to impersonate from config (if mentioned)
	var credentialContent, tokenPath string
	googlesheetConfig := GetConfig(d.Connection)

	// Return if no SpreadsheetID provided
	if *googlesheetConfig.SpreadsheetId == "" {
		return nil, errors.New("spreadsheet_id must be configured")
	}

	if googlesheetConfig.Credentials != nil {
		credentialContent = *googlesheetConfig.Credentials
	}
	if googlesheetConfig.TokenPath != nil {
		tokenPath = *googlesheetConfig.TokenPath
	}

	// If credential path provided, use domain-wide delegation
	if credentialContent != "" {
		ts, err := getTokenSource(ctx, d)
		if err != nil {
			return nil, err
		}
		opts = append(opts, option.WithTokenSource(ts))
		return opts, nil
	}

	// If token path provided, authenticate using OAuth 2.0
	if tokenPath != "" {
		opts = append(opts, option.WithCredentialsFile(tokenPath))
		return opts, nil
	}

	return nil, nil
}

// Returns a JWT TokenSource using the configuration and the HTTP client from the provided context.
func getTokenSource(ctx context.Context, d *plugin.Plugin) (oauth2.TokenSource, error) {
	// Note: based on https://developers.google.com/admin-sdk/directory/v1/guides/delegation#go

	// have we already created and cached the token?
	cacheKey := "googlesheets.token_source"
	if ts, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return ts.(oauth2.TokenSource), nil
	}

	// Get credential file path, and user to impersonate from config (if mentioned)
	var impersonateUser string
	googlesheetConfig := GetConfig(d.Connection)

	// Read credential from JSON string, or from the given path
	credentialContent, err := pathOrContents(*googlesheetConfig.Credentials)
	if err != nil {
		return nil, err
	}

	if googlesheetConfig.ImpersonatedUserEmail != nil {
		impersonateUser = *googlesheetConfig.ImpersonatedUserEmail
	}

	// Return error, since impersonation required to authenticate using domain-wide delegation
	if impersonateUser == "" {
		return nil, errors.New("impersonated_user_email must be configured")
	}

	// Authorize the request
	config, err := google.JWTConfigFromJSON(
		[]byte(credentialContent),
		sheets.SpreadsheetsReadonlyScope,
	)
	if err != nil {
		return nil, err
	}
	config.Subject = impersonateUser

	ts := config.TokenSource(ctx)

	// cache the token source
	d.ConnectionManager.Cache.Set(cacheKey, ts)

	return ts, nil
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
