package dht

// SendJoin sends a request for a node to join the network (even if it is not this exact node)
// - This allows nodes to forward a message, and not care about the response
func SendJoin(addr DhtAddress, node ChordNode, hops int) (token string) {
	msg := EncodeWireMessage(ChordWireMessage{"join", node, hops})

	resp, _ := sendMessage(addr, msg)

	r := JoinedResponse{}

	DecodeStruct([]byte(resp), &r)

	return "whatever"
}

// Join a chord table
// - It returns a string of node ids
func (self *ChordTable) RequestJoin() {
	if len(self.seeds) > 0 {
		addr := self.seeds[0]
		nonce := self.randomizer.GetToken()
		_ = nonce
		SendJoin(addr, self.getNode(), 50)
	}
}

// handleJoin tries to find a place in the request chain for a node
func (self *ChordTable) handleJoin(joining ChordWireMessage) (JoinedResponse, error) {

	// You know... we gotta stop sometime
	if joining.Hops < 0 {
		self.Println("HOP LIMIT EXCEEDED")
		return JoinedResponse{}, nil
	}

	switch {
	case self.NeedsBootstrap():
		self.Prev = &joining.Source
		self.Next = &joining.Source
		self.Println("%d <-> %d", self.Id, self.Next.Id)
		SendJoin(joining.Source.GetAddress(), self.GetNode(), joining.Hops-1)
		return JoinedResponse{}, nil

	case self.Next.Id == joining.Source.Id:
		return JoinedResponse{}, nil

	case self.IsSuccessor(joining.Source):
		self.Next = &joining.Source
		self.Println("%d |-> %d", self.Id, self.Next.Id)
		SendJoin(joining.Source.GetAddress(), self.GetNode(), joining.Hops-1)
		SendJoin(self.Next.GetAddress(), self.GetNode(), joining.Hops-1)

	case self.IsPredecessor(joining.Source):
		self.Prev = &joining.Source
		self.Println("%d |<- %d", self.Id, self.Next.Id)

	default:
		self.Println("%d ?~ %d", self.Id, joining.Source.Id)
		SendJoin(self.Next.GetAddress(), joining.Source, joining.Hops-1)
	}

	// XXX: Start figuring out ring structure
	return JoinedResponse{}, nil
}
