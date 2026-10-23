package mockrepo

import (
	"context"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// Coverage for the Batch 5 entities (provider_account, provider_network,
// circuit_group, circuit_group_assignment, virtual_circuit,
// virtual_circuit_type, virtual_circuit_termination): MockRepo delegates
// forward to their Func fields and NilRepo returns empty/sentinel values. Both
// implementations are exercised so the batch's repository methods count toward
// the coverage gate.

func TestMockRepo_Batch5Delegates(t *testing.T) { //nolint:gocognit,gocyclo,funlen,maintidx
	t.Parallel()

	t.Run("ListProviderAccounts", func(t *testing.T) {
		got, err := (&MockRepo{
			ListProviderAccountsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderAccount], error) {
				return &domain.PaginatedResponse[domain.ProviderAccount]{Count: 1}, nil
			},
		}).ListProviderAccounts(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListProviderAccounts = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateProviderAccount", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateProviderAccountFunc: func(
				_ context.Context, _ string, in domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: 1, Name: in.Name}, nil
			},
		}).CreateProviderAccount(context.Background(), "t", domain.ProviderAccountWrite{Name: "ACC-1"})
		if err != nil || got == nil || got.Name != "ACC-1" {
			t.Errorf("CreateProviderAccount = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateProviderAccount", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateProviderAccountFunc: func(
				_ context.Context, _ string, id int, in domain.ProviderAccountWrite,
			) (*domain.ProviderAccount, error) {
				return &domain.ProviderAccount{ID: id, Name: in.Name}, nil
			},
		}).UpdateProviderAccount(context.Background(), "t", 7, domain.ProviderAccountWrite{Name: "ACC-2"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateProviderAccount = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteProviderAccount", func(t *testing.T) {
		err := (&MockRepo{
			DeleteProviderAccountFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteProviderAccount(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteProviderAccount = %v", err)
		}
	})

	t.Run("ListProviderNetworks", func(t *testing.T) {
		got, err := (&MockRepo{
			ListProviderNetworksFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ProviderNetwork], error) {
				return &domain.PaginatedResponse[domain.ProviderNetwork]{Count: 1}, nil
			},
		}).ListProviderNetworks(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListProviderNetworks = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateProviderNetwork", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateProviderNetworkFunc: func(
				_ context.Context, _ string, in domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return &domain.ProviderNetwork{ID: 1, Name: in.Name}, nil
			},
		}).CreateProviderNetwork(context.Background(), "t", domain.ProviderNetworkWrite{Name: "Net A"})
		if err != nil || got == nil || got.Name != "Net A" {
			t.Errorf("CreateProviderNetwork = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateProviderNetwork", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateProviderNetworkFunc: func(
				_ context.Context, _ string, id int, in domain.ProviderNetworkWrite,
			) (*domain.ProviderNetwork, error) {
				return &domain.ProviderNetwork{ID: id, Name: in.Name}, nil
			},
		}).UpdateProviderNetwork(context.Background(), "t", 7, domain.ProviderNetworkWrite{Name: "Net B"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateProviderNetwork = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteProviderNetwork", func(t *testing.T) {
		err := (&MockRepo{
			DeleteProviderNetworkFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteProviderNetwork(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteProviderNetwork = %v", err)
		}
	})

	t.Run("ListCircuitGroups", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCircuitGroupsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroup], error) {
				return &domain.PaginatedResponse[domain.CircuitGroup]{Count: 1}, nil
			},
		}).ListCircuitGroups(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListCircuitGroups = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateCircuitGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCircuitGroupFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateCircuitGroup(context.Background(), "t", domain.CircuitGroupWrite{Name: "Grp"})
		if err != nil || got == nil || got.Name != "Grp" {
			t.Errorf("CreateCircuitGroup = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateCircuitGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCircuitGroupFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupWrite,
			) (*domain.CircuitGroup, error) {
				return &domain.CircuitGroup{ID: id, Name: in.Name}, nil
			},
		}).UpdateCircuitGroup(context.Background(), "t", 7, domain.CircuitGroupWrite{Name: "Grp2"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuitGroup = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteCircuitGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCircuitGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCircuitGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuitGroup = %v", err)
		}
	})

	t.Run("ListCircuitGroupAssignments", func(t *testing.T) {
		got, err := (&MockRepo{
			ListCircuitGroupAssignmentsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.CircuitGroupAssignment], error) {
				return &domain.PaginatedResponse[domain.CircuitGroupAssignment]{Count: 1}, nil
			},
		}).ListCircuitGroupAssignments(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListCircuitGroupAssignments = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateCircuitGroupAssignment", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return &domain.CircuitGroupAssignment{ID: 1, Group: &domain.Nested{ID: in.Group}}, nil
			},
		}).CreateCircuitGroupAssignment(context.Background(), "t", domain.CircuitGroupAssignmentWrite{Group: 5})
		if err != nil || got == nil || got.Group == nil || got.Group.ID != 5 {
			t.Errorf("CreateCircuitGroupAssignment = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateCircuitGroupAssignment", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateCircuitGroupAssignmentFunc: func(
				_ context.Context, _ string, id int, in domain.CircuitGroupAssignmentWrite,
			) (*domain.CircuitGroupAssignment, error) {
				return &domain.CircuitGroupAssignment{ID: id, Group: &domain.Nested{ID: in.Group}}, nil
			},
		}).UpdateCircuitGroupAssignment(context.Background(), "t", 7, domain.CircuitGroupAssignmentWrite{Group: 9})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateCircuitGroupAssignment = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteCircuitGroupAssignment", func(t *testing.T) {
		err := (&MockRepo{
			DeleteCircuitGroupAssignmentFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteCircuitGroupAssignment(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteCircuitGroupAssignment = %v", err)
		}
	})

	t.Run("ListVirtualCircuits", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualCircuitsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuit], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuit]{Count: 1}, nil
			},
		}).ListVirtualCircuits(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualCircuits = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateVirtualCircuit", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualCircuitFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: 1, CID: in.CID}, nil
			},
		}).CreateVirtualCircuit(context.Background(), "t", domain.VirtualCircuitWrite{CID: "VC-1"})
		if err != nil || got == nil || got.CID != "VC-1" {
			t.Errorf("CreateVirtualCircuit = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateVirtualCircuit", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualCircuitFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitWrite,
			) (*domain.VirtualCircuit, error) {
				return &domain.VirtualCircuit{ID: id, CID: in.CID}, nil
			},
		}).UpdateVirtualCircuit(context.Background(), "t", 7, domain.VirtualCircuitWrite{CID: "VC-2"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualCircuit = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteVirtualCircuit", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualCircuitFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualCircuit(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualCircuit = %v", err)
		}
	})

	t.Run("ListVirtualCircuitTypes", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualCircuitTypesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitType], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitType]{Count: 1}, nil
			},
		}).ListVirtualCircuitTypes(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualCircuitTypes = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateVirtualCircuitType", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: 1, Name: in.Name}, nil
			},
		}).CreateVirtualCircuitType(context.Background(), "t", domain.VirtualCircuitTypeWrite{Name: "VPLS"})
		if err != nil || got == nil || got.Name != "VPLS" {
			t.Errorf("CreateVirtualCircuitType = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateVirtualCircuitType", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualCircuitTypeFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualCircuitTypeWrite,
			) (*domain.VirtualCircuitType, error) {
				return &domain.VirtualCircuitType{ID: id, Name: in.Name}, nil
			},
		}).UpdateVirtualCircuitType(context.Background(), "t", 7, domain.VirtualCircuitTypeWrite{Name: "EVPL"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualCircuitType = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteVirtualCircuitType", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualCircuitTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualCircuitType(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualCircuitType = %v", err)
		}
	})

	t.Run("ListVirtualCircuitTerminations", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualCircuitTerminationsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualCircuitTermination], error) {
				return &domain.PaginatedResponse[domain.VirtualCircuitTermination]{Count: 1}, nil
			},
		}).ListVirtualCircuitTerminations(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualCircuitTerminations = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateVirtualCircuitTermination", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, in domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return &domain.VirtualCircuitTermination{ID: 1, VirtualCircuit: &domain.Nested{ID: in.VirtualCircuit}}, nil
			},
		}).CreateVirtualCircuitTermination(
			context.Background(), "t", domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
		)
		if err != nil || got == nil || got.VirtualCircuit == nil || got.VirtualCircuit.ID != 5 {
			t.Errorf("CreateVirtualCircuitTermination = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateVirtualCircuitTermination", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualCircuitTerminationFunc: func(
				_ context.Context, _ string, id int, _ domain.VirtualCircuitTerminationWrite,
			) (*domain.VirtualCircuitTermination, error) {
				return &domain.VirtualCircuitTermination{ID: id}, nil
			},
		}).UpdateVirtualCircuitTermination(context.Background(), "t", 7, domain.VirtualCircuitTerminationWrite{})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualCircuitTermination = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteVirtualCircuitTermination", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualCircuitTerminationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualCircuitTermination(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualCircuitTermination = %v", err)
		}
	})
}

//nolint:funlen,gocognit,gocyclo // Sequential assertions across all 7 entities.
func TestNilRepo_Batch5(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &NilRepo{}

	if r, err := repo.ListProviderAccounts(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListProviderAccounts = (%v, %v), want count 0", r, err)
	}
	if a, err := repo.CreateProviderAccount(
		ctx, "t", domain.ProviderAccountWrite{Name: "ACC-1"},
	); err != nil || a.ID != 1 {
		t.Errorf("CreateProviderAccount = (%v, %v)", a, err)
	}
	if a, err := repo.UpdateProviderAccount(
		ctx, "t", 7, domain.ProviderAccountWrite{Name: "ACC-2"},
	); err != nil || a.ID != 7 {
		t.Errorf("UpdateProviderAccount = (%v, %v)", a, err)
	}
	if err := repo.DeleteProviderAccount(ctx, "t", 3); err != nil {
		t.Errorf("DeleteProviderAccount = %v", err)
	}

	if r, err := repo.ListProviderNetworks(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListProviderNetworks = (%v, %v), want count 0", r, err)
	}
	if n, err := repo.CreateProviderNetwork(
		ctx, "t", domain.ProviderNetworkWrite{Name: "Net A"},
	); err != nil || n.ID != 1 {
		t.Errorf("CreateProviderNetwork = (%v, %v)", n, err)
	}
	if n, err := repo.UpdateProviderNetwork(
		ctx, "t", 7, domain.ProviderNetworkWrite{Name: "Net B"},
	); err != nil || n.ID != 7 {
		t.Errorf("UpdateProviderNetwork = (%v, %v)", n, err)
	}
	if err := repo.DeleteProviderNetwork(ctx, "t", 3); err != nil {
		t.Errorf("DeleteProviderNetwork = %v", err)
	}

	if r, err := repo.ListCircuitGroups(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListCircuitGroups = (%v, %v), want count 0", r, err)
	}
	if g, err := repo.CreateCircuitGroup(ctx, "t", domain.CircuitGroupWrite{Name: "Grp"}); err != nil || g.ID != 1 {
		t.Errorf("CreateCircuitGroup = (%v, %v)", g, err)
	}
	if g, err := repo.UpdateCircuitGroup(ctx, "t", 7, domain.CircuitGroupWrite{Name: "Grp2"}); err != nil || g.ID != 7 {
		t.Errorf("UpdateCircuitGroup = (%v, %v)", g, err)
	}
	if err := repo.DeleteCircuitGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCircuitGroup = %v", err)
	}

	if r, err := repo.ListCircuitGroupAssignments(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListCircuitGroupAssignments = (%v, %v), want count 0", r, err)
	}
	if a, err := repo.CreateCircuitGroupAssignment(
		ctx, "t", domain.CircuitGroupAssignmentWrite{Group: 5},
	); err != nil || a.ID != 1 {
		t.Errorf("CreateCircuitGroupAssignment = (%v, %v)", a, err)
	}
	if a, err := repo.UpdateCircuitGroupAssignment(
		ctx, "t", 7, domain.CircuitGroupAssignmentWrite{},
	); err != nil || a.ID != 7 {
		t.Errorf("UpdateCircuitGroupAssignment = (%v, %v)", a, err)
	}
	if err := repo.DeleteCircuitGroupAssignment(ctx, "t", 3); err != nil {
		t.Errorf("DeleteCircuitGroupAssignment = %v", err)
	}

	if r, err := repo.ListVirtualCircuits(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListVirtualCircuits = (%v, %v), want count 0", r, err)
	}
	if v, err := repo.CreateVirtualCircuit(
		ctx, "t", domain.VirtualCircuitWrite{CID: "VC-1"},
	); err != nil || v.ID != 1 {
		t.Errorf("CreateVirtualCircuit = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualCircuit(
		ctx, "t", 7, domain.VirtualCircuitWrite{CID: "VC-2"},
	); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualCircuit = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualCircuit(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualCircuit = %v", err)
	}

	if r, err := repo.ListVirtualCircuitTypes(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListVirtualCircuitTypes = (%v, %v), want count 0", r, err)
	}
	if v, err := repo.CreateVirtualCircuitType(
		ctx, "t", domain.VirtualCircuitTypeWrite{Name: "VPLS"},
	); err != nil || v.ID != 1 {
		t.Errorf("CreateVirtualCircuitType = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualCircuitType(
		ctx, "t", 7, domain.VirtualCircuitTypeWrite{Name: "EVPL"},
	); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualCircuitType = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualCircuitType(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualCircuitType = %v", err)
	}

	if r, err := repo.ListVirtualCircuitTerminations(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListVirtualCircuitTerminations = (%v, %v), want count 0", r, err)
	}
	if v, err := repo.CreateVirtualCircuitTermination(
		ctx, "t", domain.VirtualCircuitTerminationWrite{VirtualCircuit: 5},
	); err != nil || v.ID != 1 {
		t.Errorf("CreateVirtualCircuitTermination = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualCircuitTermination(
		ctx, "t", 7, domain.VirtualCircuitTerminationWrite{},
	); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualCircuitTermination = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualCircuitTermination(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualCircuitTermination = %v", err)
	}
}
