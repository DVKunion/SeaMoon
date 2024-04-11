declare namespace Service {
  type Proxy = {
    id: number,
    name: string,
    type: string,
    status: number,
    status_message: string,
    listen_address: string,
    listen_port: string
    conn: number,
    speed_up: number,
    speed_down: number,
    lag: number,
    in_bound: number,
    out_bound: number,
    created_at: string,
    updated_at: string,
  }
}
