/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

// textCmd represents the text command
var textCmd = &cobra.Command{
	Use:   "text",
	Short: "爬取内容存为text文件",
	Long:  `爬取内容存为text文件`,
	Run: func(cmd *cobra.Command, args []string) {
		c := colly.NewCollector()
		c.MaxDepth = 3
		c.SetRequestTimeout(30 * time.Second)

		// 创建文件
		file, err := os.Create("novel.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// 设置文件写入器
		writer := bufio.NewWriter(file)

		c.OnHTML("h1", func(e *colly.HTMLElement) {
			title := e.Text
			// 写入标题到文件
			writer.WriteString("Title: " + title + "\n")
		})
		c.OnHTML("dd", func(e *colly.HTMLElement) {
			link := e.ChildAttr("a", "href")
			c.Visit(e.Request.AbsoluteURL(link))
		})
		c.OnHTML("#content", func(e *colly.HTMLElement) {
			chapterText := strings.TrimSpace(e.Text)
			lines := strings.Split(chapterText, "\n")

			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					_, err := writer.WriteString(line + "\n")
					if err != nil {
						log.Println("Failed to write to file:", err)
					}
				}
			}

		})

		c.OnError(func(r *colly.Response, err error) {
			log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		})

		// 将缓冲区中的内容写入文件
		defer writer.Flush()

		c.Visit("http://www.xswang.org/book/529")

	},
}

var RequestAddress string // 抓取地址
var ProcessNum int        // 并发数目
var OutputFile string     // 保存的文件

func init() {
	rootCmd.AddCommand(textCmd)

	textCmd.Flags().StringVarP(&RequestAddress, "url", "u", "", "抓取地址")
	textCmd.Flags().IntVarP(&ProcessNum, "processNum", "p", 1, "并发数目")
	textCmd.Flags().StringVarP(&OutputFile, "outputFile", "f", "result.txt", "保存的文件")

}
