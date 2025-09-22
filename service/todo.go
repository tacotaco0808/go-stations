package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	//ここでDB接続をもつインスタンスを生成
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
// service ではおもにデータベースデータベースとの処理
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	//DBにデータを入れつつ、TODOmodelをレスポンスで返す
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	var todo model.TODO
	res,err := s.db.ExecContext(ctx,insert,subject,description)
	if err != nil{
		return nil,err
	}
	id,err := res.LastInsertId()
	if err != nil{
		return nil, err
	}
	todo.ID = id
	err = s.db.QueryRowContext(ctx,confirm,id).Scan(
		&todo.Subject,
		&todo.Description,
        &todo.CreatedAt,
        &todo.UpdatedAt,
	)
	if err != nil{
		return  nil,err
	}
	
	return &todo, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)

	if size == 0{
		return []*model.TODO{},nil
	}

	var rows *sql.Rows
	var err error
	if prevID == 0{//prevIDなし
		rows,err = s.db.QueryContext(ctx,read,size)
		
	}else{//prevIDあり
		rows,err = s.db.QueryContext(ctx,readWithID,prevID,size)
		
	}
	if err != nil{
		return nil,err
	}
	defer rows.Close() // rowsの使用が終わったら自動でClose

	
	todos := []*model.TODO{}//参照が配列に保存される。もし値型であると、保存した時点のものがコピーされる
	for rows.Next(){
		var todo model.TODO
		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil{
			return nil,err
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	//DB書き換え処理
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	res,err := s.db.ExecContext(ctx,update,subject,description,id)
	if err != nil {
		return nil,err
	}
	rows,err := res.RowsAffected()
	if err != nil {
		return nil,err
	}
	if rows == 0 {
		return nil,&model.ErrNotFound{}
	}
	//書き換え後のデータを整理する処理
	var todo model.TODO
	todo.ID = id
	err = s.db.QueryRowContext(ctx,confirm,id).Scan(
		&todo.Subject,
		&todo.Description,
        &todo.CreatedAt,
        &todo.UpdatedAt,
	)
	if err != nil {
		return nil,err
	}
	return &todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
