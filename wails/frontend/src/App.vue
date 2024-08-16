<template>

  <div>
    <el-row>
      <el-col :span="8" style="background-color: rgba(25, 48, 80, 1);height: 100vh">
        <el-row justify="center">
          <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/>
        </el-row>
        <div style="height: 10%"/>

        <el-row justify="center">
          <button @click="go('/')" class="btn-22" style="width: 70%"><span>Home</span></button>
        </el-row>
        <el-row justify="center">
          <button @click="go('/info')" class="btn-22" style="width: 70%"><span>Info</span></button>
        </el-row>
      </el-col>


      <el-col :span="16">
        <div class="grid-content ep-bg-purple">
          <router-view/>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import {useRouter} from 'vue-router'

const router = useRouter()
let go = (path: string) => {
  router.push(path)
}

</script>

<style>

@import 'style/button.css';

#logo {
  display: block;
  width: 70px;
  height: 50px;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-origin: content-box;
}

.el-row {
  margin-bottom: 20px;
}

.el-row:last-child {
  margin-bottom: 0;
}

.el-col {
  border-radius: 4px;
}

.grid-content {
  min-height: 36px;
}

</style>
