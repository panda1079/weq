module core

go 1.18

require (
	controller v0.0.0-00010101000000-000000000000
	routes v0.0.0
)

replace (
	controller => ../app/controller
	routes => ../app/routes
)
