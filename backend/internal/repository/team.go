package repository

import (
	"context"

	"github.com/jackc/pgx/v5"

	"uzbekona.dev/backend/internal/model"
)

type TeamRepo struct {
	db DBTX
}

var teamColumns = `t.id, t.name, t.role, t.bio, t.photo_id, ` + mediaJSON("t.photo_id") + `,
	t.socials, t.is_published, t.position, t.created_at, t.updated_at`

func scanTeamMember(row pgx.Row) (*model.TeamMember, error) {
	var t model.TeamMember
	err := row.Scan(&t.ID, &t.Name, &t.Role, &t.Bio, &t.PhotoID, &t.Photo,
		&t.Socials, &t.IsPublished, &t.Position, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	t.Socials = orEmpty(t.Socials)
	return &t, nil
}

func (r *TeamRepo) list(ctx context.Context, publishedOnly bool) ([]model.TeamMember, error) {
	sql := `SELECT ` + teamColumns + ` FROM team_members t`
	if publishedOnly {
		sql += ` WHERE t.is_published`
	}
	rows, err := r.db.Query(ctx, sql+` ORDER BY t.position, t.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []model.TeamMember{}
	for rows.Next() {
		t, err := scanTeamMember(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return out, rows.Err()
}

func (r *TeamRepo) ListPublished(ctx context.Context) ([]model.TeamMember, error) {
	return r.list(ctx, true)
}

func (r *TeamRepo) ListAll(ctx context.Context) ([]model.TeamMember, error) {
	return r.list(ctx, false)
}

func (r *TeamRepo) GetByID(ctx context.Context, id int64) (*model.TeamMember, error) {
	t, err := scanTeamMember(r.db.QueryRow(ctx, `SELECT `+teamColumns+` FROM team_members t WHERE t.id = $1`, id))
	return t, notFound(err)
}

func (r *TeamRepo) Create(ctx context.Context, in *model.TeamMemberInput) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO team_members (name, role, bio, photo_id, socials, is_published, position)
		VALUES ($1,$2,$3,$4,$5,$6,(SELECT COALESCE(max(position), 0) + 1 FROM team_members)) RETURNING id`,
		in.Name, in.Role, in.Bio, in.PhotoID, orEmpty(in.Socials), in.IsPublished,
	).Scan(&id)
	return id, err
}

func (r *TeamRepo) Update(ctx context.Context, id int64, in *model.TeamMemberInput) error {
	tag, err := r.db.Exec(ctx, `UPDATE team_members SET name = $2, role = $3, bio = $4, photo_id = $5,
		socials = $6, is_published = $7 WHERE id = $1`,
		id, in.Name, in.Role, in.Bio, in.PhotoID, orEmpty(in.Socials), in.IsPublished)
	if err == nil && tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return err
}

func (r *TeamRepo) Reorder(ctx context.Context, ids []int64) error {
	return reorder(ctx, r.db, "team_members", ids)
}

func (r *TeamRepo) Delete(ctx context.Context, id int64) error {
	return deleteByID(ctx, r.db, "team_members", id)
}
