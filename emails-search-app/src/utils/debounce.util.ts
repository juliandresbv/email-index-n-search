export const debounce = <F extends (...args: any[]) => any | Function>(
  fn: F,
  timeout: number
): ((...args: Parameters<F>) => Promise<ReturnType<F>>) => {
  let timeoutId: NodeJS.Timeout

  return (...args: Parameters<F>) => {
    clearTimeout(timeoutId)

    return new Promise<ReturnType<F>>((resolve, reject) => {
      timeoutId = setTimeout(async () => {
        try {
          const result = await fn(...args)

          resolve(result)
        } catch (error) {
          reject(error)
        }
      }, timeout)
    })
  }
}
