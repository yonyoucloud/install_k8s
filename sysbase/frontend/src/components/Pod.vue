<template>
  <div class="pod">
    <div style="display: inline-block;">Pod列表</div>
    <div style="display: inline-block;float:right;">
      <el-button type="primary" @click="podForm = podFormInit;listK8sCluster();resourceListPod();createPodDialog = true">创建</el-button>
    </div>
    <el-divider></el-divider>
    <el-table
      :data="tablePodData.filter(data => !search || data.Name.toLowerCase().includes(search.toLowerCase()) || data.Code.toLowerCase().includes(search.toLowerCase()))"
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
        width="180">
      </el-table-column>
      <el-table-column
        prop="Code"
        label="代码号"
        align="center"
        width="100">
      </el-table-column>
      <el-table-column
        prop="Domain"
        label="根域名"
        align="center"
        width="100">
      </el-table-column>
      <el-table-column
        prop="Cap"
        label="容量"
        align="center"
        width="80">
      </el-table-column>
      <el-table-column
        prop="Iaas"
        label="IaaS服务商"
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
        width="240">
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
            @click="handleShowListResource(scope.$index, scope.row)">查看资源</el-button>
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
      :title="(podForm.ID?'编辑':'创建')+'Pod'"
      :visible.sync="createPodDialog"
      width="40%">
      <el-form ref="form" :model="podForm" label-width="100px">
        <el-form-item label="名称">
          <el-tooltip placement="top">
            <div slot="content">Pod名，规范：中国大陆-阿里云-100</div>
            <el-input v-model="podForm.Name"></el-input>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="代码号">
          <el-tooltip placement="top">
            <div slot="content">Pod的代码号，作为二级域名，和domain字段作为Pod的唯一访问入口</div>
            <el-input v-model="podForm.Code"></el-input>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="根域名">
          <el-tooltip placement="top">
            <div slot="content">Pod的根域名，和code字段作为Pod的唯一访问入口</div>
            <el-input v-model="podForm.Domain"></el-input>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="容量">
          <el-tooltip placement="top">
            <div slot="content">Pod可以容纳的租户个数</div>
            <el-input v-model="podForm.Cap" type="number"></el-input>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="IaaS服务商">
          <el-radio-group v-model="podForm.Iaas">
            <el-radio label="aliyun">阿里云</el-radio>
            <el-radio label="huaweiyun">华为云</el-radio>
            <el-radio label="tencent">腾讯云</el-radio>
            <el-radio label="amazon">亚马逊</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="选择K8S集群">
          <el-tooltip placement="top">
            <div slot="content">一个Pod只能对应一个K8S集群，一个K8S集群可用对应多个Pod</div>
            <el-select v-model="podForm.K8sClusterID" filterable placeholder="请选择">
              <el-option
                v-for="item in k8sClusterList"
                :key="item.ID"
                :label="item.Name"
                :value="item.ID">
              </el-option>
            </el-select>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="添加资源">
          <el-tooltip placement="top">
            <div slot="content">一个Pod能添加多个非vps资源，非vps资源可以被多个Pod共用</div>
            <el-select v-model="podForm.ResourceID" multiple filterable placeholder="请选择">
              <el-option
                v-for="item in resourceList"
                :key="item.ID"
                :label="item.ID + '-' + item.Name + '-' + item.Category + '-' + item.Scope + '-' + item.Host"
                :value="item.ID">
              </el-option>
            </el-select>
          </el-tooltip>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="createPodDialog = false">取消</el-button>
        <el-button type="primary" @click="createPod(podForm.ID)">确定</el-button>
      </span>
    </el-dialog>
    <el-dialog
      :title="(name?name:'')+'-资源信息'"
      :visible.sync="listResourceDialog"
      width="85%">
      <div class="resource">
        <el-table
        :data="tableListResourceData"
        style="width: 100%;"
        :row-class-name="tableRowClassName">
          <el-table-column
            prop="ID"
            label="ID"
            align="center"
            width="50">
          </el-table-column>
          <el-table-column
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
            align="center">
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/tools/request'
import moment from 'moment'

const IAASES = {
  'aliyun': '阿里云',
  'huaweiyun': '华为云',
  'tencent': '腾讯云',
  'amazon': '亚马逊'
}

const IAASES_BACK = {
  '阿里云': 'aliyun',
  '华为云': 'huaweiyun',
  '腾讯云': 'tencent',
  '亚马逊': 'amazon',
}

export default {
  name: 'Pod',
  components: {
    // HelloWorld
  },

  data() {
    return {
      podForm: {},
      podFormInit: {
        Name: '',
        Code: '',
        Domain: '',
        Cap: 100,
        Iaas: 'aliyun',
      },
      tablePodData: [],
      search: '',
      createPodDialog: false,
      k8sClusterList: [],
      resourceList: [],
      listResourceDialog: false,
      tableListResourceData: [],
      name: '',
    }
  },

  methods: {
    createPod(id) {
      if (id) {
        this.editPod(id)
        return
      }

      let formData = new FormData()
      for (let key in this.podForm) {
        formData.append(key, this.podForm[key])
      }
      request({
        url: 'api/v1/pod/create',
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
          this.listPod()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createPodDialog = false
      })
    },

    listPod() {
      request({
        url: 'api/v1/pod/list',
        method: 'get'
      }).then(response => {
        if (response.Code === 10000) {
          this.tablePodData = []
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i].Iaas = IAASES[response.Data[i].Iaas]
            response.Data[i].CreatedAt = moment(response.Data[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss')
            response.Data[i].UpdatedAt = moment(response.Data[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
            this.tablePodData.push(response.Data[i])
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
      this.podForm = Object.assign({}, row)
      this.podForm.Iaas = IAASES_BACK[row.Iaas]
      this.createPodDialog = true
      this.listK8sCluster()
      this.resourceListPod(row.ID)
    },

    editPod(id) {
      let formData = new FormData()
      formData.append('Name', this.podForm.Name)
      formData.append('Code', this.podForm.Code)
      formData.append('K8sClusterID', this.podForm.K8sClusterID)
      formData.append('Domain', this.podForm.Domain)
      formData.append('Cap', this.podForm.Cap)
      formData.append('Iaas', this.podForm.Iaas)
      if (this.podForm.ResourceID !== undefined) {
        formData.append('ResourceID', this.podForm.ResourceID)
      }
      request({
        url: 'api/v1/pod/edit/' + id,
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
          this.listPod()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createPodDialog = false
      })
    },

    handleDelete(index, row) {
      this.deletePod(row.ID)
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

    deletePod(id) {
      this.$confirm('此操作将永久删除该Pod, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        request({
          url: 'api/v1/pod/delete/' + id,
          method: 'delete'
        }).then(response => {
          if (response.Code === 10000) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.listPod()
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

    listK8sCluster() {
      request({
        url: 'api/v1/k8sCluster/list',
        method: 'get'
      }).then(response => {
        if (response.Code === 10000) {
          this.k8sClusterList = []
          for (let i = 0; i < response.Data.length; i++) {
            this.k8sClusterList.push({
              ID: response.Data[i].ID,
              Name: response.Data[i].Name
            })
          }
        }
      })
    },

    resourceListPod(podID) {
      let query = {
        podID: podID
      }
      request({
        url: 'api/v1/resource/list/pod',
        method: 'get',
        params: query
      }).then(response => {
        if (response.Code === 10000) {
          this.resourceList = response.Data
        }
      })
    },

    handleShowListResource(index, row) {
      this.name = row.Name
      this.listResourceDialog = true
      request({
        url: 'api/v1/podResource/listResource/' + row.ID,
        method: 'get',
      }).then(response => {
        if (response.Code === 10000) {
          this.tableListResourceData = []
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i].CreatedAt = moment(response.Data[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss')
            response.Data[i].UpdatedAt = moment(response.Data[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
            this.tableListResourceData.push(response.Data[i])
          }
        }
      })
    }
  },

  mounted() {
    this.listPod()
  },
}
</script>

<style>
.pod .el-table .warning-row {
  background: oldlace;
}
.pod .el-table .success-row {
  background: #f0f9eb;
}
.resource .el-table .warning-row {
  background: lavenderblush;
}
.resource .el-table .success-row {
  background: ivory;
}
</style>
