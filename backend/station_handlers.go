package main

import (
	"net/http"
	"net/url"
	"strings"

	"gorm.io/gorm"
)

func onStationGet(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get()

	if strings.TrimSpace(id) == "" {
		// id is not null
	}
}

func onStationPost(db *gorm.DB, r *http.Request) (ResponseBody, HttpError) {

}

func onStationPut(db *gorm.DB, r *http.Request) (ResponseBody, HttpError) {

}

func onStationDelete(db *gorm.DB, r *http.Request) (ResponseBody, HttpError) {

}
