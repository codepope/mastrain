package raindrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

type collectionWrapper struct {
	Items []struct {
		ID    int64  `json:"_id"`
		Title string `json:"title"`
	} `json:"items"`
}

// API responsible for connecting to Raindrop.io.

type API struct {
	client  *http.Client
	baseURL string
}

// Bookmark defines the data type to save a bookmark.
type Bookmark struct {
	// CollectionName defines in which collection the bookmark is saved.
	CollectionName string `json:"collection_name"`
	// URL defines the url of the bookmark.
	URL string `json:"url"`
}

// Raindrop defines the data structure to save a bookmark in raindrop.io.
type Raindrop struct {
	// Title defines the raindrop's title.
	Title string `json:"title"`
	// Excerpt defines the raindrop's excerpt.
	Excerpt string `json:"excerpt"`
	// Link defines the bookmark's link.
	// Link is better to be set via parsedURL.meta.canonical.
	Link string `json:"link"`
	// CollectionID defines in which collection the raindrop will save the bookmark.
	CollectionID int64 `json:"collectionId"`
}

// ParsedURL defines the result of parsed URL from Raindrop.
type ParsedURL struct {
	// Error defines the parse URL error.
	Error string `json:"error"`
	Item  struct {
		// Title defines the title of the URL.
		Title string `json:"title"`
		// Excerpt defines the excerpt of the URL.
		Excerpt string `json:"excerpt"`
		// Meta defines the metadata of the item
		Meta struct {
			Canonical string `json:"canonical"`
		} `json:"meta"`
	} `json:"item"`
}

// Collection defines the data type for a collection.
type Collection struct {
	// ID defines the collection's ID.
	ID int64 `json:"_id"`
	// Name defines collection's name.
	Name string `json:"name"`
}

type ExistsRequest struct {
	URLs []string `json:"urls"`
}

type Duplicate struct {
	ID   int64  `json:"_id"`
	Link string `json:"link"`
}

type ExistsResponse struct {
	Result       bool        `json:"result"`
	ErrorMessage string      `json:"errorMessage"`
	IDs          []int64     `json:"ids"`
	Duplicates   []Duplicate `json:"duplicates"`
}

// NewAPI creates an instance of API.
func NewAPI(baseURL, token string) *API {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	}))

	return &API{
		client:  client,
		baseURL: baseURL,
	}
}

// GetCollections gets all root collections from raindrop.io
func (a *API) GetCollections(ctx context.Context) ([]*Collection, error) {
	resp, derr := a.client.Get(a.baseURL + "/collections")
	if derr != nil {
		return []*Collection{}, derr
	}
	defer resp.Body.Close()

	var wrapper collectionWrapper
	if jerr := json.NewDecoder(resp.Body).Decode(&wrapper); jerr != nil {
		return []*Collection{}, jerr
	}

	return convertWrapperToCollections(wrapper), nil
}

// ParseURL parse an URL to get detailed information from raindrop.io.
func (a *API) ParseURL(ctx context.Context, url string) (*ParsedURL, error) {
	reqURL := fmt.Sprintf("%s/import/url/parse?url=%s", a.baseURL, url)
	resp, derr := a.client.Get(reqURL)
	if derr != nil {
		return nil, derr
	}
	defer resp.Body.Close()

	var parsed ParsedURL
	if jerr := json.NewDecoder(resp.Body).Decode(&parsed); jerr != nil {
		return nil, jerr
	}

	return &parsed, nil
}

func (a *API) GetDuplicates(ctx context.Context, urls []string) (*[]Duplicate, error) {
	reqURL := fmt.Sprintf("%s/import/url/exists", a.baseURL)

	wrappedurls := ExistsRequest{URLs: urls}

	body, err := json.Marshal(wrappedurls)
	if err != nil {
		return nil, fmt.Errorf("marshal error")
	}

	resp, err := a.client.Post(reqURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tmp ExistsResponse
	json.NewDecoder(resp.Body).Decode(&tmp)

	if !tmp.Result {
		return nil, nil
	}

	return &tmp.Duplicates, nil
}

// SaveRaindrop saves a raindrop bookmark to specific collection in raindrop.io.
func (a *API) SaveRaindrop(ctx context.Context, raindrop *Raindrop) error {
	fmt.Printf("saving raindrop: %+v\n", raindrop)

	body, merr := json.Marshal(raindrop)
	if merr != nil {
		return fmt.Errorf("marshal error")
	}

	resp, derr := a.client.Post(a.baseURL+"/raindrop", "application/json", bytes.NewBuffer(body))
	if derr != nil {
		return derr
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	var tmp interface{}
	json.NewDecoder(resp.Body).Decode(&tmp)
	return fmt.Errorf("[SaveRaindrop] errors: %v", tmp)
}

func convertWrapperToCollections(wrapper collectionWrapper) []*Collection {
	colls := make([]*Collection, len(wrapper.Items))
	for i, item := range wrapper.Items {
		colls[i] = &Collection{ID: item.ID, Name: item.Title}
	}
	return colls
}
