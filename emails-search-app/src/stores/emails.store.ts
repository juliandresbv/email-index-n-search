import { defineStore } from 'pinia'

import type { Email } from '../types/email.type'

export const useEmailsStore = defineStore('emails', {
  state: () => {
    const hits = 0
    const emails: Email[] = []
    const selectedEmail: Email | any = null

    return {
      hits,
      emails,
      selectedEmail: selectedEmail
    }
  },
  actions: {
    setEmails(emails: Email[]) {
      this.emails = emails
    },
    setHits(hits: number) {
      this.hits = hits
    },
    setSelectedEmail(selectedEmail: Email | any) {
      this.selectedEmail = selectedEmail
    }
  }
})
