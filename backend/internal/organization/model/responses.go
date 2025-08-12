package organizationmodel

type OrganizationResponse struct {
	Organization *Organization `json:"organization"`
}

type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Total         int64          `json:"total"`
}

type OrganizationStatsResponse struct {
	TotalOrganizations int64 `json:"total_organizations"`
	ActiveOrganizations int64 `json:"active_organizations"`
}