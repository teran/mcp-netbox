package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the entities added in feature/expand-entities
// Batch 5b (get_roles, get_tenants, and the Tag entity CRUD).
//
// These cover the thin relay methods in application/service.go and exercise the
// mockrepo methods they delegate to. Each method gets a success + an error
// propagation case so both branches of every wrapper are covered.

func TestNetworkService_ListRoles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				return &domain.PaginatedResponse[domain.Role]{
					Count:   1,
					Results: []domain.Role{{ID: 1, Name: "netbox"}},
				}, nil
			},
		})
		resp, err := svc.ListRoles(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRoles() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "netbox" {
			t.Errorf("resp = %+v, want count 1 with netbox", resp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Role], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListRoles(context.Background(), nil); err == nil {
			t.Fatal("ListRoles() expected error, got nil")
		}
	})
}

func TestNetworkService_ListTenants(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTenantsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				return &domain.PaginatedResponse[domain.Tenant]{
					Count:   1,
					Results: []domain.Tenant{{ID: 1, Name: "acme"}},
				}, nil
			},
		})
		resp, err := svc.ListTenants(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListTenants() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "acme" {
			t.Errorf("resp = %+v, want count 1 with acme", resp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTenantsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tenant], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListTenants(context.Background(), nil); err == nil {
			t.Fatal("ListTenants() expected error, got nil")
		}
	})
}

func TestNetworkService_ListTags(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTagsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				return &domain.PaginatedResponse[domain.Tag]{
					Count:   1,
					Results: []domain.Tag{{ID: 1, Name: "prod"}},
				}, nil
			},
		})
		resp, err := svc.ListTags(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListTags() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "prod" {
			t.Errorf("resp = %+v, want count 1 with prod", resp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTagsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Tag], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListTags(context.Background(), nil); err == nil {
			t.Fatal("ListTags() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateTag(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTagFunc: func(
				_ context.Context, _ string, in domain.TagWrite,
			) (*domain.Tag, error) {
				return &domain.Tag{ID: 1, Name: in.Name, Slug: "prod"}, nil
			},
		})
		tag, err := svc.CreateTag(context.Background(), domain.TagWrite{Name: "prod"})
		if err != nil {
			t.Fatalf("CreateTag() returned error: %v", err)
		}
		if tag.ID != 1 || tag.Name != "prod" || tag.Slug != "prod" {
			t.Errorf("tag = %+v, want id 1 name prod slug prod", tag)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTagFunc: func(
				_ context.Context, _ string, _ domain.TagWrite,
			) (*domain.Tag, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateTag(context.Background(), domain.TagWrite{}); err == nil {
			t.Fatal("CreateTag() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateTag(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTagFunc: func(
				_ context.Context, _ string, id int, in domain.TagWrite,
			) (*domain.Tag, error) {
				return &domain.Tag{ID: id, Name: in.Name, Slug: "renamed"}, nil
			},
		})
		tag, err := svc.UpdateTag(context.Background(), 7, domain.TagWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateTag() returned error: %v", err)
		}
		if tag.ID != 7 || tag.Name != "Renamed" {
			t.Errorf("tag = %+v, want id 7 name Renamed", tag)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTagFunc: func(
				_ context.Context, _ string, _ int, _ domain.TagWrite,
			) (*domain.Tag, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateTag(context.Background(), 7, domain.TagWrite{}); err == nil {
			t.Fatal("UpdateTag() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteTag(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTagFunc: func(_ context.Context, _ string, _ int) error {
				return nil
			},
		})
		if err := svc.DeleteTag(context.Background(), 3); err != nil {
			t.Fatalf("DeleteTag() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTagFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteTag(context.Background(), 3); err == nil {
			t.Fatal("DeleteTag() expected error, got nil")
		}
	})
}
