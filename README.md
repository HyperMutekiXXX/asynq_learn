# asynq_learn
学习asynq延迟队列

支付超时案例
```
asynq_learn
├─client // 生产者
├─constant
├─controller
├─model
│  ├─item
│  ├─order
│  └─user
├─server // 消费者
├─tasks // 在这里写不同的任务
└─test
```

```
/api/order/list  //订单列表
/api/order/buy   //购买,订单状态转变为等待支付,一分钟内需支付，否则状态转为已超时
/api/order/pay   //支付,订单状态转变为已支付
/api/order/cancel //取消订单,订单状态转变为取消状态
/api/item/list   //商品列表
/api/user/list   //用户列表
```
