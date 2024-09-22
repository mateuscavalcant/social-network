package models

type SearchRequest struct {
	Search string `json:"search" binding:"required"`
}
