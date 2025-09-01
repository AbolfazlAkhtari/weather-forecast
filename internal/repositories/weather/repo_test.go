package weather

import (
	"context"
	"testing"
	"time"

	"github.com/AbolfazlAkhtari/weather-forecast/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Weather{})
	require.NoError(t, err)

	return db
}

func createTestWeather() *models.Weather {
	return &models.Weather{
		CityName:    "Tehran",
		Country:     "Iran",
		Temperature: 25.5,
		Description: "Sunny",
		Humidity:    60,
		WindSpeed:   10.5,
		FetchedAt:   time.Now(),
	}
}

func TestNewRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		weather := createTestWeather()

		err := repo.Create(ctx, weather)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, weather.ID)
		assert.False(t, weather.CreatedAt.IsZero())
		assert.False(t, weather.UpdatedAt.IsZero())

		// Verify it was actually saved
		var savedWeather models.Weather
		err = db.First(&savedWeather, weather.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, weather.CityName, savedWeather.CityName)
	})

	t.Run("creation with nil weather", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
	})
}

func TestRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		weather := createTestWeather()
		err := repo.Create(ctx, weather)
		require.NoError(t, err)

		updates := map[string]interface{}{
			"temperature": 30.0,
			"description": "Hot",
			"humidity":    70,
		}

		err = repo.Update(ctx, weather.ID, updates)
		assert.NoError(t, err)

		// Verify updates
		var updatedWeather models.Weather
		err = db.First(&updatedWeather, weather.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, 30.0, updatedWeather.Temperature)
		assert.Equal(t, "Hot", updatedWeather.Description)
		assert.Equal(t, 70, updatedWeather.Humidity)
	})

	t.Run("update non-existent record", func(t *testing.T) {
		nonExistentID := uuid.New()
		updates := map[string]interface{}{
			"temperature": 30.0,
		}

		err := repo.Update(ctx, nonExistentID, updates)
		assert.NoError(t, err) // GORM doesn't return error for no rows affected
	})

	t.Run("update with empty updates map", func(t *testing.T) {
		weather := createTestWeather()
		err := repo.Create(ctx, weather)
		require.NoError(t, err)

		err = repo.Update(ctx, weather.ID, map[string]interface{}{})
		assert.NoError(t, err)
	})
}

func TestRepository_PaginatedList(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	// Create test data
	weathers := []*models.Weather{
		{CityName: "Tehran", Country: "Iran", Temperature: 25.0, Description: "Sunny", Humidity: 60, WindSpeed: 10.0, FetchedAt: time.Now()},
		{CityName: "Mashhad", Country: "Iran", Temperature: 28.0, Description: "Cloudy", Humidity: 65, WindSpeed: 12.0, FetchedAt: time.Now()},
		{CityName: "Isfahan", Country: "Iran", Temperature: 22.0, Description: "Clear", Humidity: 55, WindSpeed: 8.0, FetchedAt: time.Now()},
		{CityName: "Shiraz", Country: "Iran", Temperature: 30.0, Description: "Hot", Humidity: 70, WindSpeed: 15.0, FetchedAt: time.Now()},
		{CityName: "Tabriz", Country: "Iran", Temperature: 18.0, Description: "Cool", Humidity: 50, WindSpeed: 6.0, FetchedAt: time.Now()},
	}

	for _, w := range weathers {
		err := repo.Create(ctx, w)
		require.NoError(t, err)
		// Add some delay to ensure different timestamps
		time.Sleep(1 * time.Millisecond)
	}

	t.Run("first page", func(t *testing.T) {
		results, totalPage, count, err := repo.PaginatedList(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.Equal(t, int64(1), totalPage)
		assert.Len(t, results, 5)
		// Should be ordered by created_at desc
		assert.Equal(t, "Tabriz", results[0].CityName)
		assert.Equal(t, "Shiraz", results[1].CityName)
	})

	t.Run("page with limit", func(t *testing.T) {
		// Create more records to test pagination
		for i := 0; i < 15; i++ {
			weather := &models.Weather{
				CityName:    "City" + string(rune('A'+i)),
				Country:     "Iran",
				Temperature: float64(20 + i),
				Description: "Test",
				Humidity:    60,
				WindSpeed:   10.0,
				FetchedAt:   time.Now(),
			}
			err := repo.Create(ctx, weather)
			require.NoError(t, err)
		}

		results, totalPage, count, err := repo.PaginatedList(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(20), count)    // 5 original + 15 new
		assert.Equal(t, int64(2), totalPage) // 20 records / 10 limit = 2 pages
		assert.Len(t, results, 10)           // First page should have 10 records
	})

	t.Run("second page", func(t *testing.T) {
		results, totalPage, count, err := repo.PaginatedList(ctx, 2)

		assert.NoError(t, err)
		assert.Equal(t, int64(20), count)
		assert.Equal(t, int64(2), totalPage)
		assert.Len(t, results, 10) // Second page should have 10 records
	})

	t.Run("page beyond available data", func(t *testing.T) {
		results, totalPage, count, err := repo.PaginatedList(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, int64(20), count)
		assert.Equal(t, int64(2), totalPage)
		assert.Len(t, results, 0) // No results for page 5
	})

	t.Run("page 0", func(t *testing.T) {
		results, totalPage, count, err := repo.PaginatedList(ctx, 0)

		assert.NoError(t, err)
		assert.Equal(t, int64(20), count)
		assert.Equal(t, int64(2), totalPage)
		assert.Len(t, results, 10) // Should default to first page
	})
}

func TestRepository_LatestByCityName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("find latest by city name", func(t *testing.T) {
		// Create multiple records for the same city
		weather1 := &models.Weather{
			CityName:    "Tehran",
			Country:     "Iran",
			Temperature: 25.0,
			Description: "Sunny",
			Humidity:    60,
			WindSpeed:   10.0,
			FetchedAt:   time.Now().Add(-1 * time.Hour),
		}
		weather2 := &models.Weather{
			CityName:    "Tehran",
			Country:     "Iran",
			Temperature: 28.0,
			Description: "Hot",
			Humidity:    70,
			WindSpeed:   12.0,
			FetchedAt:   time.Now(),
		}

		err := repo.Create(ctx, weather1)
		require.NoError(t, err)
		err = repo.Create(ctx, weather2)
		require.NoError(t, err)

		result, err := repo.LatestByCityName(ctx, "Tehran")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Tehran", result.CityName)
		assert.Equal(t, 28.0, result.Temperature) // Should be the latest one
		assert.Equal(t, "Hot", result.Description)
	})

	t.Run("case insensitive search", func(t *testing.T) {
		weather := &models.Weather{
			CityName:    "Mashhad",
			Country:     "Iran",
			Temperature: 22.0,
			Description: "Cool",
			Humidity:    55,
			WindSpeed:   8.0,
			FetchedAt:   time.Now(),
		}
		err := repo.Create(ctx, weather)
		require.NoError(t, err)

		result, err := repo.LatestByCityName(ctx, "mashhad")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Mashhad", result.CityName)

		result, err = repo.LatestByCityName(ctx, "MASHHAD")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Mashhad", result.CityName)
	})

	t.Run("city not found", func(t *testing.T) {
		_, err := repo.LatestByCityName(ctx, "NonExistentCity")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})
}

func TestRepository_FindById(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("find existing record", func(t *testing.T) {
		weather := createTestWeather()
		err := repo.Create(ctx, weather)
		require.NoError(t, err)

		result, err := repo.FindById(ctx, weather.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, weather.ID, result.ID)
		assert.Equal(t, weather.CityName, result.CityName)
		assert.Equal(t, weather.Temperature, result.Temperature)
	})

	t.Run("find non-existent record", func(t *testing.T) {
		nonExistentID := uuid.New()

		_, err := repo.FindById(ctx, nonExistentID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})
}

func TestRepository_DeleteById(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("delete existing record", func(t *testing.T) {
		weather := createTestWeather()
		err := repo.Create(ctx, weather)
		require.NoError(t, err)

		err = repo.DeleteById(ctx, weather.ID)
		assert.NoError(t, err)

		// Verify it was deleted
		var deletedWeather models.Weather
		err = db.First(&deletedWeather, weather.ID).Error
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "record not found")
	})

	t.Run("delete non-existent record", func(t *testing.T) {
		nonExistentID := uuid.New()

		err := repo.DeleteById(ctx, nonExistentID)
		assert.NoError(t, err) // GORM doesn't return error for no rows affected
	})
}
