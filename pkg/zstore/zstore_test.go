package zstore_test

import (
	"errors"
	"testing"

	"github.com/zarlcorp/core/pkg/zfilesystem"
	"github.com/zarlcorp/core/pkg/zstore"
)

func TestOpenFirstRun(t *testing.T) {
	fs := zfilesystem.NewMemFS()
	s, err := zstore.Open(fs, []byte("password"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer s.Close()

	// salt and verify files should exist
	if _, err := fs.ReadFile("salt"); err != nil {
		t.Fatal("salt file not created")
	}
	if _, err := fs.ReadFile("verify"); err != nil {
		t.Fatal("verify file not created")
	}
}

func TestOpenSubsequentRun(t *testing.T) {
	fs := zfilesystem.NewMemFS()

	s1, err := zstore.Open(fs, []byte("password"))
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	s1.Close()

	s2, err := zstore.Open(fs, []byte("password"))
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	s2.Close()
}

func TestOpenWrongPassword(t *testing.T) {
	fs := zfilesystem.NewMemFS()

	s, err := zstore.Open(fs, []byte("correct"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	s.Close()

	_, err = zstore.Open(fs, []byte("wrong"))
	if !errors.Is(err, zstore.ErrWrongPassword) {
		t.Fatalf("expected ErrWrongPassword, got %v", err)
	}
}

func TestCollectionPutGet(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	type item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	col, err := zstore.NewCollection[item](s, "items")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	want := item{Name: "widget", Count: 42}
	if err := col.Put("w1", want); err != nil {
		t.Fatalf("put: %v", err)
	}

	got, err := col.Get("w1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestCollectionGetNotFound(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "things")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	_, err = col.Get("nonexistent")
	if !errors.Is(err, zstore.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCollectionDelete(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "things")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	if err := col.Put("k1", "value1"); err != nil {
		t.Fatalf("put: %v", err)
	}

	if err := col.Delete("k1"); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = col.Get("k1")
	if !errors.Is(err, zstore.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestCollectionDeleteNotFound(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "things")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	err = col.Delete("nonexistent")
	if !errors.Is(err, zstore.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCollectionList(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "names")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	inputs := map[string]string{
		"a": "alice",
		"b": "bob",
		"c": "charlie",
	}
	for id, v := range inputs {
		if err := col.Put(id, v); err != nil {
			t.Fatalf("put %s: %v", id, err)
		}
	}

	got, err := col.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(got) != len(inputs) {
		t.Fatalf("list returned %d items, want %d", len(got), len(inputs))
	}

	// verify all values present (order not guaranteed)
	found := make(map[string]bool)
	for _, v := range got {
		found[v] = true
	}
	for _, v := range inputs {
		if !found[v] {
			t.Fatalf("list missing value %q", v)
		}
	}
}

func TestCollectionListEmpty(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "empty")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	got, err := col.List()
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty list, got %d items", len(got))
	}
}

func TestCollectionLen(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[int](s, "numbers")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	assertLen := func(want int) {
		t.Helper()
		n, err := col.Len()
		if err != nil {
			t.Fatalf("len: %v", err)
		}
		if n != want {
			t.Fatalf("len = %d, want %d", n, want)
		}
	}

	assertLen(0)

	for i := range 5 {
		if err := col.Put(string(rune('a'+i)), i); err != nil {
			t.Fatalf("put: %v", err)
		}
	}
	assertLen(5)

	if err := col.Delete(string(rune('a'))); err != nil {
		t.Fatalf("delete: %v", err)
	}
	assertLen(4)
}

func TestMultipleCollectionsIsolated(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	type user struct {
		Name string `json:"name"`
	}

	type secret struct {
		Value string `json:"value"`
	}

	users, err := zstore.NewCollection[user](s, "users")
	if err != nil {
		t.Fatalf("new users collection: %v", err)
	}

	secrets, err := zstore.NewCollection[secret](s, "secrets")
	if err != nil {
		t.Fatalf("new secrets collection: %v", err)
	}

	if err := users.Put("u1", user{Name: "alice"}); err != nil {
		t.Fatalf("put user: %v", err)
	}

	if err := secrets.Put("s1", secret{Value: "hunter2"}); err != nil {
		t.Fatalf("put secret: %v", err)
	}

	// each collection only sees its own data
	uLen, err := users.Len()
	if err != nil {
		t.Fatalf("users len: %v", err)
	}
	if uLen != 1 {
		t.Fatalf("users len = %d, want 1", uLen)
	}

	sLen, err := secrets.Len()
	if err != nil {
		t.Fatalf("secrets len: %v", err)
	}
	if sLen != 1 {
		t.Fatalf("secrets len = %d, want 1", sLen)
	}

	// verify the data is correct
	u, err := users.Get("u1")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if u.Name != "alice" {
		t.Fatalf("user name = %q, want %q", u.Name, "alice")
	}

	sv, err := secrets.Get("s1")
	if err != nil {
		t.Fatalf("get secret: %v", err)
	}
	if sv.Value != "hunter2" {
		t.Fatalf("secret value = %q, want %q", sv.Value, "hunter2")
	}
}

func TestPutOverwrite(t *testing.T) {
	s := openTestStore(t)
	defer s.Close()

	col, err := zstore.NewCollection[string](s, "things")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	if err := col.Put("k1", "first"); err != nil {
		t.Fatalf("put first: %v", err)
	}

	if err := col.Put("k1", "second"); err != nil {
		t.Fatalf("put second: %v", err)
	}

	got, err := col.Get("k1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got != "second" {
		t.Fatalf("got %q, want %q", got, "second")
	}

	n, err := col.Len()
	if err != nil {
		t.Fatalf("len: %v", err)
	}
	if n != 1 {
		t.Fatalf("len = %d, want 1 after overwrite", n)
	}
}

func TestCloseErasesKeys(t *testing.T) {
	fs := zfilesystem.NewMemFS()
	s, err := zstore.Open(fs, []byte("password"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}

	// create a collection to generate a sub-key
	_, err = zstore.NewCollection[string](s, "test")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	// grab references to the key slices before closing
	masterKey := s.MasterKeyForTest()
	subKeys := s.SubKeysForTest()

	s.Close()

	// master key should be zeroed
	for i, b := range masterKey {
		if b != 0 {
			t.Fatalf("master key byte %d = %d, want 0", i, b)
		}
	}

	// sub-keys should be zeroed
	for j, sk := range subKeys {
		for i, b := range sk {
			if b != 0 {
				t.Fatalf("sub-key %d byte %d = %d, want 0", j, i, b)
			}
		}
	}
}

func TestDataPersistsAcrossOpens(t *testing.T) {
	fs := zfilesystem.NewMemFS()
	password := []byte("password")

	// first session: write data
	s1, err := zstore.Open(fs, password)
	if err != nil {
		t.Fatalf("first open: %v", err)
	}

	col1, err := zstore.NewCollection[string](s1, "notes")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	if err := col1.Put("n1", "hello"); err != nil {
		t.Fatalf("put: %v", err)
	}
	s1.Close()

	// second session: read data
	s2, err := zstore.Open(fs, password)
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	defer s2.Close()

	col2, err := zstore.NewCollection[string](s2, "notes")
	if err != nil {
		t.Fatalf("new collection: %v", err)
	}

	got, err := col2.Get("n1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if got != "hello" {
		t.Fatalf("got %q, want %q", got, "hello")
	}
}

// openTestStore creates a store with an in-memory filesystem for testing.
func openTestStore(t *testing.T) *zstore.Store {
	t.Helper()
	fs := zfilesystem.NewMemFS()
	s, err := zstore.Open(fs, []byte("test-password"))
	if err != nil {
		t.Fatalf("open test store: %v", err)
	}
	return s
}
