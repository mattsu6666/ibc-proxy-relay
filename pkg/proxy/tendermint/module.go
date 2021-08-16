package tendermint

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/datachainlab/ibc-proxy-prover/pkg/proxy"
	"github.com/datachainlab/ibc-proxy-prover/pkg/proxy/tendermint/cmd"
	"github.com/hyperledger-labs/yui-relayer/config"
	"github.com/spf13/cobra"
)

type Module struct{}

var _ config.ModuleI = (*Module)(nil)

// Name returns the name of the module
func (m Module) Name() string {
	return "proxy-tendermint"
}

// RegisterInterfaces register the module interfaces to protobuf Any.
func (m Module) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*proxy.ProxyChainConfigI)(nil),
		&ProxyChainConfig{},
	)
	registry.RegisterImplementations(
		(*proxy.ProxyChainProverConfigI)(nil),
		&ProxyChainProverConfig{},
	)
}

// GetCmd returns the command
func (m Module) GetCmd(ctx *config.Context) *cobra.Command {
	return cmd.TendermintCmd(ctx.Codec, ctx)
}
