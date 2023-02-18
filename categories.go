package helix

// SearchCategoriesParams is parameters for SearchCategories
type SearchCategoriesParams struct {
	Query string `query:"query"`
	After string `query:"after"`
	First int    `query:"first,20"` // Limit 100
}

// ManySearchCategories is the response data from SearchCategories
type ManySearchCategories struct {
	Categories []Category `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// Category describes a category from SearchCategory
type Category struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
}

// SearchCategoriesResponse is the response from SearchCategories
type SearchCategoriesResponse struct {
	ResponseCommon
	Data ManySearchCategories
}

// SearchCategories searches for Twitch categories based on the given search query
func (c *Client) SearchCategories(params *SearchCategoriesParams) (*SearchCategoriesResponse, error) {
	resp, err := c.get("/search/categories", &ManySearchCategories{}, params)
	if err != nil {
		return nil, err
	}

	categories := &SearchCategoriesResponse{}
	resp.HydrateResponseCommon(&categories.ResponseCommon)
	categories.Data.Categories = resp.Data.(*ManySearchCategories).Categories
	categories.Data.Pagination = resp.Data.(*ManySearchCategories).Pagination

	return categories, nil
}
