package cache

import (
	"crud/internal/domain"
	"testing"
)

func TestGetAll(t *testing.T) {
	cache := &RecipeCache{
		pool: make(map[string]*domain.Recipe),
	}

	cache.pool["1"] = &domain.Recipe{ID: "1", Name: "Recipe 1"}
	cache.pool["2"] = &domain.Recipe{ID: "2", Name: "Recipe 2"}
	cache.pool["3"] = &domain.Recipe{ID: "3", Name: "Recipe 3"}

	allRecipes, err := cache.GetAll()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if len(allRecipes) != 3 {
		t.Errorf("expected 3 recipes, got %d", len(allRecipes))
	}

	expectedRecipes := map[string]*domain.Recipe{
		"1": {ID: "1", Name: "Recipe 1"},
		"2": {ID: "2", Name: "Recipe 2"},
		"3": {ID: "3", Name: "Recipe 3"},
	}

	for id, expectedRecipe := range expectedRecipes {
		recipe, exists := allRecipes[id]
		if !exists {
			t.Errorf("expected recipe with ID %s not found", id)
			continue
		}
		if recipe.ID != expectedRecipe.ID || recipe.Name != expectedRecipe.Name {
			t.Errorf("recipe with ID %s: expected %+v, got %+v", id, expectedRecipe, recipe)
		}
	}
}
