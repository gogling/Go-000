### 总结
- 目前阶段，直接使用`github.com/pkg/errors`
- `you should only handle errors once. Handling an error means inspecting the error value, and making a single decision.
` 错误只处理一次，打印也算做处理

### 作业：
Q : 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么？应该怎么做请写出代码？

##### 回答：
- 虽然根据准则，对于第三方库抛出的错误应该直接上抛，但是为了兼容多种DB存储，应该屏蔽底层DB细节，封装统一NotFoundError上抛
- service层，根据业务具体情况决定如何处理RecordNotFoundError
    - 对于类似加金币这种无法降级的业务，可以直接wrap错误继续上抛
    - 对于只是读取用户文章进行展示的业务，吞掉错误，返回空对象
- demo代码见`main.go`
