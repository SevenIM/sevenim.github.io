package public

import "errors"

func GenSessionId(sender uint64, recver uint64) (string, error) {
	if sender == 0 || recver == 0 {
		return "", errors.New("sender or recver is invalid")
	}
	if sender == recver {
		return "", errors.New("sender and recver is same")
	}
	if sender > recver {
		sender, recver = recver, sender
	}
	return string(sender) + ":" + string(recver), nil
}

func APHash(sessionId string) uint32 {
	hash := uint32(0)
	for i := 0; i < len(sessionId); i++ {
		if (i & 1) == 0 {
			hash ^= ((hash << 7) ^ uint32(sessionId[i]) ^ (hash >> 3))
		} else {
			hash ^= (^((hash << 11) ^ uint32(sessionId[i]) ^ (hash >> 5)))
		}
	}
	return hash
}
