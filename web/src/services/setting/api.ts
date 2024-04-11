import request from '@/services/request'

/** 获取 SysConfig 列表 GET /api/v1/config */
export async function getSysConfig() {
  return request<{
    data: Config.SystemConfig;
  }>('/api/v1/config', {
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
    method: 'GET',
  });
}

/** 更新 SysConfig */
export async function updateSysConfig(data: Config.SystemConfig) {
  return request<{
    code?: string;
    msg?: string;
    success: boolean;
    data: Config.SystemConfig;
  }>('/api/v1/config/', {
    method: 'PUT',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
    data: data,
  });
}
