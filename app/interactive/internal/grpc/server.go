// Package itrGrpc 用来将业务暴露为GRPC接口
package itrGrpc

import (
	"context"
	itrAPI "github.com/LEILEI0628/GoWeb-MicroServices/api/interactive/v1"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/domain"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InteractiveServiceServer 将service包装为grpc，在此实现grpc相关业务操作
type InteractiveServiceServer struct {
	itrAPI.UnimplementedInteractiveServiceServer                                     // 保证项目兼容性
	svc                                          service.InteractiveServiceInterface // 核心业务逻辑一定是在 service 中的
}

func NewInteractiveServiceServer(svc service.InteractiveServiceInterface) *InteractiveServiceServer {
	return &InteractiveServiceServer{svc: svc}
}
func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *itrAPI.IncrReadCntRequest) (*itrAPI.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, request.Biz, request.BizId)
	//if err != nil {
	//	return nil, err
	//}
	return &itrAPI.IncrReadCntResponse{}, err // 调用者应秉持一旦err不为nil就不访问Response
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *itrAPI.LikeRequest) (*itrAPI.LikeResponse, error) {
	if request.GetUid() <= 0 { // 一种返回错误码的方式
		return nil, status.Error(codes.InvalidArgument, "Uid错误")
	}
	err := i.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &itrAPI.LikeResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *itrAPI.CancelLikeRequest) (*itrAPI.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &itrAPI.CancelLikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *itrAPI.CollectRequest) (*itrAPI.CollectResponse, error) {
	err := i.svc.Collect(ctx, request.GetBiz(), request.GetBizId(), request.GetCid(), request.GetUid())
	return &itrAPI.CollectResponse{}, err
}

func (i *InteractiveServiceServer) Get(ctx context.Context, request *itrAPI.GetRequest) (*itrAPI.GetResponse, error) {
	res, err := i.svc.Get(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &itrAPI.GetResponse{
		Intr: i.toDTO(res),
	}, nil
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *itrAPI.GetByIdsRequest) (*itrAPI.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*itrAPI.Interactive, len(res))
	for index, interactive := range res {
		m[index] = i.toDTO(interactive)
	}
	return &itrAPI.GetByIdsResponse{Intrs: m}, err
}

func (i *InteractiveServiceServer) mustEmbedUnimplementedInteractiveServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (i *InteractiveServiceServer) toDTO(itr domain.Interactive) *itrAPI.Interactive {
	return &itrAPI.Interactive{
		Biz:        itr.Biz,
		BizId:      itr.BizId,
		CollectCnt: itr.CollectCnt,
		Collected:  itr.Collected,
		Liked:      itr.Liked,
		LikeCnt:    itr.LikeCnt,
		ReadCnt:    itr.ReadCnt,
	}
}
