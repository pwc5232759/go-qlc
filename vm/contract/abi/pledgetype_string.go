// Code generated by "stringer -type=PledgeType"; DO NOT EDIT.

package abi

import "strconv"

const _PledgeType_name = "NetowrkVoteStorageOracle"

var _PledgeType_index = [...]uint8{0, 7, 11, 18, 24}

func (i PledgeType) String() string {
	if i >= PledgeType(len(_PledgeType_index)-1) {
		return "PledgeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _PledgeType_name[_PledgeType_index[i]:_PledgeType_index[i+1]]
}
