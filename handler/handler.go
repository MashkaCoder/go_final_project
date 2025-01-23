package handler

import dbfuncs "github.com/MashkaCoder/go_final_project/database"

type Handler struct {
	DB *dbfuncs.DB
}

func NewHandler(db *dbfuncs.DB) *Handler {
	return &Handler{DB: db}
}
