package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Koshsky/subs-service/models"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepo interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]models.Subscription, error)
	Update(ctx context.Context, id int, sub models.SubscriptionUpdate) (*models.Subscription, error)
	Delete(ctx context.Context, id int) error
	SumPrice(ctx context.Context, params models.SumPriceParams) (float64, error)
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) SumPrice(ctx context.Context, params models.SumPriceParams) (float64, error) {
	query := `
        SELECT COALESCE(SUM(price), 0)
        FROM subscriptions
        WHERE user_id = $1
    `

	args := []interface{}{params.UserID}
	argIdx := 2

	if params.Service != "" {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, params.Service)
		argIdx++
	}

	startDate := params.StartMonth.Time().Format("2006-01-02")
	query += fmt.Sprintf(" AND start_date >= $%d", argIdx)
	args = append(args, startDate)
	argIdx++

	endTime := params.EndMonth.Time()
	lastDay := time.Date(endTime.Year(), endTime.Month()+1, 0, 0, 0, 0, 0, endTime.Location())
	query += fmt.Sprintf(" AND start_date <= $%d", argIdx)
	args = append(args, lastDay.Format("2006-01-02"))

	var sum float64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&sum)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate sum: %w", err)
	}

	return sum, nil
}

func (r *PostgresRepo) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		sub.Service,

		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	if id <= 0 {
		return nil, fmt.Errorf("")
	}

	var sub models.Subscription
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM subscriptions
		WHERE id = $1`

	err := r.db.GetContext(ctx, &sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("subscription with id=%d not found", id)

		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		FROM subscriptions`

	err := r.db.SelectContext(ctx, &subs, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	return subs, nil
}

func (s *PostgresRepo) Update(ctx context.Context, id int, update models.SubscriptionUpdate) (*models.Subscription, error) {
	current, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if update.Service != nil {
		current.Service = *update.Service
	}
	if update.Price != nil {
		current.Price = *update.Price
	}
	if update.UserID != nil {
		current.UserID = *update.UserID
	}
	if update.StartDate != nil {
		current.StartDate = *update.StartDate
	}
	if update.EndDate != nil {
		current.EndDate = update.EndDate
	}

	_, err = s.db.ExecContext(ctx,
		`UPDATE subscriptions SET
            service_name = $1,
            price = $2,
            user_id = $3,
            start_date = $4,
            end_date = $5
        WHERE id = $6`,
		current.Service, current.Price, current.UserID, current.StartDate, current.EndDate, id,
	)
	if err != nil {
		return nil, err
	}

	return current, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid subscription id=%d", id)
	}

	query := "DELETE FROM subscriptions WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription with id=%d not found", id)
	}

	return nil
}
