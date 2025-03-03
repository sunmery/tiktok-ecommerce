package data

import (
	"backend/application/product/internal/biz"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/google/uuid"
)

func (d *Data) CreateProduct(ctx context.Context, product *biz.CreateProductRequest, productId uuid.UUID) error {
	now := time.Now().Format(time.RFC3339)
	image := ""
	if product.Images == nil {
		image = ""
	} else {
		image = product.Images[0].URL
	}

	doc := map[string]interface{}{
		"id":            productId.String(),
		"name":          product.Name,
		"name_suggest":  product.Name,
		"description":   product.Description,
		"price":         product.Price,
		"status":        product.Status,
		"merchant_id":   product.MerchantId.String(),
		"images":        image,
		"attributes":    product.Attributes,
		"created_at":    now,
		"updated_at":    now,
		"category_id":   product.Category.CategoryId,
		"category_name": product.Category.CategoryName,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("error encoding document: %s", err)
	}

	res, err := d.Es.Index(
		"tt_product",
		&buf,
		d.Es.Index.WithDocumentID(productId.String()),
		d.Es.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s", productId.String())
	}

	log.Printf("Indexed document ID=%s", productId.String())
	return nil
}

func (d *Data) SearchProductsByName(ctx context.Context, req *biz.SearchProductRequest) (*biz.Products, error) {
	mustConditions := []map[string]interface{}{}
	for field, value := range req.Query {
		if value != nil {
			if values, ok := value.([]interface{}); ok {
				shouldConditions := []map[string]interface{}{}
				for _, v := range values {
					if str, ok := v.(string); ok && len(str) > 0 {
						if len(str) >= 2 && str[0] == '/' && str[len(str)-1] == '/' {
							shouldConditions = append(shouldConditions, map[string]interface{}{
								"regexp": map[string]interface{}{
									field: str[1 : len(str)-1],
								},
							})
						} else if containsWildcard(str) {
							shouldConditions = append(shouldConditions, map[string]interface{}{
								"wildcard": map[string]interface{}{
									field: str,
								},
							})
						} else if len(str) >= 1 && str[0] == '#' {
							shouldConditions = append(shouldConditions, map[string]interface{}{
								"match": map[string]interface{}{
									field: str[1:], // 去掉开头的 #
								},
							})
						} else {
							shouldConditions = append(shouldConditions, map[string]interface{}{
								"term": map[string]interface{}{
									field: str,
								},
							})
						}
					}
				}
				if len(shouldConditions) > 0 {
					mustConditions = append(mustConditions, map[string]interface{}{
						"bool": map[string]interface{}{
							"should":               shouldConditions,
							"minimum_should_match": 1,
						},
					})
				}
			} else if str, ok := value.(string); ok && len(str) > 0 {
				if len(str) >= 2 && str[0] == '/' && str[len(str)-1] == '/' {
					mustConditions = append(mustConditions, map[string]interface{}{
						"regexp": map[string]interface{}{
							field: str[1 : len(str)-1],
						},
					})
				} else if containsWildcard(str) {
					mustConditions = append(mustConditions, map[string]interface{}{
						"wildcard": map[string]interface{}{
							field: str,
						},
					})
				} else if len(str) >= 1 && str[0] == '#' {
					mustConditions = append(mustConditions, map[string]interface{}{
						"match": map[string]interface{}{
							field: str[1:], // 去掉开头的 #
						},
					})
				} else {
					mustConditions = append(mustConditions, map[string]interface{}{
						"term": map[string]interface{}{
							field: str,
						},
					})
				}
			}
		}
	}
	esQuery := map[string]interface{}{
		"from": (req.Page - 1) * req.Size,
		"size": req.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	queryJSON, err := json.Marshal(esQuery)
	if err != nil {
		return nil, err
	}

	search := d.Es.Search
	resp, err := search(
		search.WithContext(ctx),
		search.WithIndex("tt_product"),
		search.WithBody(bytes.NewReader(queryJSON)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("error in search response: %s", resp.String())
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	hits := searchResult["hits"].(map[string]interface{})["hits"].([]interface{})
	results := make([]*biz.Product, 0, len(hits))
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		idStr := hit.(map[string]interface{})["_id"].(string)
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}

		results = append(results, &biz.Product{
			ID:          id,
			Name:        source["name"].(string),
			Description: source["description"].(string),
			Price:       source["price"].(float64),
			Status:      biz.ProductStatus(uint32(source["status"].(float64))),
			Images:      []*biz.ProductImage{{URL: source["images"].(string)}},
			Category: biz.CategoryInfo{
				CategoryId:   uint64(source["category_id"].(float64)),
				CategoryName: source["category_name"].(string),
			},
			Attributes: map[string]*biz.AttributeValue{},
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			MerchantId: uuid.MustParse(source["merchant_id"].(string)),
		})
	}

	return &biz.Products{
		Items: results,
	}, nil
}

// containsWildcard checks if a string contains wildcard characters.
func containsWildcard(s string) bool {
	return strings.Contains(s, "*")
}

func (d *Data) AutocompleteSearch(ctx context.Context, req *biz.AutoCompleteRequest) (*biz.AutoCompleteResponse, error) {
	// 构建查询条件
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"goods-suggest": map[string]interface{}{
				"prefix": req.Prefix,
				"completion": map[string]interface{}{
					"field": "name_suggest",
					"size":  10,
					"fuzzy": map[string]interface{}{
						"fuzziness": 2,
					},
				},
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// 执行查询操作
	search := d.Es.Search
	resp, err := search(
		search.WithContext(ctx),
		search.WithIndex("tt_product"),
		search.WithBody(bytes.NewReader(queryJSON)),
		search.WithTrackTotalHits(true),
		search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return nil, fmt.Errorf("error in search response: %s", resp.String())
	}

	var searchResult map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, err
	}

	suggestions := []string{}
	if suggest, ok := searchResult["suggest"].(map[string]interface{}); ok {
		if goodsSuggest, ok := suggest["goods-suggest"].([]interface{}); ok {
			for _, item := range goodsSuggest {
				if options, ok := item.(map[string]interface{})["options"].([]interface{}); ok {
					for _, option := range options {
						if text, ok := option.(map[string]interface{})["text"].(string); ok {
							suggestions = append(suggestions, text)
						}
					}
				}
			}
		}
	}

	return &biz.AutoCompleteResponse{
		Suggestions: suggestions,
	}, nil

}

func (d *Data) UpdateProduct(ctx context.Context, req *biz.UpdateProductRequest) (bool, error) {
	var reqmap = make(map[string]interface{})
	reqmap["name"] = req.Name
	reqmap["description"] = req.Description
	reqmap["price"] = req.Price
	reqmap["category_id"] = req.Category.CategoryId
	reqmap["category_name"] = req.Category.CategoryName

	// updateDoc := map[string]interface{}{
	// 	"doc": reqmap,
	// }

	// updateJSON, err := json.Marshal(updateDoc)
	// if err != nil {
	// 	return false, err
	// }

	// 构建查询条件
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": "ctx._source.name = params.name; ctx._source.description = params.description; ctx._source.price = params.price; ctx._source.category_id = params.category_id; ctx._source.category_name = params.category_name",
			"params": reqmap,
		},
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": req.ID.String(),
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return false, err
	}

	update := d.Es.UpdateByQuery
	resp, err := update(
		[]string{"tt_product"},
		update.WithBody(bytes.NewReader(queryJSON)),
		update.WithContext(ctx),
		update.WithPretty(),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	fmt.Println(resp)
	var respmap map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bodyBytes, &respmap)
	if err != nil {
		return false, err
	}
	if respmap["error"] != nil {
		if reason, ok := respmap["error"].(map[string]interface{}); ok {
			return false, errors.New(reason["reason"].(string))
		}
		return false, errors.New("update error")
	}

	return true, err
}

func (d *Data) DeleteProduct(ctx context.Context, req *biz.DeleteProductRequest) (bool, error) {
	// 构建查询条件
	fmt.Println("******************开始\n", req.ID.String())
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": req.ID.String(),
			},
		},
	}
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return false, err
	}
	// 执行删除操作
	deleteReq := esapi.DeleteByQueryRequest{
		Index: []string{"tt_product"},
		Body:  bytes.NewReader(queryJSON),
	}

	resp, err := deleteReq.Do(ctx, d.Es)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	fmt.Println("******************已将删除\n", resp)
	// 处理响应结果
	if resp.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return false, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			// 打印错误信息
			return false, fmt.Errorf("error deleting document: %s", e["error"].(map[string]interface{})["reason"])
		}
	}

	var respmap map[string]interface{}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(bodyBytes, &respmap)
	if err != nil {
		return false, err
	}
	if respmap["error"] != nil {
		if reason, ok := respmap["error"].(map[string]interface{}); ok {
			return false, errors.New(reason["reason"].(string))
		}
	}

	return true, nil
}
