package miner

import (
	"bitcoin_miner/hash"
	msg "bitcoin_miner/message"
	"context"
)

const (
	THRESHOLD_BLOCK_SIZE  = 100000
	THRESHOLD_LIGHT_HEAVY = 50
	BUFFERED_CHANNEL_SIZE = 10
)

type Miner struct {
	inL    chan workReq
	inH    chan workReq
	cancel context.CancelFunc
}

type workReq struct {
	Data   string
	Lower  uint64
	Upper  uint64
	workCh chan workResp
}

type workResp struct {
	Hash  uint64
	Nonce uint64
}

func (r workResp) max(other workResp) workResp {
	if r.Hash > other.Hash {
		return r
	}
	return other
}

func NewMiner(capL, capH, nWorkerL, nWorkerH int) *Miner {
	ctxt := context.Background()
	ctxt, cancel := context.WithCancel(ctxt)
	inLight := make(chan workReq, capL)
	inHeavy := make(chan workReq, capH)

	for i := 0; i < nWorkerL; i++ {
		go Worker(ctxt, i, inLight, inHeavy)
	}
	for i := 0; i < nWorkerH; i++ {
		go Worker(ctxt, i, inHeavy, inLight)
	}
	return &Miner{inLight, inHeavy, cancel}
}

func Worker(ctxt context.Context, id int, inDefault, inOther <-chan workReq) {
	for {
		select {
		case <-ctxt.Done():
			return
		case req := <-inDefault:
			computeHigherHash(req,id)
		default:
			{
				select {
				case req := <-inOther:
					computeHigherHash(req,id)
				case req := <-inDefault:
					computeHigherHash(req,id)
				}
			}
		}
	}
}

func computeHigherHash(req workReq,id int) {
	var max uint64
	var maxi uint64
	for i := req.Lower; i <= req.Upper; i++ {
		aux := hash.Hash(req.Data, i)
		if aux > max {
			max = aux
			maxi = i
		}
	}
	//fmt.Printf("[%v] [%v] computed between %v and %v\n",id, time.Now() , req.Lower, req.Upper)
	req.workCh <- workResp{max, maxi}
}

func (m *Miner) Cancel() {
	m.cancel()
}

func min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

func submitBlocks(blocks uint64, in chan<- workReq, request *msg.Message, respCh chan workResp) {
	var i uint64

	//fmt.Printf("wants to submit %v blocks.\n", blocks+1)
	for i = 0; i <= blocks; i++ {
		in <- workReq{
			request.Data,
			request.Lower + THRESHOLD_BLOCK_SIZE*i,
			min(request.Lower+THRESHOLD_BLOCK_SIZE*(i+1)-1, request.Upper),
			respCh}
		//fmt.Printf("submitted block %v.\n", i+1)
	}
	//fmt.Printf("submitted %v blocks.\n", blocks+1)
}

func (m *Miner) SubmitJob(request *msg.Message) *msg.Message {

	blocks := (request.Upper - request.Lower) / THRESHOLD_BLOCK_SIZE
	respCh := make(chan workResp, BUFFERED_CHANNEL_SIZE)


	if blocks < THRESHOLD_LIGHT_HEAVY {
		go submitBlocks(blocks, m.inL, request, respCh)
	} else {
		go submitBlocks(blocks, m.inH, request, respCh)
	}

	max := <-respCh
	var i uint64
	for i = 0; i < blocks; i++ {
		max = max.max(<-respCh)
		//fmt.Printf("[%v] got %v/%v blocks.\n", time.Now() , i+2, blocks+1)
	}
	//fmt.Printf("received all blocks\n\n")

	return msg.NewResult(max.Hash, max.Nonce, request.Lower, request.Upper)
}
