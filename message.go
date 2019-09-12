package paxos

type msgType int

const (
	Prepare msgType = iota + 1 // Send from proposer -> acceptor
	Promise                    // Send from acceptor -> proposer
	Propose                    // Send from proposer -> acceptor
	Accept                     // Send from acceptor -> learner
)

// 消息类型
type message struct {
	from   int
	to     int
	typ    msgType
	seq    int
	preSeq int // 本条消息之前的序号 。 preSeq = seq - 1。
	val    string
}

func (m *message) getProposeVal() string {
	return m.val
}

func (m *message) getProposeSeq() int {
	return m.seq
}
