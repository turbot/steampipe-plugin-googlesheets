package googlesheets

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/mitchellh/go-homedir"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns all the spreadsheets at the given ID
func getSpreadsheets(ctx context.Context, d *plugin.Plugin) ([]string, error) {
	// have we already created and cached the service?
	serviceCacheKey := "googlesheets.sheets"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.([]string), nil
	}

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
	resp, err := svc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	var spreadsheetList []string
	for _, sheet := range resp.Sheets {
		spreadsheetList = append(spreadsheetList, sheet.Properties.Title)
	}

	// cache the service
	d.ConnectionManager.Cache.Set(serviceCacheKey, resp)

	return spreadsheetList, nil
}

// Returns all the cells at the given spreadsheet
func getSpreadsheetData(ctx context.Context, d *plugin.Plugin, sheetRange string) (*sheets.ValueRange, error) {
	// have we already created and cached the service?
	serviceCacheKey := "googlesheets" + sheetRange
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sheets.ValueRange), nil
	}

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

	resp, err := svc.Spreadsheets.Values.Get(spreadsheetID, sheetRange).ValueRenderOption("FORMATTED_VALUE").Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	// cache the service
	d.ConnectionManager.Cache.Set(serviceCacheKey, resp)

	return resp, nil
}

func getSessionConfig(ctx context.Context, d *plugin.Plugin) ([]option.ClientOption, error) {
	opts := []option.ClientOption{}

	// Get credential file path, and user to impersonate from config (if mentioned)
	var credentialContent, tokenPath string
	googlesheetConfig := GetConfig(d.Connection)
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
	cacheKey := "googleworkspace.token_source"
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
