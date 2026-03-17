<template>
  <a-layout class="dashboard-layout">
    <a-layout-sider
      v-model:collapsed="collapsed"
      :trigger="null"
      collapsible
      class="dashboard-sider"
    >
      <div class="logo">Poprako</div>
      <a-menu
        v-model:selectedKeys="selectedMenuKeys"
        theme="dark"
        mode="inline"
        :items="menuItems"
        @click="handleMenuClick"
      />
    </a-layout-sider>

    <a-layout>
      <a-layout-header class="dashboard-header">
        <div class="header-left">
          <a-button type="text" @click="toggleSider">
            <template #icon>
              <MenuUnfoldOutlined v-if="collapsed" />
              <MenuFoldOutlined v-else />
            </template>
          </a-button>
          <span class="header-title">高阶协作仪表盘</span>
        </div>
        <a-space>
          <a-button :loading="loading" @click="refreshAllData"
            >刷新数据</a-button
          >
          <a-button danger @click="logout">退出登录</a-button>
        </a-space>
      </a-layout-header>

      <a-layout-content class="dashboard-content">
        <a-row :gutter="16" class="summary-row">
          <a-col
            v-for="summaryCard in summaryCards"
            :key="summaryCard.label"
            :xs="24"
            :md="12"
            :xl="6"
          >
            <a-card class="summary-card" :bordered="false">
              <a-statistic
                :title="summaryCard.label"
                :value="summaryCard.value"
              />
            </a-card>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :xl="14">
            <a-card title="我的汉化组" :bordered="false" class="panel-card">
              <a-table
                :columns="teamColumns"
                :data-source="teams"
                row-key="id"
                :pagination="false"
                size="middle"
              />
            </a-card>
          </a-col>
          <a-col :xs="24" :xl="10">
            <a-card title="我的分配" :bordered="false" class="panel-card">
              <a-list
                :data-source="assignments"
                :locale="{ emptyText: '暂无分配' }"
              >
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta :description="`章节：${item.chapter_id}`">
                      <template #title>
                        <a-space>
                          <span>{{ item.id }}</span>
                          <a-tag :color="getRoleTagColor(item.role)">{{
                            item.role
                          }}</a-tag>
                        </a-space>
                      </template>
                    </a-list-item-meta>
                  </a-list-item>
                </template>
              </a-list>
            </a-card>
          </a-col>
        </a-row>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { message, type MenuProps, type TableColumnsType } from "ant-design-vue";
import {
  AppstoreOutlined,
  BarsOutlined,
  CloudUploadOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TeamOutlined,
} from "@ant-design/icons-vue";
import { getMyAssignments, getMyTeams } from "../api/modules";
import { useAuthStore } from "../stores/auth";
import type { AssignmentInfo, TeamInfo } from "../types/domain";

const router = useRouter();
const authStore = useAuthStore();
const collapsed = ref(false);
const loading = ref(false);
const teams = ref<TeamInfo[]>([]);
const assignments = ref<AssignmentInfo[]>([]);
const selectedMenuKeys = ref<string[]>(["dashboard"]);

const menuItems: MenuProps["items"] = [
  {
    key: "dashboard",
    label: "仪表盘",
    icon: () => h(AppstoreOutlined),
  },
  {
    key: "team",
    label: "团队概览",
    icon: () => h(TeamOutlined),
  },
  {
    key: "assignment",
    label: "任务分配",
    icon: () => h(BarsOutlined),
  },
  {
    key: "file-test",
    label: "文件测试",
    icon: () => h(CloudUploadOutlined),
  },
];

const teamColumns: TableColumnsType<TeamInfo> = [
  {
    title: "团队名称",
    dataIndex: "name",
    key: "name",
  },
  {
    title: "描述",
    dataIndex: "description",
    key: "description",
    customRender: ({ value }) => value || "-",
  },
  {
    title: "创建时间",
    dataIndex: "created_at",
    key: "created_at",
    customRender: ({ value }) => formatDate(value),
  },
];

/**
 * 生成首页统计卡片数据。
 */
const summaryCards = computed(() => [
  {
    label: "我的团队数",
    value: teams.value.length,
  },
  {
    label: "我的分配数",
    value: assignments.value.length,
  },
  {
    label: "待处理任务",
    value: assignments.value.filter(
      (assignmentInfo) => assignmentInfo.role !== "done",
    ).length,
  },
  {
    label: "已完成任务",
    value: assignments.value.filter(
      (assignmentInfo) => assignmentInfo.role === "done",
    ).length,
  },
]);

/**
 * 切换侧边栏折叠状态。
 */
function toggleSider(): void {
  collapsed.value = !collapsed.value;
}

/**
 * 处理菜单点击行为。
 * 当前示例页面仅保留仪表盘入口，其他菜单通过提示说明。
 */
function handleMenuClick(menuInfo: { key: string }): void {
  selectedMenuKeys.value = [menuInfo.key];

  if (menuInfo.key === "dashboard") {
    router.push("/dashboard");
    return;
  }

  if (menuInfo.key === "file-test") {
    router.push("/file-test");
    return;
  }

  message.info("该菜单页可按同样模式继续扩展");
}

/**
 * 拉取仪表盘依赖数据。
 * 该方法串行加载团队与分配数据，保证页面统计一致。
 */
async function refreshAllData(): Promise<void> {
  loading.value = true;
  try {
    const [myTeams, myAssignments] = await Promise.all([
      getMyTeams(),
      getMyAssignments({
        offset: 0,
        limit: 20,
      }),
    ]);
    teams.value = myTeams;
    assignments.value = myAssignments;
    message.success("数据已刷新");
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "加载失败";
    message.error(errorMessage);
  } finally {
    loading.value = false;
  }
}

/**
 * 执行退出登录并返回登录页。
 */
async function logout(): Promise<void> {
  authStore.clearAccessToken();
  await router.push("/login");
}

/**
 * 根据角色返回不同标签颜色。
 */
function getRoleTagColor(role: string): string {
  if (role === "done") {
    return "green";
  }
  if (role === "reviewer") {
    return "blue";
  }
  return "orange";
}

/**
 * 格式化时间字符串。
 */
function formatDate(rawTime: string | undefined): string {
  if (!rawTime) {
    return "-";
  }
  const date = new Date(rawTime);
  if (Number.isNaN(date.getTime())) {
    return "-";
  }
  return date.toLocaleString("zh-CN", { hour12: false });
}

/**
 * 页面挂载时主动拉取数据。
 */
onMounted(async () => {
  await refreshAllData();
});
</script>

<style scoped lang="scss">
/* 仪表盘基础布局。 */
.dashboard-layout {
  min-height: 100vh;
}

/* 侧边栏背景统一走全局变量。 */
.dashboard-sider {
  background: linear-gradient(
    180deg,
    var(--dashboard-sider-start),
    var(--dashboard-sider-end)
  );
}

.logo {
  height: 64px;
  color: var(--logo-text);
  font-size: 20px;
  font-weight: 700;
  display: grid;
  place-items: center;
  letter-spacing: 0.08em;
}

.dashboard-header {
  background: var(--dashboard-header-bg);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--dashboard-header-border);
  padding: 0 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
}

.dashboard-content {
  padding: 20px;
}

.summary-row {
  margin-bottom: 16px;
}

.summary-card {
  border-radius: 14px;
  background: var(--surface);
}

.panel-card {
  border-radius: 16px;
  margin-bottom: 16px;
  background: var(--panel-bg);
}

:global(html[data-theme="dark"]) {
  /* 暗黑模式颜色全部来自全局变量，页面仅声明文字行为。 */
  .summary-card,
  .panel-card {
    background: var(--panel-bg);
    color: var(--text-primary);
  }

  .header-title {
    color: var(--text-primary);
  }
}
</style>
