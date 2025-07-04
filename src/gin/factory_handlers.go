package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunargon/bolt-tui/src/bolt"
)

type FactoryHandler struct {
	factory *bolt.Factory
}

func NewFactoryHandler(factory *bolt.Factory) *FactoryHandler {
	return &FactoryHandler{factory: factory}
}

// GetDatabases returns all databases managed by the factory
func (h *FactoryHandler) GetDatabases(c *gin.Context) {
	databases, err := h.factory.GetDatabases()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"databases": databases})
}

func (h *FactoryHandler) getDatabase(c *gin.Context, name string) *bolt.DB {
	db, err := h.factory.Get(name)
	if err != nil {
		if db == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "database not found"})
			return nil
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil
	}

	return db
}

// GetBuckets returns all buckets in a specific database
func (h *FactoryHandler) GetBuckets(c *gin.Context) {
	dbName := c.Param("name")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name is required"})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	buckets, err := db.GetBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}

// GetKeysInBucket returns all keys in a specific bucket
func (h *FactoryHandler) GetKeysInBucket(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	keys, err := db.GetKeysInBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// GetValue returns the value for a specific key in a bucket
func (h *FactoryHandler) GetValue(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if dbName == "" || bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name, bucket name, and key are required"})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	value, err := db.GetValue(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": string(value)})
}

// CreateBucket creates a new bucket in a database
func (h *FactoryHandler) CreateBucket(c *gin.Context) {
	var req struct {
		BucketName string `json:"bucket_name" binding:"required"`
	}

	dbName := c.Param("name")
	if dbName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.CreateBucket(req.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "bucket created successfully"})
}

// DeleteBucket deletes a bucket from a database
func (h *FactoryHandler) DeleteBucket(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.DeleteBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bucket deleted successfully"})
}

// PutValue puts a value for a key in a bucket
func (h *FactoryHandler) PutValue(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	dbName := c.Param("name")
	bucketName := c.Param("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.PutValue(bucketName, req.Key, []byte(req.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "value stored successfully"})
}

// DeleteValue deletes a key from a bucket
func (h *FactoryHandler) DeleteValue(c *gin.Context) {
	dbName := c.Param("name")
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if dbName == "" || bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name, bucket name, and key are required"})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.DeleteValue(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "key deleted successfully"})
}

// RenameBucket renames a bucket in a database
func (h *FactoryHandler) RenameBucket(c *gin.Context) {
	var req struct {
		NewName string `json:"new_name" binding:"required"`
	}

	dbName := c.Param("name")
	bucketName := c.Param("bucket")

	if dbName == "" || bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name and bucket name are required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.RenameBucket(bucketName, req.NewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bucket renamed successfully"})
}

// RenameKey renames a key within a bucket
func (h *FactoryHandler) RenameKey(c *gin.Context) {
	var req struct {
		NewKey string `json:"new_key" binding:"required"`
	}

	dbName := c.Param("name")
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if dbName == "" || bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "database name, bucket name, and key are required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.getDatabase(c, dbName)
	if db == nil {
		return
	}

	err := db.RenameKey(bucketName, key, req.NewKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "key renamed successfully"})
}
