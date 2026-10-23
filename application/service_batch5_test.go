package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the entities added in feature/expand-entities
// Batch 5 (provider_account, provider_network, circuit_group,
// circuit_group_assignment, virtual_circuit, virtual_circuit_type,
// virtual_circuit_termination).
//
// These cover the thin relay methods in application/service.go and exercise the
// mockrepo methods they delegate to. Each entity gets success + error
// propagation cases so both branches of every wrapper are covered.

func TestNetworkService_ListProviderAccounts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListProviderAccountsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				return &domain.PaginatedResponse[domain.ProviderAccount]{Count: 1}, nil
			},
		})
		resp, err := svc.ListProviderAccounts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListProviderAccounts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListProviderAccountsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListProviderAccounts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListProviderAccounts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateProviderAccount(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, in domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: 1, Name: in.Name}, nil
			},
		})
		acc, err := svc.CreateProviderAccount(context.Background(), domain.ProviderAccountWrite{Name: "ACC-1"})
		if err != nil {
			t.Fatalf("CreateProviderAccount() returned error: %v", err)
		}
		if acc.ID != 1 || acc.Name != "ACC-1" {
			t.Errorf("account = %+v, want id 1 name ACC-1", acc)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateProviderAccount(context.Background(), domain.ProviderAccountWrite{})
		if err == nil {
			t.Fatal("CreateProviderAccount() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateProviderAccount(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderAccountFunc: func(
				_ context.Context, _ string, id int, in domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: id, Name: in.Name}, nil
			},
		})
		acc, err := svc.UpdateProviderAccount(context.Background(), 7, domain.ProviderAccountWrite{Name: "ACC-2"})
		if err != nil {
			t.Fatalf("UpdateProviderAccount() returned error: %v", err)
		}
		if acc.ID != 7 || acc.Name != "ACC-2" {
			t.Errorf("account = %+v, want id 7 name ACC-2", acc)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderAccountFunc: func(
				_ context.Context, _ string, _ int, _ domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateProviderAccount(context.Background(), 7, domain.ProviderAccountWrite{})
		if err == nil {
			t.Fatal("UpdateProviderAccount() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteProviderAccount(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderAccountFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteProviderAccount(context.Background(), 3); err != nil {
			t.Fatalf("DeleteProviderAccount() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderAccountFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteProviderAccount(context.Background(), 3); err == nil {
			t.Fatal("DeleteProviderAccount() expected error, got nil")
		}
	})
}

func TestNetworkService_ListProviderNetworks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListProviderNetworksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
				return &domain.PaginatedResponse[domain.ProviderNetwork]{Count: 1}, nil
			},
		})
		resp, err := svc.ListProviderNetworks(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListProviderNetworks() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListProviderNetworksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListProviderNetworks(context.Background(), nil)
		if err == nil {
			t.Fatal("ListProviderNetworks() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateProviderNetwork(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderNetworkFunc: func(
				_ context.Context, _ string, in domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return &domain.ProviderNetwork{ID: 1, Name: in.Name}, nil
			},
		})
		net, err := svc.CreateProviderNetwork(context.Background(), domain.ProviderNetworkWrite{Name: "Net A"})
		if err != nil {
			t.Fatalf("CreateProviderNetwork() returned error: %v", err)
		}
		if net.ID != 1 || net.Name != "Net A" {
			t.Errorf("network = %+v, want id 1 name Net A", net)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateProviderNetworkFunc: func(
				_ context.Context, _ string, _ domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateProviderNetwork(context.Background(), domain.ProviderNetworkWrite{})
		if err == nil {
			t.Fatal("CreateProviderNetwork() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateProviderNetwork(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderNetworkFunc: func(
				_ context.Context, _ string, id int, in domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return &domain.ProviderNetwork{ID: id, Name: in.Name}, nil
			},
		})
		net, err := svc.UpdateProviderNetwork(context.Background(), 7, domain.ProviderNetworkWrite{Name: "Net B"})
		if err != nil {
			t.Fatalf("UpdateProviderNetwork() returned error: %v", err)
		}
		if net.ID != 7 || net.Name != "Net B" {
			t.Errorf("network = %+v, want id 7 name Net B", net)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateProviderNetworkFunc: func(
				_ context.Context, _ string, _ int, _ domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateProviderNetwork(context.Background(), 7, domain.ProviderNetworkWrite{})
		if err == nil {
			t.Fatal("UpdateProviderNetwork() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteProviderNetwork(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderNetworkFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteProviderNetwork(context.Background(), 3); err != nil {
			t.Fatalf("DeleteProviderNetwork() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteProviderNetworkFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteProviderNetwork(context.Background(), 3); err == nil {
			t.Fatal("DeleteProviderNetwork() expected error, got nil")
		}
	})
}

func TestNetworkService_ListCircuitGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				return &domain.PaginatedResponse[domain.CircuitGroup]{Count: 1}, nil
			},
		})
		resp, err := svc.ListCircuitGroups(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCircuitGroups() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListCircuitGroups(context.Background(), nil)
		if err == nil {
			t.Fatal("ListCircuitGroups() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateCircuitGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: 1, Name: in.Name}, nil
			},
		})
		cg, err := svc.CreateCircuitGroup(context.Background(), domain.CircuitGroupWrite{Name: "Grp"})
		if err != nil {
			t.Fatalf("CreateCircuitGroup() returned error: %v", err)
		}
		if cg.ID != 1 || cg.Name != "Grp" {
			t.Errorf("group = %+v, want id 1 name Grp", cg)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCircuitGroup(context.Background(), domain.CircuitGroupWrite{})
		if err == nil {
			t.Fatal("CreateCircuitGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateCircuitGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitGroupFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: id, Name: in.Name}, nil
			},
		})
		cg, err := svc.UpdateCircuitGroup(context.Background(), 7, domain.CircuitGroupWrite{Name: "Grp2"})
		if err != nil {
			t.Fatalf("UpdateCircuitGroup() returned error: %v", err)
		}
		if cg.ID != 7 || cg.Name != "Grp2" {
			t.Errorf("group = %+v, want id 7 name Grp2", cg)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCircuitGroup(context.Background(), 7, domain.CircuitGroupWrite{})
		if err == nil {
			t.Fatal("UpdateCircuitGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteCircuitGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteCircuitGroup(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCircuitGroup() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitGroupFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteCircuitGroup(context.Background(), 3); err == nil {
			t.Fatal("DeleteCircuitGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_ListCircuitGroupAssignments(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.CircuitGroupAssignment]{Count: 1}, nil
			},
		})
		resp, err := svc.ListCircuitGroupAssignments(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListCircuitGroupAssignments() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListCircuitGroupAssignments(context.Background(), nil)
		if err == nil {
			t.Fatal("ListCircuitGroupAssignments() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateCircuitGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return &domain.CircuitGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
			},
		})
		asg, err := svc.CreateCircuitGroupAssignment(
			context.Background(), domain.CircuitGroupAssignmentWrite{Group: 5, MemberType: "circuits.circuit", MemberID: 6},
		)
		if err != nil {
			t.Fatalf("CreateCircuitGroupAssignment() returned error: %v", err)
		}
		if asg.ID != 1 || asg.Group == nil || asg.Group.ID != 5 {
			t.Errorf("assignment = %+v, want id 1 group 5", asg)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, _ domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateCircuitGroupAssignment(context.Background(), domain.CircuitGroupAssignmentWrite{})
		if err == nil {
			t.Fatal("CreateCircuitGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateCircuitGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return &domain.CircuitGroupAssignment{ID: id, Group: &domain.Nested{ID: in.Group}}, nil
			},
		})
		asg, err := svc.UpdateCircuitGroupAssignment(context.Background(), 7, domain.CircuitGroupAssignmentWrite{Group: 9})
		if err != nil {
			t.Fatalf("UpdateCircuitGroupAssignment() returned error: %v", err)
		}
		if asg.ID != 7 || asg.Group == nil || asg.Group.ID != 9 {
			t.Errorf("assignment = %+v, want id 7 group 9", asg)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, _ int, _ domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateCircuitGroupAssignment(context.Background(), 7, domain.CircuitGroupAssignmentWrite{})
		if err == nil {
			t.Fatal("UpdateCircuitGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteCircuitGroupAssignment(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteCircuitGroupAssignment(context.Background(), 3); err != nil {
			t.Fatalf("DeleteCircuitGroupAssignment() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteCircuitGroupAssignmentFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteCircuitGroupAssignment(context.Background(), 3); err == nil {
			t.Fatal("DeleteCircuitGroupAssignment() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualCircuits(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuit]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVirtualCircuits(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuits() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListVirtualCircuits(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVirtualCircuits() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVirtualCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: 1, CID: in.CID}, nil
			},
		})
		vc, err := svc.CreateVirtualCircuit(context.Background(), domain.VirtualCircuitWrite{CID: "VC-1"})
		if err != nil {
			t.Fatalf("CreateVirtualCircuit() returned error: %v", err)
		}
		if vc.ID != 1 || vc.CID != "VC-1" {
			t.Errorf("vc = %+v, want id 1 cid VC-1", vc)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualCircuit(context.Background(), domain.VirtualCircuitWrite{})
		if err == nil {
			t.Fatal("CreateVirtualCircuit() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVirtualCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: id, CID: in.CID}, nil
			},
		})
		vc, err := svc.UpdateVirtualCircuit(context.Background(), 7, domain.VirtualCircuitWrite{CID: "VC-2"})
		if err != nil {
			t.Fatalf("UpdateVirtualCircuit() returned error: %v", err)
		}
		if vc.ID != 7 || vc.CID != "VC-2" {
			t.Errorf("vc = %+v, want id 7 cid VC-2", vc)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualCircuit(context.Background(), 7, domain.VirtualCircuitWrite{})
		if err == nil {
			t.Fatal("UpdateVirtualCircuit() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVirtualCircuit(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVirtualCircuit(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualCircuit() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVirtualCircuit(context.Background(), 3); err == nil {
			t.Fatal("DeleteVirtualCircuit() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualCircuitTypes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitType]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVirtualCircuitTypes(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuitTypes() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListVirtualCircuitTypes(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVirtualCircuitTypes() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVirtualCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: 1, Name: in.Name}, nil
			},
		})
		vct, err := svc.CreateVirtualCircuitType(context.Background(), domain.VirtualCircuitTypeWrite{Name: "VPLS"})
		if err != nil {
			t.Fatalf("CreateVirtualCircuitType() returned error: %v", err)
		}
		if vct.ID != 1 || vct.Name != "VPLS" {
			t.Errorf("vct = %+v, want id 1 name VPLS", vct)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualCircuitType(context.Background(), domain.VirtualCircuitTypeWrite{})
		if err == nil {
			t.Fatal("CreateVirtualCircuitType() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVirtualCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: id, Name: in.Name}, nil
			},
		})
		vct, err := svc.UpdateVirtualCircuitType(context.Background(), 7, domain.VirtualCircuitTypeWrite{Name: "EVPL"})
		if err != nil {
			t.Fatalf("UpdateVirtualCircuitType() returned error: %v", err)
		}
		if vct.ID != 7 || vct.Name != "EVPL" {
			t.Errorf("vct = %+v, want id 7 name EVPL", vct)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualCircuitType(context.Background(), 7, domain.VirtualCircuitTypeWrite{})
		if err == nil {
			t.Fatal("UpdateVirtualCircuitType() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVirtualCircuitType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVirtualCircuitType(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualCircuitType() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTypeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVirtualCircuitType(context.Background(), 3); err == nil {
			t.Fatal("DeleteVirtualCircuitType() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualCircuitTerminations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitTermination]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVirtualCircuitTerminations(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualCircuitTerminations() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListVirtualCircuitTerminations(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVirtualCircuitTerminations() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVirtualCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return &domain.VirtualCircuitTermination{
					ID: 1, VirtualCircuit: &domain.Nested{ID: in.VirtualCircuit},
				}, nil
			},
		})
		term, err := svc.CreateVirtualCircuitTermination(
			context.Background(), domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
		)
		if err != nil {
			t.Fatalf("CreateVirtualCircuitTermination() returned error: %v", err)
		}
		if term.ID != 1 || term.VirtualCircuit == nil || term.VirtualCircuit.ID != 5 {
			t.Errorf("term = %+v, want id 1 vc 5", term)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualCircuitTermination(context.Background(), domain.VirtualCircuitTerminationWrite{})
		if err == nil {
			t.Fatal("CreateVirtualCircuitTermination() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVirtualCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return &domain.VirtualCircuitTermination{ID: id}, nil
			},
		})
		term, err := svc.UpdateVirtualCircuitTermination(
			context.Background(), 7, domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
		)
		if err != nil {
			t.Fatalf("UpdateVirtualCircuitTermination() returned error: %v", err)
		}
		if term.ID != 7 {
			t.Errorf("term = %+v, want id 7", term)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualCircuitTermination(context.Background(), 7, domain.VirtualCircuitTerminationWrite{})
		if err == nil {
			t.Fatal("UpdateVirtualCircuitTermination() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVirtualCircuitTermination(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVirtualCircuitTermination(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualCircuitTermination() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualCircuitTerminationFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVirtualCircuitTermination(context.Background(), 3); err == nil {
			t.Fatal("DeleteVirtualCircuitTermination() expected error, got nil")
		}
	})
}
