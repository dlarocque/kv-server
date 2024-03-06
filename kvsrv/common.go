package kvsrv

import "github.com/google/uuid"

// Put or Append
type PutAppendArgs struct {
	Key   string
	Value string
	Id    uuid.UUID
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key string
	Id  uuid.UUID
}

type GetReply struct {
	Value string
}
