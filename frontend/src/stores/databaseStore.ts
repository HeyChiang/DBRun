import { defineStore } from 'pinia';
import { ref, Ref } from 'vue';
import { GetAllCredentials, DeleteCredentials, InsertCredentials,  UpdateCredentials } from '@/../wailsjs/go/api/SQLiteAPI';
import { ListDatabasesByConfig, SetTableVOCacheByTableID } from '@/../wailsjs/go/api/MetadatasAPI';
import { metadata, service, models } from '@/../wailsjs/go/models';


// 扩展DBCredentials类型，添加dbs和tables属性
interface EnhancedDBCredentials extends metadata.Credentials {
  dbs?: models.DatabaseInfoVO[];
  display?: models.Display;
}

export const useDatabaseStore = defineStore('database', () => {
  const dbLinks: Ref<EnhancedDBCredentials[]> = ref([]);

   const addDatabase = async (dbLink: EnhancedDBCredentials): Promise<void> => {
     try {
       const credentials:metadata.Credentials =  await InsertCredentials(dbLink);
       
       console.log("添加数据库成功", credentials);

       dbLink.id = credentials.id;

       // 获取数据库表信息
       const dbInfo: models.DBInfoVO = await ListDatabasesByConfig(dbLink);
       
       dbLink.dbs = dbInfo.dbs;
       dbLinks.value.push(dbLink);

     } catch (err) {
       console.error("添加数据库错误:", err);
       throw err;
     }
   };

   const removeDatabase = async (id: number): Promise<void> => {
     try {
       await DeleteCredentials(id);
       dbLinks.value = dbLinks.value.filter(db => db.id !== id);
     } catch (err) {
       console.error("删除数据库错误:", err);
       throw err;
     }
   };

   const refreshDatabases = async (): Promise<EnhancedDBCredentials[] | undefined> => {
     try {
       const allCredentials: EnhancedDBCredentials[] = await GetAllCredentials();
       
       if(!allCredentials){
         return allCredentials;
       }

       dbLinks.value = allCredentials;

       // 创建所有请求的 Promise 数组
       const promises = allCredentials.map(link => {
         return ListDatabasesByConfig(link)
           .then(dbInfo => {
             link.dbs = dbInfo.dbs;
             console.log("完成获取数据表信息", link);
           })
           .catch(err => {
             console.error(`获取数据表信息错误 for link ${link.id}:`, err);
             link.dbs = []; // 设置空数组表示获取失败
           });
       });

       // 并行执行所有请求
       await Promise.all(promises);
       
       return dbLinks.value;
     } catch (err) {
       console.error("刷新数据库信息错误:", err);
       throw err;
     }
   };


   const updateDBCredentials = async (db: EnhancedDBCredentials): Promise<void> => {
     try {
       await UpdateCredentials(db);
       const index = dbLinks.value.findIndex(item => item.id === db.id);
       if (index !== -1) {
         // Get updated database info
         const dbInfo = await ListDatabasesByConfig(db);
         db.dbs = dbInfo.dbs;
         
         // 创建新的EnhancedDBCredentials对象
         const updatedDb = metadata.Credentials.createFrom(db) as EnhancedDBCredentials;
         updatedDb.dbs = db.dbs;
         
         // Update the database in the store
         dbLinks.value[index] = updatedDb;
       }
     } catch (err) {
       console.error("更新数据库错误:", err);
       throw err;
     }
   };


  const updateTableRemark = async (
    tableInfo: models.TableInfoVO,
    tableId?: number
  ): Promise<void> => {
    const tableCacheVO: service.TableCacheVO = {
      remark: tableInfo.remark || '',
      version: 0,
      fieldMap: {}
    };

    if (tableInfo.fields && tableInfo.fields.length > 0) {
      for (const f of tableInfo.fields) {
        if (f.remark && f.name) {
          tableCacheVO.fieldMap![f.name] = f.remark;
        }
      }
    }

    const id = tableId ?? (tableInfo as any)?.id;
    if (!id) {
      console.warn('updateTableRemark skipped: tableId/id not provided.');
      return;
    }
    await SetTableVOCacheByTableID(id, tableCacheVO);
  };

  return {
    dbLinks,
    addDatabase,
    removeDatabase,
    refreshDatabases,
    updateDBCredentials,
    updateTableRemark
  };
});
