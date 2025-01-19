export class FunctionCallError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'FunctionCallError'
  }
}
