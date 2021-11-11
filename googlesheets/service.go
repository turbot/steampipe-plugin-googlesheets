package googlesheets

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// Returns all the cells of a sheet in given spreadsheet
func getSpreadsheetHeadersMap(ctx context.Context, d *plugin.Plugin, sheetNames []string) (map[string][]string, error) {
	// have we already created and cached the token?
	cacheKey := "googlesheets.headers_map"
	if ts, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return ts.(map[string][]string), nil
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

	resp, err := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Ranges(sheetNames...).Fields(googleapi.Field("sheets(properties.title,data(rowData),merges)")).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	headersMap := map[string][]string{}
	for _, data := range resp.Sheets {
		if len(data.Data) == 0 || len(data.Data[0].RowData) == 0 || len(data.Data[0].RowData[0].Values) == 0 || data.Data[0].RowData[0].Values[0].UserEnteredValue == nil {
			continue
		}

		var spreadsheetHeaders []string
		for _, i := range data.Data {
			rowData := i.RowData[0].Values

			for col_count, value := range rowData {
				mergeRow, mergeColumn, parentRow, parentColumn := findMergeCells(data.Merges, int64(1), int64(col_count+1))
				if mergeRow != nil && mergeColumn != nil {
					parentData := i.RowData[*parentRow-1].Values[*parentColumn-1]
					spreadsheetHeaders[len(spreadsheetHeaders)-1] = fmt.Sprintf("%s [%s]", spreadsheetHeaders[len(spreadsheetHeaders)-1], intToLetters(col_count))
					spreadsheetHeaders = append(spreadsheetHeaders, fmt.Sprintf("%s [%s]", parentData.FormattedValue, intToLetters(col_count+1)))
				} else if value.FormattedValue == "" {
					columnName := intToLetters(col_count + 1) // since index in for is zero-based
					spreadsheetHeaders = append(spreadsheetHeaders, columnName)
				} else {
					if helpers.StringSliceContains(spreadsheetHeaders, value.FormattedValue) {
						spreadsheetHeaders = append(spreadsheetHeaders, fmt.Sprintf("%s [%s]", value.FormattedValue, intToLetters(6)))
					} else {
						spreadsheetHeaders = append(spreadsheetHeaders, value.FormattedValue)
					}
				}
			}
		}
		headersMap[data.Properties.Title] = spreadsheetHeaders
	}

	// cache the token source
	d.ConnectionManager.Cache.Set(cacheKey, headersMap)

	return headersMap, nil
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

	resp := svc.Spreadsheets.Get(spreadsheetID).IncludeGridData(true).Fields(googleapi.Field("sheets(properties.title,data(rowData),merges)"))
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
