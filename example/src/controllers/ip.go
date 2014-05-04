package controllers

import (
	"models"
	"web/context"
)

func GetIp(ctx *context.Context, ip string) []byte {
	return_map := make(map[string]interface{})
	ipInfo, err := models.GetIp(ip)
	if err != nil {
		return_map["error"] = err.Error()
		return ctx.Response.JsonError(return_map)
	}
	return ctx.Response.JsonSuccess(ipInfo)
}
