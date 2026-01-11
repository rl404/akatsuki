package api

import (
	"context"

	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// GetAnimeByID to get anime by id.
func (api *API) GetAnimeByID(ctx context.Context, req *schema.GetAnimeByIDRequest) (*schema.GetAnimeByIDResponse, error) {
	anime, code, err := api.service.GetAnimeByID(ctx, req.GetId())
	if err != nil {
		return nil, utils.ResponseWithGRPC(code, stack.Wrap(ctx, err))
	}
	return &schema.GetAnimeByIDResponse{
		Data: api.animeFromService(anime),
	}, nil
}
