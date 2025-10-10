// 行与行可以乱序，行内不能截断
// 性能尽可能好
// 程序可以正常退出
```go
input:=getData()// slice
data:=make(chan,1000000)
info:=make([]string,1000)
func input(id int64){
    for(i:=id;i<=id+1000;i++){
        data<-input[i]
    }
}
func main(){
	ans:=""
	cg:=sync.Group{}
    for i:=0;i<=len(input);i+=1000{
		cg.Add(1)
        go func(){
			defer cg.Done()
            input(i)
        }   
    }
    cg1:=sync.Group{}
	for i:=0;i<=1000;i++{
		cg1.Add(1)
	    go func(){
            defer cg1.Done()
			for{
			    select{
				    case info<-data:
					info[i]+=info
                }   	
            }       
        }       	
    }
    cg.Wait()
	close(data)
	cg1.Wait()
	for i:=0;i<1000;i++{
		ans+=info[i]
    }   
	log.Info(ans)
}

```