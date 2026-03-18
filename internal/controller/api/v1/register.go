package v1

import "github.com/wb-go/wbf/ginext"

const (
	V1GroupURI = "/v1"
)

type APIV1 struct {
	service ServiceI
}

func New(service ServiceI) *APIV1 {
	return &APIV1{
		service: service,
	}
}

func (v1 *APIV1) RegisterHandlers(group *ginext.RouterGroup) {
	v1Group := group.Group(V1GroupURI)

	v1.registerNotificationHandlers(v1Group)
}
