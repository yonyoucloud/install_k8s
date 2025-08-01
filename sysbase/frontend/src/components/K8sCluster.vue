<template>
  <div class="k8sCluster">
    <div style="display: inline-block;">K8sCluster列表</div>
    <div style="display: inline-block;float:right;">
      <el-button type="primary" @click="k8sClusterForm = k8sClusterFormInit;resourceListK8sCluster();createK8sClusterDialog = true">创建</el-button>
    </div>
    <el-divider></el-divider>
    <el-table
      :data="tableK8sClusterData.filter(data => !search || data.Name.toLowerCase().includes(search.toLowerCase()))"
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
      :title="(k8sClusterForm.ID?'编辑':'创建')+'K8sCluster'"
      :visible.sync="createK8sClusterDialog"
      width="40%">
      <el-form ref="form" :model="k8sClusterForm" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="k8sClusterForm.Name"></el-input>
        </el-form-item>
        <el-form-item label="添加资源">
          <el-select size="medium" v-model="k8sClusterForm.ResourceID" multiple filterable placeholder="请选择">
            <el-option
              v-for="item in resourceList"
              :key="item.ID"
              :label="item.ID + '-' + item.Name + '-' + item.Category + '-' + item.Host"
              :value="item.ID">
            </el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="createK8sClusterDialog = false">取消</el-button>
        <el-button type="primary" @click="createK8sCluster(k8sClusterForm.ID)">确定</el-button>
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
      <div class="install-k8s-button">
        <font style="font-size:30px;">安装K8s集群：</font><br />
        <el-button type="danger" @click="install('test')">测试</el-button>
        <el-button type="danger" @click="install('all')">一键安装</el-button>
        <el-button type="primary" @click="install('base')">1.基础安装</el-button>
        <!-- <el-button type="primary" @click="install('kernel')">2.升级内核</el-button> -->
        <el-button type="primary" @click="install('dns')">2.安装Dns</el-button>
        <el-button type="primary" @click="install('basebin')">3.安装可执行文件</el-button>
        <el-button type="primary" @click="install('containerd')">4.安装Containerd</el-button>
        <el-button type="primary" @click="install('registry')">5.安装私有镜像仓库</el-button>
        <el-button type="primary" @click="install('etcd')">6.安装Etcd集群</el-button>
        <el-button type="primary" @click="install('master')">7.安装Master集群</el-button>
        <el-button type="primary" @click="install('node')">8.安装Node集群</el-button>
        <el-button type="primary" @click="install('registrycrt')">9.安装私有镜像仓库证书</el-button>
        <el-button type="primary" @click="install('lvs')">10.安装Lvs</el-button>
        <el-button type="primary" @click="install('finish')">11.完成安装</el-button>
        <br />
        <el-button type="primary" @click="install('newnode')">1.安装新节点</el-button>
        <el-button type="primary" @click="install('newetcd')">2.安装新Etcd</el-button>
        <el-button type="primary" @click="install('newmaster')">3.安装新Master</el-button>
        <br />
        <el-button type="primary" @click="install('sslmaster')">1.更新Master证书</el-button>
        <el-button type="primary" @click="install('ssletcd')">2.更新Etcd证书</el-button>
        <el-button type="primary" @click="install('sslnode')">3.更新Node证书</el-button>
      </div>
      <div class="install-k8s-button" style="margin-top:10px;">
        <font style="font-size:30px;">服务管理：</font><br />
        <el-radio-group v-model="service" style="margin-right:30px;">
          <el-radio v-for="(item, key) in serviceData" :key="key" :label="key">{{item}}</el-radio>
        </el-radio-group>
        <font style="font-size:13px;font-weight: bolder;">只针对Node服务：</font>
        <el-radio-group v-model="doContainerd" style="margin-right:30px;">
          <el-radio label="true">包括Containerd</el-radio>
          <el-radio label="false">不包括Containerd</el-radio>
        </el-radio-group>
        <el-button type="primary" @click="install('start')">启动</el-button>
        <el-button type="primary" @click="install('restart')">重启</el-button>
        <el-button type="danger" @click="install('stop')">停止</el-button>
      </div>
    </el-dialog>
    <el-dialog
      :title="(name?name:'')+'-K8S集群安装-'+installLogDesc"
      :visible.sync="installLogDialog"
      :close-on-click-modal=false
      @opened="openedInstallLogDialog"
      width="85%">
      <div id="terminal" class="infinite-list-wrapper" style="overflow:auto;">
      <!-- <div class="infinite-list-wrapper" style="overflow:auto">
        <ul class="list"
          infinite-scroll-disabled="disabled">
          <li v-for="(item, index) in installLog" :key="index" class="list-item">{{ item }}</li>
          <li id="terminal" class="list-item"></li>
          <li class="loading">
            <span v-if="loading" style="color:goldenrod">正在安装...<i class="el-icon-loading"></i></span>
            <span v-if="installLog.length > 0 && !loading" style="color:green">安装完成<i class="el-icon-success"></i></span>
          </li>
        </ul>
      </div> -->
        <!-- <ul class="list"
          infinite-scroll-disabled="disabled">
          <li v-for="(item, index) in installLog" :key="index" class="list-item">{{ item }}</li>
          <li id="terminal" class="list-item"></li>
          <li class="loading">
            <span v-if="loading" style="color:goldenrod">正在安装...<i class="el-icon-loading"></i></span>
            <span v-if="installLog.length > 0 && !loading" style="color:green">安装完成<i class="el-icon-success"></i></span>
          </li>
        </ul> -->
      </div>
    </el-dialog>
  </div>
</template>

<script>
import request from '@/tools/request'
import { stream } from '@/tools/stream'
import moment from 'moment'
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'

export default {
  name: 'K8sCluster',
  components: {
    // HelloWorld
  },

  data() {
    return {
      term: null,
      socket: null,
      k8sClusterForm: {},
      k8sClusterFormInit: {
        Name: '',
      },
      tableK8sClusterData: [],
      search: '',
      createK8sClusterDialog: false,
      resourceList: [],
      listResourceDialog: false,
      tableListResourceData: [],
      id: 0,
      name: '',
      loading: false,
      installLog: [],
      installLogDialog: false,
      installLogDesc: '',
      service: 'servicePublish',
      serviceData: {
        servicePublish: '发布服务',
        serviceEtcd: 'Etcd服务',
        serviceMaster: 'Master服务',
        serviceNode: 'Node服务',
        serviceDns: 'Dns服务',
      },
      doContainerd: "true",
    }
  },

  methods: {
    initXterm() {
      this.term = new Terminal({
        fontSize: 14,
        rendererType: 'canvas', //渲染类型
        rows: 35, //行数
        convertEol: true, //启用时，光标将设置为下一行的开头
        // scrollback: 10, //终端中的回滚量
        disableStdin: true, //是否应禁用输入
        cursorStyle: 'underline', //光标样式
        cursorBlink: false, //光标闪烁
        theme: {
          // foreground: 'yellow', //字体
          // background: '#060101', //背景色
          cursor: 'help' //设置光标
        }
      })

      this.term.open(document.getElementById('terminal'))
      const fitAddon = new FitAddon()
      this.term.loadAddon(fitAddon)
      fitAddon.fit()

      // 支持输入与粘贴方法
      let _this = this; //一定要重新定义一个this，不然this指向会出问题
      this.term.onData(function(key) {
        let order = ['stdin', key] //这里key值是你输入的值，数据格式一定要找后端要！！！！
        _this.socket.onsend(JSON.stringify(order)) //转换为字符串
      })
    },
    init(url) {
      // 实例化socket
      this.socket = new WebSocket(url)
      // 监听socket连接
      this.socket.onopen = this.open
      // 监听socket错误信息
      this.socket.onerror = this.error
      // 监听socket消息
      this.socket.onmessage = this.getMessage
      // 发送socket消息
      this.socket.onsend = this.send
    },
    open: function() {
      console.log("socket连接成功")
      this.initXterm()
    },
    error: function() {
      console.log("连接错误")
    },
    close: function() {
      this.socket.close()
      console.log("socket已经关闭")
    },
    getMessage: function(msg) {
      this.term.write(JSON.parse(msg.data)[1])
    },
    send: function(order) {
      this.socket.send(order)
    },

    createK8sCluster(id) {
      if (id) {
        this.editK8sCluster(id)
        return
      }

      let formData = new FormData()
      for (let key in this.k8sClusterForm) {
        formData.append(key, this.k8sClusterForm[key])
      }
      request({
        url: 'api/v1/k8sCluster/create',
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
          this.listK8sCluster()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createK8sClusterDialog = false
      })
    },

    listK8sCluster() {
      request({
        url: 'api/v1/k8sCluster/list',
        method: 'get'
      }).then(response => {
        if (response.Code === 10000) {
          this.tableK8sClusterData = []
          for (let i = 0; i < response.Data.length; i++) {
            response.Data[i].CreatedAt = moment(response.Data[i].CreatedAt).format('YYYY-MM-DD HH:mm:ss')
            response.Data[i].UpdatedAt = moment(response.Data[i].UpdatedAt).format('YYYY-MM-DD HH:mm:ss')
            this.tableK8sClusterData.push(response.Data[i])
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
      this.resourceListK8sCluster()
      this.k8sClusterForm = Object.assign({}, row)
      this.createK8sClusterDialog = true
    },

    editK8sCluster(id) {
      let formData = new FormData()
      formData.append('Name', this.k8sClusterForm.Name)
      if (this.k8sClusterForm.ResourceID !== undefined) {
        formData.append('ResourceID', this.k8sClusterForm.ResourceID)
      }
      request({
        url: 'api/v1/k8sCluster/edit/' + id,
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
          this.listK8sCluster()
        } else {
          this.$message({
            message: response.Msg,
            type: 'error'
          })
        }
        this.createK8sClusterDialog = false
      })
    },

    handleDelete(index, row) {
      this.deleteK8sCluster(row.ID)
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

    deleteK8sCluster(id) {
      this.$confirm('此操作将永久删除该K8sCluster, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        request({
          url: 'api/v1/k8sCluster/delete/' + id,
          method: 'delete'
        }).then(response => {
          if (response.Code === 10000) {
            this.$message({
              message: '删除成功',
              type: 'success'
            })
            this.listK8sCluster()
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

    resourceListK8sCluster() {
      request({
        url: 'api/v1/resource/list/k8sCluster',
        method: 'get',
      }).then(response => {
        if (response.Code === 10000) {
          this.resourceList = response.Data
        }
      })
    },

    handleShowListResource(index, row) {
      this.id = row.ID
      this.name = row.Name
      this.listResourceDialog = true
      request({
        url: 'api/v1/k8sClusterResource/listResource/' + row.ID,
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
    },

    install(doWhat) {
      let url = '',
          title = ''
      switch(doWhat) {
        case 'test':
          title = '测试'
          url = 'api/v1/installK8s/installTest?k8s_cluster_id=' + this.id
          break
        case 'all':
          title = '一键安装'
          url = 'api/v1/installK8s/installAll?k8s_cluster_id=' + this.id
          break
        case 'base':
          title = '基础安装'
          url = 'api/v1/installK8s/installBase?k8s_cluster_id=' + this.id
          break
        case 'kernel':
          title = '升级内核'
          url = 'api/v1/installK8s/updateKernel?k8s_cluster_id=' + this.id
          break
        case 'dns':
          title = '安装Dns'
          url = 'api/v1/installK8s/installDns?k8s_cluster_id=' + this.id
          break
        case 'basebin':
          title = '安装可执行文件'
          url = 'api/v1/installK8s/InstallBaseBin?k8s_cluster_id=' + this.id
          break
        case 'containerd':
          title = '安装Containerd'
          url = 'api/v1/installK8s/InstallContainerd?k8s_cluster_id=' + this.id
          break
        case 'registry':
          title = '安装私有镜像库'
          url = 'api/v1/installK8s/installRegistry?k8s_cluster_id=' + this.id
          break
        case 'etcd':
          title = '安装Etcd集群'
          url = 'api/v1/installK8s/installEtcd?k8s_cluster_id=' + this.id
          break
        case 'master':
          title = '安装Master集群'
          url = 'api/v1/installK8s/installMaster?k8s_cluster_id=' + this.id
          break
        case 'node':
          title = '安装Node集群'
          url = 'api/v1/installK8s/installNode?k8s_cluster_id=' + this.id
          break
        case 'registrycrt':
          title = '安装私有镜像仓库证书'
          url = 'api/v1/installK8s/InstallContainerdCrt?k8s_cluster_id=' + this.id
          break
        case 'lvs':
          title = '安装Lvs'
          url = 'api/v1/installK8s/installLvs?k8s_cluster_id=' + this.id
          break
        case 'finish':
          title = '完成安装'
          url = 'api/v1/installK8s/finishInstall?k8s_cluster_id=' + this.id
          break
        case 'newnode':
          title = '安装新节点'
          url = 'api/v1/installK8s/newnodeInstall?k8s_cluster_id=' + this.id
          break
        case 'newetcd':
          title = '安装新Etcd'
          url = 'api/v1/installK8s/newetcdInstall?k8s_cluster_id=' + this.id
          break
        case 'newmaster':
          title = '安装新Master'
          url = 'api/v1/installK8s/newmasterInstall?k8s_cluster_id=' + this.id
          break
        case 'sslmaster':
          title = '更新Master证书'
          url = 'api/v1/installK8s/updateSslMaster?k8s_cluster_id=' + this.id
          break
        case 'ssletcd':
          title = '更新Etcd证书'
          url = 'api/v1/installK8s/updateSslEtcd?k8s_cluster_id=' + this.id
          break
        case 'sslnode':
          title = '更新Node证书'
          url = 'api/v1/installK8s/updateSslNode?k8s_cluster_id=' + this.id
          break
        case 'start':
          title = this.serviceData[this.service]
          url = 'api/v1/installK8s/' + this.service + '?k8s_cluster_id=' + this.id + '&do_what=start&do_containerd=' + this.doContainerd
          break
        case 'restart':
          title = this.serviceData[this.service]
          url = 'api/v1/installK8s/' + this.service + '?k8s_cluster_id=' + this.id + '&do_what=restart&do_containerd=' + this.doContainerd
          break
        case 'stop':
          title = this.serviceData[this.service]
          url = 'api/v1/installK8s/' + this.service + '?k8s_cluster_id=' + this.id + '&do_what=stop&do_containerd=' + this.doContainerd
          break
        default:
          return
      }
      this.initInstallLog(title)
      setTimeout(() => {
        let obj = stream(url)
        this.addEventListener(obj)
      }, 500)
    },

    initInstallLog(desc) {
      this.listResourceDialog = false
      this.installLogDialog = true
      this.installLog = []
      this.installLogDesc = desc
    },

    openedInstallLogDialog() {
      if (this.term === null) {
        this.initXterm()
      } else {
        this.term.clear()
      }
    },

    addEventListener(obj) {
      obj.addEventListener('message', (event) => {
        this.loading = true
        let msg = event.data
        // 参考 gin sse-encode.go 处理
        msg = msg.replaceAll("\\r", "\r").replaceAll("\ndata:", "\n")
        // console.log(JSON.stringify(msg))
        this.term.writeln(msg)
      })
      obj.addEventListener('close', () => {
        this.loading = false
        obj.close()
      })
    },

  },

  mounted() {
    this.listK8sCluster()
    // let url = 'ws://***********'
    // this.init(url)
  },
}
</script>

<style>
.k8sCluster {
  margin-bottom: 50px;
}
.k8sCluster .el-table .warning-row {
  background: aliceblue;
}
.k8sCluster .el-table .success-row {
  background: antiquewhite;
}
.resource .el-table .warning-row {
  background: lavenderblush;
}
.resource .el-table .success-row {
  background: ivory;
}
.infinite-list-wrapper .list-item {
  display: flex;
}
.infinite-list-wrapper .list {
  padding: 0;
  margin: 0;
  list-style: none;
}
.infinite-list-wrapper .list .loading {
  text-align: center;
  font-size: xx-large;
}
.install-k8s-button button {
  margin-top: 10px;
}
</style>
