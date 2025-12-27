<template>
  <Dialog 
    v-model:visible="visible" 
    modal 
    header="反馈建议" 
    :style="{ width: '50vw' }"
    :breakpoints="{ '960px': '75vw', '641px': '90vw' }"
    :closable="!isSubmitting"
  >
    <form @submit.prevent="submitFeedback" class="feedback-form">
      <div class="form-group">
        <label for="title">问题标题</label>
        <InputText id="title" v-model="feedback.title" required class="w-full" />
      </div>
      
      <div class="form-group">
        <label for="content">问题内容</label>
        <Textarea
          id="content"
          v-model="feedback.content"
          rows="6"
          required
          class="w-full"
          placeholder="请详细描述您遇到的问题，包括：
1. 问题发生的具体场景
2. 重现问题的具体步骤
3. 期望的正确行为
4. 实际发生的行为"
        />
      </div>

      <div class="form-group no-label">
        <FileUpload
          :multiple="true"
          accept="image/*"
          :maxFileSize="5000000"
          @select="onFileSelect"
          @remove="onFileRemove"
          :auto="true"
          chooseLabel="添加截图"
          :chooseIcon="'pi pi-images'"
          cancelLabel="取消"
          :customUpload="true"
          @uploader="customUpload"
          :showUploadButton="false"
          :showCancelButton="false"
          invalidFileSizeMessage="文件大小不能超过 5MB"
          invalidFileTypeMessage="只能上传图片文件"
          :class="{ 'p-error': imageError }"
        >
          <template #empty>
            <div class="upload-placeholder">
              <i class="pi pi-images text-2xl mb-2 text-gray-600"></i>
              <p class="upload-text">点击添加或拖放图片到此处</p>
              <p class="upload-hint">最多可上传3张截图（jpg、png），每张不超过5MB</p>
            </div>
          </template>
        </FileUpload>
        <small v-if="imageError" class="p-error block mt-1">{{ imageError }}</small>
      </div>
      
      <div class="form-group">
        <label for="contact">联系方式</label>
        <InputText 
          id="contact" 
          v-model="feedback.contact"
          class="w-full"
          placeholder="请留下您的邮箱或其他联系方式，以便我们回复您"
        />
      </div>

      <div class="dialog-footer">
        <Button 
          label="取消" 
          @click="closeDialog" 
          text 
          :disabled="isSubmitting"
        />
        <Button 
          type="submit" 
          label="提交反馈" 
          :loading="isSubmitting"
          :disabled="isSubmitting"
        />
      </div>
    </form>
  </Dialog>
  <Toast position="top-center" />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';
import Button from 'primevue/button';
import FileUpload from 'primevue/fileupload';
import Toast from 'primevue/toast';
import { useToast } from 'primevue/usetoast';

const visible = ref(false);
const feedback = ref({
  title: '',
  content: '',
  contact: '',
  images: [] as File[]
});
const imageError = ref('');
const isSubmitting = ref(false);
const toast = useToast();

const closeDialog = () => {
  if (isSubmitting.value) return;
  
  visible.value = false;
  feedback.value = {
    title: '',
    content: '',
    contact: '',
    images: []
  };
  imageError.value = '';
};

const onFileSelect = (event: any) => {
  const files = event.files;
  
  // 检查文件总数是否超过3张
  if (feedback.value.images.length + files.length > 3) {
    imageError.value = '最多只能上传3张图片';
    return;
  }

  // 检查每个文件是否为图片
  for (const file of files) {
    if (!file.type.startsWith('image/')) {
      imageError.value = '只能上传图片文件';
      return;
    }
  }

  imageError.value = '';
  feedback.value.images.push(...files);
};

const onFileRemove = (event: any) => {
  const removedFile = event.file;
  feedback.value.images = feedback.value.images.filter(
    file => file.name !== removedFile.name
  );
  imageError.value = '';
};

const customUpload = (event: any) => {
  // 这里不做实际上传，只是将文件保存到本地状态中
  // 实际提交时再一起处理
};

const submitFeedback = async () => {
  try {
    isSubmitting.value = true;

    // 创建 FormData 对象
    const formData = new FormData();
    formData.append('title', feedback.value.title);
    formData.append('content', feedback.value.content);
    formData.append('contact', feedback.value.contact);

    // 添加图片文件
    feedback.value.images.forEach((file, index) => {
      formData.append(`image${index + 1}`, file);
    });

    // 发送请求
    const response = await fetch('example.com/feedback', {
      method: 'POST',
      body: formData
    });

    if (!response.ok) {
      throw new Error('提交失败，请稍后重试');
    }

    const result = await response.json();
    
    // 显示成功提示
    toast.add({
      severity: 'success',
      summary: '提交成功',
      detail: '感谢您的反馈，我们会尽快处理',
      life: 3000
    });

    closeDialog();
  } catch (error) {
    // 显示错误提示
    toast.add({
      severity: 'error',
      summary: '提交失败',
      detail: error instanceof Error ? error.message : '请稍后重试',
      life: 5000
    });
  } finally {
    isSubmitting.value = false;
  }
};

// 对外暴露打开弹窗的方法
defineExpose({
  open: () => {
    visible.value = true;
  }
});
</script>

<style scoped>
.feedback-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group.no-label {
  margin-top: -0.5rem;
}

label {
  font-weight: 500;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 1rem;
}

:deep(.p-fileupload) {
  width: 100%;
}

:deep(.p-fileupload-content) {
  padding: 1.5rem;
  border-style: dashed;
  border-radius: 8px;
  border-color: var(--surface-border);
  background-color: var(--surface-ground);
  transition: all 0.2s;
}

:deep(.p-fileupload-content:hover) {
  border-color: var(--primary-color);
  background-color: var(--surface-hover);
}

:deep(.p-fileupload.p-error .p-fileupload-content) {
  border-color: var(--red-500);
}

:deep(.p-button.p-fileupload-choose) {
  background: var(--primary-color);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  text-align: center;
}

.upload-text {
  color: var(--text-color);
  font-size: 0.95rem;
  margin: 0.5rem 0;
}

.upload-hint {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
  margin: 0;
}

:deep(.p-fileupload-row) {
  margin: 0.5rem 0;
}

:deep(.p-fileupload-filename) {
  color: var(--text-color);
  font-size: 0.9rem;
}
</style>
