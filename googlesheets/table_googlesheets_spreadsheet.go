package googlesheets

import (
	"context"

	"google.golang.org/api/drive/v3"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableGoogleSheetsSpreadsheet(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "googlesheets_spreadsheet",
		Description: "Retrieve the metadata of given spreadsheet.",
		List: &plugin.ListConfig{
			Hydrate: listGoogleSheetSpreadsheet,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the spreadsheet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id"),
			},
			{
				Name:        "name",
				Description: "The name of the spreadsheet.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_time",
				Description: "The time at which the spreadsheet was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "drive_id",
				Description: "ID of the shared drive the spreadsheet resides in.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "copy_requires_writer_permission",
				Description: "Indicates whether the options to copy, print, or download this spreadsheet, should be disabled for readers and commenters.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "explicitly_trashed",
				Description: "Indicates whether the spreadsheet has been explicitly trashed, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "has_thumbnail",
				Description: "Indicates whether the spreadsheet has a thumbnail, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "icon_link",
				Description: "A static, unauthenticated link to the spreadsheets's icon.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "is_app_authorized",
				Description: "Indicates whether the spreadsheet was created or opened by the requesting app, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "modified_by_me",
				Description: "Indicates whether the spreadsheet has been modified by you, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "modified_by_me_time",
				Description: "The last time the spreadsheet was modified by the user.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_time",
				Description: "The last time the spreadsheet was modified by anyone.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "owned_by_me",
				Description: "Indicates whether the spreadsheet owns by you, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "quota_bytes_used",
				Description: "The number of storage quota bytes used by the spreadsheet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "shared",
				Description: "Indicates whether the spreadsheet has been shared, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "starred",
				Description: "Indicates whether the user has starred the spreadsheet, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "thumbnail_link",
				Description: "A short-lived link to the spreadsheet's thumbnail, if available.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "thumbnail_version",
				Description: "The thumbnail version for use in thumbnail cache invalidation.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "trashed",
				Description: "Indicates whether the spreadsheet has been trashed, either explicitly or from a trashed parent folder, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "trashed_time",
				Description: "The time that the item was trashed.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "version",
				Description: "A monotonically increasing version number for the spreadsheet.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "viewed_by_me",
				Description: "Indicates whether the spreadsheet has been viewed by this user, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "viewed_by_me_time",
				Description: "The last time the spreadsheet was viewed by the user.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "viewers_can_copy_content",
				Description: "Indicates whether the spreadsheet has been viewed by this user, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "web_view_link",
				Description: "A link for opening the spreadsheet in a relevant Google editor or viewer in a browser.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "writers_can_share",
				Description: "Indicates whether users with only writer permission can modify the spreadsheet's permissions, or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "capabilities",
				Description: "Specifies a set of capabilities the current user has on this spreadsheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "export_links",
				Description: "Links for exporting Docs Editors files to specific formats.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "last_modifying_user",
				Description: "Specifies the details of last user to modify the spreadsheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "link_share_metadata",
				Description: "Specifies details about the link URLs that clients are using to refer to this item.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "owners",
				Description: "Specifies the owner of this spreadsheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "permission_ids",
				Description: "A list of permission IDs for users with access to this spreadsheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "permissions",
				Description: "The full list of permissions for the spreadsheet.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "properties",
				Description: "A collection of arbitrary key-value pairs which are visible to all apps.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "spaces",
				Description: "The list of spaces which contain the spreadsheet.",
				Type:        proto.ColumnType_JSON,
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

func listGoogleSheetSpreadsheet(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create client
	opts, err := getSessionConfigStatic(ctx, d)
	if err != nil {
		return nil, err
	}

	// Create service
	svc, err := drive.NewService(ctx, opts...)
	if err != nil {
		plugin.Logger(ctx).Error("listGoogleSheetSpreadsheet", "connection_error", err)
		return nil, err
	}

	spreadsheetID := getSpreadsheetIDStatic(ctx, d)

	resp, err := svc.Files.Get(spreadsheetID).Fields("*").Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	d.StreamListItem(ctx, resp)

	return nil, nil
}
