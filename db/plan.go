package db

import (
	"context"
	"database/sql"
	"foodie/core"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/rs/xid"
)

// InsertPlan inserts a new plan into the database.
func InsertPlan(
	ctx context.Context,
	db *sql.DB,
	uid xid.ID,
	pc core.PlanCore,
) (*core.Plan, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	pl := core.Plan{
		ID:        xid.New(),
		UserID:    uid,
		CreatedAt: time.Now(),
		PlanCore:  pc,
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Insert("plans").SetMap(map[string]interface{}{
			"plans.id":          pl.ID,
			"plans.user_id":     pl.UserID,
			"plans.name":        pl.Name,
			"plans.description": pl.Description,
			"plans.created_at":  pl.CreatedAt,
		}),
	)
	if err != nil {
		return nil, err
	}

	for _, pr := range pc.Recipies {
		pr.PlanID = pl.ID

		if err := upsertPlanRecipy(
			ctx,
			tx,
			pr,
		); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &pl, nil
}

// GetPlans retrieves all plans.
func GetPlans(
	ctx context.Context,
	qc squirrel.QueryerContext,
) ([]core.Plan, error) {

	return selectPlans(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb
		},
	)
}

// GetPlansByUserID retrieves plans by the user id.
func GetPlansByUserID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	uid xid.ID,
) ([]core.Plan, error) {

	return selectPlans(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"plans.user_id": uid},
			)
		},
	)
}

// GetPlanByID retrieves a plan by its id.
func GetPlanByID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) (*core.Plan, error) {

	pp, err := selectPlans(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"plans.id": id},
			)
		},
	)
	if err != nil {
		return nil, err
	}

	if len(pp) == 0 {
		return nil, ErrNotFound
	}

	return &pp[0], nil
}

// UpdatePlanByID updates an existing plan by its id. An updated plan
// is returned.
func UpdatePlanByID(
	ctx context.Context,
	db *sql.DB,
	id xid.ID,
	pc core.PlanCore,
) (*core.Plan, error) {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil
	}
	defer tx.Rollback()

	if err := deletePlanRecipes(
		ctx,
		tx,
		id,
	); err != nil {
		return nil, err
	}

	for _, pr := range pc.Recipies {
		pr.PlanID = id

		if err := upsertPlanRecipy(
			ctx,
			tx,
			pr,
		); err != nil {
			return nil, err
		}
	}

	_, err = squirrel.ExecContextWith(
		ctx,
		tx,
		squirrel.Update("plans").SetMap(map[string]interface{}{
			"plans.name":        pc.Name,
			"plans.description": pc.Description,
		}).Where(
			squirrel.Eq{"plans.id": id},
		),
	)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	pl, err := GetPlanByID(ctx, db, id)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

// DeletePlanByID deletes a plan by its id.
func DeletePlanByID(
	ctx context.Context,
	ec squirrel.ExecerContext,
	id xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("plans").Where(
			squirrel.Eq{"plans.id": id},
		),
	)
	return err
}

// selectPlans selects all plans by the provided decorator function.
func selectPlans(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.Plan, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"plans.id",
			"plans.user_id",
			"plans.name",
			"plans.description",
			"plans.created_at",
		).From("plans"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pp := make([]core.Plan, 0)
	for rows.Next() {
		var pl core.Plan
		if err := rows.Scan(
			&pl.ID,
			&pl.UserID,
			&pl.Name,
			&pl.Description,
			&pl.CreatedAt,
		); err != nil {
			return nil, err
		}

		prs, err := getPlanRecipesByPlanID(ctx, qc, pl.ID)
		if err != nil {
			return nil, err
		}

		pl.Recipies = prs
		pp = append(pp, pl)
	}

	return pp, nil
}

// GetPlanRecipesByRecipyID selects plan recipes by the recipy id.
func GetPlanRecipesByRecipyID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) ([]core.PlanRecipy, error) {

	return selectPlanRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"plan_recipes.recipy_id": id},
			)
		},
	)
}

// getPlanRecipesByPlanID selects plan recipes by the plan id.
func getPlanRecipesByPlanID(
	ctx context.Context,
	qc squirrel.QueryerContext,
	id xid.ID,
) ([]core.PlanRecipy, error) {

	return selectPlanRecipes(
		ctx,
		qc,
		func(sb squirrel.SelectBuilder) squirrel.SelectBuilder {
			return sb.Where(
				squirrel.Eq{"plan_recipes.plan_id": id},
			)
		},
	)
}

// deletePlanRecipes deletes all plan recipes.
func deletePlanRecipes(
	ctx context.Context,
	ec squirrel.ExecerContext,
	pid xid.ID,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Delete("plan_recipes").Where(
			squirrel.Eq{"plan_recipes.plan_id": pid},
		),
	)
	return err
}

// upsertPlanRecipy upserts plan recipy.
func upsertPlanRecipy(
	ctx context.Context,
	ec squirrel.ExecerContext,
	pr core.PlanRecipy,
) error {

	_, err := squirrel.ExecContextWith(
		ctx,
		ec,
		squirrel.Insert("plan_recipes").SetMap(map[string]interface{}{
			"plan_recipes.plan_id":   pr.PlanID,
			"plan_recipes.recipy_id": pr.RecipyID,
			"plan_recipes.quantity":  pr.Quantity,
		}).Suffix("ON DUPLICATE KEY UPDATE plan_recipes.quantity = VALUES(plan_recipes.quantity)"),
	)
	return err
}

// selectPlanRecipes selects all plan recipes by the provided decorator
// function.
func selectPlanRecipes(
	ctx context.Context,
	qc squirrel.QueryerContext,
	dec func(squirrel.SelectBuilder) squirrel.SelectBuilder,
) ([]core.PlanRecipy, error) {

	rows, err := squirrel.QueryContextWith(ctx, qc, dec(squirrel.
		Select(
			"plan_recipes.plan_id",
			"plan_recipes.recipy_id",
			"plan_recipes.quantity",
		).From("plan_recipes"),
	))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]core.PlanRecipy, 0)
	for rows.Next() {
		var pr core.PlanRecipy
		if err := rows.Scan(
			&pr.PlanID,
			&pr.RecipyID,
			&pr.Quantity,
		); err != nil {
			return nil, err
		}

		prs = append(prs, pr)
	}

	return prs, nil
}
