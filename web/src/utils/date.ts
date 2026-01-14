import { format, parseISO, formatDuration, intervalToDuration } from 'date-fns'

export const formatDateTime = (dateString: string | null | undefined): string => {
  if (!dateString) return 'N/A'
  try {
    return format(parseISO(dateString), 'yyyy-MM-dd HH:mm:ss')
  } catch {
    return dateString
  }
}

export const formatDate = (input: string | Date): string => {
  try {
    const date = typeof input === 'string' ? parseISO(input) : input
    return format(date, 'yyyy-MM-dd')
  } catch {
    return typeof input === 'string' ? input : input.toISOString().split('T')[0]
  }
}

export const formatTime = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'HH:mm:ss')
  } catch {
    return dateString
  }
}

export const calculateDuration = (startTime: string, endTime?: string | null): string => {
  try {
    const start = parseISO(startTime)
    // 当endTime为null时，使用当前时间的UTC字符串，确保与start_time的时区一致
    const end = endTime ? parseISO(endTime) : new Date()

    const duration = intervalToDuration({ start, end })

    // 对于跨月记录，显示更完整的格式
    if (duration.months || (duration.days && duration.days > 7)) {
      const result = formatDuration(duration, {
        format: ['months', 'days', 'hours'],
      })
      return result || '0 hours'
    }

    // 对于普通记录，显示小时和分钟
    const result = formatDuration(duration, {
      format: ['hours', 'minutes'],
    })

    return result || '0 minutes'
  } catch {
    return 'Invalid duration'
  }
}

export const getCurrentDateTime = (): string => {
  return new Date().toISOString()
}

export const formatDateTimeLocal = (date: Date = new Date()): string => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')

  return `${year}-${month}-${day}T${hours}:${minutes}`
}

// 将UTC时间字符串转换为本地datetime-local格式
export const formatUTCToLocal = (utcString: string): string => {
  try {
    const utcDate = new Date(utcString)
    return formatDateTimeLocal(utcDate)
  } catch {
    return ''
  }
}

// 将本地datetime-local值转换为UTC ISO字符串
export const formatLocalToUTC = (localDateTimeString: string): string => {
  try {
    // datetime-local输入的值应该被当作用户的本地时间
    // JavaScript的Date构造函数会正确处理这个转换
    const localDate = new Date(localDateTimeString)
    return localDate.toISOString()
  } catch {
    return new Date().toISOString()
  }
}
