package kvsrv

import "github.com/google/uuid"

// Put or Append
type PutAppendArgs struct {
	Key     string
	Value   string
	ClerkID uuid.UUID
	Idx     int
}

type PutAppendReply struct {
	Value string
}

type GetArgs struct {
	Key     string
	ClerkID uuid.UUID
	Idx     int
}

type GetReply struct {
	Value string
}
