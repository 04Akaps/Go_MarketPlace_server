package gRpcUtils

import (
	"google.golang.org/grpc/keepalive"
	"time"
)

func GetKeepAliveClientParameters(time, timeout time.Duration) keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                time,
		Timeout:             timeout,
		PermitWithoutStream: false,
	}
}

func GetKeepAliveEnforcement(minTime time.Duration) keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             minTime,
		PermitWithoutStream: false,
	}
}

func GetKeepAliveServerParameters(maxConnIdle, maxConnAge, MaxConnAgeGrace, time, timeout time.Duration) keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     maxConnIdle,
		MaxConnectionAge:      maxConnAge,
		MaxConnectionAgeGrace: MaxConnAgeGrace,
		Time:                  time,
		Timeout:               timeout,
	}
}
