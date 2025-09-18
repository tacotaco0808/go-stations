package model

import "time"

type (
	// A TODO expresses ...
	TODO struct {
		ID          int
		Subject     string
		Description string
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	// A CreateTODORequest expresses ...
	CreateTODORequest struct{}
	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct{}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)
