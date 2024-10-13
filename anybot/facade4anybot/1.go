package facade4anybot

import (
	"github.com/strongo/delaying"
)

func InitDelaying(mustRegisterFunc func(key string, i any) delaying.Delayer) {
	//delayerSetUserReferrer = mustRegisterFunc("delayedSetUserReferrer", delayedSetUserReferrer)
}

//var delayerSetUserReferrer delaying.Delayer
