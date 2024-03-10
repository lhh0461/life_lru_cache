# LifeLru

## 介绍 
通过泛形实现的支持key过期时间LruCache，当Cache容量满时，会找到一个过期的key剔除，如果都没过期，则找到最近没使用的来剔除。

## 使用方法
```
    cache := NewLRUCache[int, int](3)

    //设置数据
    cache.Set(1, 1, time.Minute)

    //获取数据
    value, ok := cache.Get(1)

```