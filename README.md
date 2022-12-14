```
├─.idea
│  │  .gitignore
│  │  dataSources.local.xml
│  │  dataSources.xml
│  │  modules.xml
│  │  user-center.iml
│  │  workspace.xml
│  │  
│  └─dataSources
│      │  b7a7de30-85f6-4de8-90f1-706e2dd5b366.xml
│      │  
│      └─b7a7de30-85f6-4de8-90f1-706e2dd5b366
│          └─storage_v2
│              └─_src_
│                  └─schema
│                          information_schema.FNRwLQ.meta
│                          mysql.osA4Bg.meta
│                          performance_schema.kIw0nw.meta
│                          sys.zb4BAA.meta
│                          
└─server
    ├─api
    │      interceptor.go
    │      issue.go
    │      user.go
    │      
    ├─cache
    │      redis.go
    │      
    ├─cmd
    │      main.go
    │      
    ├─conf
    │      conf.go
    │      config.ini
    │      
    ├─dao
    │      dao.go
    │      user.go
    │      
    ├─model
    │      token.go
    │      user.go
    │      
    ├─pbfile
    │  ├─cert
    │  │      ca.crt
    │  │      ca.csr
    │  │      ca.key
    │  │      openssl.cnf
    │  │      server.csr
    │  │      server.key
    │  │      server.pem
    │  │      
    │  ├─pb
    │  │      user.pb.go
    │  │      user_grpc.pb.go
    │  │      
    │  └─proto
    │          user.proto
    │          
    ├─service
    │      token.go
    │      user.go
    │      
    └─tool
            check.go
            check_test.go
            resp.go
            trie.go
            
```