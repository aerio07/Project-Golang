package routes

import (
	"project_uas/app/service"
	"project_uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, svc *service.UserService) {
    base := "/api/v1/users"

    app.Get(base,
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.GetUsers,
    )

    app.Get(base+"/:id",
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.GetUser,
    )

    app.Post(base,
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.CreateUser,
    )

    app.Put(base+"/:id",
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.UpdateUser,
    )

    app.Delete(base+"/:id",
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.DeleteUser,
    )

    app.Put(base+"/:id/role",
        middleware.JWTMiddleware,
        middleware.RequirePermission("user:manage"),
        svc.AssignRole,
    )
}
