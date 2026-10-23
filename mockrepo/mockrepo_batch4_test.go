package mockrepo

import (
	"context"
	"testing"

	"github.com/teran/mcp-netbox/domain"
)

// Coverage for the Batch 4 entities (contact_group, contact_role, virtual_disk,
// virtual_machine_type): MockRepo delegates forward to their Func fields and
// NilRepo returns empty/sentinel values. Both implementations are exercised so
// the batch's repository methods count toward the coverage gate.

func TestMockRepo_Batch4Delegates(t *testing.T) { //nolint:gocognit,gocyclo,funlen,maintidx
	t.Parallel()

	t.Run("ListContactGroups", func(t *testing.T) {
		got, err := (&MockRepo{
			ListContactGroupsFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ContactGroup], error) {
				return &domain.PaginatedResponse[domain.ContactGroup]{Count: 1}, nil
			},
		}).ListContactGroups(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListContactGroups = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateContactGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateContactGroupFunc: func(
				_ context.Context, _ string, in domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: 1, Name: in.Name}, nil
			},
		}).CreateContactGroup(context.Background(), "t", domain.ContactGroupWrite{Name: "Admins"})
		if err != nil || got == nil || got.Name != "Admins" {
			t.Errorf("CreateContactGroup = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateContactGroup", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateContactGroupFunc: func(
				_ context.Context, _ string, id int, in domain.ContactGroupWrite,
			) (*domain.ContactGroup, error) {
				return &domain.ContactGroup{ID: id, Name: in.Name}, nil
			},
		}).UpdateContactGroup(context.Background(), "t", 7, domain.ContactGroupWrite{Name: "Ops"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateContactGroup = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteContactGroup", func(t *testing.T) {
		err := (&MockRepo{
			DeleteContactGroupFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteContactGroup(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteContactGroup = %v", err)
		}
	})

	t.Run("ListContactRoles", func(t *testing.T) {
		got, err := (&MockRepo{
			ListContactRolesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.ContactRole], error) {
				return &domain.PaginatedResponse[domain.ContactRole]{Count: 1}, nil
			},
		}).ListContactRoles(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListContactRoles = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateContactRole", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateContactRoleFunc: func(
				_ context.Context, _ string, in domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return &domain.ContactRole{ID: 1, Name: in.Name}, nil
			},
		}).CreateContactRole(context.Background(), "t", domain.ContactRoleWrite{Name: "Admin"})
		if err != nil || got == nil || got.Name != "Admin" {
			t.Errorf("CreateContactRole = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateContactRole", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateContactRoleFunc: func(
				_ context.Context, _ string, id int, in domain.ContactRoleWrite,
			) (*domain.ContactRole, error) {
				return &domain.ContactRole{ID: id, Name: in.Name}, nil
			},
		}).UpdateContactRole(context.Background(), "t", 7, domain.ContactRoleWrite{Name: "Operator"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateContactRole = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteContactRole", func(t *testing.T) {
		err := (&MockRepo{
			DeleteContactRoleFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteContactRole(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteContactRole = %v", err)
		}
	})

	t.Run("ListVirtualDisks", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualDisksFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualDisk], error) {
				return &domain.PaginatedResponse[domain.VirtualDisk]{Count: 1}, nil
			},
		}).ListVirtualDisks(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualDisks = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateVirtualDisk", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualDiskFunc: func(
				_ context.Context, _ string, in domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: 1, Name: in.Name}, nil
			},
		}).CreateVirtualDisk(context.Background(), "t", domain.VirtualDiskWrite{Name: "Data"})
		if err != nil || got == nil || got.Name != "Data" {
			t.Errorf("CreateVirtualDisk = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateVirtualDisk", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualDiskFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualDiskWrite,
			) (*domain.VirtualDisk, error) {
				return &domain.VirtualDisk{ID: id, Name: in.Name}, nil
			},
		}).UpdateVirtualDisk(context.Background(), "t", 7, domain.VirtualDiskWrite{Name: "Renamed"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualDisk = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteVirtualDisk", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualDiskFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualDisk(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualDisk = %v", err)
		}
	})

	t.Run("ListVirtualMachineTypes", func(t *testing.T) {
		got, err := (&MockRepo{
			ListVirtualMachineTypesFunc: func(
				context.Context, string, map[string]string,
			) (*domain.PaginatedResponse[domain.VirtualMachineType], error) {
				return &domain.PaginatedResponse[domain.VirtualMachineType]{Count: 1}, nil
			},
		}).ListVirtualMachineTypes(context.Background(), "t", nil)
		if err != nil || got.Count != 1 {
			t.Errorf("ListVirtualMachineTypes = (%v, %v), want count 1", got, err)
		}
	})

	t.Run("CreateVirtualMachineType", func(t *testing.T) {
		got, err := (&MockRepo{
			CreateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, in domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return &domain.VirtualMachineType{ID: 1, Name: in.Name}, nil
			},
		}).CreateVirtualMachineType(context.Background(), "t", domain.VirtualMachineTypeWrite{Name: "Standard"})
		if err != nil || got == nil || got.Name != "Standard" {
			t.Errorf("CreateVirtualMachineType = (%v, %v)", got, err)
		}
	})

	t.Run("UpdateVirtualMachineType", func(t *testing.T) {
		got, err := (&MockRepo{
			UpdateVirtualMachineTypeFunc: func(
				_ context.Context, _ string, id int, in domain.VirtualMachineTypeWrite,
			) (*domain.VirtualMachineType, error) {
				return &domain.VirtualMachineType{ID: id, Name: in.Name}, nil
			},
		}).UpdateVirtualMachineType(context.Background(), "t", 7, domain.VirtualMachineTypeWrite{Name: "Tiny"})
		if err != nil || got == nil || got.ID != 7 {
			t.Errorf("UpdateVirtualMachineType = (%v, %v)", got, err)
		}
	})

	t.Run("DeleteVirtualMachineType", func(t *testing.T) {
		err := (&MockRepo{
			DeleteVirtualMachineTypeFunc: func(_ context.Context, _ string, id int) error {
				if id != 3 {
					t.Errorf("id = %d, want 3", id)
				}

				return nil
			},
		}).DeleteVirtualMachineType(context.Background(), "t", 3)
		if err != nil {
			t.Errorf("DeleteVirtualMachineType = %v", err)
		}
	})
}

func TestNilRepo_Batch4(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := &NilRepo{}

	if r, err := repo.ListContactGroups(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListContactGroups = (%v, %v), want count 0", r, err)
	}
	if g, err := repo.CreateContactGroup(ctx, "t", domain.ContactGroupWrite{Name: "Admins"}); err != nil || g.ID != 1 {
		t.Errorf("CreateContactGroup = (%v, %v)", g, err)
	}
	if g, err := repo.UpdateContactGroup(ctx, "t", 7, domain.ContactGroupWrite{Name: "Ops"}); err != nil || g.ID != 7 {
		t.Errorf("UpdateContactGroup = (%v, %v)", g, err)
	}
	if err := repo.DeleteContactGroup(ctx, "t", 3); err != nil {
		t.Errorf("DeleteContactGroup = %v", err)
	}

	if r, err := repo.ListContactRoles(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListContactRoles = (%v, %v), want count 0", r, err)
	}
	if c, err := repo.CreateContactRole(ctx, "t", domain.ContactRoleWrite{Name: "Admin"}); err != nil || c.ID != 1 {
		t.Errorf("CreateContactRole = (%v, %v)", c, err)
	}
	if c, err := repo.UpdateContactRole(ctx, "t", 7, domain.ContactRoleWrite{Name: "Operator"}); err != nil || c.ID != 7 {
		t.Errorf("UpdateContactRole = (%v, %v)", c, err)
	}
	if err := repo.DeleteContactRole(ctx, "t", 3); err != nil {
		t.Errorf("DeleteContactRole = %v", err)
	}

	if r, err := repo.ListVirtualDisks(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListVirtualDisks = (%v, %v), want count 0", r, err)
	}
	if d, err := repo.CreateVirtualDisk(ctx, "t", domain.VirtualDiskWrite{Name: "Data"}); err != nil || d.ID != 1 {
		t.Errorf("CreateVirtualDisk = (%v, %v)", d, err)
	}
	if d, err := repo.UpdateVirtualDisk(ctx, "t", 7, domain.VirtualDiskWrite{Name: "Renamed"}); err != nil || d.ID != 7 {
		t.Errorf("UpdateVirtualDisk = (%v, %v)", d, err)
	}
	if err := repo.DeleteVirtualDisk(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualDisk = %v", err)
	}

	if r, err := repo.ListVirtualMachineTypes(ctx, "t", nil); err != nil || r.Count != 0 {
		t.Errorf("ListVirtualMachineTypes = (%v, %v), want count 0", r, err)
	}
	if v, err := repo.CreateVirtualMachineType(
		ctx, "t", domain.VirtualMachineTypeWrite{Name: "Standard"},
	); err != nil || v.ID != 1 {
		t.Errorf("CreateVirtualMachineType = (%v, %v)", v, err)
	}
	if v, err := repo.UpdateVirtualMachineType(
		ctx, "t", 7, domain.VirtualMachineTypeWrite{Name: "Tiny"},
	); err != nil || v.ID != 7 {
		t.Errorf("UpdateVirtualMachineType = (%v, %v)", v, err)
	}
	if err := repo.DeleteVirtualMachineType(ctx, "t", 3); err != nil {
		t.Errorf("DeleteVirtualMachineType = %v", err)
	}
}
