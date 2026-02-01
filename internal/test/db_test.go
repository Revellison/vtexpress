package db_test

import (
	"os"
	"testing"
	"vtgui/internal/modules"
)

var simulated_api = "e00b5ad1-54c5-4048-b910-9fa65145664e"

func TestDBWorker(t *testing.T) {
	os.Remove("test.db")

	w, err := modules.NewStorage("test.db")
	if err != nil {
		t.Fatal(err)
	}

	// сеткей
	if err := w.SetKey("vt_api_key", simulated_api); err != nil {
		t.Fatal(err)
	}

	// геткей аеаеаеае
	val, err := w.GetKey("vt_api_key")
	if err != nil {
		t.Fatal(err)
	}

	if val != simulated_api {
		t.Fatalf("expected %s, got %s", simulated_api, val)
	}

}
