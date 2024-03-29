import request from '@/services/request'
import {FormValueType} from "@/pages/function/components/CreateForm";
import {toNumber} from "lodash";

export async function getServerlessTunnel(page: number, size: number) {
  return request<{
    success: boolean;
    data: Serverless.Tunnel[];
  }>('/api/v1/tunnel?page=' + page + '&size=' + size, {
    method: 'GET',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function createFunctionTunnel(data: FormValueType) {
  const params = {
    "provider_id": data.provider_id,
    "name": data.tunnel_name,
    "port":data.port,
    "type": data.tunnel_type,
    "status": 1,
    "tunnel_config": {
      "region": data.region,
      "cpu": toNumber(data.cpu),
      "memory": toNumber(data.memory),
      "instance": toNumber(data.instance),
      "tunnel_auth_type": data.tunnel_auth_type,
      "tls": data.tls,
      "tor": data.tor,
    }
  }
  return request<{
    success: boolean;
    data: Serverless.Tunnel[];
  }>('/api/v1/tunnel', {
    method: 'POST',
    data: params,
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function updateFunctionTunnel(data: FormValueType) {
  return request<{
    success: boolean;
    data: Serverless.Tunnel[];
  }>('/api/v1/tunnel/' + data.id + "/", {
    method: 'PUT',
    data: data,
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}

export async function deleteFunctionTunnel(id: number | undefined) {
  return request<{
    success: boolean;
    data: Serverless.Tunnel[];
  }>('/api/v1/tunnel/' + id + "/", {
    method: 'DELETE',
    headers: {
      'Authorization': localStorage.getItem("token") || "",
    },
  });
}
