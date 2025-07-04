package gin_server

import (
	"github.com/gin-gonic/gin"
	"github.com/lunargon/bolt-tui/src/bolt"
)

// SetupRoutes configures all the routes for the BoltDB API and web interface
func SetupRoutes(r *gin.Engine, factory *bolt.Factory) {
	// Create handlers
	factoryHandler := NewFactoryHandler(factory)
	viewHandler := NewViewHandler(factory)

	// Set HTML template renderer
	r.LoadHTMLGlob("src/gin_server/views/*.html")

	// View routes
	r.GET("/", viewHandler.Index)
	r.GET("/databases", viewHandler.DatabaseList)
	r.GET("/databases/:name/buckets", viewHandler.BucketsView)
	r.GET("/databases/:name/buckets/:bucket/keys", viewHandler.KeysView)
	r.GET("/databases/:name/buckets/:bucket/keys/:key", viewHandler.KeyValueView)
	r.GET("/views/create-bucket-form", viewHandler.CreateBucketForm)
	r.GET("/views/edit-bucket-form", viewHandler.EditBucketForm)
	r.GET("/views/create-key-form", viewHandler.CreateKeyForm)
	r.GET("/views/edit-key-form", viewHandler.EditKeyForm)

	// API routes
	api := r.Group("/api/v1")
	{
		// Database management routes
		api.GET("/databases", factoryHandler.GetDatabases)

		// Bucket operations
		api.GET("/databases/:name/buckets", factoryHandler.GetBuckets)
		api.POST("/databases/:name/buckets", factoryHandler.CreateBucket)
		api.DELETE("/databases/:name/buckets/:bucket", factoryHandler.DeleteBucket)
		api.PUT("/databases/:name/buckets/:bucket", factoryHandler.RenameBucket)

		// Key-value operations
		api.GET("/databases/:name/buckets/:bucket/keys", factoryHandler.GetKeysInBucket)
		api.GET("/databases/:name/buckets/:bucket/keys/:key", factoryHandler.GetValue)
		api.POST("/databases/:name/buckets/:bucket/keys", factoryHandler.PutValue)
		api.DELETE("/databases/:name/buckets/:bucket/keys/:key", factoryHandler.DeleteValue)
		api.PUT("/databases/:name/buckets/:bucket/keys/:key", factoryHandler.RenameKey)
	}
}

// SetupSingleDBRoutes configures routes for a single database (legacy support)
func SetupSingleDBRoutes(r *gin.Engine, db *bolt.DB) {
	// Create handler for single database
	boltHandler := NewBoltHandler(db)

	// API routes for single database
	api := r.Group("/api/v1/db")
	{
		// Bucket operations
		api.GET("/buckets", boltHandler.GetBuckets)
		api.POST("/buckets", boltHandler.CreateBucket)
		api.DELETE("/buckets/:bucket", boltHandler.DeleteBucket)
		api.PUT("/buckets/:bucket", boltHandler.RenameBucket)

		// Key-value operations
		api.GET("/buckets/:bucket/keys", boltHandler.GetKeysInBucket)
		api.GET("/buckets/:bucket/keys/:key", boltHandler.GetValue)
		api.POST("/buckets/:bucket/keys", boltHandler.PutValue)
		api.DELETE("/buckets/:bucket/keys/:key", boltHandler.DeleteValue)
		api.PUT("/buckets/:bucket/keys/:key", boltHandler.RenameKey)
	}
}
