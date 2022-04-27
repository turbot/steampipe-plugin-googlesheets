package googlesheets

import (
	"context"
	"errors"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

// Returns headers of a sheet in given spreadsheet
func getSpreadsheetHeaders(ctx context.Context, d *plugin.Plugin, sheetNames []string) ([]*sheets.ValueRange, error) {
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

	resp, err := svc.Spreadsheets.Values.BatchGet(spreadsheetID).ValueRenderOption("FORMATTED_VALUE").Ranges(sheetNames...).Fields(googleapi.Field("valueRanges")).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return resp.ValueRanges, nil
}

// Returns all the sheets in a given spreadsheets
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

// Returns all the merge cells in a given sheet
func getMergeCells(ctx context.Context, d *plugin.Plugin, sheetName string) ([]*sheets.GridRange, error) {
	// To get config arguments from plugin config file
	opts, err := getSessionConfig(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := sheets.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("getMergeCells", "connection_error", err)
		return nil, err
	}

	spreadsheetID := getSpreadsheetID(ctx, d)

	resp, err := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Ranges(sheetName).Fields(googleapi.Field("sheets(merges)")).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Sheets[0].Merges) > 0 {
		return resp.Sheets[0].Merges, nil
	}

	return nil, nil
}

// Returns all the cells of a sheet in given spreadsheet
func getSpreadsheetData(ctx context.Context, d *plugin.Plugin, sheetNames []string) ([]*sheets.Sheet, error) {
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

	resp := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Fields(googleapi.Field("sheets(properties.title,data(rowData(values(formattedValue))),merges)"))
	if len(sheetNames) > 0 {
		resp.Ranges(sheetNames...)
	}
	data, err := resp.Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	return data.Sheets, nil
}

func getSessionConfig(ctx context.Context, d *plugin.Plugin) ([]option.ClientOption, error) {
	opts := []option.ClientOption{}

	// Get credential file path, and user to impersonate from config (if mentioned)
	var credentialContent, tokenPath string
	googleSheetsConfig := GetConfig(d.Connection)

	// Return if no Spreadsheet ID provided
	if googleSheetsConfig.SpreadsheetId == nil {
		return nil, errors.New("spreadsheet_id must be configured")
	}

	if googleSheetsConfig.TokenPath != nil {
		tokenPath = *googleSheetsConfig.TokenPath
	}

	if googleSheetsConfig.Credentials != nil {
		credentialContent = *googleSheetsConfig.Credentials
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
	googleSheetsConfig := GetConfig(d.Connection)

	// Read credential from JSON string, or from the given path
	credentialContent, err := pathOrContents(*googleSheetsConfig.Credentials)
	if err != nil {
		return nil, err
	}

	if googleSheetsConfig.ImpersonatedUserEmail != nil {
		impersonateUser = *googleSheetsConfig.ImpersonatedUserEmail
	}

	// Return error, since impersonation required to authenticate using domain-wide delegation
	if impersonateUser == "" {
		return nil, errors.New("impersonated_user_email must be configured")
	}

	// Authorize the request
	config, err := google.JWTConfigFromJSON(
		[]byte(credentialContent),
		sheets.DriveReadonlyScope,
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
