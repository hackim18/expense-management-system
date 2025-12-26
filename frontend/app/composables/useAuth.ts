type AuthUser = {
  id: string
  name: string
  email: string
  role: string
}

type AuthPayload = {
  access_token: string
  id: string
  name: string
  email: string
  role: string
}

export const useAuth = () => {
  const token = useState<string | null>('auth:token', () => null)
  const user = useState<AuthUser | null>('auth:user', () => null)

  const init = () => {
    if (!import.meta.client || token.value) {
      return
    }

    const raw = localStorage.getItem('auth')
    if (!raw) {
      return
    }

    try {
      const parsed = JSON.parse(raw) as { token: string; user: AuthUser }
      token.value = parsed.token
      user.value = parsed.user
    } catch {
      localStorage.removeItem('auth')
    }
  }

  const setAuth = (payload: AuthPayload) => {
    token.value = payload.access_token
    user.value = {
      id: payload.id,
      name: payload.name,
      email: payload.email,
      role: payload.role
    }

    if (import.meta.client) {
      localStorage.setItem(
        'auth',
        JSON.stringify({ token: token.value, user: user.value })
      )
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    if (import.meta.client) {
      localStorage.removeItem('auth')
    }
  }

  const isAuthenticated = computed(() => Boolean(token.value))

  return {
    token,
    user,
    isAuthenticated,
    init,
    setAuth,
    logout
  }
}
