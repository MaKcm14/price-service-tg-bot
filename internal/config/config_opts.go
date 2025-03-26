package config

import (
	"fmt"
	"os"
	"strings"
)

// ConfigOpt configs the option for the app's start.
type ConfigOpt func(st *Settings) error

func ConfigBotToken(tokenKey string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"

		val := os.Getenv(tokenKey)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", tokenKey, errVar)
		}
		st.TgBotToken = val

		return nil
	}
}

func ConfigDSN(dsnKey string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"

		val := os.Getenv(dsnKey)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", dsnKey, errVar)
		}
		st.DSN = val

		return nil
	}
}

func ConfigRedisInfo(socketKey, pwdKey string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"

		socketVal := os.Getenv(socketKey)

		if len(socketVal) == 0 {
			return fmt.Errorf("error of %s var: %v", socketKey, errVar)
		}
		st.CacheConf.Socket = socketVal

		pwdVal := os.Getenv(pwdKey)

		if len(pwdVal) == 0 {
			return fmt.Errorf("error of %s var: %v", pwdKey, errVar)
		}
		st.CacheConf.PWD = pwdVal

		return nil
	}
}

func ConfigSocket(socketKey string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"

		val := os.Getenv(socketKey)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", socketKey, errVar)
		}
		st.PriceServiceSocket = val

		return nil
	}
}

func ConfigBrokers(brokersKey string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"

		val := os.Getenv(brokersKey)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", brokersKey, errVar)
		}
		st.Brokers = strings.Split(val, " ")

		return nil
	}
}
