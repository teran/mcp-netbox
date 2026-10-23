package mockrepo

import (
	"context"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// Coverage for the Batch 5b entities (get_roles, get_tenants, Tag CRUD):
// MockRepo delegates forward to their Func fields and NilRepo returns
// empty/sentinel values. Both implementations are exercised so the batch's
// repository methods count toward the coverage gate.

func TestMockRepo_Batch5bDelegates(t *testing.T) { //nolint:funlen
	t.Parallel()

	t.Run("ListRoles", func(t *testing.T) {
		got, err := (&MockRepo{
			ListRolesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				return &domain.PaginatedResponse[domain.Role]{Count: 1}, nil
			},
		}).ListRoles(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListRoles = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("ListTenants", func(t *testing.T) {
		got, err := (&MockRepo{
			ListTenantsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				return &domain.PaginatedResponse[domain.Tenant]{Count: 1}, nil
			},
		}).ListTenants(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListTenants = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("ListTags", func(t *testing.T) {
		got, err := (&MockRepo{
			ListTagsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				return &domain.PaginatedResponse[domain.Tag]{Count: 1}, nil
			},
		}).ListTags(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListTags = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateTag", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateTagFunc: func(
				_ context.Context, _ string, in domain.TagWrite,
			) (*domain.Tag, error) {
				return &domain.Tag{ID: 1, Name: in.Name}, nil
			},
		}).CreateTag(context.Background(), "t", domain.TagWrite{Name: "prod"})
		if err != nil || got == nil || got.Name != "prod" {
			t.Errorf("CreateTag = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateTag", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateTagFunc: func(
				_ context.Context, _ string, id int, in domain.TagWrite,
			) (*domain.Tag, error) {
				return &domain.Tag{ID: id, Name: in.Name}, nil
			},
		}).UpdateTag(context.Background(), "t", 7, domain.TagWrite{Name: "Renamed"})
		if err != nil || got == nil || got.ID != 7 || got.Name != "Renamed" {
			t.Errorf("UpdateTag = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteTag", func(t *testing.T) {
		err := (&MockRepo{
			DeleteTagFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteTag(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteTag = %v", err)
		}
	})
}

func TestNilRepo_Batch5b(t *testing.T) {
	t.Parallel()

	repo := &NilRepo{}

	t.Run("ListRoles", func(t *testing.T) {
		got, err := repo.ListRoles(context.Background(), "t", nil)
		if err != nil || got == nil || len(got.Results) != 0 {
			t.Errorf("ListRoles = (%v, %v), want empty page", got, err)
		}
	})

	t.Run("ListTenants", func(t *testing.T) {
		got, err := repo.ListTenants(context.Background(), "t", nil)
		if err != nil || got == nil || len(got.Results) != 0 {
			t.Errorf("ListTenants = (%v, %v), want empty page", got, err)
		}
	})

	t.Run("ListTags", func(t *testing.T) {
		got, err := repo.ListTags(context.Background(), "t", nil)
		if err != nil || got == nil || len(got.Results) != 0 {
			t.Errorf("ListTags = (%v, %v), want empty page", got, err)
		}
	})

	t.Run("CreateTag", func(t *testing.T) {
		got, err := repo.CreateTag(context.Background(), "t", domain.TagWrite{Name: "prod"})
		if err != nil || got == nil || got.ID != 1 {
			t.Errorf("CreateTag = (%v, %v), want id 1", got, err)
		}
	})

	t.Run("UpdateTag", func(t *testing.T) {
		got, err := repo.UpdateTag(context.Background(), "t", 7, domain.TagWrite{Name: "X"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateTag = (%v, %v), want id 7", got, err)
		}
	})

	t.Run("DeleteTag", func(t *testing.T) {
		if err := repo.DeleteTag(context.Background(), "t", 3); err != nil {
			t.Errorf("DeleteTag = %v", err)
		}
	})
}
