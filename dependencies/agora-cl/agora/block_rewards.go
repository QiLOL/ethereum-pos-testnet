package agora

import (
	"github.com/pkg/errors"
	"math/big"

	"github.com/prysmaticlabs/prysm/v4/config/params"
)

// Global constants (start with a capital letter)
const (
	YearOfSecs uint64 = 365 * 24 * 60 * 60
)

// Package scope constants
const (
	firstYearValRewards uint64 = 7 * (YearOfSecs / 5)
)

type RewardConfig struct {
	SlotsPerEpoch             uint64
	SecondsPerSlot            uint64
	GweiPerBoa                uint64
	EffectiveBalanceIncrement uint64
}

// Agora rewards as defined in the white paper (in Gwei)
//
// Allocated Validator rewards are 7 coins per 5 seconds for first year
//   then reduced by 1.347% every year
func AllocatedYearlyValidatorRewards(secondsSinceGenesis uint64, GweiPerBoa uint64) uint64 {
	yearsSinceGenesis := secondsSinceGenesis / YearOfSecs

	bigGweiPerBoa := new(big.Int).SetUint64(GweiPerBoa)
	bigfirstYearValRewards := new(big.Int).SetUint64(firstYearValRewards)
	yearlyReward := new(big.Int).Mul(bigGweiPerBoa, bigfirstYearValRewards)

	for y := yearsSinceGenesis; y > 0; y-- {
		yearlyReward.Mul(yearlyReward, big.NewInt(98_653))
		yearlyReward.Div(yearlyReward, big.NewInt(100_000))
	}
	return yearlyReward.Uint64()
}

func AllocatedValidatorRewardsPerEpoch(secondsSinceGenesis uint64, cfg RewardConfig) uint64 {
	allocatedRewardsPerSecond := AllocatedYearlyValidatorRewards(secondsSinceGenesis, cfg.GweiPerBoa) / YearOfSecs
	return cfg.SecondsPerSlot * cfg.SlotsPerEpoch * allocatedRewardsPerSecond
}

func ValidatorRewardPerEpoch(secondsSinceGenesis uint64, totalBalance uint64, effectiveBalance uint64, cfg RewardConfig) (uint64, error) {
	if totalBalance <= 0 {
		return 0, errors.New("active balance can't be 0")
	}
	allocatedValidatorRewardsPerEpoch := new(big.Int).SetUint64(AllocatedValidatorRewardsPerEpoch(secondsSinceGenesis, cfg))
	increments := new(big.Int).SetUint64(effectiveBalance / cfg.EffectiveBalanceIncrement)
	bigTotalBalance := new(big.Int).SetUint64(totalBalance)

	bigEffectiveBalance := new(big.Int).SetUint64(cfg.EffectiveBalanceIncrement)
	bigEffectiveBalance.Mul(bigEffectiveBalance, increments)
	bigEffectiveBalance.Mul(bigEffectiveBalance, allocatedValidatorRewardsPerEpoch)
	bigEffectiveBalance.Div(bigEffectiveBalance, bigTotalBalance)

	return bigEffectiveBalance.Uint64(), nil
}

func MakeAgoraRewardConfig(beaconCfg *params.BeaconChainConfig) RewardConfig {
	agoraConfig := RewardConfig{
		SlotsPerEpoch:             uint64(beaconCfg.SlotsPerEpoch),
		SecondsPerSlot:            beaconCfg.SecondsPerSlot,
		GweiPerBoa:                beaconCfg.GweiPerEth,
		EffectiveBalanceIncrement: beaconCfg.EffectiveBalanceIncrement,
	}
	return agoraConfig
}
