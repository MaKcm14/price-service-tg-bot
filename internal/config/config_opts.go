package config

import (
	"fmt"
	"os"
	"strings"
)

// ConfigOpt configs the option for the app's start.
type ConfigOpt func(st *Settings) error

func configEnv(key string) (string, error) {
	const errVar = "check this var is set or its value is not empty and try again"

	val := os.Getenv(key)

	if len(val) == 0 {
		return "", fmt.Errorf("error of %s var: %v", key, errVar)
	}

	return val, nil
}

func ConfigBotToken(tokenKey string) ConfigOpt {
	return func(st *Settings) error {
		tokenVal, err := configEnv(tokenKey)

		if err != nil {
			return err
		}
		st.TgBotToken = tokenVal

		return nil
	}
}

func ConfigDSN(dsnKey string) ConfigOpt {
	return func(st *Settings) error {
		dsnVal, err := configEnv(dsnKey)

		if err != nil {
			return err
		}
		st.DSN = dsnVal

		return nil
	}
}

func ConfigRedisInfo(socketKey, pwdKey string) ConfigOpt {
	return func(st *Settings) error {
		socketVal, err := configEnv(socketKey)

		if err != nil {
			return err
		}

		st.CacheConf.Socket = socketVal
		st.CacheConf.PWD = os.Getenv(pwdKey)

		return nil
	}
}

func ConfigSocket(socketKey string) ConfigOpt {
	return func(st *Settings) error {
		socketVal, err := configEnv(socketKey)

		if err != nil {
			return err
		}
		st.PriceServiceSocket = socketVal

		return nil
	}
}

func ConfigBrokers(brokersKey string) ConfigOpt {
	return func(st *Settings) error {
		brokersVal, err := configEnv(brokersKey)

		if err != nil {
			return err
		}
		st.Brokers = strings.Split(brokersVal, " ")

		return nil
	}
}
