package kvraft

import (
	"6.824/labrpc"
	"crypto/rand"
	"math/big"
)

type Clerk struct {
	servers   []*labrpc.ClientEnd
	leaderId  int64
	clientId  int64 // generated by nrand(), it would be better to use some distributed ID generation algorithm that guarantees no conflicts
	commandId int64 // (clientId, commandId) defines a operation uniquely
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	return bigx.Int64()
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	return &Clerk{
		servers:   servers,
		leaderId:  0,
		clientId:  nrand(),
		commandId: 0,
	}
}

func (ck *Clerk) Get(key string) string {
	return ck.Command(&CommandArgs{Key: key, Op: OpGet})
}

func (ck *Clerk) Put(key string, value string) {
	ck.Command(&CommandArgs{Key: key, Value: value, Op: OpPut})
}
func (ck *Clerk) Append(key string, value string) {
	ck.Command(&CommandArgs{Key: key, Value: value, Op: OpAppend})
}

// you can send an RPC with code like this:
// ok := ck.servers[i].Call("KVServer.Command", &request, &response)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Command(request *CommandArgs) string {
	request.ClientId, request.CommandId = ck.clientId, ck.commandId
	for {
		var response CommandReply
		if !ck.servers[ck.leaderId].Call("KVServer.Command", request, &response) || response.Err == ErrWrongLeader || response.Err == ErrTimeout {
			ck.leaderId = (ck.leaderId + 1) % int64(len(ck.servers))
			continue
		}
		ck.commandId++
		return response.Value
	}
}
