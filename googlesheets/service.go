package googlesheets

import (
	"context"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func getSpreadsheetData(ctx context.Context, d *plugin.Plugin, sheetRange string) (*sheets.ValueRange, error) {
	// have we already created and cached the service?
	serviceCacheKey := "googlesheets" + sheetRange
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sheets.ValueRange), nil
	}
	
	// Create service
	opts := []option.ClientOption{}
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
