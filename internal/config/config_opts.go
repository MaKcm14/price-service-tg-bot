package config

import (
	"fmt"
	"os"
)

// ConfigOpt configs the option for the app's start.
type ConfigOpt func(st *Settings) error

func ConfigBotToken(key string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"
		tokenName := key

		val := os.Getenv(tokenName)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", tokenName, errVar)
		}
		st.TgBotToken = val

		return nil
	}
}

func ConfigDSN(key string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"
		dsnName := key

		val := os.Getenv(dsnName)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", dsnName, errVar)
		}
		st.DSN = val

		return nil
	}
}

func ConfigSocket(key string) ConfigOpt {
	return func(st *Settings) error {
		const errVar = "check this var is set or its value is not empty and try again"
		socketName := key

		val := os.Getenv(socketName)

		if len(val) == 0 {
			return fmt.Errorf("error of %s var: %v", socketName, errVar)
		}
		st.PriceServiceSocket = val

		return nil
	}
}
