<script setup>
defineProps({
  id: {
    type: String,
    required: true,
  },
  type: {
    type: String,
    required: true,
  },
  stroke: {
    type: String,
    required: false,
    default: '#b1b1b7',
  },
  fill: {
    type: String,
    required: false,
    default: '#b1b1b7',
  },
  strokeWidth: {
    type: Number,
    required: false,
    default: 0.8,
  },
  width: {
    type: Number,
    required: false,
    default: 10,
  },
  height: {
    type: Number,
    required: false,
    default: 10,
  },
})
</script>

<template>
  <svg class="vue-flow__marker vue-flow__container">
    <defs>
      <marker
        :id="id"
        class="vue-flow__arrowhead"
        viewBox="-10 -10 20 20"
        :refX="type === 'many' ? 5 : 4"
        refY="0"
        :markerWidth="width"
        :markerHeight="height"
        markerUnits="strokeWidth"
        orient="auto-start-reverse"
      >
        <!-- 多关系标记 (many) - 乌鸦脚图标 -->
        <g v-if="type === 'many'">
          <path
            :style="{
              stroke,
              strokeWidth,
            }"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
            d="M 0,0 L 4,4 M 0,0 L 4,-4"
          />
        </g>

        <!-- 一关系标记 (one) -->
        <g v-if="type === 'one'">
          <path
            :style="{
              stroke,
              strokeWidth,
            }"
            stroke-linecap="round"
            stroke-linejoin="round"
            fill="none"
            d="M 0,-4 L 0,4"
          />
        </g>
      </marker>
    </defs>
  </svg>
</template>
