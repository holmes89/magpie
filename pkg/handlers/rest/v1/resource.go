package v1

import (
	"context"

	"github.com/holmes89/magpie/magpie"
)

type ResourceService interface {
	Get(context.Context, string) (magpie.Resource, error)
}

type resourceHandler struct {
}
