package gin_server

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunargon/bolt-tui/src/bolt"
)

type ViewHandler struct {
	factory *bolt.Factory
	tmpl    *template.Template
}

func NewViewHandler(factory *bolt.Factory) *ViewHandler {
	// Load all templates
	tmpl := template.Must(template.ParseGlob("src/gin_server/views/*.html"))

	return &ViewHandler{
		factory: factory,
		tmpl:    tmpl,
	}
}

// Index serves the main application page
func (h *ViewHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{})
}

// DatabaseList serves the database list view
func (h *ViewHandler) DatabaseList(c *gin.Context) {
	// For now, return a simple list - this would need to be implemented in the factory
	databases := []string{"default"} // This should come from the factory

	c.HTML(http.StatusOK, "database_list.html", gin.H{
		"databases": databases,
	})
}

// BucketsView serves the buckets view for a specific database
func (h *ViewHandler) BucketsView(c *gin.Context) {
	dbName := c.Param("name")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name is required"})
		return
	}

	db, err := h.factory.Get(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	buckets, err := db.GetBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "buckets.html", gin.H{
		"database": dbName,
		"buckets":  buckets,
	})
}

// KeysView serves the keys view for a specific bucket
func (h *ViewHandler) KeysView(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	db, err := h.factory.Get(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	keys, err := db.GetKeysInBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "keys.html", gin.H{
		"database": dbName,
		"bucket":   bucketName,
		"keys":     keys,
	})
}

// KeyValueView serves the key value view
func (h *ViewHandler) KeyValueView(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if dbName == "" || bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name, bucket name, and key are required"})
		return
	}

	db, err := h.factory.Get(dbName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	value, err := db.GetValue(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "key_value.html", gin.H{
		"database": dbName,
		"bucket":   bucketName,
		"key":      key,
		"value":    string(value),
	})
}

// CreateBucketForm serves the create bucket form
func (h *ViewHandler) CreateBucketForm(c *gin.Context) {
	dbName := c.Query("database")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name is required"})
		return
	}

	c.HTML(http.StatusOK, "create_bucket_form.html", gin.H{
		"database": dbName,
	})
}

// EditBucketForm serves the edit bucket form
func (h *ViewHandler) EditBucketForm(c *gin.Context) {
	dbName := c.Query("database")
	bucketName := c.Query("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	c.HTML(http.StatusOK, "edit_bucket_form.html", gin.H{
		"database": dbName,
		"bucket":   bucketName,
	})
}

// CreateKeyForm serves the create key form
func (h *ViewHandler) CreateKeyForm(c *gin.Context) {
	dbName := c.Query("database")
	bucketName := c.Query("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	c.HTML(http.StatusOK, "create_key_form.html", gin.H{
		"database": dbName,
		"bucket":   bucketName,
	})
}

// EditKeyForm serves the edit key form
func (h *ViewHandler) EditKeyForm(c *gin.Context) {
	dbName := c.Query("database")
	bucketName := c.Query("bucket")
	key := c.Query("key")

	if dbName == "" || bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name, bucket name, and key are required"})
		return
	}

	c.HTML(http.StatusOK, "edit_key_form.html", gin.H{
		"database": dbName,
		"bucket":   bucketName,
		"key":      key,
	})
}
