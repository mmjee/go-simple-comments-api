package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	urllib "net/url"
	"strconv"
)

// Comment have this general Structure
type Comment struct {
	SiteID uint64
	PageURL string
	Time SCATime

	// If PGPSigned is true, then the text is likely signed, but not guaranteed to be signed. Clients will check for themselves.
	PGPSigned bool
	Text string
}

func (api *simpleCommentsAPI) createCommentForURL (ctx *gin.Context) {
	reqCtx := context.Background()
	var comment Comment
	err := ctx.BindJSON(&comment)
	if err != nil {
		return
	}

	// Validate the object itself
	{
		url, err := urllib.Parse(comment.PageURL)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidComment, "PageURL Not A URL"))
			return
		}

		if !url.IsAbs() {
			ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidComment, "PageURL not absolute"))
			return
		}
	}

	// Insert it to the database
	_, err = api.comments.InsertOne(reqCtx, comment)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, api.constructError(ErrServerFailure, "Couldn't insert comment to database."))
		return
	}
	ctx.JSON(http.StatusOK, true)
}

func (api *simpleCommentsAPI) getCommentsForURL(ctx *gin.Context) {
	reqCtx := context.Background()

	url := ctx.Query("url")
	idStr := ctx.Param("id")
	if len(url) == 0 || len(idStr) == 0 {
		ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidQuery, "Either URL or ID missing."))
		return
	}

	id, err := strconv.ParseUint(idStr, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidQuery, "Couldn't parse id"))
		return
	}

	pageStr, limitStr := ctx.DefaultQuery("page", "1"), ctx.DefaultQuery("limit", "50")

	// Wonder what people will do with 2^64 pages of comments, and good luck with sending hexadecimal or binary base queries
	page, err := strconv.ParseUint(pageStr, 0, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidPagination, "Couldn't parse page"))
		return
	}

	// Good luck with sending hexadecimal or binary base queries, but 255 is the max limit so none can potentially DDoS
	// the server by sending 2^64 limit queries that thrash the CPU
	limit, err := strconv.ParseUint(limitStr, 0, 8)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, api.constructError(ErrInvalidPagination, "Couldn't parse limit"))
		return
	}

	findOpts := options.Find().SetLimit(int64(limit)).SetSkip(int64(page * limit))
	cursor, err := api.comments.Find(reqCtx, bson.M{
		"SiteID": id,
		"PageURL": url,
	}, findOpts)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, api.constructError(ErrServerFailure, "Couldn't acquire documents from MongoDB"))
		return
	}

	var results []bson.M
	if err = cursor.All(reqCtx, &results); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, api.constructError(ErrServerFailure, "Couldn't iterate through documents"))
		return
	}

	ctx.JSON(http.StatusOK, results)
	return
}
