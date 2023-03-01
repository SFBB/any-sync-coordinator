package spacestatus

import "github.com/anytypeio/any-sync/commonspace/object/tree/treechangeproto"

type ChangeVerifier interface {
	Verify(rawDelete *treechangeproto.RawTreeChangeWithId, identity []byte, peerId string) (err error)
}

type changeVerifier struct {
}

func (c *changeVerifier) Verify(rawDelete *treechangeproto.RawTreeChangeWithId, identity []byte, peerId string) (err error) {
	return
}
