package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the entity added in feature/expand-entities
// Batch 5c (the CustomField entity CRUD).
//
// These cover the thin relay methods in application/service.go and exercise the
// mockrepo methods they delegate to. Each method gets a success + an error
// propagation case so both branches of every wrapper are covered.

func TestNetworkService_ListCustomFields(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCustomFieldsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				return &domain.PaginatedResponse[domain.CustomField]{
					Count:   1,
					Results: []domain.CustomField{{ID: 1, Name: "env", Label: "Environment"}},
				}, nil
			},
		})
		resp, err := svc.ListCustomFields(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCustomFields() returned error: %v", err)
		}
		if resp.Count != 1 || len(resp.Results) != 1 || resp.Results[0].Name != "env" {
			t.Errorf("resp = %+v, want count 1 with env", resp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCustomFieldsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CustomField], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListCustomFields(context.Background(), nil); err == nil {
			t.Fatal("ListCustomFields() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateCustomField(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCustomFieldFunc: func(
				_ context.Context, _ string, in domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return &domain.CustomField{ID: 1, Name: in.Name, Label: in.Label, Type: "text"}, nil
			},
		})
		cf, err := svc.CreateCustomField(context.Background(), domain.CustomFieldWrite{Name: "env", Label: "Environment"})
		if err != nil {
			t.Fatalf("CreateCustomField() returned error: %v", err)
		}
		if cf.ID != 1 || cf.Name != "env" || cf.Label != "Environment" {
			t.Errorf("cf = %+v, want id 1 name env label Environment", cf)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCustomFieldFunc: func(
				_ context.Context, _ string, _ domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateCustomField(context.Background(), domain.CustomFieldWrite{}); err == nil {
			t.Fatal("CreateCustomField() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateCustomField(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCustomFieldFunc: func(
				_ context.Context, _ string, id int, in domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return &domain.CustomField{ID: id, Name: in.Name, Label: "Renamed"}, nil
			},
		})
		cf, err := svc.UpdateCustomField(context.Background(), 7, domain.CustomFieldWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateCustomField() returned error: %v", err)
		}
		if cf.ID != 7 || cf.Name != "Renamed" {
			t.Errorf("cf = %+v, want id 7 name Renamed", cf)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCustomFieldFunc: func(
				_ context.Context, _ string, _ int, _ domain.CustomFieldWrite,
			) (*domain.CustomField, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateCustomField(context.Background(), 7, domain.CustomFieldWrite{}); err == nil {
			t.Fatal("UpdateCustomField() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteCustomField(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCustomFieldFunc: func(_ context.Context, _ string, _ int) error {
				return nil
			},
		})
		if err := svc.DeleteCustomField(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCustomField() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCustomFieldFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteCustomField(context.Background(), 3); err == nil {
			t.Fatal("DeleteCustomField() expected error, got nil")
		}
	})
}
