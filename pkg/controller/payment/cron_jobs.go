package payment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vesicash/payment-ms/cronjobs"
	"github.com/vesicash/payment-ms/utility"
)

func (base *Controller) StartCronJob(c *gin.Context) {
	var (
		req cronjobs.StartCronJobRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	if req.IntervalNumber != 0 && req.IntervalBase != "" {
		err := cronjobs.UpdateCronJobInterval(base.ExtReq, base.Db, req.Name, req.IntervalNumber, req.IntervalBase)
		if err != nil {
			rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
			c.JSON(http.StatusBadRequest, rd)
			return
		}
	}

	cronjobs.StartCronJob(base.ExtReq, base.Db, req.Name)

	rd := utility.BuildSuccessResponse(http.StatusOK, "started cron job", nil)
	c.JSON(http.StatusOK, rd)

}
func (base *Controller) StartCronJobsBulk(c *gin.Context) {
	var (
		reqSlice struct {
			Jobs []cronjobs.StartCronJobRequest `json:"jobs" validate:"required"`
		}
	)

	err := c.ShouldBind(&reqSlice)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&reqSlice)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	for _, req := range reqSlice.Jobs {
		if req.IntervalNumber != 0 && req.IntervalBase != "" {
			err := cronjobs.UpdateCronJobInterval(base.ExtReq, base.Db, req.Name, req.IntervalNumber, req.IntervalBase)
			if err != nil {
				rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
				c.JSON(http.StatusBadRequest, rd)
				return
			}
		}

		cronjobs.StartCronJob(base.ExtReq, base.Db, req.Name)
	}

	rd := utility.BuildSuccessResponse(http.StatusOK, "started cron jobs", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) StopCronJob(c *gin.Context) {
	var (
		req struct {
			Name string `json:"name"  validate:"required"`
		}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	go cronjobs.StopCronJob(req.Name)

	rd := utility.BuildSuccessResponse(http.StatusOK, "stopped cron job", nil)
	c.JSON(http.StatusOK, rd)

}

func (base *Controller) UpdateCronJobInterval(c *gin.Context) {
	var (
		req cronjobs.UpdateCronJobRequest
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = cronjobs.UpdateCronJobInterval(base.ExtReq, base.Db, req.Name, req.IntervalNumber, req.IntervalBase)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	go cronjobs.RestartCronJob(base.ExtReq, base.Db, req.Name)

	rd := utility.BuildSuccessResponse(http.StatusOK, "updated", nil)
	c.JSON(http.StatusOK, rd)

}
