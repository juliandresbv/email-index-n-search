import type { Email } from './email.type'

interface EmailSearch {
  hits: number
  emails: Email[]
}

export interface EmailsSearchResponse {
  data: EmailSearch
}
