import { createRouter, createWebHashHistory } from 'vue-router'
import Relation from '@/views/Relation.vue'
import ProjectSelect from '@/views/ProjectSelect.vue'
import FlowCanvas from '@/components/flow/FlowCanvas.vue'
import NewProject from '@/views/NewProject.vue' // 新增导入

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'projectSelect',
      component: ProjectSelect
    },
    {
      path: '/new-project',
      name: 'newProject',
      component: NewProject
    },
    {
      path: '/app',
      name: 'app',
      component: Relation,
      children: [
        {
          path: ':pageId',
          name: 'flowPage',
          component: FlowCanvas,
          props: true
        }
      ]
    }
  ]
})

export default router
