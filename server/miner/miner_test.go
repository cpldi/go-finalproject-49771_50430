package miner

import (
	"bitcoin_miner/hash"
	"bitcoin_miner/message"
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type Randomizer string

var randomString Randomizer = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "abcdefghijklmnopqrstuvwxyz"

func (r Randomizer) rand() string{
	length := rand.Uint32() % 100
	buf := make([]byte, length)
	var i uint32
	for; i < length; i++ {
		buf[i] = randomString[rand.Intn(len(randomString))]
	}

	return string(buf)
}

func TestRandMiner(t *testing.T) {

	//test setup
	rand.Seed(time.Now().UnixNano())
	var size uint64 = 6000000
	nTests := 30
	if testing.Short() { // Use -short for short test
		size = 200
		nTests = 5
	}
	ctx := context.Background()
	m := NewMiner(ctx,50,50,5,5)

	//test start here
	//It runs nTests in parallel
	for i:= 0 ; i < nTests ; i++ {
		t.Run(strconv.Itoa(i),func(t *testing.T) {
			t.Parallel()
			msg := message.NewRequest(randomString.rand(), 0, rand.Uint64()%size)
			result := m.SubmitJob(msg)
			expect := serialResult(msg)
			if *result != *expect {
				t.Fatalf("Expecting %v got %v\n", expect, result)
			}
		})
	}
}


func serialResult(req *message.Message) *message.Message{
	i := req.Lower + 1
	maxHash := hash.Hash(req.Data, req.Lower)
	maxNouce := req.Lower

	for  ; i <= req.Upper;i++ {
		h := hash.Hash(req.Data,i)
		if maxHash < h {
			maxHash = h
			maxNouce = i
		}
	}

	return message.NewResult(maxHash,maxNouce, req.Lower, req.Upper)

}