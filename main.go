package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/stevebargelt/trainNewCatModel/config"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.1/customvision/training"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	viper.SetDefault("HTTP_RETRY_COUNT", 20)

	var configuration config.Configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	projectID, err := makeUUID(configuration.ProjectID)
	// projectIDDirection, err := makeUUID(configuration.ProjectIDDirection)
	//iterationID, err := makeUUID(configuration.IterationID)
	// iterationIDDirection, err := makeUUID(configuration.IterationIDDirection)
	ctx := context.Background()

	trainer := training.New(configuration.TrainingKey, configuration.TrainingEndpoint)

	var emptyUUIDS []uuid.UUID
	allImages := true
	trainer.DeleteImages(ctx, projectID, emptyUUIDS, &allImages, &allImages)

	// fmt.Printf("iterationID = %s\n", trainer.GetIteration)
	// Make two tags in the new project
	bearTag, err := trainer.CreateTag(ctx, projectID, "Bear", "Bear cat tag", string(training.Regular))
	if err != nil {
		panic(err)
	}

	fmt.Printf("BearTag: %v\n", bearTag)
	naraTag, _ := trainer.CreateTag(ctx, projectID, "Nara", "Nara cat tag", string(training.Regular))
	negativeTag, _ := trainer.CreateTag(ctx, projectID, "Negative", "Negative", string(training.Negative))

	fmt.Println("Adding images...")
	bearImages, err := ioutil.ReadDir(path.Join(configuration.TrainingFolder, "bear"))
	if err != nil {
		fmt.Println("Error finding Sample images")
	}

	naraImages, err := ioutil.ReadDir(path.Join(configuration.TrainingFolder, "nara"))
	if err != nil {
		fmt.Println("Error finding Sample images")
	}

	negativeImages, err := ioutil.ReadDir(path.Join(configuration.TrainingFolder, "negative"))
	if err != nil {
		fmt.Println("Error finding Sample images")
	}

	for _, file := range bearImages {
		fmt.Printf("File: %s\n", path.Join(configuration.TrainingFolder, "/bear", file.Name()))
		imageFile, _ := ioutil.ReadFile(path.Join(configuration.TrainingFolder, "/bear", file.Name()))
		imageData := ioutil.NopCloser(bytes.NewReader(imageFile))

		trainer.CreateImagesFromData(ctx, projectID, imageData, []uuid.UUID{*bearTag.ID})
	}

	for _, file := range naraImages {
		imageFile, _ := ioutil.ReadFile(path.Join(configuration.TrainingFolder, "/nara", file.Name()))
		imageData := ioutil.NopCloser(bytes.NewReader(imageFile))
		trainer.CreateImagesFromData(ctx, projectID, imageData, []uuid.UUID{*naraTag.ID})
	}

	for _, file := range negativeImages {
		imageFile, _ := ioutil.ReadFile(path.Join(configuration.TrainingFolder, "/negative", file.Name()))
		imageData := ioutil.NopCloser(bytes.NewReader(imageFile))
		trainer.CreateImagesFromData(ctx, projectID, imageData, []uuid.UUID{*negativeTag.ID})
	}

	fmt.Println("Training...")

	var trainingHours int32 = 2
	var forceTrain bool = false
	iteration, _ := trainer.TrainProject(ctx, projectID, "Advanced", &trainingHours, &forceTrain, "steve@bargelt.com")
	for {
		if *iteration.Status != "Training" {
			break
		}
		fmt.Println("Training status: " + *iteration.Status)
		time.Sleep(1 * time.Second)
		iteration, _ = trainer.GetIteration(ctx, projectID, *iteration.ID)
	}
	fmt.Println("Training status: " + *iteration.Status)

	// TODO: is this the right resourceID??? it's under training in the portal
	trainer.PublishIteration(ctx, projectID, *iteration.ID, "AutoGen-Iteration", configuration.TrainingResourceID)

	// fmt.Println("Predicting...")
	// predictor := prediction.New(configuration.PredictionKey, configuration.PredictionEndpointURL)

	// testImageData, _ := ioutil.ReadFile(path.Join(configuration.TestingFolder, "Test", "test_image.jpg"))
	// results, _ := predictor.ClassifyImage(ctx, *project.ID, iteration_publish_name, ioutil.NopCloser(bytes.NewReader(testImageData)), "")

	// for _, prediction := range *results.Predictions {
	// 	fmt.Printf("\t%s: %.2f%%", *prediction.TagName, *prediction.Probability*100)
	// 	fmt.Println("")
	// }
}

func makeUUID(source string) (uuid.UUID, error) {

	result, err := uuid.FromString(source)
	if err != nil {
		fmt.Printf("Something went wrong creating UUID: %s", err)
	}
	return result, err

}
