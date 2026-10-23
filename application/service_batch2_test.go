package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the DCIM entities added in feature/expand-entities
// Batch 2. These cover the thin relay methods in application/service.go (and
// exercise the mockrepo methods they delegate to). Each entity gets success +
// error-propagation cases so both branches of the wrapper are covered.

func TestNetworkService_ListFrontPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFrontPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				return &domain.PaginatedResponse[domain.FrontPort]{Count: 1}, nil
			},
		})
		resp, err := svc.ListFrontPorts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListFrontPorts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListFrontPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.FrontPort], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListFrontPorts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListFrontPorts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateFrontPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFrontPortFunc: func(
				_ context.Context, _ string, in domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: 1, Name: in.Name}, nil
			},
		})
		fp, err := svc.CreateFrontPort(context.Background(), domain.FrontPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateFrontPort() returned error: %v", err)
		}
		if fp.ID != 1 || fp.Name != "A" {
			t.Errorf("front port = %+v, want id 1 name A", fp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateFrontPortFunc: func(
				_ context.Context, _ string, _ domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateFrontPort(context.Background(), domain.FrontPortWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateFrontPort() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateFrontPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFrontPortFunc: func(
				_ context.Context, _ string, id int, in domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return &domain.FrontPort{ID: id, Name: in.Name}, nil
			},
		})
		fp, err := svc.UpdateFrontPort(context.Background(), 7, domain.FrontPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateFrontPort() returned error: %v", err)
		}
		if fp.ID != 7 {
			t.Errorf("front port ID = %d, want 7", fp.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateFrontPortFunc: func(
				_ context.Context, _ string, _ int, _ domain.FrontPortWrite,
			) (*domain.FrontPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateFrontPort(context.Background(), 7, domain.FrontPortWrite{})
		if err == nil {
			t.Fatal("UpdateFrontPort() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteFrontPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFrontPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteFrontPort(context.Background(), 3); err != nil {
			t.Fatalf("DeleteFrontPort() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteFrontPortFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteFrontPort(context.Background(), 3); err == nil {
			t.Fatal("DeleteFrontPort() expected error, got nil")
		}
	})
}

func TestNetworkService_ListRearPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRearPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RearPort], error) {
				return &domain.PaginatedResponse[domain.RearPort]{Count: 1}, nil
			},
		})
		resp, err := svc.ListRearPorts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRearPorts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRearPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RearPort], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListRearPorts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListRearPorts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateRearPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRearPortFunc: func(
				_ context.Context, _ string, in domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return &domain.RearPort{ID: 1, Name: in.Name}, nil
			},
		})
		rp, err := svc.CreateRearPort(context.Background(), domain.RearPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateRearPort() returned error: %v", err)
		}
		if rp.ID != 1 || rp.Name != "A" {
			t.Errorf("rear port = %+v, want id 1 name A", rp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRearPortFunc: func(
				_ context.Context, _ string, _ domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateRearPort(context.Background(), domain.RearPortWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateRearPort() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateRearPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRearPortFunc: func(
				_ context.Context, _ string, id int, in domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return &domain.RearPort{ID: id, Name: in.Name}, nil
			},
		})
		rp, err := svc.UpdateRearPort(context.Background(), 7, domain.RearPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateRearPort() returned error: %v", err)
		}
		if rp.ID != 7 {
			t.Errorf("rear port ID = %d, want 7", rp.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRearPortFunc: func(
				_ context.Context, _ string, _ int, _ domain.RearPortWrite,
			) (*domain.RearPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateRearPort(context.Background(), 7, domain.RearPortWrite{})
		if err == nil {
			t.Fatal("UpdateRearPort() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteRearPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRearPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteRearPort(context.Background(), 3); err != nil {
			t.Fatalf("DeleteRearPort() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRearPortFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteRearPort(context.Background(), 3); err == nil {
			t.Fatal("DeleteRearPort() expected error, got nil")
		}
	})
}

func TestNetworkService_ListConsolePorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListConsolePortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				return &domain.PaginatedResponse[domain.ConsolePort]{Count: 1}, nil
			},
		})
		resp, err := svc.ListConsolePorts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListConsolePorts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListConsolePortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsolePort], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListConsolePorts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListConsolePorts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateConsolePort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateConsolePortFunc: func(
				_ context.Context, _ string, in domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: 1, Name: in.Name}, nil
			},
		})
		cp, err := svc.CreateConsolePort(context.Background(), domain.ConsolePortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateConsolePort() returned error: %v", err)
		}
		if cp.ID != 1 || cp.Name != "A" {
			t.Errorf("console port = %+v, want id 1 name A", cp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateConsolePortFunc: func(
				_ context.Context, _ string, _ domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateConsolePort(context.Background(), domain.ConsolePortWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateConsolePort() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateConsolePort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateConsolePortFunc: func(
				_ context.Context, _ string, id int, in domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return &domain.ConsolePort{ID: id, Name: in.Name}, nil
			},
		})
		cp, err := svc.UpdateConsolePort(context.Background(), 7, domain.ConsolePortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateConsolePort() returned error: %v", err)
		}
		if cp.ID != 7 {
			t.Errorf("console port ID = %d, want 7", cp.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateConsolePortFunc: func(
				_ context.Context, _ string, _ int, _ domain.ConsolePortWrite,
			) (*domain.ConsolePort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateConsolePort(context.Background(), 7, domain.ConsolePortWrite{})
		if err == nil {
			t.Fatal("UpdateConsolePort() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteConsolePort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteConsolePortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteConsolePort(context.Background(), 3); err != nil {
			t.Fatalf("DeleteConsolePort() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteConsolePortFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteConsolePort(context.Background(), 3); err == nil {
			t.Fatal("DeleteConsolePort() expected error, got nil")
		}
	})
}

func TestNetworkService_ListConsoleServerPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListConsoleServerPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
				return &domain.PaginatedResponse[domain.ConsoleServerPort]{Count: 1}, nil
			},
		})
		resp, err := svc.ListConsoleServerPorts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListConsoleServerPorts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListConsoleServerPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ConsoleServerPort], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListConsoleServerPorts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListConsoleServerPorts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateConsoleServerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateConsoleServerPortFunc: func(
				_ context.Context, _ string, in domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return &domain.ConsoleServerPort{ID: 1, Name: in.Name}, nil
			},
		})
		csp, err := svc.CreateConsoleServerPort(context.Background(), domain.ConsoleServerPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateConsoleServerPort() returned error: %v", err)
		}
		if csp.ID != 1 || csp.Name != "A" {
			t.Errorf("console server port = %+v, want id 1 name A", csp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateConsoleServerPortFunc: func(
				_ context.Context, _ string, _ domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateConsoleServerPort(context.Background(), domain.ConsoleServerPortWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateConsoleServerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateConsoleServerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateConsoleServerPortFunc: func(
				_ context.Context, _ string, id int, in domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return &domain.ConsoleServerPort{ID: id, Name: in.Name}, nil
			},
		})
		csp, err := svc.UpdateConsoleServerPort(context.Background(), 7, domain.ConsoleServerPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateConsoleServerPort() returned error: %v", err)
		}
		if csp.ID != 7 {
			t.Errorf("console server port ID = %d, want 7", csp.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateConsoleServerPortFunc: func(
				_ context.Context, _ string, _ int, _ domain.ConsoleServerPortWrite,
			) (*domain.ConsoleServerPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateConsoleServerPort(context.Background(), 7, domain.ConsoleServerPortWrite{})
		if err == nil {
			t.Fatal("UpdateConsoleServerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteConsoleServerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteConsoleServerPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteConsoleServerPort(context.Background(), 3); err != nil {
			t.Fatalf("DeleteConsoleServerPort() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteConsoleServerPortFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteConsoleServerPort(context.Background(), 3); err == nil {
			t.Fatal("DeleteConsoleServerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_ListPowerPorts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPort], error) {
				return &domain.PaginatedResponse[domain.PowerPort]{Count: 1}, nil
			},
		})
		resp, err := svc.ListPowerPorts(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListPowerPorts() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerPortsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerPort], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListPowerPorts(context.Background(), nil)
		if err == nil {
			t.Fatal("ListPowerPorts() expected error, got nil")
		}
	})
}

func TestNetworkService_CreatePowerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerPortFunc: func(
				_ context.Context, _ string, in domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return &domain.PowerPort{ID: 1, Name: in.Name}, nil
			},
		})
		pp, err := svc.CreatePowerPort(context.Background(), domain.PowerPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreatePowerPort() returned error: %v", err)
		}
		if pp.ID != 1 || pp.Name != "A" {
			t.Errorf("power port = %+v, want id 1 name A", pp)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerPortFunc: func(
				_ context.Context, _ string, _ domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreatePowerPort(context.Background(), domain.PowerPortWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreatePowerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdatePowerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerPortFunc: func(
				_ context.Context, _ string, id int, in domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return &domain.PowerPort{ID: id, Name: in.Name}, nil
			},
		})
		pp, err := svc.UpdatePowerPort(context.Background(), 7, domain.PowerPortWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdatePowerPort() returned error: %v", err)
		}
		if pp.ID != 7 {
			t.Errorf("power port ID = %d, want 7", pp.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerPortFunc: func(
				_ context.Context, _ string, _ int, _ domain.PowerPortWrite,
			) (*domain.PowerPort, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdatePowerPort(context.Background(), 7, domain.PowerPortWrite{})
		if err == nil {
			t.Fatal("UpdatePowerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_DeletePowerPort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerPortFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeletePowerPort(context.Background(), 3); err != nil {
			t.Fatalf("DeletePowerPort() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerPortFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeletePowerPort(context.Background(), 3); err == nil {
			t.Fatal("DeletePowerPort() expected error, got nil")
		}
	})
}

func TestNetworkService_ListPowerOutlets(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerOutletsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
				return &domain.PaginatedResponse[domain.PowerOutlet]{Count: 1}, nil
			},
		})
		resp, err := svc.ListPowerOutlets(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListPowerOutlets() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerOutletsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerOutlet], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListPowerOutlets(context.Background(), nil)
		if err == nil {
			t.Fatal("ListPowerOutlets() expected error, got nil")
		}
	})
}

func TestNetworkService_CreatePowerOutlet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerOutletFunc: func(
				_ context.Context, _ string, in domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return &domain.PowerOutlet{ID: 1, Name: in.Name}, nil
			},
		})
		po, err := svc.CreatePowerOutlet(context.Background(), domain.PowerOutletWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreatePowerOutlet() returned error: %v", err)
		}
		if po.ID != 1 || po.Name != "A" {
			t.Errorf("power outlet = %+v, want id 1 name A", po)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerOutletFunc: func(
				_ context.Context, _ string, _ domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreatePowerOutlet(context.Background(), domain.PowerOutletWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreatePowerOutlet() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdatePowerOutlet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerOutletFunc: func(
				_ context.Context, _ string, id int, in domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return &domain.PowerOutlet{ID: id, Name: in.Name}, nil
			},
		})
		po, err := svc.UpdatePowerOutlet(context.Background(), 7, domain.PowerOutletWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdatePowerOutlet() returned error: %v", err)
		}
		if po.ID != 7 {
			t.Errorf("power outlet ID = %d, want 7", po.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerOutletFunc: func(
				_ context.Context, _ string, _ int, _ domain.PowerOutletWrite,
			) (*domain.PowerOutlet, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdatePowerOutlet(context.Background(), 7, domain.PowerOutletWrite{})
		if err == nil {
			t.Fatal("UpdatePowerOutlet() expected error, got nil")
		}
	})
}

func TestNetworkService_DeletePowerOutlet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerOutletFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeletePowerOutlet(context.Background(), 3); err != nil {
			t.Fatalf("DeletePowerOutlet() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerOutletFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeletePowerOutlet(context.Background(), 3); err == nil {
			t.Fatal("DeletePowerOutlet() expected error, got nil")
		}
	})
}

func TestNetworkService_ListPowerFeeds(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerFeedsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				return &domain.PaginatedResponse[domain.PowerFeed]{Count: 1}, nil
			},
		})
		resp, err := svc.ListPowerFeeds(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListPowerFeeds() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListPowerFeedsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.PowerFeed], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListPowerFeeds(context.Background(), nil)
		if err == nil {
			t.Fatal("ListPowerFeeds() expected error, got nil")
		}
	})
}

func TestNetworkService_CreatePowerFeed(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerFeedFunc: func(
				_ context.Context, _ string, in domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: 1, Name: in.Name}, nil
			},
		})
		pf, err := svc.CreatePowerFeed(context.Background(), domain.PowerFeedWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreatePowerFeed() returned error: %v", err)
		}
		if pf.ID != 1 || pf.Name != "A" {
			t.Errorf("power feed = %+v, want id 1 name A", pf)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreatePowerFeedFunc: func(
				_ context.Context, _ string, _ domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreatePowerFeed(context.Background(), domain.PowerFeedWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreatePowerFeed() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdatePowerFeed(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerFeedFunc: func(
				_ context.Context, _ string, id int, in domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return &domain.PowerFeed{ID: id, Name: in.Name}, nil
			},
		})
		pf, err := svc.UpdatePowerFeed(context.Background(), 7, domain.PowerFeedWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdatePowerFeed() returned error: %v", err)
		}
		if pf.ID != 7 {
			t.Errorf("power feed ID = %d, want 7", pf.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdatePowerFeedFunc: func(
				_ context.Context, _ string, _ int, _ domain.PowerFeedWrite,
			) (*domain.PowerFeed, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdatePowerFeed(context.Background(), 7, domain.PowerFeedWrite{})
		if err == nil {
			t.Fatal("UpdatePowerFeed() expected error, got nil")
		}
	})
}

func TestNetworkService_DeletePowerFeed(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerFeedFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeletePowerFeed(context.Background(), 3); err != nil {
			t.Fatalf("DeletePowerFeed() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeletePowerFeedFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeletePowerFeed(context.Background(), 3); err == nil {
			t.Fatal("DeletePowerFeed() expected error, got nil")
		}
	})
}

func TestNetworkService_ListDeviceBays(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListDeviceBaysFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.DeviceBay], error) {
				return &domain.PaginatedResponse[domain.DeviceBay]{Count: 1}, nil
			},
		})
		resp, err := svc.ListDeviceBays(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListDeviceBays() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListDeviceBaysFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.DeviceBay], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListDeviceBays(context.Background(), nil)
		if err == nil {
			t.Fatal("ListDeviceBays() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateDeviceBay(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceBayFunc: func(
				_ context.Context, _ string, in domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return &domain.DeviceBay{ID: 1, Name: in.Name}, nil
			},
		})
		db, err := svc.CreateDeviceBay(context.Background(), domain.DeviceBayWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateDeviceBay() returned error: %v", err)
		}
		if db.ID != 1 || db.Name != "A" {
			t.Errorf("device bay = %+v, want id 1 name A", db)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateDeviceBayFunc: func(
				_ context.Context, _ string, _ domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateDeviceBay(context.Background(), domain.DeviceBayWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateDeviceBay() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateDeviceBay(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceBayFunc: func(
				_ context.Context, _ string, id int, in domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return &domain.DeviceBay{ID: id, Name: in.Name}, nil
			},
		})
		db, err := svc.UpdateDeviceBay(context.Background(), 7, domain.DeviceBayWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateDeviceBay() returned error: %v", err)
		}
		if db.ID != 7 {
			t.Errorf("device bay ID = %d, want 7", db.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateDeviceBayFunc: func(
				_ context.Context, _ string, _ int, _ domain.DeviceBayWrite,
			) (*domain.DeviceBay, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateDeviceBay(context.Background(), 7, domain.DeviceBayWrite{})
		if err == nil {
			t.Fatal("UpdateDeviceBay() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteDeviceBay(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceBayFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteDeviceBay(context.Background(), 3); err != nil {
			t.Fatalf("DeleteDeviceBay() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteDeviceBayFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteDeviceBay(context.Background(), 3); err == nil {
			t.Fatal("DeleteDeviceBay() expected error, got nil")
		}
	})
}

func TestNetworkService_ListRackReservations(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRackReservationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				return &domain.PaginatedResponse[domain.RackReservation]{Count: 1}, nil
			},
		})
		resp, err := svc.ListRackReservations(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListRackReservations() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListRackReservationsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.RackReservation], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListRackReservations(context.Background(), nil)
		if err == nil {
			t.Fatal("ListRackReservations() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateRackReservation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRackReservationFunc: func(
				_ context.Context, _ string, in domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: 1, Description: in.Description}, nil
			},
		})
		rr, err := svc.CreateRackReservation(context.Background(), domain.RackReservationWrite{Description: "reserved"})
		if err != nil {
			t.Fatalf("CreateRackReservation() returned error: %v", err)
		}
		if rr.ID != 1 || rr.Description != "reserved" {
			t.Errorf("rack reservation = %+v, want id 1 reserved", rr)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateRackReservationFunc: func(
				_ context.Context, _ string, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateRackReservation(context.Background(), domain.RackReservationWrite{})
		if err == nil {
			t.Fatal("CreateRackReservation() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateRackReservation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRackReservationFunc: func(
				_ context.Context, _ string, id int, in domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return &domain.RackReservation{ID: id, Description: in.Description}, nil
			},
		})
		rr, err := svc.UpdateRackReservation(context.Background(), 7, domain.RackReservationWrite{Description: "x"})
		if err != nil {
			t.Fatalf("UpdateRackReservation() returned error: %v", err)
		}
		if rr.ID != 7 {
			t.Errorf("rack reservation ID = %d, want 7", rr.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateRackReservationFunc: func(
				_ context.Context, _ string, _ int, _ domain.RackReservationWrite,
			) (*domain.RackReservation, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateRackReservation(context.Background(), 7, domain.RackReservationWrite{})
		if err == nil {
			t.Fatal("UpdateRackReservation() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteRackReservation(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRackReservationFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteRackReservation(context.Background(), 3); err != nil {
			t.Fatalf("DeleteRackReservation() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteRackReservationFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteRackReservation(context.Background(), 3); err == nil {
			t.Fatal("DeleteRackReservation() expected error, got nil")
		}
	})
}

func TestNetworkService_ListInventoryItems(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListInventoryItemsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItem], error) {
				return &domain.PaginatedResponse[domain.InventoryItem]{Count: 1}, nil
			},
		})
		resp, err := svc.ListInventoryItems(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListInventoryItems() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListInventoryItemsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.InventoryItem], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListInventoryItems(context.Background(), nil)
		if err == nil {
			t.Fatal("ListInventoryItems() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateInventoryItem(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateInventoryItemFunc: func(
				_ context.Context, _ string, in domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{ID: 1, Name: in.Name}, nil
			},
		})
		ii, err := svc.CreateInventoryItem(context.Background(), domain.InventoryItemWrite{Name: "A"})
		if err != nil {
			t.Fatalf("CreateInventoryItem() returned error: %v", err)
		}
		if ii.ID != 1 || ii.Name != "A" {
			t.Errorf("inventory item = %+v, want id 1 name A", ii)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateInventoryItemFunc: func(
				_ context.Context, _ string, _ domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateInventoryItem(context.Background(), domain.InventoryItemWrite{Name: "A"})
		if err == nil {
			t.Fatal("CreateInventoryItem() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateInventoryItem(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateInventoryItemFunc: func(
				_ context.Context, _ string, id int, in domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return &domain.InventoryItem{ID: id, Name: in.Name}, nil
			},
		})
		ii, err := svc.UpdateInventoryItem(context.Background(), 7, domain.InventoryItemWrite{Name: "A"})
		if err != nil {
			t.Fatalf("UpdateInventoryItem() returned error: %v", err)
		}
		if ii.ID != 7 {
			t.Errorf("inventory item ID = %d, want 7", ii.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateInventoryItemFunc: func(
				_ context.Context, _ string, _ int, _ domain.InventoryItemWrite,
			) (*domain.InventoryItem, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateInventoryItem(context.Background(), 7, domain.InventoryItemWrite{})
		if err == nil {
			t.Fatal("UpdateInventoryItem() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteInventoryItem(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteInventoryItemFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteInventoryItem(context.Background(), 3); err != nil {
			t.Fatalf("DeleteInventoryItem() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteInventoryItemFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteInventoryItem(context.Background(), 3); err == nil {
			t.Fatal("DeleteInventoryItem() expected error, got nil")
		}
	})
}

func TestNetworkService_ListMACAddresses(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListMACAddressesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				return &domain.PaginatedResponse[domain.MACAddress]{Count: 1}, nil
			},
		})
		resp, err := svc.ListMACAddresses(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListMACAddresses() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListMACAddressesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.MACAddress], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListMACAddresses(context.Background(), nil)
		if err == nil {
			t.Fatal("ListMACAddresses() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateMACAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateMACAddressFunc: func(
				_ context.Context, _ string, in domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: 1, Address: in.Address}, nil
			},
		})
		ma, err := svc.CreateMACAddress(context.Background(), domain.MACAddressWrite{Address: "00:11:22:33:44:55"})
		if err != nil {
			t.Fatalf("CreateMACAddress() returned error: %v", err)
		}
		if ma.ID != 1 || ma.Address != "00:11:22:33:44:55" {
			t.Errorf("MAC address = %+v, want id 1 with MAC", ma)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateMACAddressFunc: func(
				_ context.Context, _ string, _ domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateMACAddress(context.Background(), domain.MACAddressWrite{Address: "00:11:22:33:44:55"})
		if err == nil {
			t.Fatal("CreateMACAddress() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateMACAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateMACAddressFunc: func(
				_ context.Context, _ string, id int, in domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return &domain.MACAddress{ID: id, Address: in.Address}, nil
			},
		})
		ma, err := svc.UpdateMACAddress(context.Background(), 7, domain.MACAddressWrite{Address: "00:11:22:33:44:55"})
		if err != nil {
			t.Fatalf("UpdateMACAddress() returned error: %v", err)
		}
		if ma.ID != 7 {
			t.Errorf("MAC address ID = %d, want 7", ma.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateMACAddressFunc: func(
				_ context.Context, _ string, _ int, _ domain.MACAddressWrite,
			) (*domain.MACAddress, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateMACAddress(context.Background(), 7, domain.MACAddressWrite{})
		if err == nil {
			t.Fatal("UpdateMACAddress() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteMACAddress(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteMACAddressFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteMACAddress(context.Background(), 3); err != nil {
			t.Fatalf("DeleteMACAddress() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteMACAddressFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteMACAddress(context.Background(), 3); err == nil {
			t.Fatal("DeleteMACAddress() expected error, got nil")
		}
	})
}
