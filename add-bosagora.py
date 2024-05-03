import json

# Read the JSON data
with open('./network/genesis.json', 'r') as f:
    data = json.load(f)

'''
	  "bosagoraBlock": 1,
	  "bosagora": {
		"commonsBudget": "0x71D208bfd49375285301343C719e1EA087c87b43",
		"commonsBudgetReward": 120000000000000000000,
		"lastCommonsBudgetRewardBlock": 16202748
	  },
'''

data['config']['bosagoraBlock'] = 1
data['config']['bosagora'] = {
    "commonsBudget": "0x71D208bfd49375285301343C719e1EA087c87b43",
    "commonsBudgetReward": 120000000000000000000,
    "lastCommonsBudgetRewardBlock": 16202748
}

# Write the data back to the file
with open('./network/genesis.json', 'w') as f:
    json.dump(data, f, indent=4)