package main

import (
	"fmt"
	"testing"
)

func TestErrorHandle(t *testing.T) {
	// 读取到数据
	article, err := readArticle(1)
	if err == nil {
		fmt.Println("read article: ", article)
	}

	// 可以降级的业务,读到空数据
	article1, err1 := readArticle(2)
	if err1 == nil {
		fmt.Println("read article: ", article1)
	}

	// 不可降级的业务,读到数据
	article2, err2 := editArticle(1)
	if err2 == nil {
		fmt.Println("edit article: ", article2)
	}

	// 不可降级的业务,产生错误
	_, err3 := editArticle(2)
	if err3 != nil {
		// 打印日志
		fmt.Printf("article not found error:%+v \n", err3)
	}
}
