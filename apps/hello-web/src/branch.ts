export const cookieName = "dev_branch"

export function validRoute(value: string): boolean {
  return /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/.test(value)
}

export function selectedRoute(cookie = document.cookie): string {
  const prefix = `${cookieName}=`
  const value = cookie
    .split(";")
    .map((part) => part.trim())
    .find((part) => part.startsWith(prefix))
    ?.slice(prefix.length)
  return value && validRoute(value) ? decodeURIComponent(value) : ""
}

export function selectRoute(value: string): void {
  if (!validRoute(value)) throw new Error("invalid route key")
  document.cookie = `${cookieName}=${encodeURIComponent(value)}; Path=/; Max-Age=28800; SameSite=Lax; Secure`
}

export function clearRoute(): void {
  document.cookie = `${cookieName}=; Path=/; Max-Age=0; SameSite=Lax; Secure`
}
