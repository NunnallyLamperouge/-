import { useUserStore } from "@/stores/user";
import { get } from "@/utility/request";
import { message } from "ant-design-vue";
import { createRouter, createWebHistory } from "vue-router";
import KTBGTemplateEditor from "../views/tests/KTBGTemplateEditor/KTBGTemplateEditor.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: () => import("@/views/Home.vue"),
      meta: {
        title: "首页 - 研途智伴",
      },
    },
    {
      path: "/auth/login",
      name: "login",
      component: () => import("@/views/auth/Login.vue"),
      meta: {
        title: "登录 - 研途智伴",
      },
    },
    {
      path: "/auth/register",
      name: "register",
      component: () => import("@/views/auth/Register.vue"),
      meta: {
        title: "注册 - 研途智伴",
      },
    },
    {
      path: "/features/document-editor",
      name: "documentEditor",
      component: () => import("@/views/features/DocumentEditor.vue"),
      meta: {
        title: "文档编辑器 - 研途智伴",
        requiresAuth: true,
      },
    },
    {
      path: "/student",
      component: () => import("@/layouts/StudentLayout.vue"),
      children: [
        {
          path: "dashboard",
          component: () => import("@/views/student/Dashboard.vue"),
        },
        {
          path: "modeling",
          component: () => import("@/views/student/Modeling.vue"),
        },
        {
          path: "editor",
          component: () => import("@/views/student/Editing.vue"),
        },
        {
          path: "chat",
          component: () => import("@/views/ChatGroups.vue"),
        },{
          path: "select-teacher",
          component: () => import("@/views/student/Select.vue"),
        },{
          path: "notifications",
          component: () => import("@/views/Notifications.vue"),
        },{
          path: "format",
          component: () => import("@/views/student/Virtual.vue"),
        },{
          path: "missions",
          component: () => import("@/views/student/Missions.vue"),
        },{
          path:"template/history/:id",
          component: () => import('@/views/student/StudentTemplateView/StudentTemplate.vue')
        },{
          path:"documents",
          component: () => import('@/views/student/Documents.vue')
        }
      ],
    },
    {
      path: "/template/student/:id",
      component: () =>
        import(
          "@/views/student/StudentTemplateEditor/StudentTemplateEditor.vue"
        ),
    },
    {
      path: "/teacher",
      component: () => import("@/views/teacher/TeacherLayout.vue"),
      children: [
        {
          path: "dashboard",
          component: () => import("@/views/teacher/StudentManage.vue"),
        },
        {
          path: "chat",
          component: () => import("@/views/ChatGroups.vue"),
        },{
          path: "templates",
          component: () => import("@/views/student/Editing.vue"),
        },{
          path: "notifications",
          component: () => import("@/views/Notifications.vue"),
        },{
          path: "review",
          component: () => import("@/views/teacher/Review.vue"),
        },{
          path: "missions",
          component: () => import("@/views/teacher/Mission.vue"),
        }
        
      ],
    },{
      path: "/view-template/:id",
      component: () => import("@/views/teacher/TeacherTemplateView/TeacherTemplate.vue"),
      meta: {
        title: "查看学生论文 - 研途智伴",
        requiresAuth: true,
      },
    },
    {
      path: "/template",
      component: () =>
        import(
          "@/views/TemplateEditor/UniversalTemplateEditor/UniversalTemplateEditor.vue"
        ),
    },
    {
      path: "/test",
      children: [
        {
          path: "write1",
          component: () => import("@/views/student/templates/rws.vue"),
        },
        {
          path: "write2",
          component: () =>
            import("@/views/tests/RWSTemplateEditor/CustomTest.vue"),
        },
        {
          path: "write3",
          component: () => import("@/views/student/templates/ktbg.vue"),
        },
        {
          path: "write4",
          component: () =>
            import("@/views/tests/KTBGTemplateEditor/KTBGTemplateEditor.vue"),
          meta: {
            title: "开题报告模板编辑器",
          },
        },
      ],
    },
    {
      path: "/ktbg-template-editor",
      name: "KTBGTemplateEditor",
      component: KTBGTemplateEditor,
      meta: {
        title: "开题报告模板编辑器",
      },
    },
  ],
});

// 路由守卫逻辑
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // 处理无需登录的路由
  if (to.path === "/" || to.path.startsWith("/auth")) {
    next();
    return;
  }

  // 处理需要管理员权限的路由
  if (to.matched.some((record) => record.meta.requiresAdmin)) {
    if (userStore.isLogin && userStore.user?.permission >= 1) {
      next();
    } else {
      message.info("此页面需要管理员权限，请使用管理员账号登录！");
      setTimeout(() => {
        next("/auth/login");
      }, 2000);
    }
    return;
  }

  // 处理普通需登录路由
  if (userStore.isLogin) {
    next();
  } else {
    get(
      "/api/user/me",
      {},
      (message, data) => {
        userStore.login(data);
        next();
      },
      (messager, data) => {
        message.info("您还尚未登录，请先登录！");
        setTimeout(() => {
          next("/auth/login");
        }, 2000);
      },
      (messager, data) => {
        message.info("您还尚未登录，请先登录！");
        setTimeout(() => {
          next("/auth/login");
        }, 2000);
      }
    );
  }
});
export default router;
