package spacestatus

import (
	"context"
	"github.com/anytypeio/any-sync/commonspace/object/tree/treechangeproto"
	"github.com/anytypeio/any-sync/util/periodicsync"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type DelSender interface {
	Delete(ctx context.Context, spaceId string, raw *treechangeproto.RawTreeChangeWithId) (err error)
}

type pendingSpacesQuery struct {
	deletionPeriod time.Duration
}

func (d pendingSpacesQuery) toMap() bson.M {
	return bson.M{"$and": bson.A{
		bson.D{{"status", SpaceStatusDeletionPending}},
		bson.D{{"deletion_date", bson.M{
			"$gt":  time.Time{},
			"$lte": time.Now().Add(-d.deletionPeriod)}}}}}
}

type statusEntry struct {
	SpaceId         string    `bson:"_id"`
	Identity        []byte    `bson:"identity"`
	DeletionPayload []byte    `bson:"deletionChange"`
	DeletionDate    time.Time `bson:"deletionDate"`
}

type spaceDeleter struct {
	spaces         *mongo.Collection
	runSeconds     int
	deletionPeriod time.Duration
	loop           periodicsync.PeriodicSync
	delSender      DelSender
}

func newSpaceDeleter(runSeconds int, deletionPeriod time.Duration) *spaceDeleter {
	return &spaceDeleter{
		deletionPeriod: deletionPeriod,
		runSeconds:     runSeconds,
	}
}

func (s *spaceDeleter) Run(spaces *mongo.Collection, delSender DelSender) {
	s.delSender = delSender
	s.spaces = spaces
	s.loop = periodicsync.NewPeriodicSync(s.runSeconds, time.Second*10, s.delete, log)
	s.loop.Run()
}

func (s *spaceDeleter) delete(ctx context.Context) (err error) {
	query := pendingSpacesQuery{s.deletionPeriod}.toMap()
	cur, err := s.spaces.Find(ctx, query)
	if err != nil {
		return
	}
	err = s.processEntry(ctx, cur)
	if err != nil {
		return
	}
	for cur.Next(ctx) {
		err = s.processEntry(ctx, cur)
		if err != nil {
			return
		}
	}
	return
}

func (s *spaceDeleter) processEntry(ctx context.Context, cur *mongo.Cursor) (err error) {
	entry := &statusEntry{}
	err = cur.Decode(entry)
	if err != nil {
		return
	}
	raw := &treechangeproto.RawTreeChangeWithId{}
	err = proto.Unmarshal(entry.DeletionPayload, raw)
	if err != nil {
		return
	}
	op := modifyStatusOp{}
	op.Set.DeletionPayload = entry.DeletionPayload
	op.Set.Status = SpaceStatusDeletionStarted
	op.Set.DeletionDate = entry.DeletionDate
	res := s.spaces.FindOneAndUpdate(ctx, findStatusQuery{
		SpaceId:  entry.SpaceId,
		Status:   SpaceStatusDeletionPending,
		Identity: entry.Identity,
	}, op)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil
		}
		return res.Err()
	}
	err = s.delSender.Delete(ctx, entry.SpaceId, raw)
	// TODO: if this is an error related to contents of change that would never be accepted, we should remove the deletion status altogether
	if err != nil {
		op.Set.Status = SpaceStatusDeletionPending
	} else {
		op.Set.Status = SpaceStatusDeleted
	}
	s.spaces.FindOneAndUpdate(ctx, findStatusQuery{
		SpaceId:  entry.SpaceId,
		Status:   SpaceStatusDeletionStarted,
		Identity: entry.Identity,
	}, op)
	return nil
}

func (s *spaceDeleter) Close() {
	s.loop.Close()
}
