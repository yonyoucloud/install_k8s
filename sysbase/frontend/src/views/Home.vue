<template>
  <div class="home">
    <el-row :gutter="20">
      <el-col :span="12">
        <div class="grid-content bg-purple">
          <Resource />
          <K8sCluster />
          <Pod />
        </div>
      </el-col>
      <el-col :span="12">
        <div class="grid-content bg-purple">
          <Tenant />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import Tenant from '@/components/Tenant.vue'
import Resource from '@/components/Resource.vue'
import K8sCluster from '@/components/K8sCluster.vue'
import Pod from '@/components/Pod.vue'

export default {
  name: 'Home',
  components: {
    Tenant,
    Resource,
    K8sCluster,
    Pod
  }
}
</script>

<style>
.el-row {
  margin-bottom: 20px;
  &:last-child {
    margin-bottom: 0;
  }
}
.el-col {
  border-radius: 4px;
}
.bg-purple-dark {
  background: #99a9bf;
}
.bg-purple {
  border: 1px solid #ebebeb;
  border-radius: 3px;
  transition: .2s;
}
.bg-purple-light {
  background: #e5e9f2;
}
.grid-content {
  border-radius: 4px;
  min-height: 36px;
  padding: 20px;
  text-align: left;
}
.row-bg {
  padding: 10px 0;
  background-color: #f9fafc;
}
</style>
