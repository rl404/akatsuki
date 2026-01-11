package api

import (
	"context"

	"github.com/rl404/akatsuki/internal/delivery/grpc/schema"
	"github.com/rl404/akatsuki/internal/utils"
	"github.com/rl404/fairy/errors/stack"
)

// GetUserAnimeRelations to get user anime relations.
func (api *API) GetUserAnimeRelations(ctx context.Context, req *schema.GetUserAnimeRelationsRequest) (*schema.GetUserAnimeRelationsResponse, error) {
	relations, code, err := api.service.GetUserAnimeRelations(ctx, req.GetUsername())
	if err != nil {
		return nil, utils.ResponseWithGRPC(code, stack.Wrap(ctx, err))
	}
	return &schema.GetUserAnimeRelationsResponse{
		Data: api.userAnimeRelationFromService(relations),
	}, nil
}
