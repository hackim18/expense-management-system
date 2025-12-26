const formatWithDots = (value: string) => {
  const parts: string[] = []
  let remainder = value
  while (remainder.length > 3) {
    parts.unshift(remainder.slice(-3))
    remainder = remainder.slice(0, -3)
  }
  if (remainder) {
    parts.unshift(remainder)
  }
  return parts.join('.')
}

export const useIdr = () => {
  const format = (amount: number) => {
    const safe = Number.isFinite(amount) ? Math.floor(amount) : 0
    return `Rp ${formatWithDots(Math.abs(safe).toString())}`
  }

  const formatInput = (raw: string) => {
    const digits = raw.replace(/\D/g, '')
    const numeric = digits ? parseInt(digits, 10) : 0
    return {
      formatted: digits ? formatWithDots(digits) : '',
      value: numeric
    }
  }

  return { format, formatInput }
}
