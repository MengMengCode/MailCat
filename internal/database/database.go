package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"mailcat/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS emails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_address TEXT NOT NULL,
		to_address TEXT NOT NULL,
		subject TEXT,
		body TEXT,
		html_body TEXT,
		headers TEXT,
		raw_email TEXT,
		received_at DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_emails_from ON emails(from_address);
	CREATE INDEX IF NOT EXISTS idx_emails_to ON emails(to_address);
	CREATE INDEX IF NOT EXISTS idx_emails_created_at ON emails(created_at);
	`

	_, err := db.conn.Exec(query)
	if err != nil {
		return err
	}

	// 添加raw_email列（如果不存在）
	alterQuery := `ALTER TABLE emails ADD COLUMN raw_email TEXT;`
	db.conn.Exec(alterQuery) // 忽略错误，因为列可能已存在

	return nil
}

func (db *DB) SaveEmail(emailReq *models.EmailRequest) (*models.Email, error) {
	headersJSON, err := json.Marshal(emailReq.Headers)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal headers: %w", err)
	}

	query := `
	INSERT INTO emails (from_address, to_address, subject, body, html_body, headers, raw_email, received_at, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now()
	result, err := db.conn.Exec(query,
		emailReq.From,
		emailReq.To,
		emailReq.Subject,
		emailReq.Body,
		emailReq.HTMLBody,
		string(headersJSON),
		emailReq.RawEmail,
		now,
		now,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert email: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return db.GetEmailByID(int(id))
}

func (db *DB) GetEmailByID(id int) (*models.Email, error) {
	// 首先尝试查询包含raw_email的完整记录
	query := `
	SELECT id, from_address, to_address, subject, body, html_body, headers,
	       COALESCE(raw_email, '') as raw_email, received_at, created_at
	FROM emails WHERE id = ?
	`

	row := db.conn.QueryRow(query, id)
	email := &models.Email{}

	err := row.Scan(
		&email.ID,
		&email.From,
		&email.To,
		&email.Subject,
		&email.Body,
		&email.HTMLBody,
		&email.Headers,
		&email.RawEmail,
		&email.ReceivedAt,
		&email.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan email: %w", err)
	}

	return email, nil
}

func (db *DB) GetEmails(page, limit int) (*models.EmailListResponse, error) {
	offset := (page - 1) * limit

	// Get total count
	countQuery := "SELECT COUNT(*) FROM emails"
	var total int
	err := db.conn.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get emails
	query := `
	SELECT id, from_address, to_address, subject, body, html_body, headers,
	       COALESCE(raw_email, '') as raw_email, received_at, created_at
	FROM emails
	ORDER BY created_at DESC
	LIMIT ? OFFSET ?
	`

	rows, err := db.conn.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query emails: %w", err)
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		email := models.Email{}
		err := rows.Scan(
			&email.ID,
			&email.From,
			&email.To,
			&email.Subject,
			&email.Body,
			&email.HTMLBody,
			&email.Headers,
			&email.RawEmail,
			&email.ReceivedAt,
			&email.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan email: %w", err)
		}
		emails = append(emails, email)
	}

	return &models.EmailListResponse{
		Emails: emails,
		Total:  total,
		Page:   page,
		Limit:  limit,
	}, nil
}

// GetEmailStats 获取邮件统计信息
func (db *DB) GetEmailStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 获取总邮件数
	var totalEmails int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM emails").Scan(&totalEmails)
	if err != nil {
		return nil, fmt.Errorf("failed to get total emails: %w", err)
	}
	stats["total_emails"] = totalEmails
	
	// 获取今日邮件数
	var todayEmails int
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM emails
		WHERE DATE(created_at) = DATE('now', 'localtime')
	`).Scan(&todayEmails)
	if err != nil {
		return nil, fmt.Errorf("failed to get today emails: %w", err)
	}
	stats["today_emails"] = todayEmails
	
	// 获取最近7天的邮件统计（用于图表）
	// 使用CTE生成最近7天的日期，确保所有日期都显示
	rows, err := db.conn.Query(`
		WITH RECURSIVE dates(date) AS (
			SELECT date('now', '-6 days')
			UNION ALL
			SELECT date(date, '+1 day')
			FROM dates
			WHERE date < date('now')
		)
		SELECT dates.date, COALESCE(COUNT(emails.id), 0) as count
		FROM dates
		LEFT JOIN emails ON DATE(emails.created_at) = dates.date
		GROUP BY dates.date
		ORDER BY dates.date ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %w", err)
	}
	defer rows.Close()
	
	var weeklyStats []map[string]interface{}
	for rows.Next() {
		var date string
		var count int
		err := rows.Scan(&date, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weekly stats: %w", err)
		}
		weeklyStats = append(weeklyStats, map[string]interface{}{
			"date":  date,
			"count": count,
		})
	}
	stats["weekly_stats"] = weeklyStats
	
	return stats, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}