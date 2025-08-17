package models

import (
	"time"
)

type Email struct {
	ID          int       `json:"id" db:"id"`
	From        string    `json:"from" db:"from_address"`
	To          string    `json:"to" db:"to_address"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	HTMLBody    string    `json:"html_body" db:"html_body"`
	Headers     string    `json:"headers" db:"headers"`
	RawEmail    string    `json:"raw_email" db:"raw_email"`
	ReceivedAt  time.Time `json:"received_at" db:"received_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type EmailRequest struct {
	From     string            `json:"from" binding:"required"`
	To       string            `json:"to" binding:"required"`
	Subject  string            `json:"subject"`
	Body     string            `json:"body"`
	HTMLBody string            `json:"html_body"`
	Headers  map[string]string `json:"headers"`
	RawEmail string            `json:"raw_email"`
}

type EmailListResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}