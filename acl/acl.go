package acl

import (
	"context"
	"time"

	commonaccount "github.com/anyproto/any-sync/accountservice"
	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/any-sync/app/logger"
	"github.com/anyproto/any-sync/app/ocache"
	"github.com/anyproto/any-sync/commonspace/object/acl/list"
	"github.com/anyproto/any-sync/consensus/consensusclient"
	"github.com/anyproto/any-sync/consensus/consensusproto"
	"github.com/anyproto/any-sync/metric"
	"github.com/prometheus/client_golang/prometheus"
)

const CName = "coordinator.acl"

var log = logger.NewNamed(CName)

func New() Acl {
	return &aclService{}
}

type Acl interface {
	AddRecord(ctx context.Context, spaceId string, rec *consensusproto.RawRecord) (result *consensusproto.RawRecordWithId, err error)
	RecordsAfter(ctx context.Context, spaceId, aclHead string) (result []*consensusproto.RawRecordWithId, err error)
	app.ComponentRunnable
}

type aclService struct {
	consService    consensusclient.Service
	cache          ocache.OCache
	accountService commonaccount.Service
}

func (as *aclService) Init(a *app.App) (err error) {
	as.consService = app.MustComponent[consensusclient.Service](a)
	as.accountService = app.MustComponent[commonaccount.Service](a)

	var metricReg *prometheus.Registry
	if m := a.Component(metric.CName); m != nil {
		metricReg = m.(metric.Metric).Registry()
	}
	as.cache = ocache.New(as.loadObject,
		ocache.WithTTL(5*time.Minute),
		ocache.WithLogger(log.Sugar()),
		ocache.WithPrometheus(metricReg, "coordinator", "acl"),
	)
	return
}

func (as *aclService) Name() (name string) {
	return CName
}

func (as *aclService) loadObject(ctx context.Context, id string) (ocache.Object, error) {
	return as.newAclObject(ctx, id)
}

func (as *aclService) get(ctx context.Context, spaceId string) (list.AclList, error) {
	obj, err := as.cache.Get(ctx, spaceId)
	if err != nil {
		return nil, err
	}
	aObj := obj.(*aclObject)
	aObj.lastUsage.Store(time.Now())
	return aObj.AclList, nil
}

func (as *aclService) AddRecord(ctx context.Context, spaceId string, rec *consensusproto.RawRecord) (result *consensusproto.RawRecordWithId, err error) {
	acl, err := as.get(ctx, spaceId)
	if err != nil {
		return nil, err
	}
	acl.RLock()
	defer acl.RUnlock()
	err = acl.ValidateRawRecord(rec)
	if err != nil {
		return
	}

	return as.consService.AddRecord(ctx, spaceId, rec)
}

func (as *aclService) RecordsAfter(ctx context.Context, spaceId, aclHead string) (result []*consensusproto.RawRecordWithId, err error) {
	acl, err := as.get(ctx, spaceId)
	if err != nil {
		return nil, err
	}
	acl.RLock()
	defer acl.RUnlock()
	return acl.RecordsAfter(ctx, aclHead)
}

func (as *aclService) Run(ctx context.Context) (err error) {
	return
}

func (as *aclService) Close(ctx context.Context) (err error) {
	return as.cache.Close()
}
