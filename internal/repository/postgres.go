package repository

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/failfilin/otus_hw/internal/repository/sqlcdb"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var pgDB *sql.DB
var pgQueries *sqlcdb.Queries

// InitPostgres connects to PostgreSQL, runs goose migrations, and loads all
// collections into the in-memory SafeSlices used by existing handlers.
func InitPostgres() error {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/otus?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	pgDB = db
	pgQueries = sqlcdb.New(db)

	if err := runMigrations(db); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL:", dsn)
	return nil
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, "migrations")
}

// LoadCollectionsFromPostgres fills the global SafeSlices from PostgreSQL.
func LoadCollectionsFromPostgres() {
	ctx := context.Background()

	if rows, err := pgQueries.ListRestaurants(ctx); err == nil {
		RestSlice.Mu.Lock()
		RestSlice.Slice = make([]models.Restaurant, 0, len(rows))
		for _, r := range rows {
			RestSlice.Slice = append(RestSlice.Slice, pgRestaurantToModel(r))
		}
		RestSlice.Mu.Unlock()
	} else {
		log.Printf("Warning: failed to load restaurants from postgres: %v", err)
	}

	if rows, err := pgQueries.ListMenus(ctx); err == nil {
		MenuSlice.Mu.Lock()
		MenuSlice.Slice = make([]models.Menu, 0, len(rows))
		for _, r := range rows {
			MenuSlice.Slice = append(MenuSlice.Slice, pgMenuRowToModel(r))
		}
		MenuSlice.Mu.Unlock()
	} else {
		log.Printf("Warning: failed to load menus from postgres: %v", err)
	}

	if rows, err := pgQueries.ListDishes(ctx); err == nil {
		DishSlice.Mu.Lock()
		DishSlice.Slice = make([]models.Dish, 0, len(rows))
		for _, r := range rows {
			DishSlice.Slice = append(DishSlice.Slice, pgDishToModel(r))
		}
		DishSlice.Mu.Unlock()
	} else {
		log.Printf("Warning: failed to load dishes from postgres: %v", err)
	}
}

// --- Restaurant CRUD ---

func PgInsertRestaurant(r models.Restaurant) error {
	if pgQueries == nil {
		return nil
	}
	_, err := pgQueries.CreateRestaurant(context.Background(), sqlcdb.CreateRestaurantParams{
		ID:      r.Id,
		Name:    r.Name,
		Logo:    r.Logo,
		Active:  r.Active,
		Phone:   sql.NullString{String: r.Phone, Valid: r.Phone != ""},
		Address: sql.NullString{String: r.Address, Valid: r.Address != ""},
	})
	return err
}

func PgUpdateRestaurant(r models.Restaurant) error {
	if pgQueries == nil {
		return nil
	}
	_, err := pgQueries.UpdateRestaurant(context.Background(), sqlcdb.UpdateRestaurantParams{
		ID:      r.Id,
		Name:    r.Name,
		Logo:    r.Logo,
		Active:  r.Active,
		Phone:   sql.NullString{String: r.Phone, Valid: r.Phone != ""},
		Address: sql.NullString{String: r.Address, Valid: r.Address != ""},
	})
	return err
}

func PgDeleteRestaurant(id uuid.UUID) error {
	if pgQueries == nil {
		return nil
	}
	return pgQueries.DeleteRestaurant(context.Background(), id)
}

// --- Menu CRUD ---

func PgInsertMenu(m models.Menu) error {
	if pgQueries == nil {
		return nil
	}
	_, err := pgQueries.CreateMenu(context.Background(), sqlcdb.CreateMenuParams{
		ID:           m.Id,
		RestaurantID: m.RestaurantId,
		Name:         m.Name,
		Description:  m.Description,
		Active:       m.Active,
	})
	return err
}

func PgDeleteMenu(id uuid.UUID) error {
	if pgQueries == nil {
		return nil
	}
	return pgQueries.DeleteMenu(context.Background(), id)
}

// --- Dish CRUD ---

func PgInsertDish(d models.Dish) (int32, error) {
	if pgQueries == nil {
		return 0, nil
	}
	row, err := pgQueries.CreateDish(context.Background(), sqlcdb.CreateDishParams{
		MenuID:        d.MenuID,
		Name:          d.Name,
		Compound:      d.Compound,
		Price:         fmt.Sprintf("%.2f", d.Price),
		Category:      d.Category,
		Calories:      d.Macros.Calories,
		Proteins:      d.Macros.Proteins,
		Fats:          d.Macros.Fats,
		Carbohydrates: d.Macros.Carbohydrates,
		WeightGrams:   d.WeightGrams,
	})
	if err != nil {
		return 0, err
	}
	return row.ID, nil
}

func PgDeleteDish(id int32) error {
	if pgQueries == nil {
		return nil
	}
	return pgQueries.DeleteDish(context.Background(), id)
}

// --- Converters ---

func pgRestaurantToModel(r sqlcdb.Restaurant) models.Restaurant {
	return models.Restaurant{
		Id:        r.ID,
		Name:      r.Name,
		Logo:      r.Logo,
		Active:    r.Active,
		Phone:     r.Phone.String,
		Address:   r.Address.String,
		CreatedAt: r.CreatedAt,
	}
}

func pgMenuRowToModel(r sqlcdb.ListMenusRow) models.Menu {
	return models.Menu{
		Id:           r.ID,
		RestaurantId: r.RestaurantID,
		Name:         r.Name,
		Description:  r.Description,
		Active:       r.Active,
		CreatedAt:    r.CreatedAt,
	}
}

func pgDishToModel(r sqlcdb.Dish) models.Dish {
	var price float32
	fmt.Sscanf(r.Price, "%f", &price)
	return models.Dish{
		Id:          int(r.ID),
		MenuID:      r.MenuID,
		Name:        r.Name,
		Compound:    r.Compound,
		Price:       price,
		Category:    r.Category,
		WeightGrams: r.WeightGrams,
		CreatedAt:   r.CreatedAt,
		Macros: models.Macros{
			Calories:      r.Calories,
			Proteins:      r.Proteins,
			Fats:          r.Fats,
			Carbohydrates: r.Carbohydrates,
		},
	}
}
