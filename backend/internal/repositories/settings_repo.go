package repositories

import (
	"context"
	"database/sql"
)

type SettingsRepo struct{ db *sql.DB }

func NewSettingsRepo(db *sql.DB) *SettingsRepo { return &SettingsRepo{db: db} }

func (r *SettingsRepo) GetAll(ctx context.Context) (map[string]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT setting_key, setting_value FROM dbo.site_settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]string)
	for rows.Next() {
		var k string
		var v sql.NullString
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		if v.Valid {
			out[k] = v.String
		} else {
			out[k] = ""
		}
	}
	return out, rows.Err()
}

func (r *SettingsRepo) Set(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx, `
		MERGE dbo.site_settings AS t
		USING (SELECT @p1 AS setting_key, @p2 AS setting_value) AS s
		ON t.setting_key = s.setting_key
		WHEN MATCHED THEN UPDATE SET setting_value = s.setting_value, updated_at = SYSUTCDATETIME()
		WHEN NOT MATCHED THEN INSERT (setting_key, setting_value) VALUES (s.setting_key, s.setting_value);`,
		key, value)
	return err
}

func (r *SettingsRepo) SetMany(ctx context.Context, pairs map[string]string) error {
	for k, v := range pairs {
		if err := r.Set(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}
