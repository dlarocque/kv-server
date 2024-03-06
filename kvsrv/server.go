package kvsrv

import (
	"log"
	"sync"

	"github.com/google/uuid"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type Operation struct {
	idx          int // Monotonically increasing
	preAppendIdx int
}

type KVServer struct {
	mu sync.Mutex

	db              map[string]string
	clerkOperations map[uuid.UUID]Operation
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	reply.Value = kv.db[args.Key]
	kv.mu.Unlock()
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	prev, ok := kv.clerkOperations[args.ClerkID]
	if ok && args.Idx == prev.idx {
		reply.Value = ""
	} else {
		reply.Value = kv.db[args.Key]
		kv.db[args.Key] = args.Value
		kv.clerkOperations[args.ClerkID] = Operation{
			idx:          args.Idx,
			preAppendIdx: -1,
		}
	}
	kv.mu.Unlock()
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	prev, ok := kv.clerkOperations[args.ClerkID]
	if ok && prev.idx == args.Idx {
		reply.Value = kv.db[args.Key][:prev.preAppendIdx]
	} else {
		reply.Value = kv.db[args.Key]
		preAppendIdx := len(reply.Value)
		kv.db[args.Key] += args.Value
		kv.clerkOperations[args.ClerkID] = Operation{
			idx:          args.Idx,
			preAppendIdx: preAppendIdx,
		}
	}
	kv.mu.Unlock()
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.db = make(map[string]string)
	kv.clerkOperations = make(map[uuid.UUID]Operation)
	return kv
}
