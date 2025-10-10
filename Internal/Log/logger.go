package logger

import "go.uber.org/zap"

var Log *zap.Logger

// to init the global logger
func InitLogger(env string){
	var err error

	// checking the environment (development or production)
	if env == "development" {
		Log, err = zap.NewDevelopment()
	} else if env == "production" {
		Log, err = zap.NewProduction()
	}

	if err != nil {
		panic("Failed to Initialize Logger: " + err.Error())
	}

	defer Log.Sync()

	Log.Info("Logger Initialized ",zap.String("Environment",env))
}
