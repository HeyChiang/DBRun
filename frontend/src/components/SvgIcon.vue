<template>
  <div>
    <component 
      class="svg-icon-wrapper" 
      :is="currentIcon" 
      v-if="currentIcon" 
      :style="{ width: width, height: height }"
    />
  </div>
</template>

<script setup>
import { shallowRef, watchEffect } from 'vue';

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  disabled: {
    type: Boolean,
    default: false
  },
  width: {
    type: String,
    default: '1.3rem'
  },
  height: {
    type: String,
    default: '1.3rem'
  }
});

const currentIcon = shallowRef(null);

// Import icons using relative paths
const modules = import.meta.glob('../assets/icons/*.svg', { eager: true });

watchEffect(() => {
  try {
    const iconPath = props.disabled
        ? `../assets/icons/${props.name}-disabled.svg`
        : `../assets/icons/${props.name}.svg`;

    if (modules[iconPath]) {
      currentIcon.value = modules[iconPath].default;
    } else {
      console.error(`Icon not found: ${iconPath}`);
    }
  } catch (error) {
    console.error(`Failed to load icon: ${props.name}, error:`, error);
  }
});
</script>

<style>
.svg-icon-wrapper {
  overflow: hidden;
  outline: none; /* Remove outline on click */
  border: none;  /* Remove border on click */
}

.disabled {
  opacity: 0.5; /* Example style for disabled state */
}
</style>
