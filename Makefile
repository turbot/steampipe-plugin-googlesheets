STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/googlesheets@latest/steampipe-plugin-googlesheets.plugin -tags "${BUILD_TAGS}" *.go
