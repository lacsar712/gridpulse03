package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Record struct {
	ID        int64     `json:"id"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	OwnerID   int64     `json:"owner_id"`
	Bytes     int       `json:"bytes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Revision struct {
	ID        int64     `json:"id"`
	RecordID  int64     `json:"record_id"`
	Body      string    `json:"body"`
	Editor    string    `json:"editor"`
	CreatedAt time.Time `json:"created_at"`
}

type Attachment struct {
	ID        int64     `json:"id"`
	RecordID  int64     `json:"record_id"`
	Name      string    `json:"name"`
	SHA       string    `json:"sha"`
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}

type Audit struct {
	ID        int64     `json:"id"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

type ExportFile struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
