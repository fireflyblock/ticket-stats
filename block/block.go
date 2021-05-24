package block

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/gen"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/fireflyblock/ticket-stats/client"
	logging "github.com/ipfs/go-log/v2"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type blockInfo struct {
	Epoch abi.ChainEpoch
	Time  string
}

type mineBlockNum struct {
	mineNum float64
	winNum  float64
	lock    sync.Mutex
}

type QueryParams struct {
	Start string
	End   string
	Miner string
}

var groupNum abi.ChainEpoch = 5

var wait sync.WaitGroup
var blockLog = logging.Logger("block-log")

func (query *QueryParams) GetTicketStatsByTime() error {
	//totalBlockNum:=halfBlockNum*t
	if err:=client.CreateLotusClient();err!=nil{
		return err
	}
	if err:=client.CreateLotusSignClient();err!=nil{
		return err
	}

	from, to, err := calculateBlock(query.End, query.Start)
	if err != nil {
		blockLog.Errorf("calculateBlock error:%+v", err)
		return err
	}

	return query.exportResult(calculateTicketStats(from, to, query.Miner))
}

func (query *QueryParams) GetTicketStatsByEpoch() error {
	if err:=client.CreateLotusClient();err!=nil{
		return err
	}
	if err:=client.CreateLotusSignClient();err!=nil{
		return err
	}
	from, err := strconv.Atoi(query.End)
	if err != nil {
		return err
	}
	to, err := strconv.Atoi(query.Start)
	if err != nil {
		return err
	}

	return query.exportResult(calculateTicketStats(abi.ChainEpoch(from), abi.ChainEpoch(to), query.Miner))
}

func calculateBlock(start, end string) (abi.ChainEpoch, abi.ChainEpoch, error) {
	tipset, err := client.Client.ChainHead(context.Background())
	if err != nil {
		blockLog.Errorf("calculateBlock get chain head error:%+v", err)
		return 0, 0, err
	}
	t := time.Unix(int64(tipset.MinTimestamp()), 0)
	//h,m,_:=t.Clock()
	//totalNum:=h*120+m*2
	var startTime time.Time
	var endTime time.Time
	startTime, err = time.ParseInLocation("2006-01-02T15:04:05", start, time.Local)
	if err != nil {
		blockLog.Errorf("calculateBlock parse start time error:%+v", err)
		return 0, 0, err
	}
	//startTime = startTime.AddDate(0, 0, 1)

	endTime, err = time.ParseInLocation("2006-01-02T15:04:05", end, time.Local)
	if err != nil {
		blockLog.Errorf("calculateBlock parse end time error:%+v", err)
		return 0, 0, err
	}

	blockLog.Infof("calculateBlock parse start time :%+v end time :%+v", startTime.String(), endTime.String())

	from := tipset.Height() - abi.ChainEpoch(int64(t.Sub(startTime).Minutes())*2)
	to := tipset.Height() - abi.ChainEpoch(int64(t.Sub(endTime).Minutes())*2)
	return from, to, nil
}

func calculateTicketStats(begin, end abi.ChainEpoch, miner string) (*mineBlockNum, []blockInfo, []blockInfo) {
	counter := new(mineBlockNum)
	minerAddr, err := address.NewFromString(miner)
	if err != nil {
		blockLog.Errorf("NewFromString err:", err)
		return counter, nil, nil
	}
	missed := make([]blockInfo, 0)
	mined := make([]blockInfo, 0)
	wins := make([]blockInfo, 0)
	//将begin end分组
	for {
		if begin-end <= groupNum {
			wait.Add(1)
			go do(begin, end, counter, minerAddr, &missed, &mined, &wins)
			break
		}
		wait.Add(1)
		go do(begin, begin-groupNum, counter, minerAddr, &missed, &mined, &wins)
		//是否sleep
		begin -= groupNum
	}
	wait.Wait()
	blockLog.Infof("mined:%+v    missed:%+v   wins:%+v", mined, missed, wins)
	return counter, missed, mined
}

func do(begin, end abi.ChainEpoch, counter *mineBlockNum, miner address.Address, missed, mined, wins *[]blockInfo) {

	defer wait.Done()
	for i := begin; i > end; i-- {
		start := time.Now()
		tipset, err := client.Client.ChainGetTipSetByHeight(context.Background(), i, types.NewTipSetKey())
		if err != nil {
			blockLog.Errorf("calculateBlock get chain head error:%+v", err)
			return
		}
		if tipset.Height() != i {
			continue
		}
		//time
		tStr := time.Unix(int64(tipset.MinTimestamp()), 0).String()
		//flag
		flag := true
		//计算mined
		for _, b := range tipset.Blocks() {
			if b.Miner.String() == miner.String() {
				counter.lock.Lock()
				counter.mineNum++
				counter.lock.Unlock()
				flag = false
				*mined = append(*mined, blockInfo{i, tStr})
			}
		}
		blockLog.Infof("mined:%+v", time.Now().Sub(start))
		end := time.Now()
		//计算winner
		if calculateWiner(i-1, miner) {
			counter.lock.Lock()
			counter.winNum++
			counter.lock.Unlock()
			if flag {
				*missed = append(*missed, blockInfo{i, tStr})
			}
			*wins = append(*wins, blockInfo{i, tStr})
		}
		blockLog.Infof("winer:%+v", time.Now().Sub(end))
	}

}

func calculateWiner(h abi.ChainEpoch, miner address.Address) bool {
	ctx := context.Background()
	round := h + 1
	tp, err := client.Client.ChainGetTipSetByHeight(ctx, h, types.NewTipSetKey())
	if err != nil {
		blockLog.Errorf("ChainGetTipSetByHeight err:%+v", err)
		return false
	}

	mbi, err := client.Client.MinerGetBaseInfo(ctx, miner, round, tp.Key())
	if err != nil {
		blockLog.Errorf("MinerGetBaseInfo err:%+v", err)
		return false
	}

	if mbi == nil {

		return false
	}
	if !mbi.EligibleForMining {
		// slashed or just have no power yet
		return false
	}

	beaconPrev := mbi.PrevBeaconEntry
	bvals := mbi.BeaconEntries

	rbase := beaconPrev
	if len(bvals) > 0 {
		rbase = bvals[len(bvals)-1]
	}

	p, err := gen.IsRoundWinner(ctx, tp, round, miner, rbase, mbi, client.SignClient)
	if err != nil {
		blockLog.Errorf("IsRoundWinner err:%+v", err)
		return false
	}

	if p == nil {
		return false
	}
	return true
}

func (query *QueryParams) exportResult(counter *mineBlockNum, missed, mined []blockInfo) error {
	sort.Slice(missed, func(i, j int) bool {
		return missed[i].Epoch < missed[j].Epoch
	})
	sort.Slice(mined, func(i, j int) bool {
		return mined[i].Epoch < mined[j].Epoch
	})
	if counter.winNum == 0 {
		return nil
	} else {
		fileName := os.Getenv("OUTPUT_FILE")
		data := []byte(fmt.Sprintf("%s Ticket stats:\n", query.Miner))
		data = append(data, []byte(fmt.Sprintf("      Start:%s\n", query.Start))...)
		data = append(data, []byte(fmt.Sprintf("      End:%s\n", query.End))...)
		data = append(data, []byte(fmt.Sprintf("      Mined block number:%d\n", len(mined)))...)
		data = append(data, []byte(fmt.Sprintf("      Missed block number:%d\n", len(missed)))...)
		if counter.winNum != 0 {
			data = append(data, []byte(fmt.Sprintf("      Stats:%.2f%% \n", counter.mineNum/counter.winNum*100))...)
		} else {
			data = append(data, []byte(fmt.Sprintf("      Stats:%.2f%% \n", 100))...)
		}
		data = append(data, []byte(fmt.Sprintf("      Mined block :%+v\n", mined))...)
		data = append(data, []byte(fmt.Sprintf("      Missed block :%+v\n", missed))...)
		if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
			return err
		}
		return nil
	}
}
