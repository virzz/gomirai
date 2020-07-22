package message

// Chain 消息链
type Chain struct {
	Msg []Message
}

// GenChain 生成消息链
func GenChain(args ...Message) Chain {
	return Chain{
		Msg: args,
	}
}
