package main

import (
	"database/sql"
	"github.com/pkg/errors"
)

func main() {
	// handle错误，打印日志，返回错误码
	_, _ = editArticle(1)
	_, _ = readArticle(1)
}

// ------------------- service ------------------- //
// service层，根据业务具体情况决定如何处理RecordNotFoundError
// 1. 对于类似加金币这种无法降级的业务，可以直接wrap错误继续上抛
// 2. 对于只是读取用户文章进行展示的业务，吞掉错误，返回空对象

// 无法降级的业务，把错误wrap上抛
func editArticle(id int) (Article, error) {
	article, err := getArticle(id)

	if err != nil {
		if errors.Is(err, RecordNotFoundError) {
			return article, errors.Wrapf(err, "can not edit an not exist article : %d", id)
		}
	}

	// 继续业务逻辑

	return article, nil
}

// 可以降级的业务，返回mock或者空数据
func readArticle(id int) (Article, error) {
	article, err := getArticle(id)

	if err != nil {
		if errors.Is(err, RecordNotFoundError) {
			return Article{}, nil
		}
	}

	return article, nil
}

// ------------------- dao ------------------- //
// 虽然根据准则，对于第三方库抛出的错误应该直接上抛
// 但是为了兼容多种DB存储，应该屏蔽底层DB细节
// 封装统一notfoundError上抛

type Article struct {
	ID  int
	Cnt string
}

var RecordNotFoundError = errors.New("record not found")

func getArticle(id int) (Article, error) {
	// 找不到的逻辑
	if id != 1 {
		err := sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, RecordNotFoundError
		}
	}

	// 找到的逻辑
	return Article{1, "cnt"}, nil
}
