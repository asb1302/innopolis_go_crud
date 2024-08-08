package service

import (
	"context"
	"crud/internal/domain"
	"crud/internal/repository/cache"
	"sync"
	"testing"
)

func setupTestCache(t *testing.T) *cache.RecipeCache {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	cache, err := cache.RecipeCacheInit(ctx, &wg)
	if err != nil {
		t.Fatalf("Error initializing cache: %v", err)
	}

	cache.Set("1", &domain.Recipe{ID: "1", Name: "Recipe 1"})
	cache.Set("2", &domain.Recipe{ID: "2", Name: "Recipe 2"})
	cache.Set("3", &domain.Recipe{ID: "3", Name: "Recipe 3"})
	cache.Set("4", &domain.Recipe{ID: "4", Name: "Recipe 4"})
	cache.Set("5", &domain.Recipe{ID: "5", Name: "Recipe 5"})

	return cache
}

func TestGetPaginated(t *testing.T) {
	cache := setupTestCache(t)
	Init(cache)

	tests := []struct {
		page     int
		limit    int
		expected []domain.Recipe
	}{
		{
			page:  1,
			limit: 2,
			expected: []domain.Recipe{
				{ID: "1", Name: "Recipe 1"},
				{ID: "2", Name: "Recipe 2"},
			},
		},
		{
			page:  2,
			limit: 2,
			expected: []domain.Recipe{
				{ID: "3", Name: "Recipe 3"},
				{ID: "4", Name: "Recipe 4"},
			},
		},
		{
			page:  3,
			limit: 2,
			expected: []domain.Recipe{
				{ID: "5", Name: "Recipe 5"},
			},
		},
		{
			page:     4,
			limit:    2,
			expected: []domain.Recipe{},
		},
		{
			page:     -1,
			limit:    2,
			expected: nil, // Ожидается ошибка
		},
		{
			page:     1,
			limit:    -2,
			expected: nil, // Ожидается ошибка
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			pageRecipes, err := GetPaginated(test.page, test.limit)

			if (err != nil) != (test.expected == nil) {
				t.Errorf("error: %v", err)
				return
			}

			if err == nil {
				if len(pageRecipes) != len(test.expected) {
					t.Errorf("expected %d recipes, got %d", len(test.expected), len(pageRecipes))
				}

				for i, recipe := range pageRecipes {
					if recipe.ID != test.expected[i].ID || recipe.Name != test.expected[i].Name {
						t.Errorf("recipe %d: expected %+v, got %+v", i, test.expected[i], recipe)
					}
				}
			}
		})
	}
}
