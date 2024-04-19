package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Key string `form:"key" binding:"required"`
}

type Response struct {
	Data       []map[string]interface{} `json:"data"`
	RequestKey string                   `json:"request_key"`
	Error      string                   `json:"error"`
}

// @Summary Get data
// @Description Get data
// @Tags sample
// @Accept json
// @Body {object} Request
// @Produce json
// @Success 200 {object} Response
// @Router /sample [get]
func SampleHandler(c *gin.Context, ENV_VALUE string) {
	// define response struct
	var res Response

	// bind request to Request struct
	var req Request
	if err := c.ShouldBindQuery(&req); err != nil {
		res.Error = fmt.Sprintf("Invalid request%s", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Do something to get data
	res.Data = append(res.Data, map[string]interface{}{"ENV_KEY": ENV_VALUE})
	res.RequestKey = req.Key

	c.JSON(http.StatusOK, res)

}
