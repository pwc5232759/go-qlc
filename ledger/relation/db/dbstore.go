package db

import "io"

type TableName string

const (
	TableBlockHash    TableName = "blockhash"
	TableBlockMessage TableName = "blockmessage"
)

const LikeSign = "_like_"

type Column string

const (
	ColumnHash      Column = "hash"
	ColumnType      Column = "type"
	ColumnAddress   Column = "address"
	ColumnSender    Column = "sender"
	ColumnReceiver  Column = "receiver"
	ColumnMessage   Column = "message"
	ColumnTimestamp Column = "timestamp"
	ColumnNoNeed    Column = ""
)

type DbStore interface {
	io.Closer
	Create(table TableName, condition map[Column]interface{}) error
	Read(table TableName, condition map[Column]interface{}, offset int, limit int, order Column, dest interface{}) error
	Update(table TableName, condition map[Column]interface{}) error
	Delete(table TableName, condition map[Column]interface{}) error
	Count(table TableName, dest interface{}) error
	Group(table TableName, column Column, dest interface{}) error
}
