package endpoint

import (
	"context"
	"net/http"

	"github.com/charmingruby/pipo/service/api/internal/core/model"
	"github.com/charmingruby/pipo/service/api/internal/core/service"
	"github.com/gin-gonic/gin"
)

// IngestRawDataRequest is the request for the IngestRawData endpoint.
type IngestRawDataRequest struct {
	// FilePath is the path to the file to be ingested.
	FilePath string `json:"file_path" binding:"required"`
	// Records is the number of records to be ingested.
	Records int `json:"records"   binding:"required"`
}

// IngestRawDataResponse is the response for the IngestRawData endpoint.
type IngestRawDataResponse struct {
	// IngestedData is the data that was ingested.
	IngestedData []model.RawSentiment `json:"ingested_data"`
	// Errors is the errors that occurred during the ingestion.
	Errors []error `json:"errors"`
}

func (e *Endpoint) makeIngestRawDataEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req IngestRawDataRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		op, err := e.service.IngestRawData(context.Background(), service.IngestRawDataInput{
			FilePath: req.FilePath,
			Records:  req.Records,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"message": "Raw data ingested",
			"data": IngestRawDataResponse{
				IngestedData: op.IngestedData,
				Errors:       op.Errors,
			},
		})
	}
}
