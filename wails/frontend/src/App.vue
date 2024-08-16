<template>
  <el-row>
    <el-col :span="4" style="background-color: aqua;">
      <div>
        <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png"/>
      </div>
      <router-link to="/">
        <div class="grid-content ep-bg-purple">home</div>
      </router-link>

      <router-link to="/info">
        <div class="grid-content ep-bg-purple">33333</div>
      </router-link>

    </el-col>
    <el-col :span="16">
      <div class="grid-content ep-bg-purple">
        <router-view/>
      </div>
    </el-col>
  </el-row>
</template>

<script lang="ts" setup>
</script>

<style>
#logo {
  display: block;
  width: 70px;
  height: 50px;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
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
  height: 100vh;
}

.grid-content {
  border-radius: 4px;
  min-height: 36px;
}


</style>
