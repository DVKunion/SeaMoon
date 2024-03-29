declare namespace Auth {
  type User = {
    ID: number,
    Username: string,
    Password: string
  }

  type Response = {
    code?: number,
    msg?: string,
    success: boolean,
    total?: number,
    data?: any,
  }
}

