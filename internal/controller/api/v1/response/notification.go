package response

import "github.com/Oleska1601/WBDelayedNotifier/internal/models"

type GetNotificationStatusResponse struct {
	Status models.Status `json:"status"`
}

func ToGetNotificationStatusResponse(status models.Status) *GetNotificationStatusResponse {
	return &GetNotificationStatusResponse{
		Status: status,
	}
}

type CreateNotificationResponse struct {
	ID int `json:"id"`
}

func ToCreateNotificationResponse(id int) *CreateNotificationResponse {
	return &CreateNotificationResponse{
		ID: id,
	}
}
