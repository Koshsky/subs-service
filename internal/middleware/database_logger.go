package middleware

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

func DatabaseLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()

		if isSuccess(status) {
			id, exists := c.Get("db_affected_id")
			idStr := ""
			if exists {
				idStr = fmt.Sprintf(" (ID: %v)", id)
			}

			rowsAffected, rowsExists := c.Get("db_rows_affected")
			rowsStr := ""
			if rowsExists {
				rowsStr = fmt.Sprintf(" [rows: %d]", rowsAffected)
			}

			switch method {
			case "POST":
				log.Printf("%s %s %s%s%s %s",
					color.CyanString("[DB]"),
					color.GreenString("CREATE"),
					path,
					idStr,
					rowsStr,
					color.GreenString("success"),
				)
			case "PUT":
				log.Printf("%s %s %s%s%s %s",
					color.CyanString("[DB]"),
					color.BlueString("UPDATE"),
					path,
					idStr,
					rowsStr,
					color.GreenString("success"),
				)
			case "DELETE":
				log.Printf("%s %s %s%s%s %s",
					color.CyanString("[DB]"),
					color.RedString("DELETE"),
					path,
					idStr,
					rowsStr,
					color.GreenString("success"),
				)
			case "GET":
				if path == "/subscriptions/total" {
					total, exists := c.Get("db_total_value")
					if exists {
						log.Printf("%s %s %s %v %s",
							color.CyanString("[DB]"),
							color.MagentaString("CALCULATE"),
							path,
							total,
							color.GreenString("success"),
						)
					}
				}
			}
		} else {
			errorMsg, exists := c.Get("db_error")
			if exists {
				log.Printf("%s %s %s %s: %v",
					color.CyanString("[DB]"),
					color.RedString("ERROR"),
					method,
					path,
					color.RedString("%v", errorMsg),
				)
			}
		}
	}
}
