package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lunargon/bolt-tui/src/bolt"
)

type BoltHandler struct {
	db *bolt.DB
}

func NewBoltHandler(db *bolt.DB) *BoltHandler {
	return &BoltHandler{db: db}
}

// GetBuckets returns all buckets in the database
func (h *BoltHandler) GetBuckets(c *gin.Context) {
	buckets, err := h.db.GetBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"buckets": buckets})
}

// GetKeysInBucket returns all keys in a specific bucket
func (h *BoltHandler) GetKeysInBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name is required"})
		return
	}

	keys, err := h.db.GetKeysInBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// GetValue returns the value for a specific key in a bucket
func (h *BoltHandler) GetValue(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name and key are required"})
		return
	}

	value, err := h.db.GetValue(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": string(value)})
}

// CreateBucket creates a new bucket
func (h *BoltHandler) CreateBucket(c *gin.Context) {
	var req struct {
		BucketName string `json:"bucket_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.CreateBucket(req.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "bucket created successfully"})
}

// DeleteBucket deletes a bucket
func (h *BoltHandler) DeleteBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name is required"})
		return
	}

	err := h.db.DeleteBucket(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bucket deleted successfully"})
}

// PutValue puts a value for a key in a bucket
func (h *BoltHandler) PutValue(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}

	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.PutValue(bucketName, req.Key, []byte(req.Value))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "value stored successfully"})
}

// DeleteValue deletes a key from a bucket
func (h *BoltHandler) DeleteValue(c *gin.Context) {
	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name and key are required"})
		return
	}

	err := h.db.DeleteValue(bucketName, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "key deleted successfully"})
}

// RenameBucket renames a bucket
func (h *BoltHandler) RenameBucket(c *gin.Context) {
	var req struct {
		NewName string `json:"new_name" binding:"required"`
	}

	bucketName := c.Param("bucket")
	if bucketName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.RenameBucket(bucketName, req.NewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bucket renamed successfully"})
}

// RenameKey renames a key within a bucket
func (h *BoltHandler) RenameKey(c *gin.Context) {
	var req struct {
		NewKey string `json:"new_key" binding:"required"`
	}

	bucketName := c.Param("bucket")
	key := c.Param("key")

	if bucketName == "" || key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bucket name and key are required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.db.RenameKey(bucketName, key, req.NewKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "key renamed successfully"})
}
