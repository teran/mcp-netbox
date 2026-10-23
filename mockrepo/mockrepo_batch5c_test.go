package mockrepo

import (
	"context"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// Coverage for the Batch 5c entity (CustomField CRUD): MockRepo delegates
// forward to their Func fields and NilRepo returns empty/sentinel values. Both
// implementations are exercised so the batch's repository methods count toward
// the coverage gate.

func TestMockRepo_Batch5cDelegates(t *testing.T) {
	t.Parallel()

	t.Run("ListCustomFields", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCustomFieldsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				return &domain.PaginatedResponse[domain.CustomField]{Count: 1}, nil
			},
		}).ListCustomFields(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListCustomFields = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateCustomField", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCustomFieldFunc: func(
				_ context.Context, _ string, in domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return &domain.CustomField{ID: 1, Name: in.Name, Label: in.Label}, nil
			},
		}).CreateCustomField(context.Background(), "t", domain.CustomFieldWrite{Name: "env", Label: "Environment"})
		if err != nil || got == nil || got.Name != "env" || got.Label != "Environment" {
			t.Errorf("CreateCustomField = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateCustomField", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCustomFieldFunc: func(
				_ context.Context, _ string, id int, in domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return &domain.CustomField{ID: id, Name: in.Name, Label: "Renamed"}, nil
			},
		}).UpdateCustomField(context.Background(), "t", 7, domain.CustomFieldWrite{Name: "Renamed"})
		if err != nil || got == nil || got.ID != 7 || got.Name != "Renamed" {
			t.Errorf("UpdateCustomField = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteCustomField", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCustomFieldFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCustomField(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCustomField = %v", err)
		}
	})
}

func TestNilRepo_Batch5c(t *testing.T) {
	t.Parallel()

	repo := &NilRepo{}

	t.Run("ListCustomFields", func(t *testing.T) {
		got, err := repo.ListCustomFields(context.Background(), "t", nil)
		if err != nil || got == nil || len(got.Results) != 0 {
			t.Errorf("ListCustomFields = (%v, %v), want empty page", got, err)
		}
	})

	t.Run("CreateCustomField", func(t *testing.T) {
		got, err := repo.CreateCustomField(context.Background(), "t", domain.CustomFieldWrite{Name: "env"})
		if err != nil || got == nil || got.ID != 1 || got.Name != "env" {
			t.Errorf("CreateCustomField = (%v, %v), want id 1 name env", got, err)
		}
	})

	t.Run("UpdateCustomField", func(t *testing.T) {
		got, err := repo.UpdateCustomField(context.Background(), "t", 7, domain.CustomFieldWrite{Name: "X"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCustomField = (%v, %v), want id 7", got, err)
		}
	})

	t.Run("DeleteCustomField", func(t *testing.T) {
		if err := repo.DeleteCustomField(context.Background(), "t", 3); err != nil {
			t.Errorf("DeleteCustomField = %v", err)
		}
	})
}
