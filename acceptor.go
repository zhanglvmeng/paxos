package paxos

import "log"

//Create a accetor and also assign learning IDs into acceptor.
//Acceptor: Will response request from proposer, promise the first and largest seq number propose.
//          After proposer reach the majority promise.  Acceptor will pass the proposal value to learner to confirn and choose.
func NewAcceptor(id int, nt nodeNetwork, learners ...int) acceptor {
	newAccptor := acceptor{id: id, nt: nt}
	newAccptor.learners = learners
	return newAccptor
}

type acceptor struct {
	id         int
	learners   []int
	acceptMsg  message  // accept 消息
	promiseMsg message  // 对 prepare 消息的响应
	nt         nodeNetwork
}

//Acceptor process detail logic.
func (a *acceptor) run() {
	// acceptor 一直在等待。。。
	for {
		//	log.Println("acceptor:", a.id, " wait to recev msg")
		m := a.nt.recev()
		if m == nil {
			continue
		}

		//	log.Println("acceptor:", a.id, " recev message ", *m)
		switch m.typ {
		case Prepare:
			// 接收到prepare消息
			promiseMsg := a.recevPrepare(*m)
			// 发送promise 响应消息
			a.nt.send(*promiseMsg)
			continue
		case Propose:
			// 接收到propose消息
			accepted := a.recevPropose(*m)
			if accepted {
				// 如果接收了Propose消息，则给learners发Accept消息。
				for _, lId := range a.learners {
					m.from = a.id
					m.to = lId
					m.typ = Accept
					a.nt.send(*m)
				}
			}
		default:
			log.Fatalln("Unsupport message in accpetor ID:", a.id)
		}
	}
	log.Println("accetor :", a.id, " leave.")
}

//After acceptor receive prepare message.
//It will check  prepare number and return acceptor if it is bigest one.
func (a *acceptor) recevPrepare(prepare message) *message {
	// 如果已经节后到更大编号的消息，则直接退出。
	if a.promiseMsg.getProposeSeq() >= prepare.getProposeSeq() {
		log.Println("ID:", a.id, "Already accept bigger one")
		return nil
	}
	// 当前的消息时最大编号的消息，则组装promise 响应消息，
	log.Println("ID:", a.id, " Promise")
	prepare.to = prepare.from
	prepare.from = a.id
	prepare.typ = Promise
	a.acceptMsg = prepare
	return &prepare
}

//Recev Propose only check if acceptor already accept bigger propose before.
//Otherwise, will just forward this message out and change its type to "Accept" to learning later.
func (a *acceptor) recevPropose(proposeMsg message) bool {
	//Already accept message is identical with previous promise message
	log.Println("accept:check propose. ", a.acceptMsg.getProposeSeq(), proposeMsg.getProposeSeq())
	// 如果acceptor  接收到的消息编号跟 发过来的promose消息编号不一致，则不处理。
	if a.acceptMsg.getProposeSeq() > proposeMsg.getProposeSeq() || a.acceptMsg.getProposeSeq() < proposeMsg.getProposeSeq() {
		log.Println("ID:", a.id, " acceptor not take propose:", proposeMsg.val)
		return false
	}
	log.Println("ID:", a.id, " Accept")
	return true
}
