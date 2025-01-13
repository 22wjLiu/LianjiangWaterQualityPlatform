<template>
  <div class="form-container">
    <Program2 class="program"></Program2>
  </div>
</template>

<script>
// import Program1 from '@/components/fortcase/program1.vue'
import Program2 from '@/components/fortcase/program2.vue'
export default {
  data() {
    return {
      program: 'Program1',
      btnType: 'primary',
      btnType2: ''
    }
  },
  methods: {},
  components: {
    Program2
  }
}
</script>

<style lang="less" scoped>
.form-container {
  width: 30%;
  height: 100%;
  overflow: auto;
  border: 1px solid #ccc;
  .program {
    width: 80%;
    padding: 30px;
  }
}
</style>
