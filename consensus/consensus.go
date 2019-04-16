package consensus

import (
	"github.com/bluele/gcache"
	"github.com/qlcchain/go-qlc/common/event"
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/config"
	"github.com/qlcchain/go-qlc/ledger"
	"github.com/qlcchain/go-qlc/ledger/process"
	"github.com/qlcchain/go-qlc/log"
	"go.uber.org/zap"
)

type ConsensusAlgorithm interface {
	Init()
	Start()
	Stop()
	ProcessMsg(msgType MsgType, result process.ProcessResult, bs *BlockSource, p interface{})
}

type Consensus struct {
	ca       ConsensusAlgorithm
	recv     *Receiver
	logger   *zap.SugaredLogger
	cache    gcache.Cache
	ledger   *ledger.Ledger
	verifier *process.LedgerVerifier
}

func NewConsensus(ca ConsensusAlgorithm, cfg *config.Config, eb event.EventBus) *Consensus {
	l := ledger.NewLedger(cfg.LedgerDir(), eb)

	return &Consensus{
		ca:       ca,
		recv:     NewReceiver(eb),
		logger:   log.NewLogger("consensus"),
		cache:    gcache.New(msgCacheSize).LRU().Expiration(msgCacheExpirationTime).Build(),
		ledger:   l,
		verifier: process.NewLedgerVerifier(l),
	}
}

func (c *Consensus) Init() {
	c.ca.Init()
	c.recv.init(c)
}

func (c *Consensus) Start() {
	go c.ca.Start()

	err := c.recv.start()
	if err != nil {
		c.logger.Error("receiver start failed")
	}
}

func (c *Consensus) Stop() {
	err := c.recv.stop()
	if err != nil {
		c.logger.Error("receiver stop failed")
	}

	c.ca.Stop()
}

func (c *Consensus) processed(hash types.Hash) bool {
	return c.cache.Has(hash)
}

func (c *Consensus) processedUpdate(hash types.Hash) {
	err := c.cache.Set(hash, "")
	if err != nil {
		c.logger.Errorf("Set cache error [%s] for block [%s] with publish message", err, hash)
	}
}
