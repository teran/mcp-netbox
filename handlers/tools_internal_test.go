package handlers

import (
	"testing"
)

func TestPaginationParams(t *testing.T) {
	t.Parallel()

	t.Run("default values when page=0 and pageSize=0", func(t *testing.T) {
		params := paginationParams(0, 0)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})

	t.Run("page=1 and pageSize=25", func(t *testing.T) {
		params := paginationParams(1, 25)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})

	t.Run("page=3 and pageSize=10", func(t *testing.T) {
		params := paginationParams(3, 10)
		if params["offset"] != "20" {
			t.Errorf("offset = %q, want %q", params["offset"], "20")
		}
		if params["limit"] != "10" {
			t.Errorf("limit = %q, want %q", params["limit"], "10")
		}
	})

	t.Run("pageSize capped at 100", func(t *testing.T) {
		params := paginationParams(1, 200)
		if params["limit"] != "100" {
			t.Errorf("limit = %q, want %q", params["limit"], "100")
		}
	})

	t.Run("negative page defaults to 1", func(t *testing.T) {
		params := paginationParams(-5, 25)
		if params["offset"] != "0" {
			t.Errorf("offset = %q, want %q", params["offset"], "0")
		}
	})

	t.Run("negative pageSize defaults to 25", func(t *testing.T) {
		params := paginationParams(1, -5)
		if params["limit"] != "25" {
			t.Errorf("limit = %q, want %q", params["limit"], "25")
		}
	})
}

func TestAddIntParam(t *testing.T) {
	t.Parallel()

	t.Run("non-zero value added", func(t *testing.T) {
		m := make(map[string]string)
		addIntParam(m, "family", 4)
		if m["family"] != "4" {
			t.Errorf("family = %q, want %q", m["family"], "4")
		}
	})

	t.Run("zero value skipped", func(t *testing.T) {
		m := make(map[string]string)
		addIntParam(m, "family", 0)
		if _, ok := m["family"]; ok {
			t.Error("family param should not be set for zero value")
		}
	})
}

func TestAddParam(t *testing.T) {
	t.Parallel()

	t.Run("non-empty value added", func(t *testing.T) {
		m := make(map[string]string)
		addParam(m, "site", "dc1")
		if m["site"] != "dc1" {
			t.Errorf("site = %q, want %q", m["site"], "dc1")
		}
	})

	t.Run("empty value skipped", func(t *testing.T) {
		m := make(map[string]string)
		addParam(m, "site", "")
		if _, ok := m["site"]; ok {
			t.Error("site param should not be set for empty value")
		}
	})
}
