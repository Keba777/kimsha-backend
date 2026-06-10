package pagination

import "github.com/gofiber/fiber/v2"

type Meta struct {
	Page    int   `json:"page"`
	Limit   int   `json:"limit"`
	Total   int64 `json:"total"`
	Pages   int   `json:"pages"`
}

type Params struct {
	Page  int
	Limit int
}

func Parse(c *fiber.Ctx) Params {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return Params{Page: page, Limit: limit}
}

func (p Params) Offset() int { return (p.Page - 1) * p.Limit }

func NewMeta(p Params, total int64) Meta {
	pages := int(total) / p.Limit
	if int(total)%p.Limit > 0 {
		pages++
	}
	return Meta{Page: p.Page, Limit: p.Limit, Total: total, Pages: pages}
}
