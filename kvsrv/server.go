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

type KVServer struct {
	mu sync.Mutex

	db            map[string]string
	getReqs       map[uuid.UUID]GetReply // Table storing request IDs we've received
	putAppendReqs map[uuid.UUID]PutAppendReply
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	kv.mu.Lock()
	prevReply, ok := kv.getReqs[args.Id]
	if ok {
		reply.Value = prevReply.Value
	} else {
		reply.Value = kv.db[args.Key]
		kv.getReqs[args.Id] = *reply
	}
	kv.mu.Unlock()
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	prevReply, ok := kv.putAppendReqs[args.Id]
	if ok {
		reply.Value = prevReply.Value
	} else {
		reply.Value = kv.db[args.Key]
		kv.db[args.Key] = args.Value
		kv.putAppendReqs[args.Id] = *reply
	}
	kv.mu.Unlock()
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	kv.mu.Lock()
	prevReply, ok := kv.putAppendReqs[args.Id]
	if ok {
		reply.Value = prevReply.Value
	} else {
		reply.Value = kv.db[args.Key]
		kv.db[args.Key] += args.Value
		kv.putAppendReqs[args.Id] = *reply
	}
	kv.mu.Unlock()
}

func StartKVServer() *KVServer {
	kv := new(KVServer)
	kv.db = make(map[string]string)
	kv.getReqs = make(map[uuid.UUID]GetReply)
	kv.putAppendReqs = make(map[uuid.UUID]PutAppendReply)
	return kv
}
