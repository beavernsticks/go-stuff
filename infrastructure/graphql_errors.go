package shared_infrastructure

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	bsgostuff_domain "github.com/beavernsticks/go-stuff/domain"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GraphQLErrorHandler(ctx context.Context, err error) *gqlerror.Error {
	// Преобразуем ошибку в gqlerror
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)

	// Обрабатываем доменные ошибки

	switch {
	case bsgostuff_domain.IsNotFoundError(err):
		gqlErr.Extensions = map[string]interface{}{
			"code":       "NOT_FOUND",
			"httpStatus": http.StatusNotFound,
		}
	case bsgostuff_domain.IsForbiddenError(err):
		gqlErr.Extensions = map[string]interface{}{
			"code":       "FORBIDDEN",
			"httpStatus": http.StatusForbidden,
		}
	case bsgostuff_domain.IsInvalidArgumentError(err):
		gqlErr.Extensions = map[string]interface{}{
			"code":       "BAD_REQUEST",
			"httpStatus": http.StatusBadRequest,
		}
	case bsgostuff_domain.IsDuplicateError(err):
		gqlErr.Extensions = map[string]interface{}{
			"code":       "CONFLICT",
			"httpStatus": http.StatusConflict,
		}
	default:
		gqlErr.Extensions = map[string]interface{}{
			"code":       "INTERNAL_ERROR",
			"httpStatus": http.StatusInternalServerError,
		}
	}

	slog.ErrorContext(ctx, "GraphQL error",
		slog.Any("error", err),
		slog.Any("extensions", gqlErr.Extensions),
	)

	return gqlErr
}
