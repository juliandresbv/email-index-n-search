import type { Store } from 'pinia'

import type { EmailsSearchResponse } from '../../types'

export const updateEmailsStoreState = (
  emailsData: EmailsSearchResponse | null,
  emailsStore: Store | any
) => {
  if (!emailsData) {
    emailsStore.$reset()

    return
  }

  const { emails, hits: resHits } = emailsData.data

  if (!emails || emails?.length <= 0 || !resHits || resHits <= 0) {
    emailsStore.$reset()

    return
  }

  emailsStore.setEmails(emails)
  emailsStore.setHits(resHits)
}
