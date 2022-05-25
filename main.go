package main

import (
	"github.com/angelorc/sinfonia-indexer/config"
	"github.com/angelorc/sinfonia-indexer/db"
	"github.com/angelorc/sinfonia-indexer/server"
	"github.com/angelorc/sinfonia-indexer/utility"
)

func main() {
	/*
	 * Start the program
	 */
	utility.ShowLogo()

	/**
	 * Connect to db
	 */
	defaultDB := db.Database{
		DataBaseRefName: "sinfonia-indexer",
		URL:             config.GetSecret("MONGO_URI"),
		DataBaseName:    config.GetSecret("MONGO_DATABASE"),
		RetryWrites:     config.GetSecret("MONGO_RETRYWRITES"),
	}
	defaultDB.Init()
	defer defaultDB.Disconnect()

	/**
	 * Create custom mongo indexes
	 */

	/*
	 * Start http server
	 */
	server.New()
}
