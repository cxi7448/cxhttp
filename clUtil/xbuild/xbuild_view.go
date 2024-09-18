package xbuild

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clFile"
	"os"
)

func BuildView(table, path, lang string) error {
	if path == "" {
		path = "src/views"
	}
	islang := lang == "1"
	info := GenTable(table)
	if info == nil {
		return fmt.Errorf("生成表格数据失败!")
	}
	js_script := fmt.Sprintf(`<script setup lang="ts">
import { getCurrentInstance, onMounted} from 'vue'
import {ElMessage} from "element-plus";
%v
const search:any = ref({})
const proxy = getCurrentInstance()?.proxy
const tableData:any = ref([])
const getList = ()=>{
  let param = search.value
  param.pageid = pager.value.pageid -1
  param.pcount = pager.value.pcount
  proxy?.$http.post("%v_list",param).then((res:any)=>{
    if (res.code == 0){
      tableData.value = res.data.list
      if (pager.value.pageid == 1){
        pager.value.total = res.data.total
      }
    }else{
      ElMessage({message:res.msg,type:"error"})
    }
  })
}
const dialogPop:any = ref(false)
const dialogTitle:any = ref("编辑")
const isAdd = ref(true)
const form:any = ref({%v})
const showDialog = (row:any)=>{
  form.value = {
	%v
  }
  dialogTitle.value = "添加"
  if (row){
    isAdd.value = false
    form.value = row
    dialogTitle.value = "编辑:"+row.id
  }else{
    isAdd.value = true
  }
  dialogPop.value = true
}

const onSubmit = ()=>{
  let ac = "%v_add"
  if (!isAdd.value){
    ac = "%v_edit"
  }
  proxy?.$http.post(ac,form.value).then((res:any)=>{
    if (res.code == 0){
      dialogPop.value = false
      getList()
    }else{
      ElMessage({message:res.msg,type:"error"})
    }
  })
}

const selectRows:any = ref([])
const selectChange = (rows)=>{
  selectRows.value = rows
}
const onDeleteMulti = ()=>{
	if(selectRows.value.length == 0){
		 ElMessage({message:"请选择要删除的记录",type:"error"})
	}
	let ids = []
	for(const row of selectRows.value){
		ids.push(row.id)
	}
	onDelete(ids.join(","))
}
const onDelete = (ids)=>{
  let ac = "%v_delete"
  proxy?.$http.post(ac,{ids:ids}).then((res:any)=>{
    if (res.code == 0){
      getList()
    }else{
      ElMessage({message:res.msg,type:"error"})
    }
  })
}
onMounted(()=>{
  getInit()
  getList()
})

const getInit = ()=>{
  proxy?.$http.post("getInit", {table:"%v"}).then((res:any)=>{
    if (res.code == 0){
      console.log(res.data)
    }else{
      ElMessage({message:res.msg,type:"error"})
    }
  })
}

const pager:any = ref({
  pageid:1,
  pcount:10,
  total:0
})
</script>
`, info.GenScript(islang), info.Name, info.GetFormStr(islang), info.GetFormStr(islang), info.Name, info.Name, info.Name, info.Name)
	html_template := fmt.Sprintf(`
<template>
  <div class="search_box">
    <el-button @click="getList">查询</el-button>
	<el-button @click="showDialog()">添加</el-button>
	<el-popconfirm v-if="selectRows.length > 0" title="确认删除?" @confirm="onDeleteMulti()">
      <template #reference>
        <el-button type="danger" icon="Delete">批量删除</el-button>
      </template>
    </el-popconfirm>
  </div>
  <el-table @selectionChange="selectChange" :data="tableData" :header-cell-style="{ background: '#F7F8FA' }">
	<el-table-column type="selection" width="55" />
   	%v
    <el-table-column label="操作" align="left" width="150px" :show-overflow-tooltip="true">
      <template #default="scope">
        <el-button-group>
          <el-button type="primary" @click="showDialog(scope.row)" icon="Edit"></el-button>
			<el-popconfirm title="确认删除?" @confirm="onDelete(scope.row.id)">
            <template #reference>
              <el-button type="danger" icon="Delete"></el-button>
            </template>
          </el-popconfirm>
        </el-button-group>
      </template>
    </el-table-column>
  </el-table>
  <Pager :pager="pager" @query="getList"></Pager>

  <el-dialog v-model="dialogPop" width="700px" :title="dialogTitle">
    <el-form label-width="100px">
      %v
      <el-form-item label="" ><el-button @click="onSubmit">提交</el-button><el-button @click="dialogPop = false">取消</el-button></el-form-item>
    </el-form>
  </el-dialog>

</template>
`, info.ElTableColumn(), info.ElFormItem(islang))
	var content = fmt.Sprintf(" %v \n %v", html_template, js_script)
	folder := fmt.Sprintf("%v/%v", path, info.Name)
	os.MkdirAll(folder, 0700)
	vue_file := fmt.Sprintf("%v/index.vue", folder)
	// 创建模型文件
	if !clFile.IsFile(vue_file) {
		// 自动生成，存在就不生成了
		os.WriteFile(vue_file, []byte(content), 0700)
	}
	return nil
}
