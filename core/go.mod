module core

go 1.18

require (
	routes v0.0.0
)

replace (
    routes => ../app/routes
    controller => ../app/controller
)