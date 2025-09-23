package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
        read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
        readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC`
    )

    var rows *sql.Rows
    var err error
    if size == 0 {
        // LIMIT句なしで全件取得
        if prevID == 0 {
            rows, err = s.db.QueryContext(ctx, read)
        } else {
            rows, err = s.db.QueryContext(ctx, readWithID, prevID)
        }
    } else {
        // LIMIT句あり
        if prevID == 0 {
            rows, err = s.db.QueryContext(ctx, read+" LIMIT ?", size)
        } else {
            rows, err = s.db.QueryContext(ctx, readWithID+" LIMIT ?", prevID, size)
        }
    }
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    todos := []*model.TODO{}
    for rows.Next() {
        var todo model.TODO
        err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
        if err != nil {
            return nil, err
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
	const deleteFmt = `DELETE FROM todos WHERE id IN (%s)`//?%sで可変長?を複数並べる
	if len(ids) == 0{
		return nil
	}
	placeholders := strings.Repeat("?,",len(ids))
	placeholders = placeholders[:len(placeholders)-1]//最後の文字より一個前まで代入
	query := fmt.Sprintf(deleteFmt,placeholders)//元のdeleteFmtクエリの%sに"?,?,?"を代入

	args := make([]interface{},len(ids))//[]int64{1,2,3} -→ []interface{1,2,3}の変換
	for i,id := range ids{
		args[i] = id
	}

	row,err := s.db.ExecContext(ctx,query,args...)//args...で[]interface{1,2,3}の1,2,3の値が展開される
	if err != nil{
		return err
	}
	affected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0{
		return &model.ErrNotFound{}
	}
	return nil
}
