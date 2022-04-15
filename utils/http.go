package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/funte/xmlymft/common"
)

type HelloResponse struct {
	Error   string `json:"err"`
	Message string `json:"message"`
}

func HTTPGetHelloResponse(url string) (*HelloResponse, error) {
	result := new(HelloResponse)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type SearchAlbumResponse struct {
	Error   string                   `json:"err"`
	Message string                   `json:"message"`
	Data    common.SearchAlbumResult `json:"data"`
}

func HTTPGetSearchAlbumResponse(url string) (*SearchAlbumResponse, error) {
	result := new(SearchAlbumResponse)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type QueryPlayListResponse struct {
	Error   string                     `json:"err"`
	Message string                     `json:"message"`
	Data    common.QueryPlayListResult `json:"data"`
}

func HTTPGetQueryPlayListResponse(url string) (*QueryPlayListResponse, error) {
	result := new(QueryPlayListResponse)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type QueryTrackAddressResponse struct {
	Error   string                         `json:"err"`
	Message string                         `json:"message"`
	Data    common.QueryTrackAddressResult `json:"data"`
}

func HTTPGetQueryTrackAddressResponse(url string) (*QueryTrackAddressResponse, error) {
	result := new(QueryTrackAddressResponse)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// // Requires go-1.18
// // Type contracts for HTTP server response.
// type HTTPResponseType interface {
// 	HelloResponse | SearchAlbumsResponse | QueryPlayListResponse | QueryTrackAddressResponse
// }
// func HTTPGet[R HTTPResponseType](url string) (*R, error) {
// 	result := new(R)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(data, &result)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
