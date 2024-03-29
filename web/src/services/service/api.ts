import request from "@/services/request";
import {FormValueType} from "@/pages/service/components/CreateForm";
import {toNumber} from "lodash";

export async function getServiceProxy(page: number, size: number) {
  return request<{
    success: boolean;
    data: Service.Proxy[];
  }>('/api/v1/proxy?page=' + page + '&size=' + size, {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function createServiceProxy(data: FormValueType) {
  const params = {
    "name": data.name,
    "type": data.type,
    "listen_address": data.listen_address,
    "listen_port": data.listen_port,
    "status": data.status,
  }
  // 说明是 tunnel 关联
  if (data.tunnel_id !== undefined && data.tunnel_id !== 0) {
    params["tunnel_id"] = data.tunnel_id;
  }
  // 说明是 provider 关联
  if (data.provider_id !== undefined && data.provider_id !== 0) {
    params["tunnel_id"] = 0
    params["tunnel_create_api"] = {
      "provider_id": data.provider_id,
      "port": data.port.toString(),
      "name": data.tunnel_name,
      "type": data.tunnel_type,
      "status": 1,
      "tunnel_config": {
        "cpu": toNumber(data.cpu),
        "memory": toNumber(data.memory),
        "instance": toNumber(data.instance),
        "tunnel_auth_type": data.tunnel_auth_type,
        "tls": data.tls,
        "tor": data.tor,
      }
    }
  }

  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Service.Proxy[];
  }>('/api/v1/proxy', {
    method: 'POST',
    data: params,
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function updateServiceProxy(data: FormValueType) {
  const params = {
    "id": data.id,
    "type": data.type,
    "listen_address": data.listen_address,
    "listen_port": data.listen_port,
    "status": data.status,
  }
  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Service.Proxy[];
  }>('/api/v1/proxy/' + data.id + '/', {
    method: 'PUT',
    data: params,
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function deleteServiceProxy(data: FormValueType) {
  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Service.Proxy[];
  }>('/api/v1/proxy/' + data.id + '/', {
    method: 'DELETE',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function speedServiceProxy(data: FormValueType) {
  return request<{
    success: boolean;
    msg?: string;
    code?: number;
    data: Service.Proxy[];
  }>('/api/v1/proxy/speed/' + data.id + '/', {
    method: 'GET',
    data: data,
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}
