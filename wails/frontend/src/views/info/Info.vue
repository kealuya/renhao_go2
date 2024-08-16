<template>
 <div>
   <el-row>

     的方式发生的
     <el-button @click="go2">go2</el-button>


   </el-row>
 </div>
</template>
<script setup lang="ts">

import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
let go2 = () => {
  router.push('/info/index2')
}
</script>

<style scoped>


</style>