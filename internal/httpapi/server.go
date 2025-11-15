package httpapi

import (
    "encoding/json"
    "be_test/internal/model"
    "be_test/internal/store"
    "net/http"
    "strings"
)

type Server struct {
    store store.ItemStore
}

func NewServer(s store.ItemStore) http.Handler {
    srv := &Server{store: s}
    mux := http.NewServeMux()
    mux.HandleFunc("/health", srv.handleHealth)
    mux.HandleFunc("/items", srv.handleItems)
    mux.HandleFunc("/items/", srv.handleItemByID)
    return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleItems(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        items, err := s.store.List()
        if err != nil {
            writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_server_error"})
            return
        }
        
        listItems := model.ListItemsResponse{
            Items: items,
            Count: len(items),
        }

        writeJSON(w, http.StatusOK, listItems)

    case http.MethodPost:
        var req model.CreateItemRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
            return
        }

        trimmed, ok := nonEmptyTrimmed(&req.Name)
        if !ok {
            writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "name must be non-empty"})
            return
        }

        item, err := s.store.Create(trimmed)
        if err != nil {
            writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_server_error"})
            return
        }

        writeJSON(w, http.StatusCreated, item)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func (s *Server) handleItemByID(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/items/")
    if id == "" || strings.Contains(id, "/") {
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
        return
    }

    switch r.Method {
    case http.MethodGet:
        item, err := s.store.Get(id)
        if err != nil {
            if err == store.ErrNotFound {
                writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
            } else {
                writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_server_error"})
            }
            return
        }
        writeJSON(w, http.StatusOK, item)

    case http.MethodPut:
        var req model.UpdateItemRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
            return
        }

        var namePtr *string
        if req.Name != nil {
            trimmed, ok := nonEmptyTrimmed(req.Name)
            if !ok {
                writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "name must be non-empty"})
                return
            }
            namePtr = &trimmed
        }

        item, err := s.store.Update(id, namePtr, req.Done)
        if err != nil {
            if err == store.ErrNotFound {
                writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
            } else {
                writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_server_error"})
            }
            return
        }
        writeJSON(w, http.StatusOK, item)

    case http.MethodDelete:
        err := s.store.Delete(id)
        if err != nil {
            if err == store.ErrNotFound {
                writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
            } else {
                writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal_server_error"})
            }
            return
        }
        w.WriteHeader(http.StatusNoContent)

    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

// writeJSON helper
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

// nonEmptyTrimmed helper
func nonEmptyTrimmed(s *string) (string, bool) {
    if s == nil {
        return "", false
    }
    trimmed := strings.TrimSpace(*s)
    return trimmed, trimmed != ""
}