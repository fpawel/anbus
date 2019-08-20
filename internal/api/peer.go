package api

type PeerSvc struct {
	P PeerNotifier
}

func NewPeerSvc(p PeerNotifier) *PeerSvc {
	return &PeerSvc{p}
}

type PeerNotifier interface {
	OnStarted()
	OnClosed()
}

func (x *PeerSvc) Init(_ struct{}, _ *struct{}) error {
	x.P.OnStarted()
	return nil
}

func (x *PeerSvc) Close(_ struct{}, _ *struct{}) error {
	x.P.OnClosed()
	return nil
}
