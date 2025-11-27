package geecache

import pb "geecache/geecache/geecachepb"

// PeerPicker select node by key, peer owns a specific key
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter get val of a key from group
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
