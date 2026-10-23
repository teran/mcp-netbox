package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the IPAM entities added in feature/expand-entities
// Batch 3: rir, aggregate, ip_range, asn, asn_range, route_target, service,
// service_template, fhrp_group, fhrp_group_assignment, vlan_translation_policy,
// vlan_translation_rule. These cover the thin relay methods in
// application/service.go (and exercise the mockrepo methods they delegate to).
// Each entity gets success + error-propagation cases so both branches of each
// wrapper are covered.

func TestNetworkService_ListRIRs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRIRsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RIR], error) {
				return &domain.PaginatedResponse[domain.RIR]{Count: 1}, nil
			},
		})
		resp, err := svc.ListRIRs(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRIRs() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRIRsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RIR], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListRIRs(context.Background(), nil); err == nil {
			t.Fatal("ListRIRs() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateRIR(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRIRFunc: func(
				_ context.Context, _ string, in domain.RIRWrite,
			) (*domain.RIR, error) {
				return &domain.RIR{ID: 1, Name: in.Name}, nil
			},
		})
		r, err := svc.CreateRIR(context.Background(), domain.RIRWrite{Name: "RIPE"})
		if err != nil {
			t.Fatalf("CreateRIR() returned error: %v", err)
		}
		if r.ID != 1 || r.Name != "RIPE" {
			t.Errorf("RIR = %+v, want id 1 name RIPE", r)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRIRFunc: func(
				_ context.Context, _ string, _ domain.RIRWrite,
			) (*domain.RIR, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateRIR(context.Background(), domain.RIRWrite{}); err == nil {
			t.Fatal("CreateRIR() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateRIR(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRIRFunc: func(
				_ context.Context, _ string, id int, _ domain.RIRWrite,
			) (*domain.RIR, error) {
				return &domain.RIR{ID: id}, nil
			},
		})
		r, err := svc.UpdateRIR(context.Background(), 7, domain.RIRWrite{})
		if err != nil {
			t.Fatalf("UpdateRIR() returned error: %v", err)
		}
		if r.ID != 7 {
			t.Errorf("RIR ID = %d, want 7", r.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRIRFunc: func(
				_ context.Context, _ string, _ int, _ domain.RIRWrite,
			) (*domain.RIR, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateRIR(context.Background(), 7, domain.RIRWrite{}); err == nil {
			t.Fatal("UpdateRIR() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteRIR(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRIRFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteRIR(context.Background(), 3); err != nil {
			t.Fatalf("DeleteRIR() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRIRFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteRIR(context.Background(), 3); err == nil {
			t.Fatal("DeleteRIR() expected error, got nil")
		}
	})
}

func TestNetworkService_ListAggregates(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListAggregatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				return &domain.PaginatedResponse[domain.Aggregate]{Count: 1}, nil
			},
		})
		resp, err := svc.ListAggregates(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListAggregates() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListAggregatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Aggregate], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListAggregates(context.Background(), nil); err == nil {
			t.Fatal("ListAggregates() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateAggregate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateAggregateFunc: func(
				_ context.Context, _ string, in domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				return &domain.Aggregate{ID: 1, Prefix: in.Prefix}, nil
			},
		})
		a, err := svc.CreateAggregate(context.Background(), domain.AggregateWrite{Prefix: "10.0.0.0/8"})
		if err != nil {
			t.Fatalf("CreateAggregate() returned error: %v", err)
		}
		if a.ID != 1 || a.Prefix != "10.0.0.0/8" {
			t.Errorf("Aggregate = %+v, want id 1 prefix 10.0.0.0/8", a)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateAggregateFunc: func(
				_ context.Context, _ string, _ domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateAggregate(context.Background(), domain.AggregateWrite{}); err == nil {
			t.Fatal("CreateAggregate() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateAggregate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateAggregateFunc: func(
				_ context.Context, _ string, id int, _ domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				return &domain.Aggregate{ID: id}, nil
			},
		})
		a, err := svc.UpdateAggregate(context.Background(), 7, domain.AggregateWrite{})
		if err != nil {
			t.Fatalf("UpdateAggregate() returned error: %v", err)
		}
		if a.ID != 7 {
			t.Errorf("Aggregate ID = %d, want 7", a.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateAggregateFunc: func(
				_ context.Context, _ string, _ int, _ domain.AggregateWrite,
			) (*domain.Aggregate, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateAggregate(context.Background(), 7, domain.AggregateWrite{}); err == nil {
			t.Fatal("UpdateAggregate() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteAggregate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteAggregateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteAggregate(context.Background(), 3); err != nil {
			t.Fatalf("DeleteAggregate() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteAggregateFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteAggregate(context.Background(), 3); err == nil {
			t.Fatal("DeleteAggregate() expected error, got nil")
		}
	})
}

func TestNetworkService_ListIPRanges(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListIPRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				return &domain.PaginatedResponse[domain.IPRange]{Count: 1}, nil
			},
		})
		resp, err := svc.ListIPRanges(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListIPRanges() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListIPRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.IPRange], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListIPRanges(context.Background(), nil); err == nil {
			t.Fatal("ListIPRanges() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateIPRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateIPRangeFunc: func(
				_ context.Context, _ string, in domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				return &domain.IPRange{ID: 1, StartAddress: in.StartAddress}, nil
			},
		})
		r, err := svc.CreateIPRange(context.Background(), domain.IPRangeWrite{StartAddress: "10.0.0.1"})
		if err != nil {
			t.Fatalf("CreateIPRange() returned error: %v", err)
		}
		if r.ID != 1 || r.StartAddress != "10.0.0.1" {
			t.Errorf("IPRange = %+v, want id 1 start 10.0.0.1", r)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateIPRangeFunc: func(
				_ context.Context, _ string, _ domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateIPRange(context.Background(), domain.IPRangeWrite{}); err == nil {
			t.Fatal("CreateIPRange() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateIPRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateIPRangeFunc: func(
				_ context.Context, _ string, id int, _ domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				return &domain.IPRange{ID: id}, nil
			},
		})
		r, err := svc.UpdateIPRange(context.Background(), 7, domain.IPRangeWrite{})
		if err != nil {
			t.Fatalf("UpdateIPRange() returned error: %v", err)
		}
		if r.ID != 7 {
			t.Errorf("IPRange ID = %d, want 7", r.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateIPRangeFunc: func(
				_ context.Context, _ string, _ int, _ domain.IPRangeWrite,
			) (*domain.IPRange, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateIPRange(context.Background(), 7, domain.IPRangeWrite{}); err == nil {
			t.Fatal("UpdateIPRange() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteIPRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteIPRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteIPRange(context.Background(), 3); err != nil {
			t.Fatalf("DeleteIPRange() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteIPRangeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteIPRange(context.Background(), 3); err == nil {
			t.Fatal("DeleteIPRange() expected error, got nil")
		}
	})
}

func TestNetworkService_ListASNs(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListASNsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASN], error) {
				return &domain.PaginatedResponse[domain.ASN]{Count: 1}, nil
			},
		})
		resp, err := svc.ListASNs(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListASNs() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListASNsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASN], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListASNs(context.Background(), nil); err == nil {
			t.Fatal("ListASNs() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateASN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateASNFunc: func(
				_ context.Context, _ string, in domain.ASNWrite,
			) (*domain.ASN, error) {
				return &domain.ASN{ID: 1, ASN: in.ASN}, nil
			},
		})
		a, err := svc.CreateASN(context.Background(), domain.ASNWrite{ASN: 64512})
		if err != nil {
			t.Fatalf("CreateASN() returned error: %v", err)
		}
		if a.ID != 1 || a.ASN != 64512 {
			t.Errorf("ASN = %+v, want id 1 asn 64512", a)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateASNFunc: func(
				_ context.Context, _ string, _ domain.ASNWrite,
			) (*domain.ASN, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateASN(context.Background(), domain.ASNWrite{}); err == nil {
			t.Fatal("CreateASN() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateASN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateASNFunc: func(
				_ context.Context, _ string, id int, _ domain.ASNWrite,
			) (*domain.ASN, error) {
				return &domain.ASN{ID: id}, nil
			},
		})
		a, err := svc.UpdateASN(context.Background(), 7, domain.ASNWrite{})
		if err != nil {
			t.Fatalf("UpdateASN() returned error: %v", err)
		}
		if a.ID != 7 {
			t.Errorf("ASN ID = %d, want 7", a.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateASNFunc: func(
				_ context.Context, _ string, _ int, _ domain.ASNWrite,
			) (*domain.ASN, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateASN(context.Background(), 7, domain.ASNWrite{}); err == nil {
			t.Fatal("UpdateASN() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteASN(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteASNFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteASN(context.Background(), 3); err != nil {
			t.Fatalf("DeleteASN() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteASNFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteASN(context.Background(), 3); err == nil {
			t.Fatal("DeleteASN() expected error, got nil")
		}
	})
}

func TestNetworkService_ListASNRanges(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListASNRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				return &domain.PaginatedResponse[domain.ASNRange]{Count: 1}, nil
			},
		})
		resp, err := svc.ListASNRanges(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListASNRanges() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListASNRangesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ASNRange], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListASNRanges(context.Background(), nil); err == nil {
			t.Fatal("ListASNRanges() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateASNRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateASNRangeFunc: func(
				_ context.Context, _ string, in domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				return &domain.ASNRange{ID: 1, Name: in.Name}, nil
			},
		})
		r, err := svc.CreateASNRange(context.Background(), domain.ASNRangeWrite{Name: "PRIVATE"})
		if err != nil {
			t.Fatalf("CreateASNRange() returned error: %v", err)
		}
		if r.ID != 1 || r.Name != "PRIVATE" {
			t.Errorf("ASNRange = %+v, want id 1 name PRIVATE", r)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateASNRangeFunc: func(
				_ context.Context, _ string, _ domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateASNRange(context.Background(), domain.ASNRangeWrite{}); err == nil {
			t.Fatal("CreateASNRange() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateASNRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateASNRangeFunc: func(
				_ context.Context, _ string, id int, _ domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				return &domain.ASNRange{ID: id}, nil
			},
		})
		r, err := svc.UpdateASNRange(context.Background(), 7, domain.ASNRangeWrite{})
		if err != nil {
			t.Fatalf("UpdateASNRange() returned error: %v", err)
		}
		if r.ID != 7 {
			t.Errorf("ASNRange ID = %d, want 7", r.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateASNRangeFunc: func(
				_ context.Context, _ string, _ int, _ domain.ASNRangeWrite,
			) (*domain.ASNRange, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateASNRange(context.Background(), 7, domain.ASNRangeWrite{}); err == nil {
			t.Fatal("UpdateASNRange() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteASNRange(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteASNRangeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteASNRange(context.Background(), 3); err != nil {
			t.Fatalf("DeleteASNRange() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteASNRangeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteASNRange(context.Background(), 3); err == nil {
			t.Fatal("DeleteASNRange() expected error, got nil")
		}
	})
}

func TestNetworkService_ListRouteTargets(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRouteTargetsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				return &domain.PaginatedResponse[domain.RouteTarget]{Count: 1}, nil
			},
		})
		resp, err := svc.ListRouteTargets(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRouteTargets() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRouteTargetsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RouteTarget], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListRouteTargets(context.Background(), nil); err == nil {
			t.Fatal("ListRouteTargets() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateRouteTarget(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRouteTargetFunc: func(
				_ context.Context, _ string, in domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return &domain.RouteTarget{ID: 1, Name: in.Name}, nil
			},
		})
		rt, err := svc.CreateRouteTarget(context.Background(), domain.RouteTargetWrite{Name: "65000:100"})
		if err != nil {
			t.Fatalf("CreateRouteTarget() returned error: %v", err)
		}
		if rt.ID != 1 || rt.Name != "65000:100" {
			t.Errorf("RouteTarget = %+v, want id 1 name 65000:100", rt)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRouteTargetFunc: func(
				_ context.Context, _ string, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateRouteTarget(context.Background(), domain.RouteTargetWrite{}); err == nil {
			t.Fatal("CreateRouteTarget() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateRouteTarget(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRouteTargetFunc: func(
				_ context.Context, _ string, id int, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return &domain.RouteTarget{ID: id}, nil
			},
		})
		rt, err := svc.UpdateRouteTarget(context.Background(), 7, domain.RouteTargetWrite{})
		if err != nil {
			t.Fatalf("UpdateRouteTarget() returned error: %v", err)
		}
		if rt.ID != 7 {
			t.Errorf("RouteTarget ID = %d, want 7", rt.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRouteTargetFunc: func(
				_ context.Context, _ string, _ int, _ domain.RouteTargetWrite,
			) (*domain.RouteTarget, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateRouteTarget(context.Background(), 7, domain.RouteTargetWrite{}); err == nil {
			t.Fatal("UpdateRouteTarget() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteRouteTarget(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRouteTargetFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteRouteTarget(context.Background(), 3); err != nil {
			t.Fatalf("DeleteRouteTarget() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRouteTargetFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteRouteTarget(context.Background(), 3); err == nil {
			t.Fatal("DeleteRouteTarget() expected error, got nil")
		}
	})
}

func TestNetworkService_ListServices(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListServicesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				return &domain.PaginatedResponse[domain.Service]{Count: 1}, nil
			},
		})
		resp, err := svc.ListServices(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListServices() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListServicesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.Service], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListServices(context.Background(), nil); err == nil {
			t.Fatal("ListServices() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateService(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateServiceFunc: func(
				_ context.Context, _ string, in domain.ServiceWrite,
			) (*domain.Service, error) {
				return &domain.Service{ID: 1, Name: in.Name}, nil
			},
		})
		s, err := svc.CreateService(context.Background(), domain.ServiceWrite{Name: "web"})
		if err != nil {
			t.Fatalf("CreateService() returned error: %v", err)
		}
		if s.ID != 1 || s.Name != "web" {
			t.Errorf("Service = %+v, want id 1 name web", s)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateServiceFunc: func(
				_ context.Context, _ string, _ domain.ServiceWrite,
			) (*domain.Service, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateService(context.Background(), domain.ServiceWrite{}); err == nil {
			t.Fatal("CreateService() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateService(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateServiceFunc: func(
				_ context.Context, _ string, id int, _ domain.ServiceWrite,
			) (*domain.Service, error) {
				return &domain.Service{ID: id}, nil
			},
		})
		s, err := svc.UpdateService(context.Background(), 7, domain.ServiceWrite{})
		if err != nil {
			t.Fatalf("UpdateService() returned error: %v", err)
		}
		if s.ID != 7 {
			t.Errorf("Service ID = %d, want 7", s.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateServiceFunc: func(
				_ context.Context, _ string, _ int, _ domain.ServiceWrite,
			) (*domain.Service, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateService(context.Background(), 7, domain.ServiceWrite{}); err == nil {
			t.Fatal("UpdateService() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteService(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteServiceFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteService(context.Background(), 3); err != nil {
			t.Fatalf("DeleteService() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteServiceFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteService(context.Background(), 3); err == nil {
			t.Fatal("DeleteService() expected error, got nil")
		}
	})
}

func TestNetworkService_ListServiceTemplates(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListServiceTemplatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				return &domain.PaginatedResponse[domain.ServiceTemplate]{Count: 1}, nil
			},
		})
		resp, err := svc.ListServiceTemplates(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListServiceTemplates() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListServiceTemplatesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ServiceTemplate], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListServiceTemplates(context.Background(), nil); err == nil {
			t.Fatal("ListServiceTemplates() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateServiceTemplate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, in domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return &domain.ServiceTemplate{ID: 1, Name: in.Name}, nil
			},
		})
		st, err := svc.CreateServiceTemplate(context.Background(), domain.ServiceTemplateWrite{Name: "http"})
		if err != nil {
			t.Fatalf("CreateServiceTemplate() returned error: %v", err)
		}
		if st.ID != 1 || st.Name != "http" {
			t.Errorf("ServiceTemplate = %+v, want id 1 name http", st)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateServiceTemplateFunc: func(
				_ context.Context, _ string, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateServiceTemplate(context.Background(), domain.ServiceTemplateWrite{}); err == nil {
			t.Fatal("CreateServiceTemplate() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateServiceTemplate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateServiceTemplateFunc: func(
				_ context.Context, _ string, id int, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return &domain.ServiceTemplate{ID: id}, nil
			},
		})
		st, err := svc.UpdateServiceTemplate(context.Background(), 7, domain.ServiceTemplateWrite{})
		if err != nil {
			t.Fatalf("UpdateServiceTemplate() returned error: %v", err)
		}
		if st.ID != 7 {
			t.Errorf("ServiceTemplate ID = %d, want 7", st.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateServiceTemplateFunc: func(
				_ context.Context, _ string, _ int, _ domain.ServiceTemplateWrite,
			) (*domain.ServiceTemplate, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateServiceTemplate(context.Background(), 7, domain.ServiceTemplateWrite{}); err == nil {
			t.Fatal("UpdateServiceTemplate() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteServiceTemplate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteServiceTemplateFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteServiceTemplate(context.Background(), 3); err != nil {
			t.Fatalf("DeleteServiceTemplate() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteServiceTemplateFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteServiceTemplate(context.Background(), 3); err == nil {
			t.Fatal("DeleteServiceTemplate() expected error, got nil")
		}
	})
}

func TestNetworkService_ListFHRPGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFHRPGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				return &domain.PaginatedResponse[domain.FHRPGroup]{Count: 1}, nil
			},
		})
		resp, err := svc.ListFHRPGroups(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListFHRPGroups() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFHRPGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroup], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListFHRPGroups(context.Background(), nil); err == nil {
			t.Fatal("ListFHRPGroups() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateFHRPGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFHRPGroupFunc: func(
				_ context.Context, _ string, in domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				return &domain.FHRPGroup{ID: 1, Name: in.Name}, nil
			},
		})
		g, err := svc.CreateFHRPGroup(context.Background(), domain.FHRPGroupWrite{Name: "gw"})
		if err != nil {
			t.Fatalf("CreateFHRPGroup() returned error: %v", err)
		}
		if g.ID != 1 || g.Name != "gw" {
			t.Errorf("FHRPGroup = %+v, want id 1 name gw", g)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFHRPGroupFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateFHRPGroup(context.Background(), domain.FHRPGroupWrite{}); err == nil {
			t.Fatal("CreateFHRPGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateFHRPGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFHRPGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				return &domain.FHRPGroup{ID: id}, nil
			},
		})
		g, err := svc.UpdateFHRPGroup(context.Background(), 7, domain.FHRPGroupWrite{})
		if err != nil {
			t.Fatalf("UpdateFHRPGroup() returned error: %v", err)
		}
		if g.ID != 7 {
			t.Errorf("FHRPGroup ID = %d, want 7", g.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFHRPGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.FHRPGroupWrite,
			) (*domain.FHRPGroup, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateFHRPGroup(context.Background(), 7, domain.FHRPGroupWrite{}); err == nil {
			t.Fatal("UpdateFHRPGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteFHRPGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFHRPGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteFHRPGroup(context.Background(), 3); err != nil {
			t.Fatalf("DeleteFHRPGroup() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFHRPGroupFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteFHRPGroup(context.Background(), 3); err == nil {
			t.Fatal("DeleteFHRPGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_ListFHRPGroupAssignments(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.FHRPGroupAssignment]{Count: 1}, nil
			},
		})
		resp, err := svc.ListFHRPGroupAssignments(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListFHRPGroupAssignments() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFHRPGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FHRPGroupAssignment], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListFHRPGroupAssignments(context.Background(), nil); err == nil {
			t.Fatal("ListFHRPGroupAssignments() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateFHRPGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return &domain.FHRPGroupAssignment{ID: 1}, nil
			},
		})
		assn, err := svc.CreateFHRPGroupAssignment(
			context.Background(), domain.FHRPGroupAssignmentWrite{Group: 5, InterfaceID: 6},
		)
		if err != nil {
			t.Fatalf("CreateFHRPGroupAssignment() returned error: %v", err)
		}
		if assn.ID != 1 {
			t.Errorf("FHRPGroupAssignment = %+v, want id 1", assn)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateFHRPGroupAssignment(context.Background(), domain.FHRPGroupAssignmentWrite{}); err == nil {
			t.Fatal("CreateFHRPGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateFHRPGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return &domain.FHRPGroupAssignment{ID: id}, nil
			},
		})
		a, err := svc.UpdateFHRPGroupAssignment(context.Background(), 7, domain.FHRPGroupAssignmentWrite{})
		if err != nil {
			t.Fatalf("UpdateFHRPGroupAssignment() returned error: %v", err)
		}
		if a.ID != 7 {
			t.Errorf("FHRPGroupAssignment ID = %d, want 7", a.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFHRPGroupAssignmentFunc: func(
				_ context.Context, _ string, _ int, _ domain.FHRPGroupAssignmentWrite,
			) (*domain.FHRPGroupAssignment, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateFHRPGroupAssignment(context.Background(), 7, domain.FHRPGroupAssignmentWrite{}); err == nil {
			t.Fatal("UpdateFHRPGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteFHRPGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFHRPGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteFHRPGroupAssignment(context.Background(), 3); err != nil {
			t.Fatalf("DeleteFHRPGroupAssignment() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFHRPGroupAssignmentFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteFHRPGroupAssignment(context.Background(), 3); err == nil {
			t.Fatal("DeleteFHRPGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVLANTranslationPolicies(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANTranslationPoliciesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationPolicy]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVLANTranslationPolicies(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVLANTranslationPolicies() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANTranslationPoliciesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationPolicy], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListVLANTranslationPolicies(context.Background(), nil); err == nil {
			t.Fatal("ListVLANTranslationPolicies() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVLANTranslationPolicy(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, in domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return &domain.VLANTranslationPolicy{ID: 1, Name: in.Name}, nil
			},
		})
		p, err := svc.CreateVLANTranslationPolicy(context.Background(), domain.VLANTranslationPolicyWrite{Name: "core"})
		if err != nil {
			t.Fatalf("CreateVLANTranslationPolicy() returned error: %v", err)
		}
		if p.ID != 1 || p.Name != "core" {
			t.Errorf("VLANTranslationPolicy = %+v, want id 1 name core", p)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateVLANTranslationPolicy(context.Background(), domain.VLANTranslationPolicyWrite{}); err == nil {
			t.Fatal("CreateVLANTranslationPolicy() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVLANTranslationPolicy(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, id int, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return &domain.VLANTranslationPolicy{ID: id}, nil
			},
		})
		p, err := svc.UpdateVLANTranslationPolicy(context.Background(), 7, domain.VLANTranslationPolicyWrite{})
		if err != nil {
			t.Fatalf("UpdateVLANTranslationPolicy() returned error: %v", err)
		}
		if p.ID != 7 {
			t.Errorf("VLANTranslationPolicy ID = %d, want 7", p.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANTranslationPolicyFunc: func(
				_ context.Context, _ string, _ int, _ domain.VLANTranslationPolicyWrite,
			) (*domain.VLANTranslationPolicy, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateVLANTranslationPolicy(
			context.Background(), 7, domain.VLANTranslationPolicyWrite{},
		); err == nil {
			t.Fatal("UpdateVLANTranslationPolicy() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVLANTranslationPolicy(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANTranslationPolicyFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVLANTranslationPolicy(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVLANTranslationPolicy() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANTranslationPolicyFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVLANTranslationPolicy(context.Background(), 3); err == nil {
			t.Fatal("DeleteVLANTranslationPolicy() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVLANTranslationRules(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANTranslationRulesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				return &domain.PaginatedResponse[domain.VLANTranslationRule]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVLANTranslationRules(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVLANTranslationRules() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVLANTranslationRulesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VLANTranslationRule], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListVLANTranslationRules(context.Background(), nil); err == nil {
			t.Fatal("ListVLANTranslationRules() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVLANTranslationRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return &domain.VLANTranslationRule{ID: 1}, nil
			},
		})
		r, err := svc.CreateVLANTranslationRule(context.Background(), domain.VLANTranslationRuleWrite{Policy: 5})
		if err != nil {
			t.Fatalf("CreateVLANTranslationRule() returned error: %v", err)
		}
		if r.ID != 1 {
			t.Errorf("VLANTranslationRule = %+v, want id 1", r)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateVLANTranslationRule(context.Background(), domain.VLANTranslationRuleWrite{}); err == nil {
			t.Fatal("CreateVLANTranslationRule() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVLANTranslationRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, id int, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return &domain.VLANTranslationRule{ID: id}, nil
			},
		})
		r, err := svc.UpdateVLANTranslationRule(context.Background(), 7, domain.VLANTranslationRuleWrite{})
		if err != nil {
			t.Fatalf("UpdateVLANTranslationRule() returned error: %v", err)
		}
		if r.ID != 7 {
			t.Errorf("VLANTranslationRule ID = %d, want 7", r.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVLANTranslationRuleFunc: func(
				_ context.Context, _ string, _ int, _ domain.VLANTranslationRuleWrite,
			) (*domain.VLANTranslationRule, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateVLANTranslationRule(context.Background(), 7, domain.VLANTranslationRuleWrite{}); err == nil {
			t.Fatal("UpdateVLANTranslationRule() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVLANTranslationRule(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANTranslationRuleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVLANTranslationRule(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVLANTranslationRule() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVLANTranslationRuleFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVLANTranslationRule(context.Background(), 3); err == nil {
			t.Fatal("DeleteVLANTranslationRule() expected error, got nil")
		}
	})
}

func TestNetworkService_ListTenantGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTenantGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.TenantGroup], error) {
				return &domain.PaginatedResponse[domain.TenantGroup]{Count: 1}, nil
			},
		})
		resp, err := svc.ListTenantGroups(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListTenantGroups() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListTenantGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.TenantGroup], error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.ListTenantGroups(context.Background(), nil); err == nil {
			t.Fatal("ListTenantGroups() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateTenantGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTenantGroupFunc: func(
				_ context.Context, _ string, in domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return &domain.TenantGroup{ID: 1, Name: in.Name}, nil
			},
		})
		g, err := svc.CreateTenantGroup(context.Background(), domain.TenantGroupWrite{Name: "prod"})
		if err != nil {
			t.Fatalf("CreateTenantGroup() returned error: %v", err)
		}
		if g.ID != 1 || g.Name != "prod" {
			t.Errorf("TenantGroup = %+v, want id 1 name prod", g)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateTenantGroupFunc: func(
				_ context.Context, _ string, _ domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.CreateTenantGroup(context.Background(), domain.TenantGroupWrite{}); err == nil {
			t.Fatal("CreateTenantGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateTenantGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTenantGroupFunc: func(
				_ context.Context, _ string, id int, _ domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return &domain.TenantGroup{ID: id}, nil
			},
		})
		g, err := svc.UpdateTenantGroup(context.Background(), 7, domain.TenantGroupWrite{})
		if err != nil {
			t.Fatalf("UpdateTenantGroup() returned error: %v", err)
		}
		if g.ID != 7 {
			t.Errorf("TenantGroup ID = %d, want 7", g.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateTenantGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.TenantGroupWrite,
			) (*domain.TenantGroup, error) {
				return nil, errors.New("boom")
			},
		})
		if _, err := svc.UpdateTenantGroup(context.Background(), 7, domain.TenantGroupWrite{}); err == nil {
			t.Fatal("UpdateTenantGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteTenantGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTenantGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteTenantGroup(context.Background(), 3); err != nil {
			t.Fatalf("DeleteTenantGroup() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteTenantGroupFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteTenantGroup(context.Background(), 3); err == nil {
			t.Fatal("DeleteTenantGroup() expected error, got nil")
		}
	})
}
