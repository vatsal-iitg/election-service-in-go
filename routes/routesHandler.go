package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	controller "github.com/vatsal-iitg/election-service-in-go/controllers"
	middleware "github.com/vatsal-iitg/election-service-in-go/middleware"
)

func RouterHandler(incomingRoutes *gin.Engine) {
	log.Println("Entered routeshandler")

	voterRoutes := incomingRoutes.Group("/voters")
	{
		voterRoutes.POST("/register", controller.RegisterVoter)
		voterRoutes.POST("/login", controller.LoginVoter)
	}

	candidateRoutes := incomingRoutes.Group("/candidates")
	{
		candidateRoutes.POST("/register", controller.RegisterCandidate)
	}

	electionOfficerRoutes := incomingRoutes.Group("/election-officer")
	electionOfficerRoutes.Use(middleware.AuthMiddlewareForOfficer())

	{
		electionOfficerRoutes.POST("/register", controller.RegisterElectionOfficer)
		electionOfficerRoutes.POST("/login", controller.LoginElectionOfficer)

		electionOfficerRoutes.Use(middleware.AuthorizedOnlyForOfficer())
		electionOfficerRoutes.POST("/create-constituency", controller.CreateConstituency)
		electionOfficerRoutes.PUT("/update-constituency/:id", controller.UpdateConstituency)
	}

}
