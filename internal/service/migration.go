// Package service app services
package service

import "go-infra/internal/container"

func createRepository(appContainer container.AppContainer) {

	_ = appContainer.Repository()

	initRepositoryMasterData(appContainer)
}

func initRepositoryMasterData(appContainer container.AppContainer) {
	_ = appContainer.Repository()

}
