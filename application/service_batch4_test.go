package application

import (
	"context"
	"errors"
	"testing"

	"github.com/teran/mcp-netbox/domain"
	"github.com/teran/mcp-netbox/mockrepo"
)

// Service-layer tests for the entities added in feature/expand-entities
// Batch 4 (contact_group, contact_role, virtual_disk, virtual_machine_type).
// These cover the thin relay methods in application/service.go and exercise the
// mockrepo methods they delegate to. Each entity gets success + error
// propagation cases so both branches of the wrapper are covered.

func TestNetworkService_ListContactGroups(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListContactGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				return &domain.PaginatedResponse[domain.ContactGroup]{Count: 1}, nil
			},
		})
		resp, err := svc.ListContactGroups(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListContactGroups() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListContactGroupsFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListContactGroups(context.Background(), nil)
		if err == nil {
			t.Fatal("ListContactGroups() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateContactGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateContactGroupFunc: func(
				_ context.Context, _ string, in domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: 1, Name: in.Name}, nil
			},
		})
		cg, err := svc.CreateContactGroup(context.Background(), domain.ContactGroupWrite{Name: "Admins"})
		if err != nil {
			t.Fatalf("CreateContactGroup() returned error: %v", err)
		}
		if cg.ID != 1 || cg.Name != "Admins" {
			t.Errorf("contact group = %+v, want id 1 name Admins", cg)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateContactGroupFunc: func(
				_ context.Context, _ string, _ domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateContactGroup(context.Background(), domain.ContactGroupWrite{})
		if err == nil {
			t.Fatal("CreateContactGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateContactGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateContactGroupFunc: func(
				_ context.Context, _ string, id int, in domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: id, Name: in.Name}, nil
			},
		})
		cg, err := svc.UpdateContactGroup(context.Background(), 7, domain.ContactGroupWrite{Name: "Ops"})
		if err != nil {
			t.Fatalf("UpdateContactGroup() returned error: %v", err)
		}
		if cg.ID != 7 {
			t.Errorf("contact group ID = %d, want 7", cg.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateContactGroupFunc: func(
				_ context.Context, _ string, _ int, _ domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateContactGroup(context.Background(), 7, domain.ContactGroupWrite{})
		if err == nil {
			t.Fatal("UpdateContactGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteContactGroup(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteContactGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteContactGroup(context.Background(), 3); err != nil {
			t.Fatalf("DeleteContactGroup() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteContactGroupFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteContactGroup(context.Background(), 3); err == nil {
			t.Fatal("DeleteContactGroup() expected error, got nil")
		}
	})
}

func TestNetworkService_ListContactRoles(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListContactRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactRole], error) {
				return &domain.PaginatedResponse[domain.ContactRole]{Count: 1}, nil
			},
		})
		resp, err := svc.ListContactRoles(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListContactRoles() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListContactRolesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.ContactRole], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListContactRoles(context.Background(), nil)
		if err == nil {
			t.Fatal("ListContactRoles() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateContactRole(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateContactRoleFunc: func(
				_ context.Context, _ string, in domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return &domain.ContactRole{ID: 1, Name: in.Name}, nil
			},
		})
		cr, err := svc.CreateContactRole(context.Background(), domain.ContactRoleWrite{Name: "Admin"})
		if err != nil {
			t.Fatalf("CreateContactRole() returned error: %v", err)
		}
		if cr.ID != 1 || cr.Name != "Admin" {
			t.Errorf("contact role = %+v, want id 1 name Admin", cr)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateContactRoleFunc: func(
				_ context.Context, _ string, _ domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateContactRole(context.Background(), domain.ContactRoleWrite{})
		if err == nil {
			t.Fatal("CreateContactRole() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateContactRole(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateContactRoleFunc: func(
				_ context.Context, _ string, id int, in domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return &domain.ContactRole{ID: id, Name: in.Name}, nil
			},
		})
		cr, err := svc.UpdateContactRole(context.Background(), 7, domain.ContactRoleWrite{Name: "Operator"})
		if err != nil {
			t.Fatalf("UpdateContactRole() returned error: %v", err)
		}
		if cr.ID != 7 {
			t.Errorf("contact role ID = %d, want 7", cr.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateContactRoleFunc: func(
				_ context.Context, _ string, _ int, _ domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateContactRole(context.Background(), 7, domain.ContactRoleWrite{})
		if err == nil {
			t.Fatal("UpdateContactRole() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteContactRole(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteContactRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteContactRole(context.Background(), 3); err != nil {
			t.Fatalf("DeleteContactRole() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteContactRoleFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteContactRole(context.Background(), 3); err == nil {
			t.Fatal("DeleteContactRole() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualDisks(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualDisksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				return &domain.PaginatedResponse[domain.VirtualDisk]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVirtualDisks(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualDisks() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualDisksFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListVirtualDisks(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVirtualDisks() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVirtualDisk(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualDiskFunc: func(
				_ context.Context, _ string, in domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: 1, Name: in.Name}, nil
			},
		})
		vd, err := svc.CreateVirtualDisk(context.Background(), domain.VirtualDiskWrite{Name: "Data"})
		if err != nil {
			t.Fatalf("CreateVirtualDisk() returned error: %v", err)
		}
		if vd.ID != 1 || vd.Name != "Data" {
			t.Errorf("virtual disk = %+v, want id 1 name Data", vd)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualDiskFunc: func(
				_ context.Context, _ string, _ domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualDisk(context.Background(), domain.VirtualDiskWrite{})
		if err == nil {
			t.Fatal("CreateVirtualDisk() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVirtualDisk(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualDiskFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: id, Name: in.Name}, nil
			},
		})
		vd, err := svc.UpdateVirtualDisk(context.Background(), 7, domain.VirtualDiskWrite{Name: "Renamed"})
		if err != nil {
			t.Fatalf("UpdateVirtualDisk() returned error: %v", err)
		}
		if vd.ID != 7 {
			t.Errorf("virtual disk ID = %d, want 7", vd.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualDiskFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualDisk(context.Background(), 7, domain.VirtualDiskWrite{})
		if err == nil {
			t.Fatal("UpdateVirtualDisk() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVirtualDisk(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualDiskFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVirtualDisk(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualDisk() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualDiskFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVirtualDisk(context.Background(), 3); err == nil {
			t.Fatal("DeleteVirtualDisk() expected error, got nil")
		}
	})
}

func TestNetworkService_ListVirtualMachineTypes(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualMachineTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
				return &domain.PaginatedResponse[domain.VirtualMachineType]{Count: 1}, nil
			},
		})
		resp, err := svc.ListVirtualMachineTypes(context.Background(), nil)
		if err != nil {
			t.Fatalf("ListVirtualMachineTypes() returned error: %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("resp.Count = %d, want 1", resp.Count)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			ListVirtualMachineTypesFunc: func(
				_ context.Context, _ string, _ map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.ListVirtualMachineTypes(context.Background(), nil)
		if err == nil {
			t.Fatal("ListVirtualMachineTypes() expected error, got nil")
		}
	})
}

func TestNetworkService_CreateVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, in domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return &domain.VirtualMachineType{ID: 1, Name: in.Name}, nil
			},
		})
		vmt, err := svc.CreateVirtualMachineType(context.Background(), domain.VirtualMachineTypeWrite{Name: "Standard"})
		if err != nil {
			t.Fatalf("CreateVirtualMachineType() returned error: %v", err)
		}
		if vmt.ID != 1 || vmt.Name != "Standard" {
			t.Errorf("virtual machine type = %+v, want id 1 name Standard", vmt)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			CreateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, _ domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.CreateVirtualMachineType(context.Background(), domain.VirtualMachineTypeWrite{})
		if err == nil {
			t.Fatal("CreateVirtualMachineType() expected error, got nil")
		}
	})
}

func TestNetworkService_UpdateVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return &domain.VirtualMachineType{ID: id, Name: in.Name}, nil
			},
		})
		vmt, err := svc.UpdateVirtualMachineType(context.Background(), 7, domain.VirtualMachineTypeWrite{Name: "Tiny"})
		if err != nil {
			t.Fatalf("UpdateVirtualMachineType() returned error: %v", err)
		}
		if vmt.ID != 7 {
			t.Errorf("virtual machine type ID = %d, want 7", vmt.ID)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			UpdateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, _ int, _ domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return nil, errors.New("boom")
			},
		})
		_, err := svc.UpdateVirtualMachineType(context.Background(), 7, domain.VirtualMachineTypeWrite{})
		if err == nil {
			t.Fatal("UpdateVirtualMachineType() expected error, got nil")
		}
	})
}

func TestNetworkService_DeleteVirtualMachineType(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualMachineTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		})
		if err := svc.DeleteVirtualMachineType(context.Background(), 3); err != nil {
			t.Fatalf("DeleteVirtualMachineType() returned error: %v", err)
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		svc := newTestService(&mockrepo.MockRepo{
			DeleteVirtualMachineTypeFunc: func(_ context.Context, _ string, _ int) error {
				return errors.New("boom")
			},
		})
		if err := svc.DeleteVirtualMachineType(context.Background(), 3); err == nil {
			t.Fatal("DeleteVirtualMachineType() expected error, got nil")
		}
	})
}
