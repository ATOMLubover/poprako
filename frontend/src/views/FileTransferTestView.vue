<template>
  <div class="file-test-page">
    <a-page-header
      title="文件上传下载测试"
      sub-title="使用 swagger 中真实接口完成预留上传、确认上传与下载验证"
      @back="goBack"
    />

    <a-row :gutter="16">
      <a-col :xs="24" :xl="12">
        <a-card
          title="上传测试（对接后端）"
          class="panel-card"
          :bordered="false"
        >
          <a-space direction="vertical" style="width: 100%" :size="12">
            <a-radio-group v-model:value="targetType">
              <a-radio-button value="user">用户头像</a-radio-button>
              <a-radio-button value="team">团队头像</a-radio-button>
            </a-radio-group>

            <a-select
              v-if="targetType === 'team'"
              v-model:value="selectedTeamID"
              placeholder="请选择团队"
              :options="teamOptions"
            />

            <a-input
              v-model:value="uploadUrl"
              placeholder="点击“获取上传URL”后自动填充"
              readonly
            />

            <a-button
              :loading="reserving"
              :disabled="!canReserve"
              @click="reserveUploadUrl"
            >
              获取上传URL
            </a-button>

            <a-upload
              :before-upload="beforeSelectFile"
              :file-list="uploadFileList"
              :max-count="1"
              @remove="removeSelectedFile"
            >
              <a-button> 选择待上传文件 </a-button>
            </a-upload>

            <a-alert
              v-if="selectedFileName"
              type="info"
              :message="`已选择文件：${selectedFileName}`"
              show-icon
            />

            <a-button
              type="primary"
              :loading="uploading"
              :disabled="!canUpload"
              @click="uploadAndConfirm"
            >
              上传并确认
            </a-button>
          </a-space>
        </a-card>
      </a-col>

      <a-col :xs="24" :xl="12">
        <a-card
          title="下载测试（对接后端）"
          class="panel-card"
          :bordered="false"
        >
          <a-space direction="vertical" style="width: 100%" :size="12">
            <a-input
              v-model:value="downloadUrl"
              placeholder="头像 URL（来自 users/mine 或 teams/mine）"
              readonly
            />

            <a-button :loading="refreshing" @click="refreshTargetData">
              刷新目标数据
            </a-button>

            <a-input
              v-model:value="downloadFileName"
              placeholder="下载文件名（可选，默认 download.bin）"
            />

            <a-space>
              <a-button
                type="primary"
                :loading="downloading"
                :disabled="!downloadUrl.trim()"
                @click="downloadFile"
              >
                下载文件
              </a-button>
              <a-button :disabled="!downloadUrl.trim()" @click="openInNewTab">
                新标签页打开
              </a-button>
            </a-space>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { message, type UploadFile, type UploadProps } from "ant-design-vue";
import {
  confirmTeamAvatarUploaded,
  confirmUserAvatarUploaded,
  getCurrentUserProfile,
  getMyTeams,
  reserveTeamAvatar,
  reserveUserAvatar,
} from "../api/modules";

const router = useRouter();

type UploadTargetType = "user" | "team";

const uploadUrl = ref("");
const downloadUrl = ref("");
const downloadFileName = ref("download.bin");

const targetType = ref<UploadTargetType>("user");
const selectedTeamID = ref("");

const currentUserID = ref("");
const currentUserAvatarURL = ref("");
const teamList = ref<Array<{ id: string; name: string; avatar_url?: string }>>(
  [],
);

const selectedFile = ref<File | null>(null);
const uploadFileList = ref<UploadFile[]>([]);
const reserving = ref(false);
const uploading = ref(false);
const downloading = ref(false);
const refreshing = ref(false);

const selectedFileName = computed(() => selectedFile.value?.name ?? "");
const teamOptions = computed(() =>
  teamList.value.map((teamInfo) => ({
    label: `${teamInfo.name} (${teamInfo.id})`,
    value: teamInfo.id,
  })),
);
const canReserve = computed(() => {
  if (selectedFile.value === null) {
    return false;
  }

  if (targetType.value === "user") {
    return currentUserID.value.length > 0;
  }

  return selectedTeamID.value.length > 0;
});
const canUpload = computed(
  () => uploadUrl.value.trim().length > 0 && selectedFile.value !== null,
);

const selectedTeamAvatarURL = computed(() => {
  const targetTeam = teamList.value.find(
    (teamInfo) => teamInfo.id === selectedTeamID.value,
  );
  return targetTeam?.avatar_url ?? "";
});

function goBack(): void {
  router.push("/dashboard");
}

function syncDownloadUrl(): void {
  if (targetType.value === "user") {
    downloadUrl.value = currentUserAvatarURL.value;
    return;
  }

  downloadUrl.value = selectedTeamAvatarURL.value;
}

async function refreshTargetData(): Promise<void> {
  refreshing.value = true;
  try {
    const [currentUserProfile, myTeams] = await Promise.all([
      getCurrentUserProfile(),
      getMyTeams({ offset: 0, limit: 50 }),
    ]);

    currentUserID.value = currentUserProfile.id;
    currentUserAvatarURL.value =
      (currentUserProfile as { avatar_url?: string }).avatar_url ?? "";

    teamList.value = myTeams.map((teamInfo) => ({
      id: teamInfo.id,
      name: teamInfo.name,
      avatar_url: (teamInfo as { avatar_url?: string }).avatar_url,
    }));

    if (!selectedTeamID.value && teamList.value.length > 0) {
      selectedTeamID.value = teamList.value[0].id;
    }

    syncDownloadUrl();
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : "刷新数据失败";
    message.error(errorMessage);
  } finally {
    refreshing.value = false;
  }
}

async function reserveUploadUrl(): Promise<void> {
  if (selectedFile.value === null) {
    message.warning("请先选择待上传文件");
    return;
  }

  const contentType =
    selectedFile.value.type.trim() || "application/octet-stream";

  reserving.value = true;
  try {
    if (targetType.value === "user") {
      const reserveResult = await reserveUserAvatar(currentUserID.value, {
        content_type: contentType,
      });
      uploadUrl.value = reserveResult.put_url;
      message.success("已获取用户头像上传 URL");
      return;
    }

    const reserveResult = await reserveTeamAvatar(selectedTeamID.value, {
      content_type: contentType,
    });
    uploadUrl.value = reserveResult.put_url;
    message.success("已获取团队头像上传 URL");
  } catch (error) {
    const errorMessage =
      error instanceof Error ? error.message : "获取上传 URL 失败";
    message.error(errorMessage);
  } finally {
    reserving.value = false;
  }
}

const beforeSelectFile: UploadProps["beforeUpload"] = (file) => {
  selectedFile.value = file;
  uploadFileList.value = [file as UploadFile];
  return false;
};

function removeSelectedFile(): boolean {
  selectedFile.value = null;
  uploadFileList.value = [];
  return true;
}

function buildCorsHint(rawURL: string): string {
  let targetOrigin = "未知域名";

  try {
    targetOrigin = new URL(rawURL).origin;
  } catch {
    // 忽略 URL 解析错误，保持默认提示文本。
  }

  return `浏览器未能访问 ${targetOrigin}，请在对象存储桶 CORS 中放行当前前端来源（Origin）、PUT/GET/OPTIONS 方法及必要请求头`;
}

async function uploadAndConfirm(): Promise<void> {
  if (!selectedFile.value || !uploadUrl.value.trim()) {
    message.warning("请先填写上传 URL 并选择文件");
    return;
  }

  uploading.value = true;
  try {
    const response = await fetch(uploadUrl.value.trim(), {
      method: "PUT",
      headers: {
        "Content-Type": selectedFile.value.type || "application/octet-stream",
      },
      body: selectedFile.value,
    });

    if (!response.ok) {
      throw new Error(`上传失败，HTTP ${response.status}`);
    }

    if (targetType.value === "user") {
      await confirmUserAvatarUploaded(currentUserID.value);
    } else {
      await confirmTeamAvatarUploaded(selectedTeamID.value);
    }

    await refreshTargetData();
    selectedFile.value = null;
    uploadFileList.value = [];
    message.success("上传并确认成功");
  } catch (error) {
    if (
      error instanceof TypeError &&
      error.message.includes("Failed to fetch")
    ) {
      message.error(buildCorsHint(uploadUrl.value.trim()));
      return;
    }

    const errorMessage = error instanceof Error ? error.message : "上传失败";
    message.error(errorMessage);
  } finally {
    uploading.value = false;
  }
}

async function downloadFile(): Promise<void> {
  if (!downloadUrl.value.trim()) {
    message.warning("请先填写下载 URL");
    return;
  }

  downloading.value = true;
  try {
    const response = await fetch(downloadUrl.value.trim(), {
      method: "GET",
    });

    if (!response.ok) {
      throw new Error(`下载失败，HTTP ${response.status}`);
    }

    const fileBlob = await response.blob();
    const objectUrl = URL.createObjectURL(fileBlob);

    const link = document.createElement("a");
    link.href = objectUrl;
    link.download = downloadFileName.value.trim() || "download.bin";
    link.click();

    URL.revokeObjectURL(objectUrl);

    message.success("下载成功");
  } catch (error) {
    if (
      error instanceof TypeError &&
      error.message.includes("Failed to fetch")
    ) {
      message.error(buildCorsHint(downloadUrl.value.trim()));
      return;
    }

    const errorMessage = error instanceof Error ? error.message : "下载失败";
    message.error(errorMessage);
  } finally {
    downloading.value = false;
  }
}

function openInNewTab(): void {
  if (!downloadUrl.value.trim()) {
    message.warning("请先填写下载 URL");
    return;
  }

  window.open(downloadUrl.value.trim(), "_blank", "noopener,noreferrer");
}

onMounted(async () => {
  await refreshTargetData();
});

watch(targetType, () => {
  syncDownloadUrl();
});

watch(selectedTeamID, () => {
  if (targetType.value === "team") {
    syncDownloadUrl();
  }
});
</script>

<style scoped lang="scss">
.file-test-page {
  min-height: 100vh;
  padding: 16px;
}

.panel-card {
  border-radius: 12px;
}
</style>
