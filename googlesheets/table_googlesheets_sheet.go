package googlesheets

import (
	"context"

	"google.golang.org/api/sheets/v4"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableGooglesheetsSheet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesheets_sheet",
		Description: "Retrieve the sheet in a given spreadsheet.",
		List: &plugin.ListConfig{
			Hydrate: listGooglesheetSheets,
		},
		Columns: []*plugin.Column{
			{
				Name:        "spreadsheet_id",
				Description: "The ID of the spreadsheet.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     spreadsheetID,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "title",
				Description: "The name of the sheet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.Title"),
			},
			{
				Name:        "sheet_id",
				Description: "The ID of the sheet.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.SheetId"),
			},
			{
				Name:        "sheet_type",
				Description: "The type of sheet. Defaults to GRID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Properties.SheetType"),
			},
			{
				Name:        "hidden",
				Description: "Indicates whether the sheet is hidden, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.Hidden"),
			},
			{
				Name:        "index",
				Description: "The index of the sheet within the spreadsheet.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Properties.Index"),
			},
			{
				Name:        "right_to_left",
				Description: "Indicates whether sheet is an RTL sheet instead of an LTR sheet.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Properties.RightToLeft"),
			},
			{
				Name:        "banded_ranges",
				Description: "The banded (alternating colors) ranges on this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "basic_filter",
				Description: "The filter on this sheet, if any.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "charts",
				Description: "The specifications of every chart on this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "column_groups",
				Description: "All column groups on this sheet, ordered by increasing range start index, then by group depth.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "conditional_formats",
				Description: "The conditional format rules in this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data",
				Description: "Data in the grid, if this is a grid sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "data_source_sheet_properties",
				Description: "Specifies the properties specific to the DATA_SOURCE sheet.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.DataSourceSheetProperties"),
			},
			{
				Name:        "developer_metadata",
				Description: "The developer metadata associated with a sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "filter_views",
				Description: "The filter views in this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "grid_properties",
				Description: "Additional properties of the sheet if this sheet is a grid.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.GridProperties"),
			},
			{
				Name:        "merges",
				Description: "The ranges that are merged together.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "protected_ranges",
				Description: "The protected ranges in this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "row_groups",
				Description: "All row groups on this sheet, ordered by increasing range start index, then by group depth.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "slicers",
				Description: "The slicers on this sheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "tab_color",
				Description: "The color of the tab in the UI.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TabColor"),
			},
			{
				Name:        "tab_color_style",
				Description: "The color of the tab in the UI.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Properties.TabColorStyle"),
			},
		},
	}
}

//// LIST FUNCTION

func listGooglesheetSheets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

	spreadsheetID := getSpreadsheetID(ctx, d.Table.Plugin)

	resp, err := svc.Spreadsheets.Get(spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	if resp.Sheets != nil {
		for _, sheet := range resp.Sheets {
			d.StreamListItem(ctx, sheet)
		}
	}

	return nil, nil
}

func spreadsheetID(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	spreadsheetID := getSpreadsheetID(ctx, d.Table.Plugin)
	return spreadsheetID, nil
}
