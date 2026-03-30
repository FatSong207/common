package gen

import (
	"gorm.io/gen"
)

// G generates a query model using the given output path and variadic list of models.
func G(output string, models ...any) {
	// Generate Query
	gTicket := gen.NewGenerator(gen.Config{
		//OutPath: "app/ticket/model/query",
		OutPath: output,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gTicket.ApplyBasic(models...)

	gTicket.Execute()
}
