package campaign

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Repository interface {
	FindAll(ctx context.Context) ([]Campaign, error)
	FindByUserID(ctx context.Context, userID string) ([]Campaign, error)
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
