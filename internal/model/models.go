package model

import "time"

type Status string

const (
	Available    Status = "available"
	NotAvailable Status = "not available"
)

const (
	Pending    Status = "pending"
	Processing Status = "processing"
	Completed  Status = "completed"
)

type Link struct {
	URL    string `json:"url"`
	Status Status `json:"status"` // avaliable, not avaliable
}

type LinkList struct {
	ID        int       `json:"id"`
	Links     *[]Link   `json:"links"`
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"` // pending, processing, completed
}

type CheckRequest struct {
	Links []string `json:"links"`
}

type CheckResponse struct {
	ID    int    `json:"links_num"`
	Links []Link `json:"results"`
}

type ReportRequest struct {
	LinksIDs []int `json:"links_num"`
}
