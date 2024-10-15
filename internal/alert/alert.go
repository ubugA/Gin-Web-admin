package alert

import (
	"fmt"

	"gin-api-admin/internal/proposal"
)

func NotifyHandler() func(msg *proposal.AlertMessage) {
	return func(msg *proposal.AlertMessage) {
		// 自行实现适合自己的告警通知方式，例如：邮件/短信/微信/企业微信/飞书 等

		// 可使用的参数如下：
		fmt.Println("---------------------------------")
		fmt.Println("自行实现适合自己的告警通知方式，信息如下：")
		fmt.Println("项目名称：", msg.ProjectName)
		fmt.Println("发生环境：", msg.Env)
		fmt.Println("当前请求的唯一ID：", msg.TraceID)
		fmt.Println("当前请求的HOST：", msg.HOST)
		fmt.Println("当前请求的URI：", msg.URI)
		fmt.Println("当前请求的Method：", msg.Method)
		fmt.Println("错误信息：", msg.ErrorMessage)
		fmt.Println("堆栈信息：", msg.ErrorStack)
		fmt.Println("发生时间：", msg.Time)
		fmt.Println("---------------------------------")

	}
}
