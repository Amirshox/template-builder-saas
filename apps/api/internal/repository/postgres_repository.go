package repository

import (
	"context"
	"fmt"
	"template-builder-api/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateOrg(ctx context.Context, name string) (*model.Org, error)
	GetOrg(ctx context.Context, id uuid.UUID) (*model.Org, error)
	CreateUser(ctx context.Context, email, name, passwordHash string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	CreateTemplate(ctx context.Context, template *model.Template) error
	GetTemplate(ctx context.Context, id uuid.UUID) (*model.Template, error)
	ListTemplates(ctx context.Context, orgID uuid.UUID) ([]model.Template, error)
	CreateTemplateVersion(ctx context.Context, version *model.TemplateVersion) error
	ListTemplateVersions(ctx context.Context, templateID uuid.UUID) ([]model.TemplateVersion, error)
	GetTemplateVersion(ctx context.Context, templateID uuid.UUID, version int) (*model.TemplateVersion, error)
	GetMaxVersion(ctx context.Context, templateID uuid.UUID) (int, error)

	CreateAsset(ctx context.Context, asset *model.Asset) error
	GetAsset(ctx context.Context, id uuid.UUID) (*model.Asset, error)

	// Jobs
	CreateJob(ctx context.Context, job *model.GenerationJob) error
	GetJob(ctx context.Context, id uuid.UUID) (*model.GenerationJob, error)
	UpdateJobStatus(ctx context.Context, id uuid.UUID, status string, outputAssetID *uuid.UUID, errMsg string) error

	// Auth
	ListMemberships(ctx context.Context, userID uuid.UUID) ([]model.Membership, error)
	CreateMembership(ctx context.Context, userID, orgID uuid.UUID, role string) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateOrg(ctx context.Context, name string) (*model.Org, error) {
	query := `INSERT INTO orgs (name) VALUES ($1) RETURNING id, name, created_at`
	row := r.db.QueryRow(ctx, query, name)

	var org model.Org
	if err := row.Scan(&org.ID, &org.Name, &org.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create org: %w", err)
	}
	return &org, nil
}

func (r *PostgresRepository) GetOrg(ctx context.Context, id uuid.UUID) (*model.Org, error) {
	query := `SELECT id, name, created_at FROM orgs WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var org model.Org
	if err := row.Scan(&org.ID, &org.Name, &org.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to get org: %w", err)
	}
	return &org, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, name, passwordHash string) (*model.User, error) {
	query := `INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3) RETURNING id, email, name, created_at`
	row := r.db.QueryRow(ctx, query, email, name, passwordHash)

	var user model.User
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &user, nil
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, name, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user model.User
	if err := row.Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *PostgresRepository) CreateTemplate(ctx context.Context, t *model.Template) error {
	query := `INSERT INTO templates (id, org_id, name, type, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, t.ID, t.OrgID, t.Name, t.Type, t.Status, t.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*model.Template, error) {
	query := `SELECT id, org_id, name, type, status, created_at FROM templates WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var t model.Template
	if err := row.Scan(&t.ID, &t.OrgID, &t.Name, &t.Type, &t.Status, &t.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	return &t, nil
}

func (r *PostgresRepository) ListTemplates(ctx context.Context, orgID uuid.UUID) ([]model.Template, error) {
	query := `SELECT id, org_id, name, type, status, created_at FROM templates WHERE org_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	var templates []model.Template
	for rows.Next() {
		var t model.Template
		if err := rows.Scan(&t.ID, &t.OrgID, &t.Name, &t.Type, &t.Status, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

func (r *PostgresRepository) CreateTemplateVersion(ctx context.Context, v *model.TemplateVersion) error {
	query := `INSERT INTO template_versions (id, template_id, version, status, template_json, schema_json, created_by, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, v.ID, v.TemplateID, v.Version, v.Status, v.TemplateJSON, v.SchemaJSON, v.CreatedBy, v.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create template version: %w", err)
	}
	return nil
}

func (r *PostgresRepository) ListTemplateVersions(ctx context.Context, templateID uuid.UUID) ([]model.TemplateVersion, error) {
	query := `SELECT id, template_id, version, status, created_by, created_at 
			  FROM template_versions WHERE template_id = $1 ORDER BY version DESC`
	rows, err := r.db.Query(ctx, query, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	defer rows.Close()

	var versions []model.TemplateVersion
	for rows.Next() {
		var v model.TemplateVersion
		// Note: We skip fetching heavy JSONs for the list view
		if err := rows.Scan(&v.ID, &v.TemplateID, &v.Version, &v.Status, &v.CreatedBy, &v.CreatedAt); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, nil
}

func (r *PostgresRepository) GetTemplateVersion(ctx context.Context, templateID uuid.UUID, version int) (*model.TemplateVersion, error) {
	query := `SELECT id, template_id, version, status, template_json, schema_json, created_at 
			  FROM template_versions WHERE template_id = $1 AND version = $2`
	row := r.db.QueryRow(ctx, query, templateID, version)

	var v model.TemplateVersion
	// Note: We might need to handle NULLs for docx_asset_id etc if we query them.
	// For MVP simplified query above ignores partial fields.
	if err := row.Scan(&v.ID, &v.TemplateID, &v.Version, &v.Status, &v.TemplateJSON, &v.SchemaJSON, &v.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to get template version: %w", err)
	}
	return &v, nil
}

func (r *PostgresRepository) GetMaxVersion(ctx context.Context, templateID uuid.UUID) (int, error) {
	query := `SELECT COALESCE(MAX(version), 0) FROM template_versions WHERE template_id = $1`
	var maxVersion int
	err := r.db.QueryRow(ctx, query, templateID).Scan(&maxVersion)
	if err != nil {
		return 0, fmt.Errorf("failed to get max version: %w", err)
	}
	return maxVersion, nil
}

func (r *PostgresRepository) CreateAsset(ctx context.Context, a *model.Asset) error {
	query := `INSERT INTO assets (id, org_id, type, filename, content_type, size_bytes, s3_key, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query, a.ID, a.OrgID, a.Type, a.Filename, a.ContentType, a.SizeBytes, a.S3Key, a.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}
	return nil
}

func (r *PostgresRepository) CreateJob(ctx context.Context, job *model.GenerationJob) error {
	query := `INSERT INTO generation_jobs (id, org_id, template_id, status, error_message, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, job.ID, job.OrgID, job.TemplateID, job.Status, job.ErrorMessage, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetJob(ctx context.Context, id uuid.UUID) (*model.GenerationJob, error) {
	query := `SELECT id, org_id, template_id, status, output_asset_id, error_message, created_at, updated_at FROM generation_jobs WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var job model.GenerationJob
	var outputAssetID *uuid.UUID
	var errMsg *string

	if err := row.Scan(&job.ID, &job.OrgID, &job.TemplateID, &job.Status, &outputAssetID, &errMsg, &job.CreatedAt, &job.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	job.OutputAssetID = outputAssetID
	if errMsg != nil {
		job.ErrorMessage = *errMsg
	}
	return &job, nil
}

func (r *PostgresRepository) UpdateJobStatus(ctx context.Context, id uuid.UUID, status string, outputAssetID *uuid.UUID, errMsg string) error {
	query := `UPDATE generation_jobs SET status = $1, output_asset_id = $2, error_message = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(ctx, query, status, outputAssetID, errMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
	}
	return nil
}

func (r *PostgresRepository) GetAsset(ctx context.Context, id uuid.UUID) (*model.Asset, error) {
	query := `SELECT id, org_id, type, filename, content_type, size_bytes, s3_key, created_at FROM assets WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var a model.Asset
	if err := row.Scan(&a.ID, &a.OrgID, &a.Type, &a.Filename, &a.ContentType, &a.SizeBytes, &a.S3Key, &a.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	return &a, nil
}

func (r *PostgresRepository) ListMemberships(ctx context.Context, userID uuid.UUID) ([]model.Membership, error) {
	// We assume there's a memberships table.
	// If not, we might need to check schema.
	// Wait, the Summary says "PostgreSQL schema with tables for ... memberships"
	query := `SELECT id, user_id, org_id, role, created_at FROM memberships WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []model.Membership
	for rows.Next() {
		var m model.Membership
		if err := rows.Scan(&m.ID, &m.UserID, &m.OrgID, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, nil
}

func (r *PostgresRepository) CreateMembership(ctx context.Context, userID, orgID uuid.UUID, role string) error {
	query := `INSERT INTO memberships (user_id, org_id, role) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, userID, orgID, role)
	return err
}
