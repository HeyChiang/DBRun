<template>
  <Dialog
    v-model:visible="visible"
    modal
    header="关于"
    :style="{ width: '40vw' }"
    :breakpoints="{ '960px': '75vw', '641px': '90vw' }"
  >
    <div class="about-content">
      <h2>DBRun</h2>
      <p class="version">版本 1.0.0</p>
      <div class="description">
        <p>DBRun 是一个强大的数据库关系可视化工具，帮助您更好地理解和管理数据库结构。</p>
      </div>
      <div class="features">
        <h3>主要特性：</h3>
        <ul>
          <li>数据库表关系可视化</li>
          <li>支持多种数据库类型</li>
          <li>直观的用户界面</li>
          <li>高效的性能表现</li>
        </ul>
      </div>
      <div class="copyright">
        <p>© 2024 DBRun. 保留所有权利。</p>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, defineExpose } from 'vue';
import Dialog from 'primevue/dialog';

const visible = ref(false);

const show = () => {
  visible.value = true;
};

defineExpose({
  show
});
</script>

<style scoped>
.about-content {
  padding: 1rem;
}

.version {
  color: #666;
  margin-bottom: 1.5rem;
}

.description {
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.features {
  margin-bottom: 2rem;
}

.features h3 {
  margin-bottom: 1rem;
}

.features ul {
  list-style-type: none;
  padding-left: 0;
}

.features li {
  margin-bottom: 0.5rem;
  padding-left: 1.5rem;
  position: relative;
}

.features li::before {
  content: "•";
  position: absolute;
  left: 0;
  color: var(--primary-color);
}

.copyright {
  margin-top: 2rem;
  padding-top: 1rem;
  border-top: 1px solid var(--surface-border);
  font-size: 0.9rem;
  color: #666;
}
</style>
