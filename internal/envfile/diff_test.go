package envfile

import (
	"testing"
)

func TestDiffNoChanges(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}
	d := Diff(old, new)
	if !d.IsEmpty() {
		t.Errorf("expected empty diff, got added=%v removed=%v changed=%v", d.Added, d.Removed, d.Changed)
	}
}

func TestDiffAdded(t *testing.T) {
	old := map[string]string{"FOO": "bar"}
	new := map[string]string{"FOO": "bar", "NEW_KEY": "new_val"}
	d := Diff(old, new)
	if len(d.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(d.Added))
	}
	if d.Added["NEW_KEY"] != "new_val" {
		t.Errorf("expected NEW_KEY=new_val, got %q", d.Added["NEW_KEY"])
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected removed or changed entries")
	}
}

func TestDiffRemoved(t *testing.T) {
	old := map[string]string{"FOO": "bar", "GONE": "bye"}
	new := map[string]string{"FOO": "bar"}
	d := Diff(old, new)
	if len(d.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(d.Removed))
	}
	if d.Removed["GONE"] != "bye" {
		t.Errorf("expected GONE=bye, got %q", d.Removed["GONE"])
	}
}

func TestDiffChanged(t *testing.T) {
	old := map[string]string{"FOO": "old_val"}
	new := map[string]string{"FOO": "new_val"}
	d := Diff(old, new)
	if len(d.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(d.Changed))
	}
	pair, ok := d.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in changed")
	}
	if pair[0] != "old_val" || pair[1] != "new_val" {
		t.Errorf("expected [old_val new_val], got %v", pair)
	}
}

func TestDiffMixed(t *testing.T) {
	old := map[string]string{"KEEP": "same", "MODIFY": "v1", "DROP": "gone"}
	new := map[string]string{"KEEP": "same", "MODIFY": "v2", "ADDED": "here"}
	d := Diff(old, new)
	if d.IsEmpty() {
		t.Fatal("expected non-empty diff")
	}
	if len(d.Added) != 1 || d.Added["ADDED"] != "here" {
		t.Errorf("unexpected added: %v", d.Added)
	}
	if len(d.Removed) != 1 || d.Removed["DROP"] != "gone" {
		t.Errorf("unexpected removed: %v", d.Removed)
	}
	if len(d.Changed) != 1 || d.Changed["MODIFY"] != ([2]string{"v1", "v2"}) {
		t.Errorf("unexpected changed: %v", d.Changed)
	}
}

func TestDiffSortedKeys(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	d := Diff(old, new)
	keys := d.SortedAdded()
	if len(keys) != 3 || keys[0] != "ALPHA" || keys[1] != "MANGO" || keys[2] != "ZEBRA" {
		t.Errorf("expected sorted keys [ALPHA MANGO ZEBRA], got %v", keys)
	}
}
