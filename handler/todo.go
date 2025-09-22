package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todosPointer, err := h.svc.ReadTODO(ctx,req.PrevID,req.Size)
	if err != nil {
		return nil, err
	}
	todos := make([]model.TODO,len(todosPointer))
	for i, value:= range todosPointer{
		if value!=nil{
			todos[i] = *value
		}
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

//routerによって自動で下記のメソッドが呼び出される
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		values := r.URL.Query()
		prevIDStr := values.Get("prev_id")
		sizeStr := values.Get("size")

		var prevID int64
		var size int64
		var err error

		if prevIDStr != "" {
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				prevID = 0 // 変換失敗時は0など適切な値
			}
		}
		if sizeStr != "" {
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				size = 10 // デフォルト値
			}
		}

		m := model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
		}
		result,err := h.Read(r.Context(),&m)
		if err != nil{
			log.Println(err)
			return 
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil{
			log.Println(err)
			return 
		}
	case http.MethodPost: //POSTリクエストの場合
		var m *model.CreateTODORequest //json形式のものをでコードしたデータが入っている
		err := json.NewDecoder(r.Body).Decode(&m) //r.bodyにはjsonデータが入ってる
		if err != nil {
			log.Println(err)
			return
		}
		if m.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//ここでresultにレスポンス用データを作成&DB操作
		result, err := h.Create(r.Context(), m) // r.Context()メタデータ情報
		if err != nil {
			log.Println(err)
			return
		}
		//ここでレスポンスを返す
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			return
		}

	case http.MethodPut:
		var m *model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}
		if m.ID == 0 || m.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		result, err := h.Update(r.Context(), m)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Println(err)
			return
		}

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}