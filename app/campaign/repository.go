package campaign

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Campaign, error)
	FindByUserID(ctx context.Context, userID string) ([]Campaign, error)
	FindByID(ctx context.Context, ID string) (Campaign, error)
	Save(ctx context.Context, campaign Campaign) (Campaign, error)
}

type repository struct {
	DB *sql.DB
}

const (
	layoutDateTime = "2006-01-02 15:04:05"
)

func NewCampaignRepository(DB *sql.DB) Repository {
	return &repository{DB}
}

func (r *repository) FindAll(ctx context.Context) ([]Campaign, error) {
	campaigns := []Campaign{}

	fmt.Println(layoutDateTime)

	sqlQuery := "SELECT id, user_id, name, short_description, description, slug, perks, goal_amount, current_amount, backer_count, created_at, updated_at FROM campaigns"

	rows, err := r.DB.QueryContext(ctx, sqlQuery)
	if err != nil {
		return campaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		campaign := Campaign{}
		var createdAt, updatedAt string

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Slug,
			&campaign.Perks,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.BackerCount,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return campaigns, err
		}

		if createdAt != "" || updatedAt != "" {
			if campaign.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
				log.Fatal(err)
			}

			if campaign.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
				log.Fatal(err)
			}
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) ([]Campaign, error) {
	campaigns := []Campaign{}

	sqlQuery := "SELECT id, user_id, name, short_description, description, slug, perks, goal_amount, current_amount, backer_count, created_at, updated_at FROM campaigns WHERE user_id = $1"

	rows, err := r.DB.QueryContext(ctx, sqlQuery, userID)
	if err != nil {
		return campaigns, err
	}

	defer rows.Close()

	for rows.Next() {
		campaign := Campaign{}
		var createdAt, updatedAt string

		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Slug,
			&campaign.Perks,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.BackerCount,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return campaigns, err
		}

		if createdAt != "" || updatedAt != "" {
			if campaign.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
				log.Fatal(err)
			}

			if campaign.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
				log.Fatal(err)
			}
		}

		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (r *repository) FindByID(ctx context.Context, ID string) (Campaign, error) {
	campaign := Campaign{}
	var createdAt, updatedAt string

	sqlQuery := "SELECT id, user_id, name, short_description, description, slug, perks, goal_amount, current_amount, backer_count, created_at, updated_at FROM campaigns WHERE id = $1"

	rows, err := r.DB.QueryContext(ctx, sqlQuery, ID)
	if err != nil {
		return campaign, err
	}

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(
			&campaign.ID,
			&campaign.UserID,
			&campaign.Name,
			&campaign.ShortDescription,
			&campaign.Description,
			&campaign.Slug,
			&campaign.Perks,
			&campaign.GoalAmount,
			&campaign.CurrentAmount,
			&campaign.BackerCount,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return campaign, err
		}
	}

	if createdAt != "" || updatedAt != "" {
		if campaign.CreatedAt, err = time.Parse(time.RFC3339, createdAt); err != nil {
			log.Fatal(err)
		}

		if campaign.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt); err != nil {
			log.Fatal(err)
		}
	}

	return campaign, nil
}

func (r *repository) Save(ctx context.Context, campaign Campaign) (Campaign, error) {
	sqlQuery := "INSERT into campaigns (id, user_id, name, short_description, description, slug, perks, goal_amount, current_amount, backer_count, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"

	stmt, err := r.DB.PrepareContext(ctx, sqlQuery)
	if err != nil {
		return campaign, err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, &campaign.ID,
		&campaign.UserID,
		&campaign.Name,
		&campaign.ShortDescription,
		&campaign.Description,
		&campaign.Slug,
		&campaign.Perks,
		&campaign.GoalAmount,
		&campaign.CurrentAmount,
		&campaign.BackerCount,
		time.Now().Format(layoutDateTime),
		time.Now().Format(layoutDateTime),
	)

	if err != nil {
		return campaign, err
	}

	log.Info("Success insert new campaign!")
	return campaign, nil
}
