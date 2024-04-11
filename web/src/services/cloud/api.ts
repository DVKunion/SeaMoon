import request from '@/services/request'
import {FormValueType} from "@/pages/provider/components/CreateForm";

export async function getCloudProvider(page: number, size: number) {
  return request<{
    success: boolean;
    data: Cloud.Provider[];
  }>('/api/v1/provider?page=' + page + '&size=' + size, {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function getActiveProvider() {
  return request<{
    success: boolean;
    data: Cloud.Provider[];
  }>('/api/v1/provider/active', {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function createProvider(data: FormValueType) {
  const params = {}
  params["cloud_auth"] = {
    "access_key": data.access_key,
    "access_secret": data.access_secret,
    "access_id": data.access_id,
    "token": data.token,
    "kube_config": data.kube_config,
  }

  params["name"] = data.name
  params["type"] = data.type
  params["desc"] = data.desc
  params["status"] = data.status
  params["regions"] = typeof data.regions === "string" ? [data.regions] : data.regions

  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Cloud.Provider[];
  }>('/api/v1/provider', {
    method: 'POST',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
    data: params,
  });
}

export async function updateProvider(data: FormValueType) {
  const params = {};

  params["cloud_auth"] = {};

  if (data.access_id !== undefined && data.access_id !== "") {
    params["cloud_auth"]["access_id"] = data.access_id;
  }

  if (data.access_key !== undefined && data.access_key !== "") {
    params["cloud_auth"]["access_key"] = data.access_key;
  }

  if (data.access_secret !== undefined && data.access_secret !== "") {
    params["cloud_auth"]["access_secret"] = data.access_secret;
  }

  if (data.token !== undefined && data.token !== "") {
    params["cloud_auth"]["token"] = data.token;
  }

  if (data.kube_config !== undefined && data.kube_config !== "") {
    params["cloud_auth"]["kube_config"] = data.kube_config;
  }

  params["id"] = data.id
  params["name"] = data.name
  params["status"] = data.status
  params["desc"] = data.desc
  params["type"] = data.type
  params["regions"] = data.regions

  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Cloud.Provider[];
  }>('/api/v1/provider/' + data.id + "/", {
    method: 'PUT',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
    data: params,
  });
}

export async function deleteProvider(id: number) {
  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Cloud.Provider[];
  }>('/api/v1/provider/' + id + '/', {
    method: 'DELETE',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function syncProvider(id: number) {
  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Cloud.Provider[];
  }>('/api/v1/provider/sync/' + id + '/', {
    method: 'PUT',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}
