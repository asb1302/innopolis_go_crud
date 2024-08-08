package handler

import (
	"crud/internal/domain"
	"crud/internal/pkg/authclient"
	"crud/internal/service"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
)

func ServerHandler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderAccessControlAllowOrigin, "*")
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodPost)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodGet)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowMethods, fasthttp.MethodDelete)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderContentType)
	ctx.Response.Header.Add(fasthttp.HeaderAccessControlAllowHeaders, fasthttp.HeaderAuthorization)

	if ctx.IsOptions() {
		return
	}

	if string(ctx.Path()) != "/ping" {
		token := ctx.Request.Header.Peek(fasthttp.HeaderAuthorization)
		if string(token) == "" || !authclient.ValidateToken(string(token)) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			log.Println("Request", string(ctx.Method()), "Unauthorized, token:", string(token))
			return
		}
	}

	switch string(ctx.Path()) {
	case "/ping":
		PingHandler(ctx)
	case "/count":
		CountHandler(ctx)
	default:
		handleRecipes(ctx)
	}
}

func PingHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)

	ctx.SetBodyString("pong")
}

func CountHandler(ctx *fasthttp.RequestCtx) {
	count, err := service.GetRecipeCount()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("Error fetching recipe count:", err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, err = ctx.Write([]byte(strconv.Itoa(count)))
	if err != nil {
		log.Println("Error writing response:", err)
	}
}

func handleRecipes(ctx *fasthttp.RequestCtx) {
	if ctx.IsGet() && ctx.QueryArgs().Has("page") {
		PaginatedHandler(ctx)
	} else if ctx.IsGet() {
		GetHandler(ctx)
	} else if ctx.IsDelete() {
		DeleteHandler(ctx)
	} else if ctx.IsPost() {
		PostHandler(ctx)
	} else {
		ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	}
}

func GetHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	rec, err := service.Get(string(id))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	marshal, err := json.Marshal(rec)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func DeleteHandler(ctx *fasthttp.RequestCtx) {
	id := ctx.QueryArgs().Peek("id")
	if len(id) == 0 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if err := service.Delete(string(id)); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func PostHandler(ctx *fasthttp.RequestCtx) {
	var rec domain.Recipe
	log.Println(string(ctx.PostBody()))
	if err := json.Unmarshal(ctx.PostBody(), &rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	if err := service.AddOrUpd(&rec); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	resp := IdResponse{ID: rec.ID}

	marshal, err := json.Marshal(resp)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

const (
	defaultLimit = 10
	maxLimit     = 10
)

func PaginatedHandler(ctx *fasthttp.RequestCtx) {
	pageStr := ctx.QueryArgs().Peek("page")
	limitStr := ctx.QueryArgs().Peek("limit")

	page, err := strconv.Atoi(string(pageStr))
	if err != nil || page < 1 {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBodyString("Invalid page number")

		return
	}

	limit, err := strconv.Atoi(string(limitStr))
	if err != nil || limit < 1 {
		limit = defaultLimit
	} else if limit > maxLimit {
		limit = maxLimit
	}

	records, err := service.GetPaginated(page, limit)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("Error getting paginated records:", err)

		return
	}

	marshal, err := json.Marshal(records)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("Error converting to json:", err)

		return
	}

	if _, err = ctx.Write(marshal); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Println("Error writing response:", err)

		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
