module "go"

go 1.18

require (
	core v0.0.0
)

replace (
    core => ./core
)