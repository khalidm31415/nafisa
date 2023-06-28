package elasticsarch_helper

import (
	"backend/dto"
	"backend/entity"
	"backend/internal_constant"
	"backend/package_helper/embeddings_helper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type IElasticsearchProfileIndex interface {
	CreateIndexIfNotExists(ctx context.Context) error
	Index(ctx context.Context, profile entity.UserProfile) error
	GetMatchingProfiles(ctx context.Context, profile entity.UserProfile) ([]ElasticSearchProfile, error)
}

type ElasticsearchProfileIndex struct {
	es         *elasticsearch.Client
	embeddings embeddings_helper.IEmbeddings
}

func NewElasticsearchProfileIndex(cfg elasticsearch.Config, embeddings embeddings_helper.IEmbeddings) IElasticsearchProfileIndex {
	// 1. Get cluster info
	r := map[string]interface{}{}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
	return &ElasticsearchProfileIndex{es: es, embeddings: embeddings}
}

func (e *ElasticsearchProfileIndex) CreateIndexIfNotExists(ctx context.Context) error {
	req := esapi.IndicesExistsRequest{
		Index: []string{internal_constant.ProfileIndexName},
	}
	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		// Handle error
		log.Fatalf("Error: %s", err.Error())
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		req2 := esapi.IndicesCreateRequest{
			Index: internal_constant.ProfileIndexName,
			Body:  strings.NewReader(Index),
		}
		res, err = req2.Do(context.Background(), e.es)
		if err != nil {
			// Handle error
			log.Fatalf("Error: %s", err.Error())
		}
		defer res.Body.Close()

		if res.IsError() {
			// Handle error response
			log.Fatalf("Error: %s", res.String())
		} else {
			fmt.Println("Index created successfully")
		}
	}
	return nil
}

func (e *ElasticsearchProfileIndex) Index(ctx context.Context, profile entity.UserProfile) error {
	profileIndex, err := dto.NewProfileIndex(profile)
	if err != nil {
		return err
	}
	summaryDenseVector, err := e.embeddings.Embed(context.Background(), profile.Summary)
	if err != nil {
		return err
	}
	profileIndex.SummaryDenseVector = summaryDenseVector

	payload, err := json.Marshal(profileIndex)
	if err != nil {
		return err
	}

	res, err := esapi.CreateRequest{
		Index:      internal_constant.ProfileIndexName,
		DocumentID: profileIndex.UserID,
		Body:       bytes.NewReader(payload),
	}.Do(ctx, e.es.Transport)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			fmt.Printf("Error: %s", err.Error())
			return err
		}
		fmt.Printf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}
	return nil
}

func (e *ElasticsearchProfileIndex) GetMatchingProfiles(ctx context.Context, profile entity.UserProfile) ([]ElasticSearchProfile, error) {
	embeddingSummary, err := e.embeddings.Embed(ctx, profile.Summary)
	if err != nil {
		return nil, err
	}
	embeddingPreferencePartnerCriteria, err := e.embeddings.Embed(ctx, profile.PreferencePartnerCriteria)
	if err != nil {
		return nil, err
	}

	genderFilter := "f"
	if profile.Gender == "f" {
		genderFilter = "m"
	}

	now := time.Now()
	currentYear, _, _ := now.Date()

	// Create the Elasticsearch query based on the input profile
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"script_score": map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"filter": []map[string]interface{}{
							{
								"term": map[string]string{
									"gender": genderFilter,
								},
							},
							{
								"range": map[string]interface{}{
									"year_born": map[string]int{
										"gte": currentYear - profile.PreferenceMaxAge,
										"lte": currentYear - profile.PreferenceMinAge,
									},
								},
							},
							{
								"terms": map[string][]string{
									"last_education": getLastEducationTerms(profile.PreferenceMinLastEducation),
								},
							},
						},
					},
				},
				"script": map[string]interface{}{
					"source": "cosineSimilarity(params.summary_vector, 'summary_dense_vector')  + cosineSimilarity(params.partner_criteria_vector, 'summary_dense_vector') + 2.0",
					"params": map[string]interface{}{
						"summary_vector":          embeddingSummary,
						"partner_criteria_vector": embeddingPreferencePartnerCriteria,
					},
				},
			},
		},
		"sort": []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"_source": map[string]interface{}{
			"excludes": []string{"summary_dense_vector"},
		},
	}

	// Convert the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Send the search request to Elasticsearch
	req := esapi.SearchRequest{
		Index: []string{"profiles"}, // Elasticsearch index name
		Body:  bytes.NewReader(queryJSON),
	}

	res, err := req.Do(context.Background(), e.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response JSON
	var response ElasticSearchResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Hits.Hits, nil
}
