package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Oleska1601/WBDelayedNotifier/internal/controller/api/v1/request"
	"github.com/Oleska1601/WBDelayedNotifier/internal/controller/api/v1/response"
	"github.com/Oleska1601/WBDelayedNotifier/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/zlog"
)

const (
	notificationsGroupURI = "/notify"
)

const (
	getNotificationStatusURI = "/:id"
	createNotificationURI    = ""
	deleteNotificationURI    = "/:id"
)

func (v1 *APIV1) registerNotificationHandlers(group *gin.RouterGroup) {
	notificationsGroup := group.Group(notificationsGroupURI)

	notificationsGroup.GET(getNotificationStatusURI, v1.getNotificationStatus)
	notificationsGroup.POST(createNotificationURI, v1.createNotification)
	notificationsGroup.DELETE(deleteNotificationURI, v1.deleteNotification)
}

// @Summary	Получение статуса уведомления
// @Description	Получение статуса конкретного уведомления
// @Tags NOTIFICATIONS API
// @Produce	json
// @Param id path integer true "ID уведомления"
// @Success	200	{object} response.GetNotificationStatusResponse
// @Failure	400	{object} map[string]string "Ошибка валидации"
// @Failure	404	{object} map[string]string "Уведомление не найдено"
// @Failure	500	{object} map[string]string "Ошибка сервера"
// @Router	/api/v1/notify/{id} [get]
func (v1 *APIV1) getNotificationStatus(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "getNotificationStatus strconv.Atoi").
			Msg("impossible to get notification status")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to get notification status"})
		return
	}

	status, err := v1.service.GetNotificationStatus(ctx, id)
	if err != nil {
		if errors.Is(err, errs.NotFoundError) {
			zlog.Logger.Error().
				Err(err).
				Int("status", http.StatusNotFound).
				Str("path", "getNotificationStatus v1.service.GetNotificationStatus").
				Int("id", id).
				Msg("impossible to get notification status")
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("impossible to get notification status: %v", err.Error())})
			return
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "getNotificationStatus v1.service.GetNotificationStatus").
			Int("id", id).
			Msg("failed to get notification status")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification status"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "getNotificationStatus").
		Int("id", id).
		Msg("get notification status successful")
	c.JSON(http.StatusOK, response.ToGetNotificationStatusResponse(status))
}

// @Summary Создание уведомления
// @Description	Создание нового уведомления
// @Tags NOTIFICATIONS API
// @Accept json
// @Produce	json
// @Param req body request.CreateNotificationRequest true "Данные для создания уведомления"
// @Success	201	{object} response.CreateNotificationResponse
// @Failure	400	{object} map[string]string "Ошибка валидации"
// @Failure	500	{object} map[string]string "Ошибка сервера"
// @Router /api/v1/notify [post]
func (v1 *APIV1) createNotification(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "createNotification c.ShouldBindJSON").
			Msg("impossible to create notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to create notification"})
		return
	}

	notification, err := req.ToModel()
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "createNotification req.ToModel").
			Msg("impossible to create notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("impossible to create notification: %v", err.Error())})
		return
	}

	id, err := v1.service.CreateNotification(ctx, notification)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusInternalServerError).
			Str("path", "createNotification v1.service.CreateNotification").
			Msg("failed to create notification")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create notification"})
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusCreated).
		Str("path", "createNotification").
		Int("id", id).
		Msg("create notification successful")
	c.JSON(http.StatusCreated, response.ToCreateNotificationResponse(id))

}

// @Summary Удаление (отмена) уведомления
// @Description	Удаление (отмена) существующего уведомления
// @Tags NOTIFICATIONS API
// @Security
// @Param id path integer true "ID уведомления"
// @Success	200
// @Failure	400	{object} map[string]string "Ошибка валидации"
// @Failure	404	{object} map[string]string "Уведомление не существует"
// @Failure	409	{object} map[string]string "Ошибка бизнес логики"
// @Failure	500	{object} map[string]string "Ошибка сервера"
// @Router /api/v1/notify/{id} [delete]
func (v1 *APIV1) deleteNotification(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zlog.Logger.Error().
			Err(err).
			Int("status", http.StatusBadRequest).
			Str("path", "deleteNotification strconv.Atoi").
			Msg("impossible to delete notification")
		c.JSON(http.StatusBadRequest, gin.H{"error": "impossible to delete notification"})
		return
	}

	if err := v1.service.DeleteNotification(ctx, id); err != nil {
		var status int
		msg := "impossible to delete notification"
		switch {
		case errors.Is(err, errs.NotFoundError):
			status = http.StatusNotFound
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("impossible to delete notification: %v", err.Error())})
		case errors.Is(err, errs.ConflictError):
			status = http.StatusConflict
			c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("impossible to delete notification: %v", err.Error())})
		default:
			status = http.StatusInternalServerError
			msg = "failed to delete notification"
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete notification"})
		}

		zlog.Logger.Error().
			Err(err).
			Int("status", status).
			Str("path", "deleteNotification v1.service.DeleteNotification").
			Int("id", id).
			Msg(msg)
		return
	}

	zlog.Logger.Info().
		Int("status", http.StatusOK).
		Str("path", "deleteNotification").
		Int("id", id).
		Msg("delete notification successful")
	c.Status(http.StatusOK)
}
