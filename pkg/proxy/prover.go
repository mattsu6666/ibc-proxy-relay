package proxy

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/modules/core/02-client/types"
	conntypes "github.com/cosmos/ibc-go/modules/core/03-connection/types"
	chantypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/modules/core/exported"
	"github.com/spf13/viper"

	"github.com/hyperledger-labs/yui-relayer/core"
)

// TODO how to initialize a chain that implements ProxiableChainI
type Prover struct {
	chain      ProxiableChainI
	prover     core.ProverI
	upstream   *Upstream
	downstream *Downstream
}

var (
	_ core.ProverI       = (*Prover)(nil)
	_ ProxyEventListener = (*Prover)(nil)
)

func NewProver(chain ProxiableChainI, prover core.ProverI, upstreamConfig *UpstreamConfig, downstreamConfig *DownstreamConfig) (*Prover, error) {
	if upstreamConfig == nil && downstreamConfig == nil {
		return nil, fmt.Errorf("either upstream or downstream must be not nil")
	} else if downstreamConfig != nil {
		prover = NewDownstreamProver(prover)
	}
	// XXX: the following params should be given from the relayer
	homePath := viper.GetString(flags.FlagHome)
	timeout := time.Minute
	pr := &Prover{
		chain:      chain,
		prover:     prover,
		upstream:   NewUpstream(upstreamConfig, chain, homePath, timeout),
		downstream: NewDownstream(downstreamConfig, chain, homePath, timeout),
	}
	if pr.upstream != nil {
		pr.chain.RegisterEventListener(pr)
	}
	return pr, nil
}

func (pr *Prover) GetUnderlyingProver() core.ProverI {
	switch prover := pr.prover.(type) {
	case *DownstreamProver:
		return prover.ProverI
	default:
		return prover
	}
}

// GetChainID returns the chain ID
func (pr *Prover) GetChainID() string {
	return pr.chain.ChainID()
}

// QueryLatestHeader returns the latest header from the chain
func (pr *Prover) QueryLatestHeader() (out core.HeaderI, err error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryLatestHeader()
	} else {
		return pr.prover.QueryLatestHeader()
	}
}

// GetLatestLightHeight returns the latest height on the light client
func (pr *Prover) GetLatestLightHeight() (int64, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.GetLatestLightHeight()
	} else {
		return pr.prover.GetLatestLightHeight()
	}
}

// CreateMsgCreateClient creates a CreateClientMsg to this chain
func (pr *Prover) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.CreateMsgCreateClient(clientID, dstHeader, signer)
	} else {
		return pr.prover.CreateMsgCreateClient(clientID, dstHeader, signer)
	}
}

// SetupHeader creates a new header based on a given header
func (pr *Prover) SetupHeader(dst core.LightClientIBCQueryierI, baseSrcHeader core.HeaderI) (core.HeaderI, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.SetupHeader(dst, baseSrcHeader)
	} else {
		return pr.prover.SetupHeader(dst, baseSrcHeader)
	}
}

// UpdateLightWithHeader updates a header on the light client and returns the header and height corresponding to the chain
func (pr *Prover) UpdateLightWithHeader() (header core.HeaderI, provableHeight int64, queryableHeight int64, err error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.UpdateLightWithHeader()
	} else {
		return pr.prover.UpdateLightWithHeader()
	}
}

// QueryClientConsensusState returns the ClientConsensusState and its proof
func (pr *Prover) QueryClientConsensusStateWithProof(height int64, dstClientConsHeight ibcexported.Height) (*clienttypes.QueryConsensusStateResponse, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyClientConsensusStateWithProof(height, dstClientConsHeight)
	} else {
		return pr.prover.QueryClientConsensusStateWithProof(height, dstClientConsHeight)
	}
}

// QueryClientStateWithProof returns the ClientState and its proof
func (pr *Prover) QueryClientStateWithProof(height int64) (*clienttypes.QueryClientStateResponse, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyClientStateWithProof(height, "TODO-upstreamClientID")
	} else {
		return pr.prover.QueryClientStateWithProof(height)
	}
}

// QueryConnectionWithProof returns the Connection and its proof
func (pr *Prover) QueryConnectionWithProof(height int64) (*conntypes.QueryConnectionResponse, error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyConnectionStateWithProof(height, "TODO-upstreamClientID")
	} else {
		return pr.prover.QueryConnectionWithProof(height)
	}
}

// QueryChannelWithProof returns the Channel and its proof
func (pr *Prover) QueryChannelWithProof(height int64) (chanRes *chantypes.QueryChannelResponse, err error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyChannelWithProof(height, "TODO-upstreamClientID")
	} else {
		return pr.prover.QueryChannelWithProof(height)
	}
}

// QueryPacketCommitmentWithProof returns the packet commitment and its proof
func (pr *Prover) QueryPacketCommitmentWithProof(height int64, seq uint64) (comRes *chantypes.QueryPacketCommitmentResponse, err error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyPacketCommitmentWithProof(height, seq, "TODO-upstreamClientID")
	} else {
		return pr.prover.QueryPacketCommitmentWithProof(height, seq)
	}
}

// QueryPacketAcknowledgementCommitmentWithProof returns the packet acknowledgement commitment and its proof
func (pr *Prover) QueryPacketAcknowledgementCommitmentWithProof(height int64, seq uint64) (ackRes *chantypes.QueryPacketAcknowledgementResponse, err error) {
	if pr.upstream != nil {
		return pr.upstream.Proxy.QueryProxyPacketAcknowledgementCommitmentWithProof(height, seq, "TODO-upstreamClientID")
	} else {
		return pr.prover.QueryPacketAcknowledgementCommitmentWithProof(height, seq)
	}
}

type DownstreamProver struct {
	core.ProverI
}

var _ core.ProverI = (*DownstreamProver)(nil)

func NewDownstreamProver(prover core.ProverI) *DownstreamProver {
	return &DownstreamProver{ProverI: prover}
}

func (p *DownstreamProver) CreateMsgCreateClient(clientID string, dstHeader core.HeaderI, signer sdk.AccAddress) (*clienttypes.MsgCreateClient, error) {
	// TODO returns a msg corredponding to MultiV client
	panic("not implemented error")
}
