declare namespace Serverless {
  type Tunnel = {
    provider_id: number,
    provider_type: number,
    id: number,
    uniq_id: string,
    created_at: string,
    updated_at: string,
    name: string,
    address: string,
    port: string,
    type: string,
    status: number,
    status_message: string,
    // 健康检查信息
    version: string,
    v2ray_version: string,
    last_check_time: string,
    tunnel_config: {
      cpu: number,
      memory: number,
      instance: number,
      tunnel_auth_type: number,
      region: string,
      tls: false,
      tor: false,
      // 级联代理配置
      cascade_proxy: boolean,
      cascade_tunnel_id: number,
      cascade_addr: string,
    }
  }
}
