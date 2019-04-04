package p2p

import (
	"math/rand"
	"sync"
	"time"

	libnet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
)

// StreamManager manages all streams
type StreamManager struct {
	mu         sync.Mutex
	allStreams *sync.Map
	node       *QlcNode
}

// NewStreamManager return a new stream manager
func NewStreamManager() *StreamManager {
	return &StreamManager{
		allStreams: new(sync.Map),
	}
}

// SetQlcService set netService
func (sm *StreamManager) SetQlcNode(node *QlcNode) {
	sm.node = node
}

// Add a new stream into the stream manager
func (sm *StreamManager) Add(s libnet.Stream) {
	stream := NewStream(s, sm.node)
	sm.AddStream(stream)
}

// AddStream into the stream manager
func (sm *StreamManager) AddStream(stream *Stream) {

	//sm.mu.Lock()
	//defer sm.mu.Unlock()

	// check & close old stream
	if v, ok := sm.allStreams.Load(stream.pid.Pretty()); ok {
		old, _ := v.(*Stream)

		sm.node.logger.Info("Removing old stream.")

		sm.allStreams.Delete(old.pid.Pretty())

		if old.stream != nil {
			old.stream.Close()
		}
	}

	sm.node.logger.Infof("Added a new stream:[%s]", stream.pid.Pretty())

	sm.allStreams.Store(stream.pid.Pretty(), stream)
	stream.StartLoop()
}

// RemoveStream from the stream manager
func (sm *StreamManager) RemoveStream(s *Stream) {
	if v, ok := sm.allStreams.Load(s.pid.Pretty()); ok {
		exist, _ := v.(*Stream)
		if s != exist {
			return
		}
		sm.node.logger.Errorf("Removing a stream:[%s]", s.pid.Pretty())
		sm.allStreams.Delete(s.pid.Pretty())
	}
}

// FindByPeerID find the stream with the given peerID
func (sm *StreamManager) FindByPeerID(peerID string) *Stream {
	if v, ok := sm.allStreams.Load(peerID); ok {
		return v.(*Stream)
	}
	return nil
}

func (sm *StreamManager) RandomPeer() (string, error) {
	allPeers := make(PeersSlice, 0)

	sm.allStreams.Range(func(key, value interface{}) bool {
		stream := value.(*Stream)
		if stream.IsConnected() {
			allPeers = append(allPeers, value)
		}
		return true
	})
	var peerID string
	rand.Seed(time.Now().Unix())
	if (len(allPeers)) == 0 {
		return "", ErrNoStream
	}
	randNum := rand.Intn(len(allPeers))
	for i, v := range allPeers {
		stream := v.(*Stream)
		if i == randNum {
			peerID = stream.pid.Pretty()
		}
	}
	return peerID, nil
}

// CloseStream with the given pid and reason
func (sm *StreamManager) CloseStream(peerID string) {
	stream := sm.FindByPeerID(peerID)
	if stream != nil {
		stream.close()
	}
}

// CreateStreamWithPeer create stream with a peer.
func (sm *StreamManager) createStreamWithPeer(pid peer.ID) {

	stream := sm.FindByPeerID(pid.Pretty())

	if stream == nil {
		stream = NewStreamFromPID(pid, sm.node)
		sm.AddStream(stream)
	}
}

// BroadcastMessage broadcast the message
func (sm *StreamManager) BroadcastMessage(messageName string, v interface{}) error {
	//	var cs []*cacheValue
	//	var c *cacheValue
	messageContent, err := marshalMessage(messageName, v)
	if err != nil {
		sm.node.logger.Error(err)
		return err
	}
	version := p2pVersion
	message := NewQlcMessage(messageContent, byte(version), messageName)
	err = sm.node.publisher.Publish(MsgTopic, message)
	if err != nil {
		return err
	}
	return nil
}

func (sm *StreamManager) PeerCounts() int {
	allPeers := make(PeersSlice, 0)
	sm.allStreams.Range(func(key, value interface{}) bool {
		stream := value.(*Stream)
		if stream.IsConnected() {
			allPeers = append(allPeers, value)
		}
		return true
	})

	return len(allPeers)
}
