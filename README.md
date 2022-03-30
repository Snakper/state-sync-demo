#mmo手游点击移动下的状态同步实现
1. 模拟网络延迟
   define.go：
   var lag = time.Millisecond * 1
2. 服务器帧率
   main.go：
   svr := NewServer(10)
3. 客户端帧率固定60
4. 客户端预测开启
   main.go：
   g.OpenForecast()
5. 对账开启
   main.go：
   g.OpenReconciliation()
