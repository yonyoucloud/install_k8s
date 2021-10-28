<template>
  <div class="resource">
    <div style="display: inline-block;">Resource列表</div>
    <div style="display: inline-block;float:right;">
      <el-button type="primary" @click="resourceForm = resourceFormInit;createResourceDialog = true">创建</el-button>
    </div>
    <el-divider></el-divider>
    <el-table
      :data="tableResourceData.filter(data => !search || data.Name.toLowerCase().includes(search.toLowerCase()) || data.Category.toLowerCase().includes(search.toLowerCase()) || data.Host.toLowerCase().includes(search.toLowerCase()))"
      style="width: 100%;"
      :row-class-name="tableRowClassName">
      <el-table-column
        fixed
        prop="ID"
        label="ID"
        align="center"
        width="50">
      </el-table-column>
      <el-table-column
        fixed
        prop="Name"
        label="名称"
        align="center"
        width="150">
      </el-table-column>
      <el-table-column
        prop="Category"
        label="资源类别"
        align="center"
        width="150">
      </el-table-column>
      <el-table-column
        prop="Scope"
        label="特定描述"
        align="center"
        width="100">
      </el-table-column>
      <el-table-column
        prop="Host"
        label="主机地址"
        align="center"
        width="120">
      </el-table-column>
      <el-table-column
        prop="Port"
        label="端口号"
        align="center"
        width="80">
      </el-table-column>
      <el-table-column
        prop="User"
        label="用户名"
        align="center"
        width="100">
      </el-table-column>
      <el-table-column
        prop="Password"
        label="密码"
        align="center"
        width="100">
      </el-table-column>
      <el-table-column
        prop="CreatedAt"
        label="创建时间"
        align="center"
        width="150">
      </el-table-column>
      <el-table-column
        prop="UpdatedAt"
        label="更新时间"
        align="center"
        width="150">
      </el-table-column>
      <el-table-column
        fixed="right"
        align="center"
        width="150">
        <!-- eslint-disable-next-line -->
        <template slot="header" slot-scope="scope">
          <el-input
            v-model="search"
            size="mini"
            placeholder="输入关键字搜索"/>
        </template>
        <template slot-scope="scope">
          <el-button
            size="mini"
            @click="handleEdit(scope.$index, scope.row)">编辑</el-button>
          <el-button
            size="mini"
            type="danger"
            @click="handleDelete(scope.$index, scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog
      :title="(resourceForm.ID?'编辑':'创建')+'Resource'"
      :visible.sync="createResourceDialog"
      width="40%">
      <el-form ref="form" :model="resourceForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="resourceForm.Name"></el-input>
        </el-form-item>
        <el-form-item label="资源类别">
          <el-select v-model="resourceForm.Category" placeholder="请选择资源类别">
            <el-option label="vps" value="vps"></el-option>
            <el-option label="mysql" value="mysql"></el-option>
            <el-option label="redis" value="redis"></el-option>
            <el-option label="mongodb" value="mongodb"></el-option>
            <el-option label="rabbitmq" value="rabbitmq"></el-option>
            <el-option label="elasticsearch" value="elasticsearch"></el-option>
            <el-option label="kafka" value="kafka"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="特定描述">
          <el-tooltip placement="top">
            <div slot="content">特定资源的特定描述，如mysql的话，可以标明是主库或是从库</div>
            <el-select v-model="resourceForm.Scope" placeholder="请选择特定描述">
              <el-option label="default" value="default"></el-option>
              <el-option label="master" value="master"></el-option>
              <el-option label="sources" value="sources"></el-option>
              <el-option label="replicas" value="replicas"></el-option>
              <el-option label="publish" value="publish"></el-option>
              <el-option label="node" value="node"></el-option>
              <el-option label="etcd" value="etcd"></el-option>
              <el-option label="etcdlb" value="etcdlb"></el-option>
              <el-option label="masterlb" value="masterlb"></el-option>
              <el-option label="lvs" value="lvs"></el-option>
              <el-option label="pridocker" value="pridocker"></el-option>
              <el-option label="pridns" value="pridns"></el-option>
              <el-option label="newnode" value="newnode"></el-option>
              <el-option label="newetcd" value="newetcd"></el-option>
              <el-option label="newmaster" value="newmaster"></el-option>
            </el-select>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="主机地址">
          <el-input v-model="resourceForm.Host"></el-input>
        </el-form-item>
        <el-form-item label="端口号">
          <el-input v-model="resourceForm.Port"></el-input>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="resourceForm.User"></el-input>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="resourceForm.Password"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="createResourceDialog = false">取消</el-button>
        <el-button type="primary" @click="createResource(resourceForm.ID)">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/tools/request'
import moment from 'moment'

export default {
  name: 'Resource',
  components: {
    // HelloWorld
  },

  data() {
    return {
      resourceForm: {},
      resourceFormInit: {
        Name: '',
        Category: 'vps',
        Scope: 'default',
        Host: '',
        Port: 0,
        User: '',
        Password: ''
      },
      tableResourceData: [],
      search: '',
      createResourceDialog: false,
    }
  },

  methods: {
    createResource(id) {
      if (id) {
        this.editResource(id)
        return
      }

      let formData = new FormData()
      for (let key in this.resourceForm) {
        formData.append(key, this.resourceForm[key])
      }
      request({
        url: 'api/v1/resource/create',
        method: 'post',
        header: {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        },
        data: formData
      }).then(response => {
        if (response.Code === 10000) {
          this.$message({
            message: '创建成功，ID为：' + response.Data.ID,
            type: 'success'
          })
          this.listResource()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createResourceDialog = false
      })
    },

    listResource() {
      request({
        url: 'api/v1/resource/list',
        method: 'get'
      }).then(response => {
        if (response.Code === 10000) {
          this.tableResourceData = []
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i].CreatedAt = moment(response.Data[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss')
            response.Data[i].UpdatedAt = moment(response.Data[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
            this.tableResourceData.push(response.Data[i])
          }
        }
      })
    },

    tableRowClassName({row, rowIndex}) {
      row
      if (rowIndex % 2 === 0) {
        return 'warning-row'
      } else if (rowIndex % 2 === 1) {
        return 'success-row'
      }
      return ''
    },

    handleEdit(index, row) {
      this.resourceForm = Object.assign({}, row)
      this.createResourceDialog = true
    },

    editResource(id) {
      let formData = new FormData()
      formData.append('Name', this.resourceForm.Name)
      formData.append('Category', this.resourceForm.Category)
      formData.append('Scope', this.resourceForm.Scope)
      formData.append('Host', this.resourceForm.Host)
      formData.append('Port', this.resourceForm.Port)
      formData.append('User', this.resourceForm.User)
      formData.append('Password', this.resourceForm.Password)
      request({
        url: 'api/v1/resource/edit/' + id,
        method: 'post',
        header: {
          headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
          }
        },
        data: formData
      }).then(response => {
        if (response.Code === 10000) {
          this.$message({
            message: '更新成功',
            type: 'success'
          })
          this.listResource()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createResourceDialog = false
      })
    },

    handleDelete(index, row) {
      this.deleteResource(row.ID)
    },

    handleClose(done) {
      this.$confirm('确认关闭？')
        // eslint-disable-next-line
        .then(_ => {
          done()
        })
        // eslint-disable-next-line
        .catch(_ => {})
    },

    deleteResource(id) {
      this.$confirm('此操作将永久删除该Resource, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        request({
          url: 'api/v1/resource/delete/' + id,
          method: 'delete'
        }).then(response => {
          if (response.Code === 10000) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.listResource()
          } else {
            this.$message({
              message: response.Msg,
              type: 'error'
            })
          }
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })       
      })
    },
  },

  mounted() {
    this.listResource()
  },
}
</script>

<style>
.resource {
  margin-bottom: 50px;
}
.resource .el-table .warning-row {
  background: lavenderblush;
}
.resource .el-table .success-row {
  background: ivory;
}
</style>
