package agora

import "testing"

func TestAllocatedYearlyValidatorRewards(t *testing.T) {
	tests := []struct {
		name                string
		secondsSinceGenesis uint64
		want                uint64
	}{
		{"beginning of year 1", 1, 44_150_400},
		{"end of year 1", YearOfSecs - 1, 44_150_400},
		{"beginning of year 2", YearOfSecs + 1, 43_555_694},
		{"end of year 2", 2*YearOfSecs - 1, 43_555_694},
		{"year 3", 3*YearOfSecs - 1, 42_968_998},
		{"year 10", 10*YearOfSecs - 1, 39_077_544},
		{"year 50", 50*YearOfSecs - 1, 22_716_364},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllocatedYearlyValidatorRewards(tt.secondsSinceGenesis, 1); got != tt.want {
				t.Errorf("AllocatedYearlyValidatorRewards() test for '%s' got %v but wanted %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestAllocatedValidatorRewardsPerEpoch(t *testing.T) {
	cfg := RewardConfig{
		SlotsPerEpoch:             32,
		SecondsPerSlot:            12,
		GweiPerBoa:                1000000000,
		EffectiveBalanceIncrement: 1000000000000,
	}

	tests := []struct {
		name                string
		secondsSinceGenesis uint64
		want                uint64
	}{
		{"beginning of year 1", 1, 537_600_000_000},
		{"end of year 1", YearOfSecs - 1, 537_600_000_000},
		{"beginning of year 2", YearOfSecs + 1, 530_358_528_000},
		{"beginning of year 3", 2*YearOfSecs + 1, 523_214_598_528},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllocatedValidatorRewardsPerEpoch(tt.secondsSinceGenesis, cfg); got != tt.want {
				t.Errorf("AllocatedValidatorRewardsPerEpoch() test for '%s' got %v but wanted %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestValidatorRewardPerEpoch(t *testing.T) {
	cfg := RewardConfig{
		SlotsPerEpoch:             32,
		SecondsPerSlot:            12,
		GweiPerBoa:                1000000000,
		EffectiveBalanceIncrement: 1000000000000,
	}

	tests := []struct {
		name             string
		totalBalance     uint64
		effectiveBalance uint64
		want             uint64
	}{
		{"all", cfg.EffectiveBalanceIncrement, cfg.EffectiveBalanceIncrement, 537_600_000_000},
		{"half", 2 * cfg.EffectiveBalanceIncrement, cfg.EffectiveBalanceIncrement, 268_800_000_000},
		{"one seventh", 7 * cfg.EffectiveBalanceIncrement, cfg.EffectiveBalanceIncrement, 76_800_000_000},
		{"one seventh, but has excess", 7 * cfg.EffectiveBalanceIncrement, cfg.EffectiveBalanceIncrement + 1, 76_800_000_000},
		{"below increment", 2 * cfg.EffectiveBalanceIncrement, cfg.EffectiveBalanceIncrement - 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ValidatorRewardPerEpoch(1, tt.totalBalance, tt.effectiveBalance, cfg); got != tt.want || err != nil {
				t.Errorf("AllocatedValidatorRewardsPerEpoch() test for '%s' got %v but wanted %v", tt.name, got, tt.want)
			}
		})
	}
}
