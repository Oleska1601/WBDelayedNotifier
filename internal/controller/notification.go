package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Oleska1601/WBDelayedNotifier/internal/controller/dto"
	"github.com/Oleska1601/WBDelayedNotifier/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
)

// GetNotificationStatusHandler godoc
// @Summary      get notification status
// @Description  get notification status by notification_id
// @Tags         notify
// @Accept       json
// @Produce      json
// @Param        notification_id   query      string  true  "notification ID"
// @Success		200					{object}	dto.GetNotificationResponse
// @Failure		400					{object}	map[string]string	"invalid notification_id"
// @Failure		500					{object}	map[string]string	"failed to get notification status"
// @Router 		/notify/{notification_id} [get]
func (s *Server) GetNotificationStatusHandler(c *gin.Context) {
	ctx := c.Request.Context()
	notificationIDStr := strings.TrimSpace(c.Query("notification_id"))
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil || notificationID <= 0 {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "GetNotificationStatusHandler strconv.ParseInt").
			Msg("invalid notification_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification_id"})
		return
	}

	notificationStatus, err := s.usecase.GetNotificationStatus(ctx, notificationID)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "GetNotificationStatusHandler s.usecase.GetNotificationStatus").
			Msg("failed to get notification status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification status"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "GetNotificationStatusHandler").
		Int64("notification_id", notificationID).
		Type("notification_status", notificationStatus).
		Msg("get notification status successful")

	c.JSON(http.StatusOK, dto.GetNotificationResponse{
		NotificationStatus: notificationStatus,
	})
}

// CreateNotificationHandler godoc
// @Summary create notification
// @Description create notification with provided params
// @Tags notify
// @Accept json
// @Produce json
// @Param        notification_id   query      string  true  "notification ID"
// @Success		200					{object}	dto.CreateNotificationResponse
// @Failure		400					{object}	map[string]string	"impossible to create notification"
// @Failure		500					{object}	map[string]string	"failed to create notification"
// @Router /notify [post]
func (s *Server) CreateNotificationHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var notificationRequest dto.CreateNotificationRequest
	if err := c.ShouldBindJSON(&notificationRequest); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "CreateNotificationHandler c.ShouldBindJSON").
			Msg("impossible to create notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to create notification"})
		return
	}
	notification, err := notificationRequest.ToModel()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "CreateNotificationHandler notificationRequest.ToModel").
			Msg("impossible to create notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to create notification"})
		return
	}
	notificationID, err := s.usecase.CreateNotification(ctx, notification)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "CreateNotificationHandler s.usecase.CreateNotification").
			Msg("failed to create notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "CreateNotificationHandler").
		Int64("notification_id", notificationID).
		Msg("create notification successful")

	c.JSON(http.StatusOK, dto.CreateNotificationResponse{
		NotificationID: notificationID,
	})
}

// DeleteNotificationHandler godoc
// @Summary delete notification
// @Description delete notification by notification_id
// @Tags notify
// @Accept json
// @Produce json
// @Param        notification_id   query      string  true  "notification ID"
// @Success		200					{string}	string "delete notification successful"
// @Failure		400					{object}	map[string]string	"invalid notification_id"
// @Failure		500					{object}	map[string]string	"failed to delete notification"
// @Router /notify/{notification_id} [delete]
func (s *Server) DeleteNotificationHandler(c *gin.Context) {
	ctx := c.Request.Context()
	notificationIDStr := strings.TrimSpace(c.Query("notification_id"))
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil || notificationID <= 0 {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "DeleteNotificationHandler strconv.ParseInt").
			Msg("invalid notification_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification_id"})
		return
	}

	err = s.usecase.UpdateNotificationStatus(ctx, notificationID, models.StatusCancelled)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "DeleteNotificationHandler s.usecase.UpdateNotificationStatus").
			Msg("failed to delete notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification"})
		return
	}
	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "DeleteNotificationHandler").
		Int64("notification_id", notificationID).
		Msg("delete notification successful")
	c.JSON(http.StatusOK, gin.H{"message": "delete notification successful"})
}
