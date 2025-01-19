type Storage = {
  token: string
  address: string
}

export const storage = {
  get<T extends keyof Storage>(key: T): Storage[T] | undefined {
    try {
      return JSON.parse(localStorage.getItem(key)!)
    } catch {
      return localStorage.getItem(key) as Storage[T]
    }
  },
  set<T extends keyof Storage>(key: T, value: Storage[T] | undefined) {
    if (value) {
      localStorage.setItem(key, JSON.stringify(value))
    }
  },
  remove<T extends keyof Storage>(key: T) {
    localStorage.removeItem(key)
  },
}
