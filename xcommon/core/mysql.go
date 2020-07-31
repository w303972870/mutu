package mtcore
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type IMysql interface {
    Connect( dbName string )
}

type IMysqlTable interface {
}

type Mysql struct {
    DbName string
    Host string
    Db * sql.DB
}

func( m * Mysql ) Connect( dbName string ) {
    if m.Db == nil {
        m.DbName = LibConfigParms.Configs.GetDatabase( dbName , ConfigDbNameKey ).(string)
        m.Host = LibConfigParms.Configs.GetDatabase( dbName ,ConfigDbHostKey ).(string)
        var err error
        m.Db , err = sql.Open( "mysql" , MtTools.Str( 
          LibConfigParms.Configs.GetDatabase( dbName , ConfigDbUserKey ).(string) , 
          ":", 
          LibConfigParms.Configs.GetDatabase( dbName , ConfigDbPwdKey ).(string) , 
          "@tcp(",
          LibConfigParms.Configs.GetDatabase( dbName , ConfigDbHostKey ).(string) , 
          ":",
          MtTools.Int2Str( LibConfigParms.Configs.GetDatabase( dbName , ConfigDbPortKey ).(int) ) , 
          ")/" , 
          LibConfigParms.Configs.GetDatabase( dbName , ConfigDbNameKey ).(string) , 
          "?charset=utf8&parseTime=true&loc=Local" ) )
        //m.db.SetConnMaxLifetime( 0 ) 
        //m.db.SetMaxIdleConns( 2 )
        if err != nil {
            MutuLogs.Error( MtTools.Str( "数据库链接失败：", dbName ) )
        }
        //关闭ONLY_FULL_GROUP_BY
        m.Db.Exec("SET sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")

        MutuLogs.Sys( MtTools.Str( "创建数据库链接成功：", dbName ) )
    }
}

