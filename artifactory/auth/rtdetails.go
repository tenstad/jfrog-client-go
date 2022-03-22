package auth

import (
	"github.com/tenstad/jfrog-client-go/artifactory"
	"github.com/tenstad/jfrog-client-go/auth"
	"github.com/tenstad/jfrog-client-go/config"
	"github.com/tenstad/jfrog-client-go/utils/log"
)

func NewArtifactoryDetails() auth.ServiceDetails {
	return &artifactoryDetails{}
}

type artifactoryDetails struct {
	auth.CommonConfigFields
}

func (rt *artifactoryDetails) GetVersion() (string, error) {
	var err error
	if rt.Version == "" {
		rt.Version, err = rt.getArtifactoryVersion()
		if err != nil {
			return "", err
		}
		log.Debug("The Artifactory version is:", rt.Version)
	}
	return rt.Version, nil
}

func (rt *artifactoryDetails) getArtifactoryVersion() (string, error) {
	cd := auth.ServiceDetails(rt)
	serviceConfig, err := config.NewConfigBuilder().
		SetServiceDetails(cd).
		SetCertificatesPath(cd.GetClientCertPath()).
		Build()
	if err != nil {
		return "", err
	}
	var sm artifactory.ArtifactoryServicesManager
	client := rt.GetClient()
	if client != nil {
		sm, err = artifactory.NewWithClient(serviceConfig, client)
	} else {
		sm, err = artifactory.New(serviceConfig)
	}
	if err != nil {
		return "", err
	}
	return sm.GetVersion()
}
