package proxy

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger-labs/yui-relayer/core"
)

type ProxyChainConfigI interface {
	proto.Message
	Build() (ProxyChainI, error)
}

type ProxyChainProverConfigI interface {
	proto.Message
	Build(ProxyChainI) (ProxyChainProverI, error)
}

var _ core.ProverConfigI = (*ProverConfig)(nil)

func (pc ProverConfig) Build(chain core.ChainI) (core.ProverI, error) {
	prover, err := pc.Prover.GetCachedValue().(core.ProverConfigI).Build(chain)
	if err != nil {
		return nil, err
	}
	return NewProver(chain, prover, pc.Upstream, pc.Downstream)
}

type UpstreamProxy struct {
	ProxyProvableChain
}

func NewUpstreamProxy(config *ProxyConfig) *UpstreamProxy {
	if config == nil {
		return nil
	}
	proxyChain, err := config.ProxyChain.GetCachedValue().(ProxyChainConfigI).Build()
	if err != nil {
		panic(err)
	}
	proxyChainProver, err := config.ProxyChainProver.GetCachedValue().(ProxyChainProverConfigI).Build(proxyChain)
	if err != nil {
		panic(err)
	}
	return &UpstreamProxy{
		ProxyProvableChain: ProxyProvableChain{ProxyChainI: proxyChain, ProxyChainProverI: proxyChainProver},
	}
}

type DownstreamProxy struct {
	ProxyProvableChain
}

func NewDownstreamProxy(config *ProxyConfig) *DownstreamProxy {
	if config == nil {
		return nil
	}
	proxyChain, err := config.ProxyChain.GetCachedValue().(ProxyChainConfigI).Build()
	if err != nil {
		panic(err)
	}
	proxyChainProver, err := config.ProxyChainProver.GetCachedValue().(ProxyChainProverConfigI).Build(proxyChain)
	if err != nil {
		panic(err)
	}
	return &DownstreamProxy{
		ProxyProvableChain: ProxyProvableChain{ProxyChainI: proxyChain, ProxyChainProverI: proxyChainProver},
	}
}

var _, _ codectypes.UnpackInterfacesMessage = (*ProverConfig)(nil), (*ProxyConfig)(nil)

func (cfg *ProverConfig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if cfg == nil {
		return nil
	}
	if err := unpacker.UnpackAny(cfg.Prover, new(core.ProverConfigI)); err != nil {
		return err
	}
	if err := cfg.Upstream.UnpackInterfaces(unpacker); err != nil {
		return err
	}
	if err := cfg.Downstream.UnpackInterfaces(unpacker); err != nil {
		return err
	}
	return nil
}

func (cfg *ProxyConfig) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if cfg == nil {
		return nil
	}
	if err := unpacker.UnpackAny(cfg.ProxyChain, new(ProxyChainConfigI)); err != nil {
		return err
	}
	if err := unpacker.UnpackAny(cfg.ProxyChainProver, new(ProxyChainProverConfigI)); err != nil {
		return err
	}
	return nil
}
