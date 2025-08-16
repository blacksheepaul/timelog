import { format, parseISO, formatDuration, intervalToDuration } from 'date-fns'

export const formatDateTime = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'yyyy-MM-dd HH:mm:ss')
  } catch {
    return dateString
  }
}

export const formatDate = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'yyyy-MM-dd')
  } catch {
    return dateString
  }
}

export const formatTime = (dateString: string): string => {
  try {
    return format(parseISO(dateString), 'HH:mm:ss')
  } catch {
    return dateString
  }
}

export const calculateDuration = (startTime: string, endTime?: string): string => {
  try {
    const start = parseISO(startTime)
    const end = endTime ? parseISO(endTime) : new Date()
    
    const duration = intervalToDuration({ start, end })
    
    return formatDuration(duration, {
      format: ['hours', 'minutes'],
      locale: { formatDistance: () => '' }
    }) || '0 minutes'
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