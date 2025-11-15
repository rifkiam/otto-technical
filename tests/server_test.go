package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "be_test/internal/httpapi"
    "be_test/internal/store"
)

type createReq struct { Name string `json:"name"` }
type updateReq struct { Name *string `json:"name,omitempty"`; Done *bool `json:"done,omitempty"` }
type itemRes struct { ID string `json:"id"`; Name string `json:"name"`; Done bool `json:"done"` }

func TestHealthOK(t *testing.T) {
    ts := httptest.NewServer(httpapi.NewServer(store.NewMemoryItemStore()))
    defer ts.Close()
    resp, err := http.Get(ts.URL + "/health")
    if err != nil { t.Fatalf("health request failed: %v", err) }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK { t.Fatalf("expected 200, got %d", resp.StatusCode) }
}

func TestCRUDSimple(t *testing.T) {
    ts := httptest.NewServer(httpapi.NewServer(store.NewMemoryItemStore()))
    defer ts.Close()

    // CREATE
    body := mustJSON(createReq{Name: "Apple"})
    r1, err := http.Post(ts.URL+"/items", "application/json", bytes.NewReader(body))
    if err != nil { t.Fatalf("create failed: %v", err) }
    defer r1.Body.Close()
    if r1.StatusCode != http.StatusCreated { t.Fatalf("expected 201, got %d", r1.StatusCode) }
    var created itemRes
    mustDecode(r1, &created)
    if created.ID == "" || created.Name != "Apple" || created.Done { t.Fatalf("invalid create result: %+v", created) }

    // LIST
    r2, err := http.Get(ts.URL + "/items")
    if err != nil { t.Fatalf("list failed: %v", err) }
    defer r2.Body.Close()
    if r2.StatusCode != http.StatusOK { t.Fatalf("expected 200, got %d", r2.StatusCode) }
    var list struct { Items []itemRes `json:"items"`; Count int `json:"count"` }
    mustDecode(r2, &list)
    if list.Count != 1 || len(list.Items) != 1 { t.Fatalf("expected 1 item, got count=%d len=%d", list.Count, len(list.Items)) }

    // GET BY ID
    r3, err := http.Get(ts.URL + "/items/" + created.ID)
    if err != nil { t.Fatalf("get failed: %v", err) }
    defer r3.Body.Close()
    if r3.StatusCode != http.StatusOK { t.Fatalf("expected 200, got %d", r3.StatusCode) }
    var got itemRes
    mustDecode(r3, &got)
    if got.ID != created.ID { t.Fatalf("id mismatch: want %s got %s", created.ID, got.ID) }

    // UPDATE
    newName := "Banana"
    done := true
    up := mustJSON(updateReq{Name: &newName, Done: &done})
    req, _ := http.NewRequest(http.MethodPut, ts.URL+"/items/"+created.ID, bytes.NewReader(up))
    req.Header.Set("Content-Type", "application/json")
    r4, err := http.DefaultClient.Do(req)
    if err != nil { t.Fatalf("update failed: %v", err) }
    defer r4.Body.Close()
    if r4.StatusCode != http.StatusOK { t.Fatalf("expected 200, got %d", r4.StatusCode) }
    var updated itemRes
    mustDecode(r4, &updated)
    if !updated.Done || updated.Name != "Banana" { t.Fatalf("update incorrect: %+v", updated) }

    // DELETE
    reqDel, _ := http.NewRequest(http.MethodDelete, ts.URL+"/items/"+created.ID, nil)
    r5, err := http.DefaultClient.Do(reqDel)
    if err != nil { t.Fatalf("delete failed: %v", err) }
    defer r5.Body.Close()
    if r5.StatusCode != http.StatusNoContent { t.Fatalf("expected 204, got %d", r5.StatusCode) }

    // Ensure 404 after delete
    r6, err := http.Get(ts.URL + "/items/" + created.ID)
    if err != nil { t.Fatalf("get after delete failed: %v", err) }
    defer r6.Body.Close()
    if r6.StatusCode != http.StatusNotFound { t.Fatalf("expected 404, got %d", r6.StatusCode) }
}

func TestValidation(t *testing.T) {
    ts := httptest.NewServer(httpapi.NewServer(store.NewMemoryItemStore()))
    defer ts.Close()

    // Create invalid name
    body := mustJSON(createReq{Name: ""})
    r1, err := http.Post(ts.URL+"/items", "application/json", bytes.NewReader(body))
    if err != nil { t.Fatalf("create invalid failed: %v", err) }
    defer r1.Body.Close()
    if r1.StatusCode != http.StatusUnprocessableEntity { t.Fatalf("expected 422, got %d", r1.StatusCode) }

    // Update invalid name
    // First create
    okBody := mustJSON(createReq{Name: "Ok"})
    r2, _ := http.Post(ts.URL+"/items", "application/json", bytes.NewReader(okBody))
    var created itemRes; mustDecode(r2, &created); r2.Body.Close()
    empty := ""
    up := mustJSON(updateReq{Name: &empty})
    req, _ := http.NewRequest(http.MethodPut, ts.URL+"/items/"+created.ID, bytes.NewReader(up))
    req.Header.Set("Content-Type", "application/json")
    r3, err := http.DefaultClient.Do(req)
    if err != nil { t.Fatalf("update invalid failed: %v", err) }
    defer r3.Body.Close()
    if r3.StatusCode != http.StatusUnprocessableEntity { t.Fatalf("expected 422, got %d", r3.StatusCode) }
}

// Helpers
func mustJSON(v any) []byte { b, _ := json.Marshal(v); return b }
func mustDecode(resp *http.Response, v any) { _ = json.NewDecoder(resp.Body).Decode(v) }