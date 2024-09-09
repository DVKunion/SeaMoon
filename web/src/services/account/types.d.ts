declare namespace Account {
  type Admin = {
    name: string
    type: string
    last_addr: string
    last_time: string
  }
  type Tunnel = {
    name: string
    type: number
  }
  type Provider = {
    id: number,
    created_at: string,
    updated_at: string,
    name: string,
    desc: string,
    type: number,
    regions: string[],
    info: {
      amount: number,
      cost: number,
    }
    status: number,
    status_message: string,
    count: number,
    max_limit: number
  }
}
