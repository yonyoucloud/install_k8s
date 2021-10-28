<template>
  <div>
    <span>开通租户</span>
    <el-divider></el-divider>
    <el-form ref="form" :model="tenantForm" label-width="100px">
      <el-form-item label="租户ID">
        <el-input v-model="tenantForm.tenantID" :disabled="true"></el-input>
      </el-form-item>
      <el-form-item label="租户名称">
        <el-input v-model="tenantForm.tenantName" :disabled="true"></el-input>
      </el-form-item>
      <el-form-item v-if="podID==0" label="选择Pod资源">
        <el-popover
          placement="right"
          width="400"
          trigger="manual"
          v-model="podVisible">
          <el-form label-position="right" label-width="100px">
            <el-form-item label="ID">
              <el-input disabled v-model="pod.ID"></el-input>
            </el-form-item>
            <el-form-item label="名称">
              <el-input disabled v-model="pod.Name"></el-input>
            </el-form-item>
            <el-form-item label="代码号">
              <el-input disabled v-model="pod.Code"></el-input>
            </el-form-item>
            <el-form-item label="根域名">
              <el-input disabled v-model="pod.Domain"></el-input>
            </el-form-item>
            <el-form-item label="容量">
              <el-input disabled v-model="pod.Cap"></el-input>
            </el-form-item>
            <el-form-item label="IaaS服务商">
              <el-input disabled v-model="pod.Iaas"></el-input>
            </el-form-item>
          </el-form>
          <el-select slot="reference" v-model="tenantForm.podID" @change="changePod" @blur="podVisible=false" placeholder="请选择Pod资源">
            <el-option v-for="(pod, index) in podData" :key="index" :label="pod.Name" :value="pod.ID"></el-option>
          </el-select>
        </el-popover>
        <el-alert
          title="特别提醒"
          type="warning"
          description="开通后不能更改对应Pod资源"
          :closable=false
          show-icon>
        </el-alert>
      </el-form-item>
      <el-form-item v-if="podID==0">
        <el-button type="primary" :loading="openLoading" @click="openTenant">开通</el-button>
      </el-form-item>
    </el-form>
    <template v-if="podID > 0 && podData[podID] !== undefined">
      <div class="address">
        <el-alert
          :title="'访问地址：' + podData[podID].Code + '.' + podData[podID].Domain"
          type="success"
          show-icon
          :closable=false>
        </el-alert>
      </div>
      <el-divider>资源信息</el-divider>
      <el-tabs v-if="podID>0" v-model="editModel" type="border-card">
        <el-tab-pane key="json" label="JSON" name="json">
          <keep-alive>
            <AceEditor :text.sync="textJson" @handleTextChange="handleTextChange" lang="json" height="800px" />
          </keep-alive>
        </el-tab-pane>
        <el-tab-pane key="yaml" label="YAML" name="yaml">
          <keep-alive>
            <AceEditor :text.sync="textYaml" @handleTextChange="handleTextChange" lang="yaml" height="800px" />
          </keep-alive>
        </el-tab-pane>
      </el-tabs>
    </template>
  </div>
</template>

<script>
import request from '@/tools/request'
import moment from 'moment'
import AceEditor from './AceEditor.vue'

const TENANTS = {
  1: '用友网络科技股份有限公司',
  2: '畅捷通信息技术有限公司',
  3: '字节跳动科技有限公司',
}

const IAASES = {
  'aliyun': '阿里云',
  'huaweiyun': '华为云',
  'tencent': '腾讯云',
  'amazon': '亚马逊'
}

export default {
  name: 'Pod',
  components: {
    AceEditor
  },
  watch: {
    data: {
      deep: true,
      handler(val) {
        this.textJson = val
        this.textYaml = val
      }
    }
  },
  data() {
    return {
      tenantID: 0,
      podID: 0,
      tenantForm: {
        tenantID: 0,
        podID: '',
        tenantName: '',
      },
      podData: {},
      podVisible: false,
      pod: {},
      openLoading: false,
      editModel: 'json',
      data: {
        'Pod资源': {},
        'K8S集群': {},
        'Mysql集群': [],
        'Redis集群': [],
      },
      textYaml: {},
      textJson: {},
    }
  },

  methods: {
    listPod() {
      request({
        url: 'api/v1/pod/list',
        method: 'get'
      }).then(response => {
        if (response.Code === 10000) {
          this.podData = {}
          for (let i = 0; i < response.Data.length; i++) {
            this.podData[response.Data[i].ID] = this.setPod(response.Data[i])
          }
        }
      })
    },

    setPod(pod) {
      pod.Iaas = IAASES[pod.Iaas]
      pod.CreatedAt = moment(pod.CreatedAt).format('YYYY-MM-DD HH:mm:ss')
      pod.UpdatedAt = moment(pod.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
      return pod
    },

    changePod(podID) {
      this.pod = this.podData[podID]
      this.podVisible = true
    },

    openTenant() {
      if (this.tenantForm.podID === '') {
        this.$message({
          message: '请选择Pod资源',
          type: 'warning'
        })
        return
      }

      this.$confirm('开通后不能更改对应Pod资源, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.openLoading = true
        let formData = new FormData()
        for (let key in this.tenantForm) {
          formData.append(key, this.tenantForm[key])
        }
        request({
          url: 'api/v1/tenantPod/open',
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
              message: '开通成功，ID为：' + response.Data.ID,
              type: 'success'
            })
            this.getByTenantID(this.tenantID)
          } else {
            this.$message({
              message: response.Msg,
              type: 'error'
            })
          }
          this.openLoading = false
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })       
      })
    },

    getByTenantID(tenantID) {
      request({
        url: 'api/v1/tenantPod/getByTenantID/' + tenantID,
        method: 'get'
      }).then(response => {
        if (response.Code === 10000 && response.Data.ID > 0) {
          this.podID = response.Data.PodID
          this.getPod(this.podID)
        }
      })
    },

    getPod(podID) {
      request({
        url: 'api/v1/pod/get/' + podID,
        method: 'get'
      }).then(response => {
        if (response.Code === 10000 && response.Data.ID > 0) {
          this.data['Pod资源'] = this.setPod(response.Data)
          this.getK8sCluster(response.Data.K8sClusterID)
          this.getListResource(podID)
        }
      })
    },

    setK8sCluster(k8sCluster) {
      k8sCluster.CreatedAt = moment(k8sCluster.CreatedAt).format('YYYY-MM-DD HH:mm:ss')
      k8sCluster.UpdatedAt = moment(k8sCluster.UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
      return k8sCluster
    },

    getK8sCluster(k8sClusterID) {
      request({
        url: 'api/v1/k8sCluster/get/' + k8sClusterID,
        method: 'get'
      }).then(response => {
        if (response.Code === 10000 && response.Data.ID > 0) {
          this.data['K8S集群'] = this.setK8sCluster(response.Data)
        }
      })
    },

    getListResource(podID) {
      request({
        url: 'api/v1/podResource/listResource/' + podID,
        method: 'get',
      }).then(response => {
        if (response.Code === 10000) {
          let mysql = [],
              redis = []
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i].CreatedAt = moment(response.Data[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss')
            response.Data[i].UpdatedAt = moment(response.Data[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
            switch (response.Data[i].Category) {
              case 'mysql':
                mysql.push(response.Data[i])
                break
              case 'redis':
                redis.push(response.Data[i])
                break
            }
          }
          this.data['Mysql集群'] = mysql
          this.data['Redis集群'] = redis
        }
      })
    },

    handleTextChange(lang, val) {
      switch (lang) {
        case 'yaml':
          this.textYaml = val
          break
        case 'json':
          this.textJson = val
          break
      }
    },
  },

  mounted() {
    this.tenantID = this.$route.query.tenant_id || 1
    this.tenantForm.tenantID = this.tenantID
    this.tenantForm.tenantName = TENANTS[this.tenantForm.tenantID]
    this.getByTenantID(this.tenantID)
    this.listPod()
  },
}
</script>

<style>
.el-form-item__content.el-form-item__content {
  line-height: 16px !important;
}
.el-alert .el-alert__title {
  font-size: 22px;
  line-height: 36px;
}
</style>
