package config

// Configuration for app
type Configuration struct {
	TrainingFolder        string `mapstructure:"TRAINING_FOLDER"`
	ProjectID             string `mapstructure:"CUSTOM_VISION_PROJECT_ID"`
	ProjectIDDirection    string `mapstructure:"CUSTOM_VISION_PROJECT_DIRECTION_ID"`
	PredictionKey         string `mapstructure:"CUSTOM_VISION_PREDICTION_KEY"`
	PredictionEndpointURL string `mapstructure:"CUSTOM_VISION_ENDPOINT"`
	IterationID           string `mapstructure:"CUSTOM_VISION_ITERATION_ID"`
	IterationIDDirection  string `mapstructure:"CUSTOM_VISION_ITERATION_DIRECTION_ID"`
	TrainingKey           string `mapstructure:"CUSTOM_VISION_TRAINING_KEY"`
	TrainingEndpoint      string `mapstructure:"CUSTOM_VISION_TRAINING_ENDPOINT"`
	TrainingResourceID    string `mapstructure:"CUSTOM_VISION_TRAINING_RESOURCEID"`
	TestingFolder         string `mapstructure:"TESTING_FOLDER"`
}
