package nodeservice

import (
	"context"
	"github.com/anytypeio/any-sync-node/nodesync/nodesyncproto"
	"github.com/anytypeio/any-sync/app"
	"github.com/anytypeio/any-sync/commonspace/object/tree/treechangeproto"
	"github.com/anytypeio/any-sync/net/pool"
	"github.com/anytypeio/any-sync/nodeconf"
)

const CName = "coordinator.nodeservice"

type NodeService interface {
	app.ComponentRunnable
	Delete(ctx context.Context, spaceId string, raw *treechangeproto.RawTreeChangeWithId) (err error)
}

func New() NodeService {
	return &nodeService{}
}

type nodeService struct {
	pool  pool.Pool
	nodes nodeconf.Service
}

func (n *nodeService) Init(a *app.App) (err error) {
	n.pool = a.MustComponent(pool.CName).(pool.Pool)
	n.nodes = a.MustComponent(nodeconf.CName).(nodeconf.Service)
	return nil
}

func (n *nodeService) Name() (name string) {
	return CName
}

func (n *nodeService) Run(ctx context.Context) (err error) {
	return
}

func (n *nodeService) Close(ctx context.Context) (err error) {
	return
}

func (n *nodeService) Delete(ctx context.Context, spaceId string, raw *treechangeproto.RawTreeChangeWithId) (err error) {
	nodeIds := n.nodes.NodeIds(spaceId)
	respPeer, err := n.pool.GetOneOf(ctx, nodeIds)
	if err != nil {
		return
	}
	nodeClient := nodesyncproto.NewDRPCNodeSyncClient(respPeer)
	_, err = nodeClient.SpaceDelete(ctx, &nodesyncproto.SpaceDeleteRequest{
		SpaceId:      spaceId,
		ChangeId:     raw.Id,
		DeleteChange: raw.RawChange,
	})
	return
}
