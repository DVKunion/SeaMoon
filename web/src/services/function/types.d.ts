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
    tunnel_config: {
      cpu: number,
      memory: number,
      instance: number,
      tunnel_auth_type: number,
      region: string,
      tls: false,
      tor: false
    }
  }
}
