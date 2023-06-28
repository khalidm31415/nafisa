package elasticsarch_helper

type ElasticSearchProfile struct {
	Score  float32 `json:"_score"`
	Source struct {
		UserID        string `json:"user_id"`
		YearBorn      int    `json:"year_born"`
		Gender        string `json:"gender"`
		LastEducation string `json:"last_education"`
		Summary       string `json:"summary"`
	} `json:"_source"`
}

type ElasticSearchResponse struct {
	Hits struct {
		Hits []ElasticSearchProfile `json:"hits"`
	} `json:"hits"`
}
