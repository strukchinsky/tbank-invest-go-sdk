package investgo

import (
	"context"
	"fmt"
	"io"
	"slices"
	"sync"

	pb "github.com/floatdrop/tbank-invest-go-sdk"

	"google.golang.org/grpc"
)

// MarketDataStream provides interface for managing subscriptions through
// one stream instance. Channels with data will receive market data from single underlying stream.
// Subscription will be automatically cancelled when no more channels are listening for it.
type MarketDataStream struct {
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	logger Logger
	stream pb.MarketDataStreamService_MarketDataStreamClient

	candleSubscriptions map[string][]*CandleSubscription
}

// CandleSubscription represents handler for candles stream.
type CandleSubscription struct {
	// InstrumentId stores unique instrument id that is used in candleSubscriptions as key
	InstrumentId string

	// Private channel that receives data from stream
	ch chan *pb.Candle
}

// Recv returns read only channel with candles market data
func (cs *CandleSubscription) Recv() <-chan *pb.Candle {
	return cs.ch
}

// NewMarketDataStream returns a new [MarketDataStream] to create subscriptions.
func (c *Client) NewMarketDataStream(ctx context.Context, logger Logger) (*MarketDataStream, error) {
	ctx, cancel := context.WithCancel(ctx)

	stream, err := pb.NewMarketDataStreamServiceClient(c.conn).MarketDataStream(ctx, grpc.EmptyCallOption{})
	if err != nil {
		cancel()
		return nil, err
	}

	mds := &MarketDataStream{
		mu:     sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
		logger: logger,
		stream: stream,

		candleSubscriptions: make(map[string][]*CandleSubscription),
	}

	go mds.listen()

	return mds, nil
}

func (m *MarketDataStream) listen() {
	defer m.shutdown()

	m.logger.Debug("Starting to listen for MarketDataStream events")
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			resp, err := m.stream.Recv()

			if err == io.EOF {
				return
			}

			if err != nil {
				m.logger.Error("Failed to Recv from MarketDataStream", "Error", err)
				return
			}

			m.consume(resp)
		}
	}
}

func (m *MarketDataStream) consume(resp *pb.MarketDataResponse) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Debug("Consuming response", "MarketDataResponse", resp)

	if cs := resp.GetSubscribeCandlesResponse(); cs != nil {
		for _, c := range cs.CandlesSubscriptions {
			if c.SubscriptionStatus != pb.SubscriptionStatus_SUBSCRIPTION_STATUS_SUCCESS {
				m.logger.Warn("SubscriptionStatus is not SUBSCRIPTION_STATUS_SUCCESS: closing subscriptions", "Id", []string{c.Figi, c.InstrumentUid})
				m.closeCandleSubscription(c.Figi)
				m.closeCandleSubscription(c.InstrumentUid)
			}
		}
	}

	if c := resp.GetCandle(); c != nil {
		m.notifyCandle(c)
	}
}

// shutdown closes all channels - should be called only from listen (writing) function
func (m *MarketDataStream) shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k := range m.candleSubscriptions {
		m.closeCandleSubscription(k)
	}
}

func (m *MarketDataStream) notifyCandle(c *pb.Candle) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, s := range m.candleSubscriptions[c.Figi] {
		s.ch <- c
	}

	for _, s := range m.candleSubscriptions[c.InstrumentUid] {
		s.ch <- c
	}
}

// SubscribeCandle creates subscriptions for candle with specified id. If subscription already
// exists – only new channel will be created.
func (m *MarketDataStream) SubscribeCandle(id string, waitingClose bool) (*CandleSubscription, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	select {
	case <-m.ctx.Done():
		return nil, fmt.Errorf("MarketDataStream is closed")
	default:
		s, newSubscription := m.createCandleSubscription(id)

		if !newSubscription {
			return s, nil
		}

		m.logger.Debug("Subscribing to candle", "Candle", id)
		return s, m.sendCandlesRequest([]string{id}, pb.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE, pb.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE, waitingClose)
	}
}

// createCandleSubscription returns new CandleSubscription with boolean flag, that indicates
// that subscription was already in map (no need to call stream.Send for it).
func (m *MarketDataStream) createCandleSubscription(id string) (*CandleSubscription, bool) {
	subs, ok := m.candleSubscriptions[id]

	ch := make(chan *pb.Candle)
	s := &CandleSubscription{InstrumentId: id, ch: ch}
	m.candleSubscriptions[id] = append(subs, s)

	return s, !ok
}

// UnsubscribeCandle closes CandleSubscription and removes subscription from subscriptions map.
// If no more subscription are present with same id – subscription will be cancelled.
func (m *MarketDataStream) UnsubscribeCandle(s *CandleSubscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.candleSubscriptions[s.InstrumentId]

	if !ok {
		return fmt.Errorf("CandleSubscription with id %s not active", s.InstrumentId)
	}

	close(s.ch)

	m.candleSubscriptions[s.InstrumentId] = slices.DeleteFunc(m.candleSubscriptions[s.InstrumentId], func(e *CandleSubscription) bool {
		return e == s
	})

	if len(m.candleSubscriptions[s.InstrumentId]) == 0 {
		m.sendCandlesRequest([]string{s.InstrumentId}, pb.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE, pb.SubscriptionAction_SUBSCRIPTION_ACTION_UNSUBSCRIBE, true)
	}

	return nil
}

func (m *MarketDataStream) closeCandleSubscription(id string) {
	for _, sub := range m.candleSubscriptions[id] {
		close(sub.ch)
	}
	m.sendCandlesRequest([]string{id}, pb.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE, pb.SubscriptionAction_SUBSCRIPTION_ACTION_UNSUBSCRIBE, true)
	delete(m.candleSubscriptions, id)
}

// Close cancels the MarketDataStream context, which will stop listening on stream
func (m *MarketDataStream) Close() {
	m.cancel()
}

func (m *MarketDataStream) sendCandlesRequest(ids []string, interval pb.SubscriptionInterval, act pb.SubscriptionAction, waitingClose bool) error {
	instruments := make([]*pb.CandleInstrument, 0, len(ids))
	for _, id := range ids {
		instruments = append(instruments, &pb.CandleInstrument{
			InstrumentId: id,
			Interval:     interval,
		})
	}

	return m.stream.Send(&pb.MarketDataRequest{
		Payload: &pb.MarketDataRequest_SubscribeCandlesRequest{
			SubscribeCandlesRequest: &pb.SubscribeCandlesRequest{
				SubscriptionAction: act,
				Instruments:        instruments,
				WaitingClose:       waitingClose,
			},
		},
	})
}
